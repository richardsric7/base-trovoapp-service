package publicmarkets

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	pm "trovo-wallet-api/internal/components/publicmarkets/models"
	"trovo-wallet-api/internal/components/publicmarkets/partners"
	"trovo-wallet-api/internal/gnosissafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ExchangePartnerOf returns the Public Markets configuration of a service
// link, refusing one that is not onboarded or is suspended.
func (e *Engine) ExchangePartnerOf(serviceLinkID string) (*pm.ExchangePartner, error) {
	var p pm.ExchangePartner
	if err := e.DB.First(&p, "service_link_id = ?", serviceLinkID).Error; err != nil {
		return nil, refuse(http.StatusForbidden, "error-not-onboarded", "", "This service link is not onboarded for Public Markets.")
	}
	if p.Status != "active" {
		return nil, refuse(http.StatusForbidden, "error-partner-suspended", "", "Public Markets access for this partner is suspended.")
	}
	return &p, nil
}

// VerifyExchangeSignature checks a request's HMAC: hex(HMAC-SHA256(secret,
// timestamp + "." + body)), within five minutes of the timestamp. The
// previous secret is accepted until it expires after a rotation.
func (e *Engine) VerifyExchangeSignature(p *pm.ExchangePartner, timestamp string, body []byte, signature string) error {
	ts, err := strconv.ParseInt(strings.TrimSpace(timestamp), 10, 64)
	if err != nil {
		return refuse(http.StatusUnauthorized, "error-invalid-signature", "X-Trovotech-Timestamp", "A Unix timestamp header (X-Trovotech-Timestamp) is required.")
	}
	if skew := e.now().Unix() - ts; skew > 300 || skew < -300 {
		return refuse(http.StatusUnauthorized, "error-invalid-signature", "X-Trovotech-Timestamp", "The request timestamp is more than 5 minutes off.")
	}
	if partners.VerifySignature(p.SigningSecret, timestamp, body, signature) {
		return nil
	}
	if p.PreviousSigningSecret != "" && p.PreviousSecretExpiresAt != nil && e.now().Before(*p.PreviousSecretExpiresAt) &&
		partners.VerifySignature(p.PreviousSigningSecret, timestamp, body, signature) {
		return nil
	}
	return refuse(http.StatusUnauthorized, "error-invalid-signature", "Authorization", "The request signature does not match.")
}

// ---------------------------------------------------------------- wallets

// ProvisionRequest is an exchange's wallet provisioning call (§6.3.1).
type ProvisionRequest struct {
	ExternalUserRef  string `json:"externalUserRef"`
	LegalName        string `json:"legalName"`
	TaxIdentifier    string `json:"taxIdentifier"`
	ResidencyCountry string `json:"residencyCountry"`
	Nationality      string `json:"nationality"`
	NDPAConsent      *bool  `json:"ndpaConsent"`
}

// WalletView is a provisioned wallet as the exchange sees it.
type WalletView struct {
	WalletID        string `json:"walletId"`
	PublicKey       string `json:"publicKey"`
	Status          string `json:"status"`
	ExternalUserRef string `json:"externalUserRef"`
}

func walletSalt(walletID string) *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256([]byte("trovo-public-markets:" + walletID)))
}

func validCountry(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) != 2 && len(s) != 3 {
		return false
	}
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

// ProvisionWallet opens an individually addressed wallet for an exchange's
// customer: a Safe owned by the Public Markets signers (deployed on its
// first outgoing move). A repeated externalUserRef returns the same
// wallet. A request missing identity data is refused with the field and
// kept on record.
func (e *Engine) ProvisionWallet(ctx context.Context, serviceLinkID string, r ProvisionRequest) (*WalletView, bool, error) {
	r.ExternalUserRef = strings.TrimSpace(r.ExternalUserRef)
	r.ResidencyCountry = strings.ToUpper(strings.TrimSpace(r.ResidencyCountry))
	r.Nationality = strings.ToUpper(strings.TrimSpace(r.Nationality))
	if r.ExternalUserRef == "" {
		return nil, false, refuse(http.StatusBadRequest, "error-missing-field", "externalUserRef", "externalUserRef is required")
	}
	var existing pm.PartnerWallet
	if e.DB.Where("service_link_id = ? AND external_user_ref = ? AND status = ?", serviceLinkID, r.ExternalUserRef, "active").First(&existing).Error == nil {
		return &WalletView{WalletID: existing.ID, PublicKey: existing.WalletAddress, Status: existing.Status, ExternalUserRef: existing.ExternalUserRef}, false, nil
	}
	var bad *Error
	switch {
	case strings.TrimSpace(r.LegalName) == "":
		bad = refuse(http.StatusBadRequest, "error-missing-field", "legalName", "legalName is required")
	case strings.TrimSpace(r.TaxIdentifier) == "":
		bad = refuse(http.StatusBadRequest, "error-missing-field", "taxIdentifier", "taxIdentifier is required")
	case !validCountry(r.ResidencyCountry):
		bad = refuse(http.StatusBadRequest, "error-invalid-field", "residencyCountry", "residencyCountry must be an ISO country code")
	case !validCountry(r.Nationality):
		bad = refuse(http.StatusBadRequest, "error-invalid-field", "nationality", "nationality must be an ISO country code")
	case r.NDPAConsent == nil || !*r.NDPAConsent:
		bad = refuse(http.StatusBadRequest, "error-missing-consent", "ndpaConsent", "NDPA consent must be recorded before provisioning")
	}
	consent := r.NDPAConsent != nil && *r.NDPAConsent
	w := pm.PartnerWallet{ID: "wlt_" + randomHex(6), ServiceLinkID: serviceLinkID, ExternalUserRef: r.ExternalUserRef, LegalName: strings.TrimSpace(r.LegalName),
		TaxIdentifier: strings.TrimSpace(r.TaxIdentifier), ResidencyCountry: r.ResidencyCountry, Nationality: r.Nationality, NDPAConsent: consent,
		CreatedAt: e.now(), UpdatedAt: e.now()}
	if bad != nil {
		w.Status, w.RejectionReason = "rejected", bad.Message
		e.DB.Create(&w)
		return nil, false, bad
	}
	if e.Chain == nil {
		return nil, false, refuse(http.StatusServiceUnavailable, "error-temporary", "", "Wallet provisioning is not available right now.")
	}
	salt := walletSalt(w.ID)
	addr, err := e.Chain.PartnerSafe(ctx, salt, false)
	if err != nil {
		e.logf("provisioning %s: %v", w.ID, err)
		return nil, false, refuse(http.StatusServiceUnavailable, "error-temporary", "", "Wallet provisioning is not available right now.")
	}
	w.WalletAddress, w.SafeSaltNonce, w.Status = addr.Hex(), salt.String(), "active"
	if err := e.DB.Create(&w).Error; err != nil {
		return nil, false, refuse(http.StatusInternalServerError, "error-temporary", "", "The wallet could not be saved.")
	}
	return &WalletView{WalletID: w.ID, PublicKey: w.WalletAddress, Status: w.Status, ExternalUserRef: w.ExternalUserRef}, true, nil
}

// PartnerWallet returns one of an exchange's active wallets.
func (e *Engine) PartnerWallet(serviceLinkID, walletID string) (*pm.PartnerWallet, error) {
	var w pm.PartnerWallet
	if err := e.DB.First(&w, "id = ? AND service_link_id = ?", strings.TrimSpace(walletID), serviceLinkID).Error; err != nil {
		return nil, refuse(http.StatusNotFound, "error-unknown-wallet", "walletId", "Unknown walletId.")
	}
	if w.Status != "active" {
		return nil, refuse(http.StatusConflict, "error-wallet-inactive", "walletId", "This wallet is %s.", w.Status)
	}
	return &w, nil
}

// ---------------------------------------------------------------- the prefunded balance

// creditExchange adds to an exchange's balance; it is idempotent on
// reference (a repeated credit is a no-op).
func (e *Engine) creditExchange(serviceLinkID string, amount decimal.Decimal, kind, reference, note, by string) error {
	return e.moveExchange(serviceLinkID, amount, kind, reference, note, by)
}

// debitExchange takes from it, refusing to go below zero.
func (e *Engine) debitExchange(serviceLinkID string, amount decimal.Decimal, kind, reference, note, by string) error {
	return e.moveExchange(serviceLinkID, amount.Neg(), kind, reference, note, by)
}

func (e *Engine) moveExchange(serviceLinkID string, amount decimal.Decimal, kind, reference, note, by string) error {
	var done int64
	e.DB.Model(&pm.ExchangeLedgerEntry{}).Where("reference = ?", reference).Count(&done)
	if done > 0 {
		return nil
	}
	for i := 0; i < 5; i++ {
		var p pm.ExchangePartner
		if err := e.DB.First(&p, "service_link_id = ?", serviceLinkID).Error; err != nil {
			return refuse(http.StatusForbidden, "error-not-onboarded", "", "Unknown exchange partner.")
		}
		next := d(p.Balance).Add(amount)
		if next.IsNegative() {
			return refuse(http.StatusPaymentRequired, "error-insufficient-balance", "amount",
				"Your Public Markets balance (%s %s) does not cover this order. Fund it from your registered funding wallet.", d(p.Balance).StringFixed(2), LoadSettings(e.DB).FundingAssetCode)
		}
		err := e.DB.Transaction(func(tx *gorm.DB) error {
			res := tx.Model(&pm.ExchangePartner{}).Where("service_link_id = ? AND balance = ?", serviceLinkID, p.Balance).
				Updates(map[string]interface{}{"balance": next.String(), "updated_at": e.now()})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return errConflict
			}
			return tx.Create(&pm.ExchangeLedgerEntry{ServiceLinkID: serviceLinkID, Kind: kind, Amount: amount.String(), BalanceAfter: next.String(),
				Reference: reference, Note: note, CreatedBy: by, CreatedAt: e.now()}).Error
		})
		if errors.Is(err, errConflict) {
			continue
		}
		return err
	}
	return errors.New("the exchange balance kept changing; try again")
}

var errConflict = errors.New("concurrent update")

// VerifyDeposit credits an exchange's deposit: a mined transfer of the
// funding stablecoin from its registered funding wallet to the treasury.
func (e *Engine) VerifyDeposit(ctx context.Context, serviceLinkID, txHash, by string) (decimal.Decimal, error) {
	p, err := e.ExchangePartnerOf(serviceLinkID)
	if err != nil {
		return decimal.Zero, err
	}
	if !common.IsHexAddress(p.FundingAddress) {
		return decimal.Zero, refuse(http.StatusConflict, "error-no-funding-address", "txHash", "No funding wallet is registered for this partner; ask Trovotech to register one.")
	}
	if e.Chain == nil {
		return decimal.Zero, refuse(http.StatusServiceUnavailable, "error-temporary", "", "Deposits cannot be verified right now.")
	}
	ft, err := e.fundingToken(ctx)
	if err != nil {
		return decimal.Zero, err
	}
	txHash = strings.ToLower(strings.TrimSpace(txHash))
	transfers, err := e.Chain.TxTransfers(ctx, txHash, ft.Address)
	if err != nil {
		return decimal.Zero, refuse(http.StatusBadRequest, "error-invalid-deposit", "txHash", "The transaction was not found or did not succeed.")
	}
	total := new(big.Int)
	for _, t := range transfers {
		if t.From == common.HexToAddress(p.FundingAddress) && t.To == e.Chain.TreasurySafe() {
			total.Add(total, t.Value)
		}
	}
	if total.Sign() == 0 {
		return decimal.Zero, refuse(http.StatusBadRequest, "error-invalid-deposit", "txHash", "The transaction moves no %s from your funding wallet to the Public Markets treasury.", ft.Code)
	}
	amount := decimal.NewFromBigInt(total, -int32(ft.Decimals))
	if err := e.creditExchange(serviceLinkID, amount, "DEPOSIT", "deposit:"+txHash, "Deposit "+txHash, by); err != nil {
		return decimal.Zero, err
	}
	return amount, nil
}

// PayExchangeWithdrawal pays part of an exchange's balance back to its
// funding wallet from the treasury.
func (e *Engine) PayExchangeWithdrawal(ctx context.Context, serviceLinkID string, amount decimal.Decimal, reference, by string) (string, error) {
	var p pm.ExchangePartner
	if err := e.DB.First(&p, "service_link_id = ?", serviceLinkID).Error; err != nil {
		return "", err
	}
	if !common.IsHexAddress(p.FundingAddress) || e.Chain == nil {
		return "", errors.New("no funding wallet or chain")
	}
	if err := e.debitExchange(serviceLinkID, amount, "WITHDRAWAL", reference, "Withdrawal to "+p.FundingAddress, by); err != nil {
		return "", err
	}
	ft, err := e.fundingToken(ctx)
	if err != nil {
		return "", err
	}
	hash, err := e.Chain.Send(ctx, e.Chain.TreasurySafe(), []gnosissafe.Call{callOf(aa.ERC20Transfer(ft.Address, common.HexToAddress(p.FundingAddress), units(amount.String(), ft.Decimals)))})
	if err != nil {
		_ = e.creditExchange(serviceLinkID, amount, "ADJUSTMENT", reference+":reversal", "Withdrawal not sent: "+trimTo(err.Error(), 200), "system")
		return "", err
	}
	e.DB.Model(&pm.ExchangeLedgerEntry{}).Where("reference = ?", reference).Update("note", "Withdrawal to "+p.FundingAddress+" (tx "+hash+")")
	return hash, e.Chain.Wait(ctx, hash)
}

// ---------------------------------------------------------------- orders

// ExchangeOrderRequest is a creation or redemption placed by an exchange
// (§6.3.2-6.3.3). A redemption takes quantity (tokens) or amount (NGN).
type ExchangeOrderRequest struct {
	WalletID         string `json:"walletId"`
	AssetCode        string `json:"assetCode"`
	Amount           string `json:"amount"`
	AmountCurrency   string `json:"amountCurrency"`
	Quantity         string `json:"quantity"`
	ExternalOrderRef string `json:"externalOrderRef"`
}

// ExchangeOrderView is what an exchange sees of an order: pending until
// it settles; internal states are not exposed.
type ExchangeOrderView struct {
	TrovotechOrderID string     `json:"trovotechOrderId"`
	Status           string     `json:"status"` // pending | settled | rejected | failed
	Type             string     `json:"type"`
	WalletID         string     `json:"walletId"`
	AssetCode        string     `json:"assetCode"`
	Quantity         string     `json:"quantity"`
	Amount           string     `json:"amount"`
	Fee              string     `json:"fee"`
	Price            string     `json:"price"`
	ExternalOrderRef string     `json:"externalOrderRef"`
	Reason           string     `json:"reason,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	SettledAt        *time.Time `json:"settledAt,omitempty"`
}

// ExchangeView renders an order for its exchange.
func ExchangeView(o *pm.Order) ExchangeOrderView {
	status := "pending"
	switch o.State {
	case pm.StateComplete:
		status = "settled"
	case pm.StateRejected, pm.StateCancelled:
		status = "rejected"
	case pm.StateFailed:
		status = "failed"
	}
	price := o.ExecutedPrice
	if price == "" {
		price = o.ReferencePrice
	}
	amount := o.Amount
	if o.Type == pm.OrderRedemption {
		amount = o.NetAmount
	}
	v := ExchangeOrderView{TrovotechOrderID: o.ID, Status: status, Type: strings.ToLower(o.Type), WalletID: o.PartnerWalletID, AssetCode: o.AssetCode,
		Quantity: o.Quantity, Amount: amount, Fee: o.Fee, Price: price, ExternalOrderRef: o.ExternalOrderRef, CreatedAt: o.CreatedAt}
	if status == "settled" {
		v.SettledAt = o.CompletedAt
	}
	if status == "rejected" || status == "failed" {
		v.Reason = o.Note
	}
	return v
}

// PlaceExchangeCreation is the exchange's creation: a thin wrapper over
// the engine every channel uses, paid from the exchange's balance.
func (e *Engine) PlaceExchangeCreation(ctx context.Context, serviceLinkID string, r ExchangeOrderRequest) (*pm.Order, error) {
	w, err := e.PartnerWallet(serviceLinkID, r.WalletID)
	if err != nil {
		return nil, err
	}
	if c := strings.ToUpper(strings.TrimSpace(r.AmountCurrency)); c != "" && c != "NGN" && c != LoadSettings(e.DB).FundingAssetCode {
		return nil, refuse(http.StatusBadRequest, "error-invalid-currency", "amountCurrency", "amountCurrency must be NGN.")
	}
	a, err := e.AssetByCode(r.AssetCode)
	if err != nil {
		return nil, err
	}
	amount := d(r.Amount)
	q, err := e.QuoteCreation(ctx, a, amount)
	if err != nil {
		return nil, err
	}
	ref := "order:" + serviceLinkID + ":" + randomHex(8)
	if err := e.debitExchange(serviceLinkID, amount, "ORDER", ref, "Creation "+a.AssetCode+" "+r.ExternalOrderRef, "exchange"); err != nil {
		return nil, err
	}
	o, err := e.Place(OrderRequest{Asset: a, Channel: pm.ChannelExchange, WalletAddress: w.WalletAddress, WalletAlias: w.ExternalUserRef,
		ServiceLinkID: serviceLinkID, PartnerWalletID: w.ID, ExternalOrderRef: strings.TrimSpace(r.ExternalOrderRef), Quote: q, State: pm.StateQueued})
	if err != nil {
		_ = e.creditExchange(serviceLinkID, amount, "REFUND", ref+":refund", "Order not saved", "system")
		return nil, err
	}
	e.DB.Model(&pm.ExchangeLedgerEntry{}).Where("reference = ?", ref).Update("note", "Creation "+o.ID)
	return o, nil
}

// PlaceExchangeRedemption is the exchange's redemption: the customer's
// tokens move from their wallet to the issuing Safe, then the order runs
// like any other; proceeds go to the exchange's balance.
func (e *Engine) PlaceExchangeRedemption(ctx context.Context, serviceLinkID string, r ExchangeOrderRequest) (*pm.Order, error) {
	w, err := e.PartnerWallet(serviceLinkID, r.WalletID)
	if err != nil {
		return nil, err
	}
	a, err := e.AssetByCode(r.AssetCode)
	if err != nil {
		return nil, err
	}
	qty := d(r.Quantity)
	if !qty.IsPositive() && d(r.Amount).IsPositive() {
		p := e.CurrentPrice(a)
		if !p.Value.IsPositive() {
			return nil, refuse(http.StatusServiceUnavailable, "error-no-price", "assetCode", "%s has no reference price yet.", a.AssetCode)
		}
		qty = d(r.Amount).Div(p.Value)
	}
	q, err := e.QuoteRedemption(ctx, a, qty)
	if err != nil {
		return nil, err
	}
	if e.Chain != nil && common.IsHexAddress(a.ContractAddress) {
		bal, err := e.Chain.BalanceOf(ctx, common.HexToAddress(a.ContractAddress), common.HexToAddress(w.WalletAddress))
		if err != nil {
			return nil, refuse(http.StatusServiceUnavailable, "error-temporary", "", "The wallet's balance could not be read.")
		}
		held := decimal.NewFromBigInt(bal, -int32(a.TokenDecimals)).Sub(e.pendingRedemptions(w.WalletAddress, a.ID))
		if held.LessThan(d(q.Quantity)) {
			return nil, refuse(http.StatusConflict, "error-insufficient-holding", "quantity", "The wallet holds %s %s available to sell.", held, a.AssetCode)
		}
	}
	return e.Place(OrderRequest{Asset: a, Channel: pm.ChannelExchange, WalletAddress: w.WalletAddress, WalletAlias: w.ExternalUserRef,
		ServiceLinkID: serviceLinkID, PartnerWalletID: w.ID, ExternalOrderRef: strings.TrimSpace(r.ExternalOrderRef), Quote: q, State: pm.StateAwaitingPayment})
}

// pendingRedemptions are tokens of a wallet already committed to open
// redemptions whose tokens have not left it yet.
func (e *Engine) pendingRedemptions(wallet, assetID string) decimal.Decimal {
	var orders []pm.Order
	e.DB.Where("wallet_address = ? AND asset_id = ? AND type = ? AND state = ?", wallet, assetID, pm.OrderRedemption, pm.StateAwaitingPayment).Find(&orders)
	total := decimal.Zero
	for _, o := range orders {
		total = total.Add(d(o.Quantity))
	}
	return total
}

// exchangeRedemptionTransfer moves an exchange customer's tokens to the
// issuing Safe (deploying their wallet's Safe the first time).
func (e *Engine) exchangeRedemptionTransfer(ctx context.Context, o *pm.Order, a *pm.Asset) error {
	if e.Chain == nil {
		return errors.New("no chain configured")
	}
	if o.PaymentTxHash == "" {
		var w pm.PartnerWallet
		if err := e.DB.First(&w, "id = ?", o.PartnerWalletID).Error; err != nil {
			return err
		}
		salt, ok := new(big.Int).SetString(w.SafeSaltNonce, 10)
		if !ok {
			return errors.New("the wallet has no Safe salt")
		}
		if !w.Deployed {
			addr, err := e.Chain.PartnerSafe(ctx, salt, true)
			if err != nil {
				return err
			}
			if !strings.EqualFold(addr.Hex(), w.WalletAddress) {
				return fmt.Errorf("the deployed Safe %s is not the wallet %s", addr.Hex(), w.WalletAddress)
			}
			now := e.now()
			e.DB.Model(&pm.PartnerWallet{}).Where("id = ?", w.ID).Updates(map[string]interface{}{"deployed": true, "deployed_at": &now})
		}
		res := e.DB.Model(&pm.Order{}).Where("id = ? AND payment_tx_hash = ''", o.ID).Update("payment_tx_hash", "pending")
		if res.RowsAffected == 0 {
			return nil
		}
		hash, err := e.Chain.Send(ctx, common.HexToAddress(w.WalletAddress), []gnosissafe.Call{callOf(aa.ERC20Transfer(common.HexToAddress(a.ContractAddress), common.HexToAddress(a.IssuingSafeAddress), units(o.Quantity, a.TokenDecimals)))})
		if err != nil {
			e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Update("payment_tx_hash", "")
			return err
		}
		e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Update("payment_tx_hash", hash)
		o.PaymentTxHash = hash
	}
	if o.PaymentTxHash == "pending" {
		return nil
	}
	if err := e.Chain.Wait(ctx, o.PaymentTxHash); err != nil {
		e.transition(o, []string{pm.StateAwaitingPayment}, pm.StateFailed, "The tokens could not be moved from the wallet: "+trimTo(err.Error(), 200), map[string]interface{}{"note": "Token transfer failed"})
		return nil
	}
	now := e.now()
	e.transition(o, []string{pm.StateAwaitingPayment}, pm.StateQueued, "Tokens moved from the wallet to the issuing Safe", map[string]interface{}{"payment_confirmed_at": &now})
	return nil
}

// ---------------------------------------------------------------- webhooks

// queueOrderWebhook tells the exchange its order settled, or did not.
func (e *Engine) queueOrderWebhook(o *pm.Order) {
	event := "creation.settled"
	switch {
	case o.State == pm.StateRejected || o.State == pm.StateFailed:
		event = "order.rejected"
	case o.Type == pm.OrderRedemption:
		event = "redemption.settled"
	}
	v := ExchangeView(o)
	payload := map[string]interface{}{"event": event, "trovotechOrderId": o.ID, "externalOrderRef": o.ExternalOrderRef, "walletId": o.PartnerWalletID,
		"assetCode": o.AssetCode, "quantity": o.Quantity, "amount": v.Amount, "price": v.Price, "timestamp": e.now().Format(time.RFC3339)}
	if event == "order.rejected" {
		payload["reason"] = o.Note
	}
	e.QueueWebhook(o.ServiceLinkID, event, o.ID, o.PartnerWalletID, o.AssetCode, payload, false)
}

// QueueWebhook records a webhook for delivery (at least once).
func (e *Engine) QueueWebhook(serviceLinkID, event, reference, walletID, assetCode string, payload map[string]interface{}, needsConfirmation bool) {
	id := "evt_" + randomHex(8)
	payload["eventId"] = id
	body, _ := json.Marshal(payload)
	now := e.now()
	e.DB.Create(&pm.WebhookDelivery{ID: id, ServiceLinkID: serviceLinkID, Event: event, Reference: reference, PartnerWalletID: walletID, AssetCode: assetCode,
		Payload: string(body), Status: pm.DeliveryPending, NextAttemptAt: &now, NeedsConfirmation: needsConfirmation, CreatedAt: now, UpdatedAt: now})
}

// DeliverWebhooks posts due webhooks, signed, with exponential backoff and
// a dead-letter after the configured attempts.
func (e *Engine) DeliverWebhooks(ctx context.Context) {
	s := LoadSettings(e.DB)
	var due []pm.WebhookDelivery
	e.DB.Where("status = ? AND next_attempt_at <= ?", pm.DeliveryPending, e.now()).Order("created_at").Limit(100).Find(&due)
	client := &http.Client{Timeout: 10 * time.Second}
	for _, w := range due {
		var p pm.ExchangePartner
		if e.DB.First(&p, "service_link_id = ?", w.ServiceLinkID).Error != nil {
			continue
		}
		code, err := e.post(ctx, client, &p, &w)
		attempts := w.Attempts + 1
		now := e.now()
		if err == nil {
			updates := map[string]interface{}{"status": pm.DeliveryDelivered, "attempts": attempts, "last_response_code": code, "delivered_at": &now, "last_error": ""}
			if w.NeedsConfirmation {
				if sla := confirmationSLA(&p, s); sla > 0 {
					due := now.Add(sla)
					updates["confirmation_due_at"] = &due
				}
			}
			e.DB.Model(&pm.WebhookDelivery{}).Where("id = ?", w.ID).Updates(updates)
			continue
		}
		updates := map[string]interface{}{"attempts": attempts, "last_response_code": code, "last_error": trimTo(err.Error(), 300)}
		if attempts >= s.WebhookMaxAttempts {
			updates["status"] = pm.DeliveryDeadLetter
			if e.GC != nil {
				e.GC.LogDiscordFailedRequest(fmt.Sprintf("[publicmarkets] webhook %s %s to %s dead-lettered: %v", w.Event, w.Reference, p.CallbackURL, err))
			}
		} else {
			next := now.Add(time.Duration(30<<uint(attempts-1)) * time.Second)
			updates["next_attempt_at"] = &next
		}
		e.DB.Model(&pm.WebhookDelivery{}).Where("id = ?", w.ID).Updates(updates)
	}
}

func confirmationSLA(p *pm.ExchangePartner, s pm.Settings) time.Duration {
	h := p.ConfirmationSLAHours
	if h == 0 {
		h = s.ConfirmationSLAHours
	}
	return time.Duration(h) * time.Hour
}

func (e *Engine) post(ctx context.Context, client *http.Client, p *pm.ExchangePartner, w *pm.WebhookDelivery) (int, error) {
	if strings.TrimSpace(p.CallbackURL) == "" {
		return 0, errors.New("no callback URL registered")
	}
	ts := strconv.FormatInt(e.now().Unix(), 10)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.CallbackURL, bytes.NewReader([]byte(w.Payload)))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Trovotech-Event", w.Event)
	req.Header.Set("X-Trovotech-Delivery", w.ID)
	req.Header.Set("X-Trovotech-Timestamp", ts)
	req.Header.Set("X-Trovotech-Signature", partners.Sign(p.SigningSecret, ts, []byte(w.Payload)))
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("callback answered %d", resp.StatusCode)
	}
	return resp.StatusCode, nil
}

// ConfirmationRequest is an exchange's receipt confirmation (§6.4.3).
type ConfirmationRequest struct {
	Event       string `json:"event"`
	WalletID    string `json:"walletId"`
	EventID     string `json:"eventId"`
	ConfirmedAt string `json:"confirmedAt"`
}

// Confirm records an exchange's confirmation of a webhook it received
// (required for dividend.paid). Without eventId the wallet's oldest
// unconfirmed event of that type is confirmed.
func (e *Engine) Confirm(serviceLinkID string, r ConfirmationRequest) (*pm.WebhookDelivery, error) {
	at := e.now()
	if t, err := time.Parse(time.RFC3339, strings.TrimSpace(r.ConfirmedAt)); err == nil {
		at = t.UTC()
	}
	q := e.DB.Where("service_link_id = ? AND needs_confirmation = ? AND confirmed_at IS NULL", serviceLinkID, true)
	if strings.TrimSpace(r.EventID) != "" {
		q = q.Where("id = ?", strings.TrimSpace(r.EventID))
	} else {
		if r.Event == "" || r.WalletID == "" {
			return nil, refuse(http.StatusBadRequest, "error-missing-field", "eventId", "eventId, or event and walletId, are required")
		}
		q = q.Where("event = ? AND partner_wallet_id = ?", r.Event, r.WalletID)
	}
	var w pm.WebhookDelivery
	if err := q.Order("created_at").First(&w).Error; err != nil {
		return nil, refuse(http.StatusNotFound, "error-nothing-to-confirm", "eventId", "No unconfirmed event matches.")
	}
	e.DB.Model(&pm.WebhookDelivery{}).Where("id = ?", w.ID).Update("confirmed_at", &at)
	w.ConfirmedAt = &at
	return &w, nil
}

// SweepConfirmationSLA escalates confirmations missing past their window.
func (e *Engine) SweepConfirmationSLA() {
	var late []pm.WebhookDelivery
	e.DB.Where("needs_confirmation = ? AND confirmed_at IS NULL AND escalated_at IS NULL AND confirmation_due_at IS NOT NULL AND confirmation_due_at < ?", true, e.now()).Limit(1000).Find(&late)
	if len(late) == 0 {
		return
	}
	now := e.now()
	ids := make([]string, 0, len(late))
	per := map[string]int{}
	for _, w := range late {
		ids = append(ids, w.ID)
		per[w.ServiceLinkID]++
	}
	e.DB.Model(&pm.WebhookDelivery{}).Where("id IN ?", ids).Update("escalated_at", &now)
	for sl, n := range per {
		if e.GC != nil {
			e.GC.LogDiscordFailedRequest(fmt.Sprintf("[publicmarkets] %d dividend confirmation(s) from service link %s past the SLA", n, sl))
		}
	}
}
