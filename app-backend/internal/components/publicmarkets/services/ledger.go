package publicmarkets

import (
	"context"
	"math/big"
	"strings"
	"time"

	pm "trovo-wallet-api/internal/components/publicmarkets/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ledgerConfirmations is how far behind the head the ledger follows.
const ledgerConfirmations = 2

// IndexLedgers follows every asset token's Transfers into the beneficial
// ownership ledger: one movement row per wallet per transfer (the
// record-date snapshots sum them) and a running balance per wallet.
func (e *Engine) IndexLedgers(ctx context.Context) {
	if e.Chain == nil {
		return
	}
	head, err := e.Chain.Head(ctx)
	if err != nil || head <= ledgerConfirmations {
		return
	}
	safeHead := head - ledgerConfirmations
	var assets []pm.Asset
	e.DB.Where("contract_address <> ''").Find(&assets)
	for i := range assets {
		if err := e.indexAsset(ctx, &assets[i], safeHead); err != nil {
			e.logf("ledger %s: %v", assets[i].AssetCode, err)
		}
	}
}

func (e *Engine) indexAsset(ctx context.Context, a *pm.Asset, safeHead uint64) error {
	if !common.IsHexAddress(a.ContractAddress) {
		return nil
	}
	if a.LedgerStartBlock == 0 {
		start, err := e.Chain.BlockAtTime(ctx, a.CreatedAt.Add(-time.Hour))
		if err != nil {
			return err
		}
		if start == 0 {
			start = 1
		}
		e.DB.Model(&pm.Asset{}).Where("id = ?", a.ID).Updates(map[string]interface{}{"ledger_start_block": start, "ledger_scanned_block": start - 1})
		a.LedgerStartBlock, a.LedgerScannedBlock = start, start-1
	}
	token := common.HexToAddress(a.ContractAddress)
	for from := a.LedgerScannedBlock + 1; from <= safeHead; {
		to := from + 1999
		if to > safeHead {
			to = safeHead
		}
		transfers, err := e.Chain.Transfers(ctx, token, from, to)
		if err != nil {
			return err
		}
		if err := e.applyTransfers(a, transfers, to); err != nil {
			return err
		}
		a.LedgerScannedBlock = to
		from = to + 1
	}
	return nil
}

// applyTransfers records a range's transfers and the scan's progress in
// one database transaction, so a crash never counts a transfer twice.
func (e *Engine) applyTransfers(a *pm.Asset, transfers []Transfer, scannedTo uint64) error {
	return e.DB.Transaction(func(tx *gorm.DB) error {
		zero := common.Address{}
		for _, t := range transfers {
			for _, leg := range []struct {
				wallet common.Address
				side   string
				delta  *big.Int
			}{{t.From, "OUT", new(big.Int).Neg(t.Value)}, {t.To, "IN", t.Value}} {
				if leg.wallet == zero {
					continue // mints come from, burns go to, the zero address
				}
				m := pm.LedgerMovement{AssetID: a.ID, WalletAddress: leg.wallet.Hex(), Delta: leg.delta.String(), BlockNumber: t.Block,
					TxHash: t.TxHash, LogIndex: t.LogIndex, Side: leg.side}
				res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&m)
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					continue
				}
				if err := e.addToLedger(tx, a, leg.wallet.Hex(), leg.delta); err != nil {
					return err
				}
			}
		}
		return tx.Model(&pm.Asset{}).Where("id = ?", a.ID).Update("ledger_scanned_block", scannedTo).Error
	})
}

func (e *Engine) addToLedger(tx *gorm.DB, a *pm.Asset, wallet string, delta *big.Int) error {
	var row pm.LedgerEntry
	if tx.Where("asset_id = ? AND wallet_address = ?", a.ID, wallet).First(&row).Error != nil {
		channel, link := e.channelOf(a, wallet)
		row = pm.LedgerEntry{AssetID: a.ID, WalletAddress: wallet, Balance: "0", Channel: channel, ServiceLinkID: link}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
	}
	bal, _ := new(big.Int).SetString(row.Balance, 10)
	if bal == nil {
		bal = new(big.Int)
	}
	bal.Add(bal, delta)
	return tx.Model(&pm.LedgerEntry{}).Where("id = ?", row.ID).Updates(map[string]interface{}{"balance": bal.String(), "updated_at": e.now()}).Error
}

// channelOf attributes a wallet: an exchange customer's, a platform Safe,
// or a Trovo App wallet.
func (e *Engine) channelOf(a *pm.Asset, wallet string) (string, string) {
	if strings.EqualFold(wallet, a.IssuingSafeAddress) || (e.Chain != nil && strings.EqualFold(wallet, e.Chain.TreasurySafe().Hex())) {
		return pm.ChannelPlatform, ""
	}
	var w pm.PartnerWallet
	if e.DB.Where("LOWER(wallet_address) = ?", strings.ToLower(wallet)).First(&w).Error == nil {
		return pm.ChannelExchange, w.ServiceLinkID
	}
	return pm.ChannelApp, ""
}

// BalancesAt is every wallet's balance of an asset at the end of block
// (base units), from the ledger's movements.
func (e *Engine) BalancesAt(a *pm.Asset, block uint64) map[string]*big.Int {
	type row struct {
		WalletAddress string
		Delta         string
	}
	var rows []row
	e.DB.Model(&pm.LedgerMovement{}).Select("wallet_address, delta").Where("asset_id = ? AND block_number <= ?", a.ID, block).Scan(&rows)
	out := map[string]*big.Int{}
	for _, r := range rows {
		v, ok := new(big.Int).SetString(r.Delta, 10)
		if !ok {
			continue
		}
		if out[r.WalletAddress] == nil {
			out[r.WalletAddress] = new(big.Int)
		}
		out[r.WalletAddress].Add(out[r.WalletAddress], v)
	}
	for k, v := range out {
		if v.Sign() <= 0 {
			delete(out, k)
		}
	}
	return out
}

// HolderView is a beneficial owner as Trovo Manager lists it.
type HolderView struct {
	WalletAddress string `json:"walletAddress"`
	Channel       string `json:"channel"`
	ServiceLinkID string `json:"serviceLinkId"`
	Balance       string `json:"balance"`
	PercentSupply string `json:"percentSupply"`
}

// Holders are an asset's beneficial owners, largest first.
func (e *Engine) Holders(a *pm.Asset, limit int) []HolderView {
	var rows []pm.LedgerEntry
	e.DB.Where("asset_id = ? AND balance <> '0'", a.ID).Find(&rows)
	total := e.LedgerTotal(a)
	out := make([]HolderView, 0, len(rows))
	for _, r := range rows {
		bal := d(r.Balance).Shift(-int32(a.TokenDecimals))
		if !bal.IsPositive() {
			continue
		}
		pct := decimal.Zero
		if total.IsPositive() {
			pct = bal.Div(total).Mul(decimal.NewFromInt(100))
		}
		out = append(out, HolderView{WalletAddress: r.WalletAddress, Channel: r.Channel, ServiceLinkID: r.ServiceLinkID, Balance: bal.String(), PercentSupply: pct.StringFixed(2)})
	}
	sortHolders(out)
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

func sortHolders(h []HolderView) {
	for i := 1; i < len(h); i++ {
		for j := i; j > 0 && d(h[j].Balance).GreaterThan(d(h[j-1].Balance)); j-- {
			h[j], h[j-1] = h[j-1], h[j]
		}
	}
}
