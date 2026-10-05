package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	swapModels "trovo-wallet-api/internal/components/swaps/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/offerbook"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Swaps fill offers on TrovoOfferBook (market/contracts): swapping token A
// for token B buys B from the open offers that sell B and accept A,
// cheapest first - the market-making offers users place with MakeOffer and
// any other offers on the book. Each fill is authorized by the platform
// (OFFER_AUTHORIZER_PRIVATE_KEY) for exactly that wallet and amount, and
// the swap is one wallet operation the user signs: the service fee and
// VAT transfers, allowing the book to take the payment, then the fills.
// Offers are found through the offer book index and re-read on-chain
// before a swap is priced. Only tokens trade on the book; the native asset
// (ETH) cannot be swapped.

// OperationSwap is a swap's wallet operation kind.
const OperationSwap = "SWAP"

// maxSwapFills bounds how many offers one swap fills (its gas).
const maxSwapFills = 8

// swapFill is one priced fill of a swap.
type swapFill struct {
	Request       aa.FillRequest
	Payment       *big.Int
	Authorization []byte
}

// swapRoute is a priced swap of source for destination.
type swapRoute struct {
	Book        common.Address
	Source      common.Address
	Destination common.Address
	SourceDec   uint8
	DestDec     uint8
	Fills       []swapFill
	Pay         *big.Int // source base units
	Receive     *big.Int // destination base units
}

// PayHuman is what the swap pays, in whole source tokens.
func (r *swapRoute) PayHuman() decimal.Decimal {
	return decimal.NewFromBigInt(r.Pay, -int32(r.SourceDec))
}

// ReceiveHuman is what the swap receives, in whole destination tokens.
func (r *swapRoute) ReceiveHuman() decimal.Decimal {
	return decimal.NewFromBigInt(r.Receive, -int32(r.DestDec))
}

// errSwapLiquidity is the low-liquidity error, saying how much is on offer.
func errSwapLiquidity(available decimal.Decimal, sourceCode, destCode string) error {
	msg := fmt.Sprintf("There is no %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destCode, sourceCode, sourceCode)
	if available.IsPositive() {
		msg = fmt.Sprintf("There is only %v %v to exchange for your %v at this time. Please reduce the quantity of %v to try again.", available, destCode, sourceCode, sourceCode)
	}
	return &tErrors.CustomError{Param: "destinationAssetCode", Err: "error-low-liquidity", ErrMessage: msg}
}

// routeSwap prices spending budget (source base units) on destination, over
// the indexed offers re-read on-chain. It does not authorize the fills.
func routeSwap(ctx context.Context, db *gorm.DB, chain *ethclient.Client, book, source, destination common.Address, budget *big.Int, sourceCode, destCode string) (*swapRoute, error) {
	candidates, err := offerbook.Candidates(db, destination, source)
	if err != nil {
		log.Printf("[routeSwap] reading offers: %v", err)
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	r := &swapRoute{Book: book, Source: source, Destination: destination, Pay: new(big.Int), Receive: new(big.Int)}
	if r.SourceDec, err = offerbook.Decimals(ctx, db, chain, source); err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	if r.DestDec, err = offerbook.Decimals(ctx, db, chain, destination); err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	left := new(big.Int).Set(budget)
	available := new(big.Int) // destination on offer at the prices seen
	short := true             // the offers ran out before the budget did
	for _, c := range candidates {
		if len(r.Fills) == maxSwapFills {
			break
		}
		// the index may lag the chain: price from the offer as it is now
		o, err := aa.ReadOffer(ctx, chain, book, c.OfferID)
		if err != nil {
			return nil, &tErrors.ErrorTemporaryServerError{}
		}
		if !o.Open || o.Remaining.Sign() == 0 || o.SellToken != destination {
			continue
		}
		price, err := aa.ReadOfferPrice(ctx, chain, book, c.OfferID, source)
		if err != nil {
			return nil, &tErrors.ErrorTemporaryServerError{}
		}
		if price.Num.Sign() == 0 {
			continue
		}
		available.Add(available, o.Remaining)
		amount := price.AmountFor(left)
		if amount.Sign() == 0 {
			// the rest does not buy a unit here (nor, cheapest first,
			// anywhere after)
			short = false
			break
		}
		if amount.Cmp(o.Remaining) >= 0 {
			amount = new(big.Int).Set(o.Remaining)
		} else {
			short = false
		}
		cost := price.Cost(amount)
		r.Fills = append(r.Fills, swapFill{Request: aa.FillRequest{OfferID: c.OfferID, PaymentToken: source, Amount: amount, MaxPayment: cost}, Payment: cost})
		r.Pay.Add(r.Pay, cost)
		r.Receive.Add(r.Receive, amount)
		left.Sub(left, cost)
		if !short {
			break
		}
	}
	if r.Receive.Sign() == 0 || (short && left.Sign() > 0) {
		return nil, errSwapLiquidity(decimal.NewFromBigInt(available, -int32(r.DestDec)), sourceCode, destCode)
	}
	return r, nil
}

// authorize has the platform authorize taker to make the route's fills,
// delivering to recipient, until deadline.
func (r *swapRoute) authorize(taker, recipient common.Address, deadline time.Time) error {
	key, err := offerAuthorizerKey()
	if err != nil {
		return err
	}
	for i := range r.Fills {
		nonce, err := aa.RandomFillNonce()
		if err != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		f := &r.Fills[i]
		f.Request.Recipient = recipient
		f.Request.Nonce = nonce
		f.Request.Deadline = big.NewInt(deadline.Unix())
		if f.Authorization, err = aa.SignFill(key, network.GetBlockchainChainID(), r.Book, f.Request, taker); err != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
	}
	return nil
}

// calls are the taker's calls making the swap: allow the book to take the
// payment, then the fills.
func (r *swapRoute) calls() ([]aa.Call, error) {
	out := []aa.Call{aa.ERC20Approve(r.Source, r.Book, r.Pay)}
	for _, f := range r.Fills {
		c, err := aa.FillCall(r.Book, f.Request, f.Authorization)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

// swapToken is a swap side's token contract, or an error for the native
// asset (which cannot trade on the book).
func swapToken(code, contract, param string) (common.Address, error) {
	if contract == "" || strings.EqualFold(code, os.Getenv("NATIVE_ASSET_CODE")) || !common.IsHexAddress(contract) {
		return common.Address{}, &tErrors.CustomError{Param: param, Err: "error-asset-not-swappable", ErrMessage: fmt.Sprintf("%v cannot be swapped. Only tokens can be swapped.", nonEmpty(code, os.Getenv("NATIVE_ASSET_CODE")))}
	}
	return common.HexToAddress(contract), nil
}

func nonEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// QuoteSwap is what spending amount (whole source tokens) on the
// destination asset receives now, in whole destination tokens.
func QuoteSwap(ctx context.Context, sourceCode, sourceContract, amount, destCode, destContract string, gc *sharedconfig.GlobalConfig) (decimal.Decimal, error) {
	book, err := offerbook.BookAddress()
	if err != nil {
		return decimal.Zero, &tErrors.CustomError{Param: "destinationAssetCode", Err: "error-swaps-not-configured", ErrMessage: "Swaps are not available at this time. Please try again later.", Code: http.StatusServiceUnavailable}
	}
	source, err := swapToken(sourceCode, sourceContract, "sourceAssetCode")
	if err != nil {
		return decimal.Zero, err
	}
	dest, err := swapToken(destCode, destContract, "destinationAssetCode")
	if err != nil {
		return decimal.Zero, err
	}
	dec, err := offerbook.Decimals(ctx, gc.DB, gc.BantuExpansionClient, source)
	if err != nil {
		return decimal.Zero, &tErrors.ErrorTemporaryServerError{}
	}
	d, err := decimal.NewFromString(amount)
	if err != nil || !d.IsPositive() {
		return decimal.Zero, &tErrors.CustomError{Param: "sourceAmount", Err: "error-invalid-amount", ErrMessage: "Invalid amount."}
	}
	r, err := routeSwap(ctx, gc.DB, gc.BantuExpansionClient, book, source, dest, d.Shift(int32(dec)).Truncate(0).BigInt(), sourceCode, destCode)
	if err != nil {
		return decimal.Zero, err
	}
	return r.ReceiveHuman(), nil
}

// swapContext is what a swap operation needs when it is submitted.
type swapContext struct {
	Fees     []sharedconfig.FeeCollection `json:"fees"`
	Receive  string                       `json:"receive"`
	Spend    string                       `json:"spend"`
	OfferIDs []string                     `json:"offerIds"`
}

// SwapSend swaps swapInfo.SourceAmount of the source asset (fees included)
// for as much of the destination asset as it buys. Like Pay it runs in two
// steps: without a signature it builds the swap as a wallet operation and
// returns it (swapInfo.Transaction) with its estimate; with the signature
// it submits it. Shared wallets with approvers get an approval request on
// commit. The caller has computed the fees (swapInfo.FeeAmount, VatAmount,
// SwapAmount).
func SwapSend(signerUser, walletOwner *userModels.User, wallet *userModels.UserWallet, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	swapInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	// step 2: the user signed the operation built in step 1
	if swapInfo.Multiparty == 0 && len(swapInfo.TransactionSignature) > 0 && len(swapInfo.Transaction) > 0 {
		rec, p, err := LoadWalletOperation(swapInfo.Transaction, wallet.ID, OperationSwap, gc)
		if err != nil {
			return err
		}
		hash, err := SignSingleOwnerOperation(ctx, rec, p, signerUser.PrimarySigner, swapInfo.TransactionSignature, gc)
		if err != nil {
			return err
		}
		swapInfo.TransactionID = hash
		recordOperationFees(rec, hash, gc)
		if rec.Context != nil {
			var c swapContext
			if json.Unmarshal([]byte(*rec.Context), &c) == nil && c.Receive != "" {
				swapInfo.SwappedEstimate = c.Receive
			}
		}
		wallet.InvalidateUserCache(gc)
		signerUser.InvalidateUserWalletCache(gc)
		signerUser.InvalidateUserCache(gc)
		return nil
	}
	// a shared wallet's initiator commits the operation built in step 1
	if swapInfo.Multiparty == 1 && swapInfo.Commit == 1 && len(swapInfo.Transaction) > 0 {
		if _, _, err := LoadWalletOperation(swapInfo.Transaction, wallet.ID, OperationSwap, gc); err != nil {
			return err
		}
		return createSwapApprovalRequest(signerUser, wallet, swapInfo, gc)
	}

	// step 1: build
	calls, sc, err := buildSwapCalls(ctx, walletOwner, wallet, swapInfo, gc)
	if err != nil {
		return err
	}
	validity := time.Duration(0)
	if swapInfo.Multiparty == 1 {
		validity = sharedWalletOperationValidity()
	}
	op, err := PrepareWalletOperation(ctx, OperationSwap, signerUser, walletOwner, wallet, calls, validity, sc, gc)
	if err != nil {
		return err
	}
	swapInfo.Transaction = op.Transaction
	swapInfo.Messages = append(swapInfo.Messages, op.Messages()...)
	swapInfo.SignatureRequired = 1
	if swapInfo.Multiparty == 1 && swapInfo.Commit == 1 {
		return createSwapApprovalRequest(signerUser, wallet, swapInfo, gc)
	}
	return nil
}

// swapAuthorizationValidity is how long a swap's fill authorizations stay
// valid: past its operation's validity.
func swapAuthorizationValidity(multiparty bool, gc *sharedconfig.GlobalConfig) time.Duration {
	v := operationBuilder(gc).DefaultValidity
	if multiparty {
		v = sharedWalletOperationValidity()
	}
	return v + 10*time.Minute
}

// buildSwapCalls prices the swap and returns the wallet's calls: the fee
// and VAT transfers, then the swap.
func buildSwapCalls(ctx context.Context, owner *userModels.User, wallet *userModels.UserWallet, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) ([]aa.Call, *swapContext, error) {
	book, err := offerbook.BookAddress()
	if err != nil {
		return nil, nil, &tErrors.CustomError{Param: "destinationAssetCode", Err: "error-swaps-not-configured", ErrMessage: "Swaps are not available at this time. Please try again later.", Code: http.StatusServiceUnavailable}
	}
	source, err := swapToken(swapInfo.SourceAssetCode, swapInfo.SourceContractAddress, "sourceAssetCode")
	if err != nil {
		return nil, nil, err
	}
	dest, err := swapToken(swapInfo.DestinationAssetCode, swapInfo.DestinationContractAddress, "destinationAssetCode")
	if err != nil {
		return nil, nil, err
	}
	client := gc.BantuExpansionClient
	sourceAsset := basetxn.CreditAsset{Code: swapInfo.SourceAssetCode, Issuer: source.Hex()}
	destAsset := basetxn.CreditAsset{Code: swapInfo.DestinationAssetCode, Issuer: dest.Hex()}

	total, err := decimal.NewFromString(swapInfo.SourceAmount)
	if err != nil || !total.IsPositive() {
		return nil, nil, &tErrors.CustomError{Param: "sourceAmount", Err: "error-invalid-amount", ErrMessage: "Invalid amount."}
	}
	_, sourceAuthorized, _, sourceBalance, _, err := network.BlockchainAccountProperties(client, wallet.ID, sourceAsset)
	if err != nil {
		return nil, nil, err
	}
	if !sourceAuthorized || sourceBalance.LessThan(total) {
		return nil, nil, &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs %v %v or you reduce the amount you want to swap.", total.Sub(sourceBalance), swapInfo.SourceAssetCode)}
	}
	// a regulated destination asset: swapping into it is this wallet's own
	// way of starting to hold it
	if network.RequiresWalletAuthorization(destAsset.GetCode()) && !network.IsWalletAuthorizedForAsset(wallet.ID, destAsset) {
		if e := network.SetWalletAssetAuthorization(wallet.ID, destAsset, true, destAsset.GetIssuer(), "auto-authorized via swap into this asset"); e != nil {
			log.Printf("[buildSwapCalls] authorizing %v for %v: %v", wallet.ID, destAsset.GetCode(), e)
			return nil, nil, &tErrors.CustomError{Param: "destination", Err: "error-could-not-authorize-wallet", ErrMessage: "Could not authorize this wallet to hold the destination asset."}
		}
	}

	decimals, err := network.AssetDecimals(ctx, client, sourceAsset)
	if err != nil {
		return nil, nil, &tErrors.ErrorTemporaryServerError{}
	}
	budget, err := baseUnits(swapInfo.SwapAmount, decimals)
	if err != nil || budget.Sign() <= 0 {
		return nil, nil, &tErrors.CustomError{Param: "sourceAmount", Err: "error-invalid-amount", ErrMessage: "The amount left to swap after fees is too small."}
	}
	route, err := routeSwap(ctx, gc.DB, client, book, source, dest, budget, swapInfo.SourceAssetCode, swapInfo.DestinationAssetCode)
	if err != nil {
		return nil, nil, err
	}
	walletAddr := common.HexToAddress(wallet.ID)
	if err := route.authorize(walletAddr, walletAddr, time.Now().Add(swapAuthorizationValidity(swapInfo.Multiparty == 1, gc))); err != nil {
		return nil, nil, err
	}

	var calls []aa.Call
	var fees []sharedconfig.FeeCollection
	serviceFee := wallet.GetSwapFee(gc)
	feeAmount, _ := decimal.NewFromString(swapInfo.FeeAmount)
	if feeAmount.IsPositive() && serviceFee.Inactive == 0 {
		feeAddr, err := feeWalletAddress(serviceFee.FeeWalletSecretKey, "swap fee wallet", gc)
		if err != nil {
			return nil, nil, err
		}
		if err := authorizeFeeWallet(feeAddr, sourceAsset, "swap fee wallet"); err != nil {
			return nil, nil, err
		}
		units, err := baseUnits(swapInfo.FeeAmount, decimals)
		if err != nil {
			return nil, nil, err
		}
		calls = append(calls, transferCall(sourceAsset, feeAddr, units))
		swapInfo.Messages = append(swapInfo.Messages, "Service fee will apply.")
		fees = append(fees, feeRecord("SWAP", wallet, owner, sourceAsset, swapInfo.FeeAmount, wallet.Alias, swapInfo.Multiparty))

		vatAmount, _ := decimal.NewFromString(swapInfo.VatAmount)
		if vatAmount.IsPositive() {
			vatAddr, err := feeWalletAddress(gc.GetVATWallet(), "VAT wallet", gc)
			if err != nil {
				return nil, nil, err
			}
			if err := authorizeFeeWallet(vatAddr, sourceAsset, "VAT wallet"); err != nil {
				return nil, nil, err
			}
			units, err := baseUnits(swapInfo.VatAmount, decimals)
			if err != nil {
				return nil, nil, err
			}
			calls = append(calls, transferCall(sourceAsset, vatAddr, units))
			fees = append(fees, feeRecord("VAT", wallet, owner, sourceAsset, swapInfo.VatAmount, wallet.Alias, swapInfo.Multiparty))
		}
	}
	swapCalls, err := route.calls()
	if err != nil {
		return nil, nil, &tErrors.ErrorTemporaryServerError{}
	}
	calls = append(calls, swapCalls...)

	swapInfo.SwappedEstimate = route.ReceiveHuman().String()
	if spent := route.PayHuman(); !spent.Equal(decimal.RequireFromString(swapInfo.SwapAmount)) {
		swapInfo.Messages = append(swapInfo.Messages, fmt.Sprintf("%v %v of it is swapped; the remaining %v %v is too little to buy more %v and stays in the wallet.",
			spent, swapInfo.SourceAssetCode, decimal.RequireFromString(swapInfo.SwapAmount).Sub(spent), swapInfo.SourceAssetCode, swapInfo.DestinationAssetCode))
	}
	sc := &swapContext{Fees: fees, Receive: swapInfo.SwappedEstimate, Spend: route.PayHuman().String()}
	for _, f := range route.Fills {
		sc.OfferIDs = append(sc.OfferIDs, f.Request.OfferID.String())
	}
	return calls, sc, nil
}

// createSwapApprovalRequest records a shared wallet's swap for its
// approvers; they sign the same operation (see approvals).
func createSwapApprovalRequest(signerUser *userModels.User, wallet *userModels.UserWallet, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) error {
	swapInfo.TransactionID = "PENDING_AUTH"
	description := fmt.Sprintf("Swap\n From: %v %v,\n To: %v,\n Est. Value After: %v", swapInfo.SourceAmount, nonEmpty(swapInfo.SourceAssetCode, os.Getenv("NATIVE_ASSET_CODE")), nonEmpty(swapInfo.DestinationAssetCode, os.Getenv("NATIVE_ASSET_CODE")), swapInfo.SwappedEstimate)
	if len(swapInfo.Memo) > 0 {
		description = fmt.Sprintf("%v\nMemo: %v", description, swapInfo.Memo)
	}
	if len(swapInfo.Messages) > 0 {
		description = fmt.Sprintf("%v\nMessages: %v", description, strings.Join(swapInfo.Messages, "\n"))
	}
	swapInfo.ReturnedDescription = description
	infoBytes, _ := json.Marshal(*swapInfo)
	info := string(infoBytes)
	pendingAuth := userModels.PendingAuth{
		ID:                     uuid.NewString(),
		Initiator:              signerUser.Username,
		InitiatorSignerAddress: signerUser.PrimarySigner,
		WalletAddress:          wallet.ID,
		TransactionType:        "SWAP",
		Description:            description,
		TransactionSource:      wallet.ID,
		ApprovalsNeeded:        wallet.NumberOfApprovalsNeeded,
		TransactionXdr:         swapInfo.Transaction,
		TransactionInfoStr:     &info,
	}
	if err := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error; err != nil {
		log.Printf("[SwapSend] saving approval request for %v: %v", wallet.ID, err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	return nil
}
