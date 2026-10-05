package publicmarkets

import (
	"context"
	"fmt"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	pm "trovo-wallet-api/internal/components/publicmarkets/models"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// Job names (System Health in Trovo Manager).
const (
	JobOrders         = "Order processing"
	JobInstructions   = "Outbound instruction retry sweep"
	JobBatches        = "Session net batches"
	JobReconciliation = "Custodian reconciliation"
	JobPrices         = "Price ingestion"
	JobWebhooks       = "Exchange webhook retry & dead-letter"
	JobConfirmations  = "Exchange confirmation SLA sweep"
	JobLedger         = "Beneficial ownership ledger"
	JobDividends      = "Dividend distribution"
)

func (e *Engine) jobRan(job, status, detail string) {
	e.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&pm.JobRun{Job: job, LastRunAt: e.now(), Status: status, Detail: trimTo(detail, 300)})
}

// Start runs the engine's loops. Each runs on one instance at a time (a
// distributed lock), so the engine scales out like the rest of app-backend.
func (e *Engine) Start() {
	kick := make(chan struct{}, 1)
	e.kick = kick
	every := func(name string, interval, stale time.Duration, fn func(ctx context.Context)) {
		go func() {
			t := time.NewTicker(interval)
			defer t.Stop()
			for {
				if sharedconfig.ShuttingDown() {
					return
				}
				sharedconfig.WithSingletonLock(e.GC, "pm-"+name, stale, func() {
					ctx, cancel := context.WithTimeout(context.Background(), stale)
					defer cancel()
					fn(ctx)
				})
				if name == "orders" {
					select {
					case <-t.C:
					case <-kick:
					}
				} else {
					<-t.C
				}
			}
		}()
	}
	every("orders", 3*time.Second, 3*time.Minute, func(ctx context.Context) {
		e.DeliverMockEvents(ctx)
		e.ProcessPendingEvents()
		e.ProcessOrders(ctx)
		var open int64
		e.DB.Model(&pm.Order{}).Where("state NOT IN ?", []string{pm.StateComplete, pm.StateRejected, pm.StateFailed, pm.StateCancelled}).Count(&open)
		e.jobRan(JobOrders, "up", fmt.Sprintf("%d open order(s)", open))
	})
	every("instructions", 20*time.Second, 3*time.Minute, func(ctx context.Context) {
		e.DispatchInstructions()
		var esc int64
		e.DB.Model(&pm.Instruction{}).Where("status = ?", pm.InstrEscalated).Count(&esc)
		status := "up"
		if esc > 0 {
			status = "degraded"
		}
		e.jobRan(JobInstructions, status, fmt.Sprintf("%d escalated", esc))
	})
	every("batches", 30*time.Second, 3*time.Minute, func(ctx context.Context) {
		e.ReleaseApprovedBatches()
		e.BuildBatches(ctx)
		e.jobRan(JobBatches, "up", "")
	})
	every("prices", time.Minute, 3*time.Minute, func(ctx context.Context) {
		stale := e.PollPrices(ctx)
		status := "up"
		if stale > 0 {
			status = "degraded"
		}
		e.jobRan(JobPrices, status, fmt.Sprintf("%d stale", stale))
	})
	every("ledger", 30*time.Second, 5*time.Minute, func(ctx context.Context) {
		e.IndexLedgers(ctx)
		e.jobRan(JobLedger, "up", "")
	})
	every("webhooks", 15*time.Second, 3*time.Minute, func(ctx context.Context) {
		e.DeliverWebhooks(ctx)
		var dlq int64
		e.DB.Model(&pm.WebhookDelivery{}).Where("status = ?", pm.DeliveryDeadLetter).Count(&dlq)
		status := "up"
		if dlq > 0 {
			status = "degraded"
		}
		e.jobRan(JobWebhooks, status, fmt.Sprintf("%d dead-lettered", dlq))
	})
	every("hourly", time.Hour, 10*time.Minute, func(ctx context.Context) {
		e.SweepConfirmationSLA()
		var late int64
		e.DB.Model(&pm.WebhookDelivery{}).Where("escalated_at IS NOT NULL AND confirmed_at IS NULL").Count(&late)
		status := "up"
		if late > 0 {
			status = "degraded"
		}
		e.jobRan(JobConfirmations, status, fmt.Sprintf("%d past SLA", late))
	})
	every("daily", time.Minute, 30*time.Minute, func(ctx context.Context) {
		e.runRequests(ctx)
		e.dailyJobs(ctx)
		e.SnapshotDueActions(ctx)
		e.AdvanceDividends(ctx)
		e.SweepFees(ctx)
	})
}

// Kick wakes the order loop now (after a new order).
func (e *Engine) Kick() {
	if e.kick == nil {
		return
	}
	select {
	case e.kick <- struct{}{}:
	default:
	}
}

// dailyJobs runs the reconciliation at its hour (WAT) and the mock
// Custodians' end-of-day position feed, once per day.
func (e *Engine) dailyJobs(ctx context.Context) {
	s := LoadSettings(e.DB)
	now := e.now().In(Lagos)
	today := now.Format("2006-01-02")
	var run pm.JobRun
	if now.Hour() >= s.ReconciliationHour && (e.DB.First(&run, "job = ?", JobReconciliation).Error != nil || run.LastRunAt.In(Lagos).Format("2006-01-02") != today) {
		n := e.ReconcileAll(ctx, "system")
		var drift int64
		e.DB.Model(&pm.ReconciliationRun{}).Where("created_at >= ? AND result <> ?", e.now().Add(-time.Hour), pm.ReconMatched).Count(&drift)
		status := "up"
		if drift > 0 {
			status = "down"
		}
		e.jobRan(JobReconciliation, status, fmt.Sprintf("%d asset(s), %d drift", n, drift))
	}
	var feed pm.JobRun
	if now.Hour() >= 18 && (e.DB.First(&feed, "job = ?", "mock-position-feed").Error != nil || feed.LastRunAt.In(Lagos).Format("2006-01-02") != today) {
		var mocks []pm.Custodian
		e.DB.Where("mode = ?", pm.ModeMock).Find(&mocks)
		for _, c := range mocks {
			if err := e.MockPositionFeed(ctx, c.Code); err != nil {
				e.logf("mock position feed %s: %v", c.Code, err)
			}
		}
		e.jobRan("mock-position-feed", "up", fmt.Sprintf("%d mock Custodian(s)", len(mocks)))
	}
}

// runRequests runs what Trovo Manager asked for: RECONCILE (an asset id or
// ALL), RESUME (reconcile and reopen a halted asset), POSITION_FEED (a mock
// Custodian's feed now) and EXCHANGE_WITHDRAWAL (serviceLinkId:amount).
func (e *Engine) runRequests(ctx context.Context) {
	var reqs []pm.JobRequest
	e.DB.Where("done_at IS NULL").Order("id").Limit(20).Find(&reqs)
	for _, r := range reqs {
		res := e.DB.Model(&pm.JobRequest{}).Where("id = ? AND done_at IS NULL", r.ID).Update("done_at", e.now())
		if res.RowsAffected == 0 {
			continue
		}
		result := "done"
		switch r.Job {
		case "RECONCILE":
			if r.Target == "" || strings.EqualFold(r.Target, "ALL") {
				result = fmt.Sprintf("%d asset(s) reconciled", e.ReconcileAll(ctx, r.RequestedBy))
			} else if a, err := e.AssetByCode(r.Target); err == nil {
				run := e.Reconcile(ctx, a, r.RequestedBy)
				result = run.Result + ": " + run.Detail
			} else {
				result = err.Error()
			}
		case "RESUME":
			if a, err := e.AssetByCode(r.Target); err == nil {
				run, ok := e.ResumeAsset(ctx, a, r.RequestedBy)
				if ok {
					result = "Reconciled and resumed"
				} else {
					result = "Not resumed: " + run.Result + " " + run.Detail
				}
			} else {
				result = err.Error()
			}
		case "POSITION_FEED":
			if err := e.MockPositionFeed(ctx, r.Target); err != nil {
				result = err.Error()
			}
		case "EXCHANGE_WITHDRAWAL":
			parts := strings.SplitN(r.Target, ":", 2)
			if len(parts) == 2 && d(parts[1]).IsPositive() {
				hash, err := e.PayExchangeWithdrawal(ctx, parts[0], d(parts[1]), fmt.Sprintf("withdrawal:%d", r.ID), r.RequestedBy)
				if err != nil {
					result = "Failed: " + err.Error()
				} else {
					result = "Paid (tx " + hash + ")"
				}
			} else {
				result = "Invalid withdrawal request"
			}
		default:
			result = "unknown job"
		}
		e.DB.Model(&pm.JobRequest{}).Where("id = ?", r.ID).Update("result", trimTo(result, 500))
	}
}

// SweepFees moves the Trovo fee of completed app orders from the treasury
// to the fee wallet, in one transfer. (Exchange orders' fees came out of
// the exchange's balance and are swept the same way.)
func (e *Engine) SweepFees(ctx context.Context) {
	s := LoadSettings(e.DB)
	if e.Chain == nil || !common.IsHexAddress(s.FeeWallet) {
		return
	}
	var orders []pm.Order
	e.DB.Where("state = ? AND fee_sweep_tx = '' AND fee <> '0' AND fee <> ''", pm.StateComplete).Limit(500).Find(&orders)
	total := decimal.Zero
	ids := make([]string, 0, len(orders))
	for _, o := range orders {
		total = total.Add(d(o.Fee))
		ids = append(ids, o.ID)
	}
	if len(ids) == 0 || total.LessThan(decimal.NewFromInt(1)) {
		return
	}
	ft, err := e.fundingToken(ctx)
	if err != nil {
		return
	}
	res := e.DB.Model(&pm.Order{}).Where("id IN ? AND fee_sweep_tx = ''", ids).Update("fee_sweep_tx", "pending")
	if res.RowsAffected != int64(len(ids)) {
		e.DB.Model(&pm.Order{}).Where("id IN ? AND fee_sweep_tx = 'pending'", ids).Update("fee_sweep_tx", "")
		return
	}
	hash, err := e.Chain.Send(ctx, e.Chain.TreasurySafe(), []gnosissafe.Call{callOf(aa.ERC20Transfer(ft.Address, common.HexToAddress(s.FeeWallet), units(total.String(), ft.Decimals)))})
	if err != nil {
		e.DB.Model(&pm.Order{}).Where("id IN ?", ids).Update("fee_sweep_tx", "")
		e.logf("fee sweep: %v", err)
		return
	}
	e.DB.Model(&pm.Order{}).Where("id IN ?", ids).Update("fee_sweep_tx", hash)
}
