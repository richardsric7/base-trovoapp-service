package publicmarkets

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	pm "trovo-wallet-api/internal/components/publicmarkets/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/gnosissafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Public Markets dividends run on their own: a corporate action fixes a
// record-date snapshot of every holder from the ownership ledger,
// withholds tax per owner, waits for approval and a funded treasury, then
// pays app holders in the funding stablecoin and credits exchanges (who
// confirm receipt of each wallet's dividend.paid).

// CorporateActionNotice is a Custodian's corporate action (§4.4.3), also
// what Trovo Manager's manual declaration records.
type CorporateActionNotice struct {
	AssetCode     string `json:"assetCode"`
	EventType     string `json:"eventType"`
	RecordDate    string `json:"recordDate"`
	PayDate       string `json:"payDate"`
	AmountPerUnit string `json:"amountPerUnit"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
	Reference     string `json:"reference"`
}

func (e *Engine) onCorporateAction(ev *pm.PartnerEvent) (string, error) {
	var n CorporateActionNotice
	if err := json.Unmarshal([]byte(ev.Payload), &n); err != nil {
		return "", refuse(http.StatusBadRequest, "error-invalid-payload", "", "Unreadable corporate action notice.")
	}
	a, err := e.AssetByCode(n.AssetCode)
	if err != nil {
		return "", err
	}
	if code := strings.TrimPrefix(ev.Source, "CUSTODIAN:"); code != ev.Source && code != e.custodianCode(a.CustodianID) {
		return "", refuse(http.StatusForbidden, "error-wrong-partner", "assetCode", "This asset is not held by you.")
	}
	ca, err := e.DeclareCorporateAction(a, n, pm.SourceCustodianCA, ev.Reference, ev.Source)
	if err != nil {
		return "", err
	}
	return "corporate action " + ca.ID + " " + ca.Status, nil
}

// DeclareCorporateAction records a corporate action. Dividends and coupons
// are distributed automatically; supply-changing events (bonus, rights,
// split) are not designed yet and are marked for manual handling.
func (e *Engine) DeclareCorporateAction(a *pm.Asset, n CorporateActionNotice, source, reference, by string) (*pm.CorporateAction, error) {
	typ := strings.ToUpper(strings.TrimSpace(n.EventType))
	if _, err := time.ParseInLocation("2006-01-02", n.RecordDate, Lagos); err != nil {
		return nil, refuse(http.StatusBadRequest, "error-invalid-field", "recordDate", "recordDate must be YYYY-MM-DD.")
	}
	status := pm.ActionAnnounced
	switch typ {
	case pm.ActionDividend, pm.ActionCoupon:
		if !d(n.AmountPerUnit).IsPositive() {
			return nil, refuse(http.StatusBadRequest, "error-invalid-field", "amountPerUnit", "amountPerUnit must be positive.")
		}
		if c := strings.ToUpper(strings.TrimSpace(n.Currency)); c != "" && c != "NGN" {
			return nil, refuse(http.StatusBadRequest, "error-invalid-field", "currency", "Only NGN distributions are supported.")
		}
	case pm.ActionBonus, pm.ActionRights, pm.ActionSplit:
		status = pm.ActionNeedsManual
	default:
		return nil, refuse(http.StatusBadRequest, "error-invalid-field", "eventType", "Unknown eventType %q.", n.EventType)
	}
	ca := pm.CorporateAction{ID: "CA-" + strings.ToUpper(randomHex(5)), AssetID: a.ID, AssetCode: a.AssetCode, EventType: typ, Description: n.Description,
		RecordDate: n.RecordDate, PayDate: n.PayDate, AmountPerUnit: d(n.AmountPerUnit).String(), Currency: "NGN", Source: source, SourceReference: reference,
		Status: status, DeclaredBy: by, CreatedAt: e.now(), UpdatedAt: e.now()}
	if ca.Description == "" {
		ca.Description = fmt.Sprintf("%s %s", a.ShortName, strings.ToLower(typ))
	}
	if err := e.DB.Create(&ca).Error; err != nil {
		return nil, err
	}
	return &ca, nil
}

// whtPercent is an owner's withholding rate: an exchange customer's from
// the tax data the exchange provisioned them with (no tax ID, or residence
// outside the asset's country, use those rates); an app user's from their
// country on record.
func (e *Engine) whtPercent(a *pm.Asset, wallet string, s pm.Settings) (decimal.Decimal, string) {
	var w pm.PartnerWallet
	if e.DB.Where("LOWER(wallet_address) = ?", strings.ToLower(wallet)).First(&w).Error == nil {
		res := w.ResidencyCountry
		switch {
		case strings.TrimSpace(w.TaxIdentifier) == "":
			return d(s.WHTMissingTaxIDPercent), res
		case !sameCountry(res, a.Country):
			return d(s.WHTNonResidentPercent), res
		}
		return d(s.WHTResidentPercent), res
	}
	var uw userModels.UserWallet
	if e.DB.Where("LOWER(id) = ?", strings.ToLower(wallet)).First(&uw).Error == nil {
		var u userModels.User
		if e.DB.Select("id, country_code").First(&u, "id = ?", uw.UserID).Error == nil && u.CountryCode != nil && *u.CountryCode != "" && !sameCountry(*u.CountryCode, a.Country) {
			return d(s.WHTNonResidentPercent), *u.CountryCode
		}
	}
	return d(s.WHTResidentPercent), a.Country
}

func sameCountry(x, y string) bool {
	norm := func(c string) string {
		c = strings.ToUpper(strings.TrimSpace(c))
		if c == "NGA" {
			return "NG"
		}
		return c
	}
	return norm(x) == norm(y)
}

func (e *Engine) usernameOf(wallet string) string {
	var uw userModels.UserWallet
	if e.DB.Where("LOWER(id) = ?", strings.ToLower(wallet)).First(&uw).Error != nil {
		return ""
	}
	var u userModels.User
	if e.DB.Select("id, username").First(&u, "id = ?", uw.UserID).Error != nil {
		return ""
	}
	return u.Username
}

// recordDateEnd is the end of the record date in Lagos.
func recordDateEnd(date string) (time.Time, error) {
	day, err := time.ParseInLocation("2006-01-02", date, Lagos)
	if err != nil {
		return time.Time{}, err
	}
	return day.AddDate(0, 0, 1).Add(-time.Second), nil
}

// SnapshotDueActions fixes the entitlements of every announced dividend
// whose record date has ended and whose ledger has been indexed past it.
func (e *Engine) SnapshotDueActions(ctx context.Context) {
	var due []pm.CorporateAction
	e.DB.Where("status = ?", pm.ActionAnnounced).Find(&due)
	for i := range due {
		if err := e.snapshot(ctx, &due[i]); err != nil {
			e.logf("snapshot %s: %v", due[i].ID, err)
		}
	}
}

func (e *Engine) snapshot(ctx context.Context, ca *pm.CorporateAction) error {
	end, err := recordDateEnd(ca.RecordDate)
	if err != nil || e.now().Before(end) {
		return err
	}
	if e.Chain == nil {
		return errors.New("no chain configured")
	}
	var a pm.Asset
	if err := e.DB.First(&a, "id = ?", ca.AssetID).Error; err != nil {
		return err
	}
	block, err := e.Chain.BlockAtTime(ctx, end)
	if err != nil {
		return err
	}
	if a.LedgerScannedBlock < block {
		return nil // the ledger has not reached the record date yet
	}
	s := LoadSettings(e.DB)
	balances := e.BalancesAt(&a, block)
	wallets := make([]string, 0, len(balances))
	for w := range balances {
		wallets = append(wallets, w)
	}
	sort.Strings(wallets)
	perUnit := d(ca.AmountPerUnit)
	eligible, retained, gross, wht, net := decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero
	var ents []pm.Entitlement
	h := sha256.New()
	for _, w := range wallets {
		qty := decimal.NewFromBigInt(balances[w], -int32(a.TokenDecimals))
		channel, link := e.channelOf(&a, w)
		en := pm.Entitlement{ID: "ENT-" + strings.ToUpper(randomHex(6)), CorporateActionID: ca.ID, AssetID: a.ID, AssetCode: a.AssetCode,
			WalletAddress: w, Channel: channel, ServiceLinkID: link, Units: qty.String(), CreatedAt: e.now(), UpdatedAt: e.now()}
		if channel == pm.ChannelPlatform {
			// inventory held by the platform: its dividend stays with Trovotech
			en.GrossAmount, en.WHTPercent, en.WHTAmount, en.NetAmount, en.Status = "0", "0", "0", "0", pm.EntitlementRetained
			en.Note = "Platform wallet: retained"
			retained = retained.Add(qty)
		} else {
			g := qty.Mul(perUnit).RoundDown(2)
			pct, res := e.whtPercent(&a, w, s)
			t := g.Mul(pct).Div(decimal.NewFromInt(100)).Round(2)
			en.GrossAmount, en.WHTPercent, en.WHTAmount, en.NetAmount, en.TaxResidency = g.String(), pct.String(), t.String(), g.Sub(t).String(), res
			en.Status = pm.EntitlementPending
			if channel == pm.ChannelExchange {
				var pw pm.PartnerWallet
				e.DB.Select("id").First(&pw, "LOWER(wallet_address) = ?", strings.ToLower(w))
				en.PartnerWalletID = pw.ID
			} else {
				en.Username = e.usernameOf(w)
			}
			eligible, gross, wht, net = eligible.Add(qty), gross.Add(g), wht.Add(t), net.Add(g.Sub(t))
		}
		fmt.Fprintf(h, "%s|%s|%s|%s|%s\n", w, en.Units, en.GrossAmount, en.WHTAmount, en.NetAmount)
		ents = append(ents, en)
	}
	checksum := hex.EncodeToString(h.Sum(nil))
	now := e.now()
	return e.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&pm.CorporateAction{}).Where("id = ? AND status = ?", ca.ID, pm.ActionAnnounced).Updates(map[string]interface{}{
			"status": pm.ActionSnapshotted, "record_block": block, "snapshot_at": &now, "snapshot_checksum": checksum,
			"holder_count": len(ents), "eligible_units": eligible.String(), "retained_units": retained.String(),
			"gross_amount": gross.String(), "wht_amount": wht.String(), "net_amount": net.String(), "approvals_required": s.ApprovalsRequired})
		if res.Error != nil || res.RowsAffected == 0 {
			return res.Error
		}
		if len(ents) == 0 {
			return nil
		}
		return tx.CreateInBatches(ents, 200).Error
	})
}

// AdvanceDividends moves distributions along: approved once enough
// approvers signed the snapshot, paying once the treasury covers it, and
// paid out batch by batch.
func (e *Engine) AdvanceDividends(ctx context.Context) {
	s := LoadSettings(e.DB)
	allowed := CSV(s.DividendApprovers)
	var snapped []pm.CorporateAction
	e.DB.Where("status = ?", pm.ActionSnapshotted).Find(&snapped)
	for _, ca := range snapped {
		var approvals []pm.DividendApproval
		e.DB.Where("corporate_action_id = ? AND checksum = ?", ca.ID, ca.SnapshotChecksum).Find(&approvals)
		n := 0
		for _, ap := range approvals {
			if len(allowed) == 0 || contains(allowed, strings.ToLower(ap.Approver)) {
				n++
			}
		}
		if n >= ca.ApprovalsRequired {
			now := e.now()
			e.DB.Model(&pm.CorporateAction{}).Where("id = ? AND status = ?", ca.ID, pm.ActionSnapshotted).Updates(map[string]interface{}{"status": pm.ActionApproved, "approved_at": &now})
		}
	}
	var approved []pm.CorporateAction
	e.DB.Where("status = ?", pm.ActionApproved).Find(&approved)
	for _, ca := range approved {
		if ca.PayDate != "" {
			if day, err := time.ParseInLocation("2006-01-02", ca.PayDate, Lagos); err == nil && e.now().Before(day) {
				continue // not before the pay date
			}
		}
		if e.Chain == nil {
			continue
		}
		free, err := e.TreasuryFreeCash(ctx)
		if err != nil {
			continue
		}
		// this distribution is itself counted in the obligations
		if free.IsNegative() {
			e.DB.Model(&pm.CorporateAction{}).Where("id = ?", ca.ID).Update("note", fmt.Sprintf("Waiting for the treasury: %s short", free.Neg().StringFixed(2)))
			continue
		}
		now := e.now()
		e.DB.Model(&pm.CorporateAction{}).Where("id = ? AND status = ?", ca.ID, pm.ActionApproved).Updates(map[string]interface{}{"status": pm.ActionPaying, "funded_at": &now, "note": ""})
	}
	var paying []pm.CorporateAction
	e.DB.Where("status = ?", pm.ActionPaying).Find(&paying)
	for i := range paying {
		if err := e.payDistribution(ctx, &paying[i]); err != nil {
			e.logf("dividend %s: %v", paying[i].ID, err)
		}
	}
}

const dividendBatchSize = 100

func (e *Engine) payDistribution(ctx context.Context, ca *pm.CorporateAction) error {
	ft, err := e.fundingToken(ctx)
	if err != nil {
		return err
	}
	var a pm.Asset
	e.DB.First(&a, "id = ?", ca.AssetID)
	// exchanges: credit their balance, then dividend.paid per wallet
	var exch []pm.Entitlement
	e.DB.Where("corporate_action_id = ? AND channel = ? AND status = ?", ca.ID, pm.ChannelExchange, pm.EntitlementPending).Find(&exch)
	for _, en := range exch {
		if err := e.creditExchange(en.ServiceLinkID, d(en.NetAmount), "DIVIDEND", "dividend:"+en.ID, fmt.Sprintf("%s %s for wallet %s", ca.AssetCode, strings.ToLower(ca.EventType), en.PartnerWalletID), "system"); err != nil {
			return err
		}
		now := e.now()
		e.DB.Model(&pm.Entitlement{}).Where("id = ?", en.ID).Updates(map[string]interface{}{"status": pm.EntitlementPaid, "paid_at": &now, "note": "Credited to the exchange balance"})
		e.QueueWebhook(en.ServiceLinkID, "dividend.paid", en.ID, en.PartnerWalletID, ca.AssetCode, map[string]interface{}{
			"event": "dividend.paid", "walletId": en.PartnerWalletID, "assetCode": ca.AssetCode, "eventType": ca.EventType,
			"amount": d(en.NetAmount).StringFixed(2), "grossAmount": d(en.GrossAmount).StringFixed(2), "currency": "NGN",
			"withholdingTaxApplied": d(en.WHTAmount).StringFixed(2), "withholdingTaxPercent": en.WHTPercent,
			"units": en.Units, "amountPerUnit": ca.AmountPerUnit, "recordDate": ca.RecordDate, "payDate": ca.PayDate,
			"corporateActionId": ca.ID, "entitlementId": en.ID, "timestamp": e.now().Format(time.RFC3339)}, true)
	}
	// app holders: a confirmed batch is settled; then the next one is sent
	var queued []pm.Entitlement
	e.DB.Where("corporate_action_id = ? AND status = ?", ca.ID, pm.EntitlementQueued).Find(&queued)
	if len(queued) > 0 {
		hash := queued[0].BatchTxHash
		if hash == "" || hash == "pending" {
			if e.now().Sub(queued[0].UpdatedAt) > 10*time.Minute {
				e.DB.Model(&pm.Entitlement{}).Where("corporate_action_id = ? AND status = ? AND batch_tx_hash IN ?", ca.ID, pm.EntitlementQueued, []string{"", "pending"}).
					Updates(map[string]interface{}{"status": pm.EntitlementFailed, "note": "The payment batch was not recorded; check the treasury before retrying"})
			}
			return nil
		}
		now := e.now()
		if err := e.Chain.Wait(ctx, hash); err != nil {
			e.DB.Model(&pm.Entitlement{}).Where("batch_tx_hash = ?", hash).Updates(map[string]interface{}{"status": pm.EntitlementFailed, "note": trimTo(err.Error(), 300)})
			return err
		}
		e.DB.Model(&pm.Entitlement{}).Where("batch_tx_hash = ?", hash).Updates(map[string]interface{}{"status": pm.EntitlementPaid, "paid_at": &now})
		for _, en := range queued {
			if en.Username != "" && e.Notify != nil {
				e.Notify(en.Username, fmt.Sprintf("%s %s paid", strings.TrimSuffix(ca.AssetCode, "-T"), strings.ToLower(ca.EventType)),
					fmt.Sprintf("₦%s (net of WHT) was paid into your wallet.", d(en.NetAmount).StringFixed(2)),
					map[string]string{"route": "publicMarketsDividend", "entitlementId": en.ID, "assetCode": ca.AssetCode})
			}
		}
		return nil
	}
	var pending []pm.Entitlement
	e.DB.Where("corporate_action_id = ? AND channel = ? AND status = ?", ca.ID, pm.ChannelApp, pm.EntitlementPending).Limit(dividendBatchSize).Find(&pending)
	if len(pending) > 0 {
		var calls []gnosissafe.Call
		ids := make([]string, 0, len(pending))
		for _, en := range pending {
			if amt := units(en.NetAmount, ft.Decimals); amt.Sign() > 0 {
				calls = append(calls, callOf(aa.ERC20Transfer(ft.Address, common.HexToAddress(en.WalletAddress), amt)))
			}
			ids = append(ids, en.ID)
		}
		res := e.DB.Model(&pm.Entitlement{}).Where("id IN ? AND status = ?", ids, pm.EntitlementPending).Updates(map[string]interface{}{"status": pm.EntitlementQueued, "batch_tx_hash": "pending", "updated_at": e.now()})
		if res.RowsAffected != int64(len(ids)) {
			return errors.New("entitlements changed while queuing")
		}
		if len(calls) == 0 {
			now := e.now()
			e.DB.Model(&pm.Entitlement{}).Where("id IN ?", ids).Updates(map[string]interface{}{"status": pm.EntitlementPaid, "paid_at": &now, "batch_tx_hash": "", "note": "Nothing to pay"})
			return nil
		}
		hash, err := e.Chain.Send(ctx, e.Chain.TreasurySafe(), calls)
		if err != nil {
			e.DB.Model(&pm.Entitlement{}).Where("id IN ?", ids).Updates(map[string]interface{}{"status": pm.EntitlementPending, "batch_tx_hash": ""})
			return err
		}
		e.DB.Model(&pm.Entitlement{}).Where("id IN ?", ids).Update("batch_tx_hash", hash)
		return nil
	}
	// withheld tax: one transfer to the WHT wallet
	if ca.WHTTxHash == "" && d(ca.WHTAmount).IsPositive() {
		s := LoadSettings(e.DB)
		if !common.IsHexAddress(s.WHTWallet) {
			e.DB.Model(&pm.CorporateAction{}).Where("id = ?", ca.ID).Update("note", "Set the withholding tax wallet in Public Markets settings to remit the tax")
			return nil
		}
		res := e.DB.Model(&pm.CorporateAction{}).Where("id = ? AND wht_tx_hash = ''", ca.ID).Update("wht_tx_hash", "pending")
		if res.RowsAffected == 0 {
			return nil
		}
		hash, err := e.Chain.Send(ctx, e.Chain.TreasurySafe(), []gnosissafe.Call{callOf(aa.ERC20Transfer(ft.Address, common.HexToAddress(s.WHTWallet), units(ca.WHTAmount, ft.Decimals)))})
		if err != nil {
			e.DB.Model(&pm.CorporateAction{}).Where("id = ?", ca.ID).Update("wht_tx_hash", "")
			return err
		}
		e.DB.Model(&pm.CorporateAction{}).Where("id = ?", ca.ID).Update("wht_tx_hash", hash)
		return e.Chain.Wait(ctx, hash)
	}
	if ca.WHTTxHash == "pending" {
		return nil
	}
	var open int64
	e.DB.Model(&pm.Entitlement{}).Where("corporate_action_id = ? AND status IN ?", ca.ID, []string{pm.EntitlementPending, pm.EntitlementQueued, pm.EntitlementFailed}).Count(&open)
	if open == 0 {
		now := e.now()
		e.DB.Model(&pm.CorporateAction{}).Where("id = ? AND status = ?", ca.ID, pm.ActionPaying).Updates(map[string]interface{}{"status": pm.ActionDistributed, "completed_at": &now})
	}
	return nil
}
