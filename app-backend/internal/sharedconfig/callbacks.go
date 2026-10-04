package sharedconfig

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Callbacks to partners (service-link login/authorization/event callbacks
// and payment notifications) are recorded before they are sent, so none is
// lost when an instance stops: SendCallback stores it and tries it at once;
// DeliverDueCallbacks (one instance at a time, see main.go) retries those
// that failed, with backoff, until callbackMaxAttempts.

const (
	CallbackPending   = "PENDING"
	CallbackDelivered = "DELIVERED"
	CallbackFailed    = "FAILED"

	callbackMaxAttempts = 10
	// the first attempt belongs to SendCallback for this long
	callbackFirstAttemptWindow = 30 * time.Second
)

// CallbackDelivery is one callback to a partner and its delivery state.
type CallbackDelivery struct {
	ID            string `gorm:"primaryKey;size:40"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	URL           string    `gorm:"type:text"`
	Body          string    `gorm:"type:text"`
	Status        string    `gorm:"size:20;index"`
	Attempts      int       `gorm:"default:0"`
	NextAttemptAt time.Time `gorm:"index"`
	LastError     string    `gorm:"type:text"`
}

var callbackHTTP = &http.Client{Timeout: 15 * time.Second}

// SendCallback POSTs body (JSON) to url: it is recorded first, then tried
// right away in the background; if that fails it is retried later.
func (gc *GlobalConfig) SendCallback(url string, body []byte) {
	cb := CallbackDelivery{ID: uuid.NewString(), URL: url, Body: string(body), Status: CallbackPending, NextAttemptAt: time.Now().Add(callbackFirstAttemptWindow)}
	if err := gc.DB.Create(&cb).Error; err != nil {
		// not recorded: still try once
		log.Printf("[SendCallback] recording callback to %v: %v", url, err)
		go func() {
			if err := postCallback(context.Background(), url, body); err != nil {
				log.Printf("[SendCallback] callback to %v failed and was not recorded for retry: %v", url, err)
			}
		}()
		return
	}
	go attemptCallback(context.Background(), gc, cb)
}

func postCallback(ctx context.Context, url string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := callbackHTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	reply, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("status %d: %s", resp.StatusCode, reply)
	}
	return nil
}

// attemptCallback makes one delivery attempt and records the outcome.
func attemptCallback(ctx context.Context, gc *GlobalConfig, cb CallbackDelivery) {
	err := postCallback(ctx, cb.URL, []byte(cb.Body))
	attempts := cb.Attempts + 1
	updates := map[string]interface{}{"attempts": attempts}
	switch {
	case err == nil:
		updates["status"] = CallbackDelivered
		updates["last_error"] = ""
	case attempts >= callbackMaxAttempts:
		updates["status"] = CallbackFailed
		updates["last_error"] = err.Error()
		log.Printf("[callbacks] giving up on callback %v to %v after %d attempts: %v", cb.ID, cb.URL, attempts, err)
	default:
		backoff := 30 * time.Second << (attempts - 1) // 30s, 1m, 2m, ...
		if backoff > time.Hour {
			backoff = time.Hour
		}
		updates["next_attempt_at"] = time.Now().Add(backoff)
		updates["last_error"] = err.Error()
	}
	if e := gc.DB.Model(&CallbackDelivery{}).Where("id = ?", cb.ID).Updates(updates).Error; e != nil {
		log.Printf("[callbacks] recording attempt of %v: %v", cb.ID, e)
	}
}

// DeliverDueCallbacks retries the callbacks whose next attempt is due.
func DeliverDueCallbacks(ctx context.Context, gc *GlobalConfig) {
	var due []CallbackDelivery
	gc.DB.Where("status = ? AND next_attempt_at <= ?", CallbackPending, time.Now()).Order("next_attempt_at").Limit(50).Find(&due)
	for _, cb := range due {
		if ctx.Err() != nil {
			return
		}
		attemptCallback(ctx, gc, cb)
	}
}
