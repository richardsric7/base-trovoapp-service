// Package engine is payout-engine's work loop. Each tick it finishes any
// Safe transaction left in flight, prepares requested schedules, checks
// requested funding, pays payouts in batches and notifies paid holders. All
// state lives in app-backend's database; tm-api changes it on admin actions
// and wakes the engine over Redis.
package engine

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"trovo-payout-engine/internal/bus"
	"trovo-payout-engine/internal/config"
	"trovo-payout-engine/internal/push"
	"trovo-payout-engine/internal/store"

	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Engine holds the engine's dependencies.
type Engine struct {
	DB      *gorm.DB
	Chain   *ethclient.Client
	ChainID *big.Int
	Cfg     *config.Config
	Signers config.SignerSource
	Push    push.Sender // nil: notifications are skipped
	Bus     *bus.Bus
}

// staleAfter is how long another instance's heartbeat may be silent before
// this one takes over.
const staleAfter = 90 * time.Second

// Run works until ctx ends. Only one instance works at a time: the others
// wait as standbys until its heartbeat goes stale.
func (e *Engine) Run(ctx context.Context) {
	wake := e.Bus.Wakeups(ctx)
	for ctx.Err() == nil {
		if e.claim() {
			if halted, reason := e.halted(); halted {
				e.activity("halted: " + reason)
			} else if err := e.Tick(ctx); err != nil {
				log.Printf("[engine] %v", err)
				e.recordError(err)
			}
		}
		select {
		case <-ctx.Done():
		case <-wake:
		case <-time.After(e.Cfg.PollInterval):
		}
	}
}

// Tick runs one round of work.
func (e *Engine) Tick(ctx context.Context) error {
	if err := e.recoverBatches(ctx); err != nil {
		return fmt.Errorf("recovering batches: %w", err)
	}
	var payouts []store.ProceedPayout
	if err := e.DB.Where("status IN ?", []string{store.StatusPrepareRequested, store.StatusPreparing, store.StatusFundingCheckRequested, store.StatusPaying}).
		Order("id").Find(&payouts).Error; err != nil {
		return err
	}
	for i := range payouts {
		if ctx.Err() != nil {
			return nil
		}
		if halted, _ := e.halted(); halted {
			return nil
		}
		p := &payouts[i]
		var err error
		switch p.Status {
		case store.StatusPrepareRequested, store.StatusPreparing:
			err = e.prepare(ctx, p)
		case store.StatusFundingCheckRequested:
			err = e.checkFunding(ctx, p)
		case store.StatusPaying:
			err = e.pay(ctx, p)
		}
		if err != nil {
			msg := fmt.Sprintf("payout %d (%s): %v", p.ID, p.Status, err)
			log.Printf("[engine] %s", msg)
			e.recordError(fmt.Errorf("%s", msg))
			e.note(p.ID, err.Error())
			e.publish(ctx, "error", p.ID, err.Error())
		}
	}
	if err := e.notifyPaid(ctx); err != nil {
		log.Printf("[engine] notifications: %v", err)
	}
	if err := e.sweep(ctx); err != nil {
		log.Printf("[engine] sweep: %v", err)
	}
	return nil
}

// claim takes or renews the single active-instance slot (the state row's
// instance + heartbeat), reporting whether this instance holds it.
func (e *Engine) claim() bool {
	e.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&store.EngineState{ID: 1})
	now := time.Now().UTC()
	res := e.DB.Model(&store.EngineState{}).
		Where("id = 1 AND (instance = ? OR instance = '' OR instance IS NULL OR heartbeat_at IS NULL OR heartbeat_at < ?)", e.Cfg.Instance, now.Add(-staleAfter)).
		Updates(map[string]interface{}{"instance": e.Cfg.Instance, "heartbeat_at": now, "version": e.Cfg.Version})
	held := res.Error == nil && res.RowsAffected == 1
	if held {
		e.Bus.Heartbeat(context.Background(), e.Cfg.Instance, staleAfter)
	}
	return held
}

func (e *Engine) halted() (bool, string) {
	var s store.EngineState
	if e.DB.First(&s, 1).Error != nil {
		return false, ""
	}
	return s.Halted, s.HaltReason
}

func (e *Engine) activity(msg string) {
	e.DB.Model(&store.EngineState{}).Where("id = 1").Update("activity", trim(msg, 500))
}

func (e *Engine) recordError(err error) {
	e.DB.Model(&store.EngineState{}).Where("id = 1").Update("last_error", trim(time.Now().UTC().Format(time.RFC3339)+" "+err.Error(), 1000))
}

// note records the engine's latest message on a payout.
func (e *Engine) note(id uint64, msg string) {
	e.DB.Model(&store.ProceedPayout{}).Where("id = ?", id).Update("note", trim(msg, 1000))
}

func (e *Engine) publish(ctx context.Context, typ string, id uint64, msg string) {
	e.Bus.Publish(ctx, bus.Event{Type: typ, PayoutID: id, Message: msg})
}

// setStatus moves a payout from one status to another, only if it is still
// in from (an admin may have paused or cancelled it meanwhile).
func (e *Engine) setStatus(id uint64, from, to string, extra map[string]interface{}) (bool, error) {
	updates := map[string]interface{}{"status": to}
	for k, v := range extra {
		updates[k] = v
	}
	res := e.DB.Model(&store.ProceedPayout{}).Where("id = ? AND status = ?", id, from).Updates(updates)
	return res.RowsAffected == 1, res.Error
}

// currentStatus re-reads a payout's status.
func (e *Engine) currentStatus(id uint64) string {
	var p store.ProceedPayout
	if e.DB.Select("status").First(&p, id).Error != nil {
		return ""
	}
	return p.Status
}

func trim(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) > n {
		return s[:n]
	}
	return s
}

func bigOf(s string) *big.Int {
	n, ok := new(big.Int).SetString(strings.TrimSpace(s), 10)
	if !ok {
		return new(big.Int)
	}
	return n
}
