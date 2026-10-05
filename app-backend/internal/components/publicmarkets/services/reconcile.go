package publicmarkets

import (
	"context"
	"fmt"
	"strings"
	"time"

	pm "trovo-wallet-api/internal/components/publicmarkets/models"

	"github.com/shopspring/decimal"
)

// Reconcile checks one asset, fail-closed:
//
//   - token supply must not exceed the units the Custodian holds (FR-G1):
//     every token is backed 1:1. Units held beyond the supply are the
//     working inventory that fast-path creations fill from;
//   - Σ BeneficialOwnershipLedger must equal the on-chain supply (FR-G3),
//     which catches internal drift rather than Custodian-side drift.
//
// A mismatch halts creation and redemption for that asset only; holders'
// balances and P2P trading are not affected.
func (e *Engine) Reconcile(ctx context.Context, a *pm.Asset, triggeredBy string) pm.ReconciliationRun {
	run := pm.ReconciliationRun{AssetID: a.ID, AssetCode: a.AssetCode, TriggeredBy: triggeredBy, CreatedAt: e.now()}
	// refresh the position from the Custodian where it can be asked
	var c pm.Custodian
	if e.DB.First(&c, "custodian_id = ?", a.CustodianID).Error == nil && e.Partners != nil {
		if client, err := e.Partners.Custodian(c); err == nil && client != nil {
			if p, err := client.Position(ctx, a.AssetCode); err == nil {
				asOf := p.AsOf
				if asOf.IsZero() {
					asOf = e.now()
				}
				e.DB.Create(&pm.CustodianPosition{AssetID: a.ID, RealUnitsHeld: p.UnitsHeld.String(), AsOf: asOf, Source: pm.PositionPull,
					Reference: "position query (" + triggeredBy + ")", CreatedAt: e.now()})
			} else {
				run.Detail = "Position query failed: " + trimTo(err.Error(), 200) + ". "
			}
		}
	}
	supply, err := e.Supply(ctx, a)
	if err != nil {
		run.Result, run.Detail = pm.ReconError, run.Detail+"Reading the token supply failed: "+trimTo(err.Error(), 200)
		e.DB.Create(&run)
		return run
	}
	run.TokenSupply = supply.String()
	pos, ok := e.Position(a)
	if !ok {
		if supply.IsZero() {
			run.Result, run.Detail = pm.ReconMatched, run.Detail+"No tokens issued and no position yet."
		} else {
			run.Result, run.Detail = pm.ReconNoPosition, run.Detail+"Tokens are issued but the Custodian has reported no position."
		}
	} else {
		units := d(pos.RealUnitsHeld)
		asOf := pos.AsOf
		run.CustodianPosition, run.PositionAsOf, run.PositionSource = units.String(), &asOf, pos.Source
		run.Inventory = decimal.Max(units.Sub(supply), decimal.Zero).String()
		if supply.GreaterThan(units) {
			run.Delta = supply.Sub(units).String()
			run.Result = pm.ReconDrift
			run.Detail += fmt.Sprintf("Token supply %s is %s above the Custodian's position %s (as of %s, %s).", supply, run.Delta, units, asOf.In(Lagos).Format("Jan 2 15:04"), strings.ToLower(pos.Source))
			if late := e.ordersSince(a, asOf); late != "" {
				run.Detail += " Settled since: " + late + "."
			}
		} else {
			run.Result = pm.ReconMatched
		}
	}
	// FR-G3: the ledger against the chain, once the ledger has caught up
	ledger := e.LedgerTotal(a)
	run.LedgerTotal = ledger.String()
	if e.Chain != nil && a.ContractAddress != "" {
		if head, err := e.Chain.Head(ctx); err == nil && a.LedgerScannedBlock+ledgerConfirmations+10 >= head {
			if !ledger.Equal(supply) {
				run.LedgerDelta = ledger.Sub(supply).String()
				if run.Result == pm.ReconMatched {
					run.Result = pm.ReconLedgerDrift
				}
				run.Detail += fmt.Sprintf(" Σ ledger %s differs from on-chain supply %s.", ledger, supply)
			}
		} else {
			run.Detail += " Ledger still catching up; FR-G3 check skipped."
		}
	}
	p := e.CurrentPrice(a)
	run.PriceStale = p.Session.Open && !p.Fresh
	if run.Result == pm.ReconDrift || run.Result == pm.ReconLedgerDrift || run.Result == pm.ReconNoPosition {
		if a.Status == pm.AssetLive {
			now := e.now()
			e.DB.Model(&pm.Asset{}).Where("id = ? AND status = ?", a.ID, pm.AssetLive).Updates(map[string]interface{}{
				"status": pm.AssetHalted, "halt_reason": trimTo("Reconciliation: "+run.Detail, 500), "halted_at": &now, "halted_by": "reconciliation"})
			run.Halted = true
			if e.GC != nil {
				e.GC.LogDiscordFailedRequest(fmt.Sprintf("[publicmarkets] %s halted by reconciliation: %s", a.AssetCode, run.Detail))
			}
		}
	}
	run.Detail = strings.TrimSpace(run.Detail)
	e.DB.Create(&run)
	return run
}

// ordersSince lists creations settled after the position's as-of time, the
// usual cause of a timing-only drift (a settlement after the Custodian's
// end-of-day cut-off).
func (e *Engine) ordersSince(a *pm.Asset, asOf time.Time) string {
	var orders []pm.Order
	e.DB.Where("asset_id = ? AND type = ? AND settled_at > ?", a.ID, pm.OrderCreation, asOf).Limit(5).Find(&orders)
	var parts []string
	for _, o := range orders {
		parts = append(parts, fmt.Sprintf("%s (+%s)", o.ID, o.Quantity))
	}
	return strings.Join(parts, ", ")
}

// ReconcileAll reconciles every asset with tokens or a position.
func (e *Engine) ReconcileAll(ctx context.Context, triggeredBy string) int {
	var assets []pm.Asset
	e.DB.Where("status IN ?", []string{pm.AssetLive, pm.AssetHalted}).Find(&assets)
	for i := range assets {
		e.Reconcile(ctx, &assets[i], triggeredBy)
	}
	return len(assets)
}

// ResumeAsset reopens a halted asset, only when a fresh reconciliation
// matches.
func (e *Engine) ResumeAsset(ctx context.Context, a *pm.Asset, by string) (pm.ReconciliationRun, bool) {
	if a.Status != pm.AssetHalted {
		return pm.ReconciliationRun{}, false
	}
	// reconcile as if live, without halting again
	a.Status = pm.AssetSetup
	run := e.Reconcile(ctx, a, by)
	if run.Result != pm.ReconMatched {
		return run, false
	}
	e.DB.Model(&pm.Asset{}).Where("id = ? AND status = ?", a.ID, pm.AssetHalted).Updates(map[string]interface{}{"status": pm.AssetLive, "halt_reason": "", "halted_at": nil, "halted_by": ""})
	return run, true
}
