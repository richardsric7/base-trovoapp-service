package aa

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// TrovoOfferBook (market/contracts): fixed-price offers between curated
// tokens. A seller escrows a token and prices it in each payment token it
// accepts; every fill needs an EIP-712 authorization from a platform
// authorizer for exactly that taker, recipient and amount.
var offerBookABI = mustABI(`[
 {"name":"createOffer","type":"function","inputs":[
   {"name":"sellToken","type":"address"},{"name":"amount","type":"uint256"},{"name":"proceedsRecipient","type":"address"},
   {"name":"paymentTokens","type":"address[]"},
   {"name":"prices","type":"tuple[]","components":[{"name":"num","type":"uint128"},{"name":"den","type":"uint128"}]}],
  "outputs":[{"name":"offerId","type":"uint256"}]},
 {"name":"setPrice","type":"function","inputs":[{"name":"offerId","type":"uint256"},{"name":"paymentToken","type":"address"},
   {"name":"price","type":"tuple","components":[{"name":"num","type":"uint128"},{"name":"den","type":"uint128"}]}],"outputs":[]},
 {"name":"cancelOffer","type":"function","inputs":[{"name":"offerId","type":"uint256"}],"outputs":[]},
 {"name":"fill","type":"function","inputs":[
   {"name":"f","type":"tuple","components":[
     {"name":"offerId","type":"uint256"},{"name":"paymentToken","type":"address"},{"name":"amount","type":"uint256"},
     {"name":"maxPayment","type":"uint256"},{"name":"recipient","type":"address"},{"name":"nonce","type":"uint256"},{"name":"deadline","type":"uint256"}]},
   {"name":"authorization","type":"bytes"}],
  "outputs":[{"name":"payment","type":"uint256"}]},
 {"name":"getOffer","type":"function","stateMutability":"view","inputs":[{"name":"offerId","type":"uint256"}],
  "outputs":[{"type":"tuple","components":[
    {"name":"seller","type":"address"},{"name":"sellToken","type":"address"},{"name":"proceedsRecipient","type":"address"},
    {"name":"remaining","type":"uint256"},{"name":"open","type":"bool"}]}]},
 {"name":"priceOf","type":"function","stateMutability":"view","inputs":[{"name":"offerId","type":"uint256"},{"name":"paymentToken","type":"address"}],
  "outputs":[{"type":"tuple","components":[{"name":"num","type":"uint128"},{"name":"den","type":"uint128"}]}]},
 {"name":"tradable","type":"function","stateMutability":"view","inputs":[{"name":"token","type":"address"}],"outputs":[{"type":"bool"}]},
 {"name":"authorizers","type":"function","stateMutability":"view","inputs":[{"name":"a","type":"address"}],"outputs":[{"type":"bool"}]},
 {"name":"OfferCreated","type":"event","anonymous":false,"inputs":[
   {"name":"offerId","type":"uint256","indexed":true},{"name":"seller","type":"address","indexed":true},
   {"name":"sellToken","type":"address","indexed":true},{"name":"amount","type":"uint256","indexed":false},
   {"name":"proceedsRecipient","type":"address","indexed":false}]},
 {"name":"Filled","type":"event","anonymous":false,"inputs":[
   {"name":"offerId","type":"uint256","indexed":true},{"name":"taker","type":"address","indexed":true},
   {"name":"recipient","type":"address","indexed":true},{"name":"paymentToken","type":"address","indexed":false},
   {"name":"amount","type":"uint256","indexed":false},{"name":"payment","type":"uint256","indexed":false},
   {"name":"nonce","type":"uint256","indexed":false}]}
]`)

// OfferPrice: buying amount base units of the sell token costs
// ceil(amount * Num / Den) base units of the payment token.
type OfferPrice struct {
	Num *big.Int
	Den *big.Int
}

// Cost is what buying amount costs at p (rounded up, as the contract does).
func (p OfferPrice) Cost(amount *big.Int) *big.Int {
	n := new(big.Int).Mul(amount, p.Num)
	q, r := new(big.Int).QuoRem(n, p.Den, new(big.Int))
	if r.Sign() > 0 {
		q.Add(q, big.NewInt(1))
	}
	return q
}

// AmountFor is the most of the sell token payment buys at p.
func (p OfferPrice) AmountFor(payment *big.Int) *big.Int {
	return new(big.Int).Quo(new(big.Int).Mul(payment, p.Den), p.Num)
}

var maxUint128 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))

// Valid reports whether the contract accepts p as a price.
func (p OfferPrice) Valid() bool {
	return p.Num != nil && p.Den != nil && p.Num.Sign() > 0 && p.Den.Sign() > 0 && p.Num.Cmp(maxUint128) <= 0 && p.Den.Cmp(maxUint128) <= 0
}

type abiPrice struct {
	Num *big.Int `abi:"num"`
	Den *big.Int `abi:"den"`
}

// CreateOfferCall escrows amount of sellToken from the caller in a new offer
// priced in each of tokens.
func CreateOfferCall(book, sellToken common.Address, amount *big.Int, proceeds common.Address, tokens []common.Address, prices []OfferPrice) (Call, error) {
	ps := make([]abiPrice, len(prices))
	for i, p := range prices {
		if !p.Valid() {
			return Call{}, fmt.Errorf("aa: invalid offer price %v/%v", p.Num, p.Den)
		}
		ps[i] = abiPrice{Num: p.Num, Den: p.Den}
	}
	data, err := offerBookABI.Pack("createOffer", sellToken, amount, proceeds, tokens, ps)
	if err != nil {
		return Call{}, err
	}
	return Call{To: book, Value: big.NewInt(0), Data: data}, nil
}

// CancelOfferCall closes an offer and returns what is left to its seller.
func CancelOfferCall(book common.Address, offerID *big.Int) Call {
	data, _ := offerBookABI.Pack("cancelOffer", offerID)
	return Call{To: book, Value: big.NewInt(0), Data: data}
}

// FillRequest is a purchase from an offer (TrovoOfferBook.FillRequest).
type FillRequest struct {
	OfferID      *big.Int       `abi:"offerId"`
	PaymentToken common.Address `abi:"paymentToken"`
	Amount       *big.Int       `abi:"amount"`
	MaxPayment   *big.Int       `abi:"maxPayment"`
	Recipient    common.Address `abi:"recipient"`
	Nonce        *big.Int       `abi:"nonce"`
	Deadline     *big.Int       `abi:"deadline"`
}

// FillCall makes the purchase f with its authorization.
func FillCall(book common.Address, f FillRequest, authorization []byte) (Call, error) {
	data, err := offerBookABI.Pack("fill", f, authorization)
	if err != nil {
		return Call{}, err
	}
	return Call{To: book, Value: big.NewInt(0), Data: data}, nil
}

// RandomFillNonce is a fresh single-use authorization nonce.
func RandomFillNonce() (*big.Int, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(b), nil
}

var (
	eip712DomainTypeHash = crypto.Keccak256Hash([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))
	fillTypeHash         = crypto.Keccak256Hash([]byte("Fill(uint256 offerId,address taker,address recipient,address paymentToken,uint256 amount,uint256 maxPayment,uint256 nonce,uint256 deadline)"))
)

func addrWord(a common.Address) []byte { return common.LeftPadBytes(a.Bytes(), 32) }

// FillDigest is the EIP-712 digest an authorizer signs for taker to make
// fill f on the offer book at book on chainID.
func FillDigest(chainID *big.Int, book common.Address, f FillRequest, taker common.Address) common.Hash {
	domain := crypto.Keccak256Hash(eip712DomainTypeHash.Bytes(), crypto.Keccak256([]byte("TrovoOfferBook")), crypto.Keccak256([]byte("1")), word(chainID), addrWord(book))
	structHash := crypto.Keccak256Hash(fillTypeHash.Bytes(), word(f.OfferID), addrWord(taker), addrWord(f.Recipient), addrWord(f.PaymentToken),
		word(f.Amount), word(f.MaxPayment), word(f.Nonce), word(f.Deadline))
	return crypto.Keccak256Hash([]byte{0x19, 0x01}, domain.Bytes(), structHash.Bytes())
}

// SignFill is an authorizer's signature on fill f for taker.
func SignFill(key *ecdsa.PrivateKey, chainID *big.Int, book common.Address, f FillRequest, taker common.Address) ([]byte, error) {
	sig, err := crypto.Sign(FillDigest(chainID, book, f, taker).Bytes(), key)
	if err != nil {
		return nil, err
	}
	sig[64] += 27
	return sig, nil
}

// Offer is an offer's on-chain state.
type Offer struct {
	Seller            common.Address
	SellToken         common.Address
	ProceedsRecipient common.Address
	Remaining         *big.Int
	Open              bool
}

// ReadOffer reads offer id from the book.
func ReadOffer(ctx context.Context, r Caller, book common.Address, id *big.Int) (*Offer, error) {
	data, _ := offerBookABI.Pack("getOffer", id)
	out, err := callRaw(ctx, r, book, data)
	if err != nil {
		return nil, err
	}
	res, err := offerBookABI.Unpack("getOffer", out)
	if err != nil {
		return nil, err
	}
	o := res[0].(struct {
		Seller            common.Address `json:"seller"`
		SellToken         common.Address `json:"sellToken"`
		ProceedsRecipient common.Address `json:"proceedsRecipient"`
		Remaining         *big.Int       `json:"remaining"`
		Open              bool           `json:"open"`
	})
	return &Offer{Seller: o.Seller, SellToken: o.SellToken, ProceedsRecipient: o.ProceedsRecipient, Remaining: o.Remaining, Open: o.Open}, nil
}

// ReadOfferPrice reads offer id's price in token (zero Num: not accepted).
func ReadOfferPrice(ctx context.Context, r Caller, book common.Address, id *big.Int, token common.Address) (OfferPrice, error) {
	data, _ := offerBookABI.Pack("priceOf", id, token)
	out, err := callRaw(ctx, r, book, data)
	if err != nil {
		return OfferPrice{}, err
	}
	res, err := offerBookABI.Unpack("priceOf", out)
	if err != nil {
		return OfferPrice{}, err
	}
	p := res[0].(struct {
		Num *big.Int `json:"num"`
		Den *big.Int `json:"den"`
	})
	return OfferPrice{Num: p.Num, Den: p.Den}, nil
}

// Tradable reports whether the book lists token.
func Tradable(ctx context.Context, r Caller, book, token common.Address) (bool, error) {
	data, _ := offerBookABI.Pack("tradable", token)
	out, err := callRaw(ctx, r, book, data)
	if err != nil {
		return false, err
	}
	res, err := offerBookABI.Unpack("tradable", out)
	if err != nil {
		return false, err
	}
	return res[0].(bool), nil
}

// IsAuthorizer reports whether the book accepts a's fill authorizations.
func IsAuthorizer(ctx context.Context, r Caller, book, a common.Address) (bool, error) {
	data, _ := offerBookABI.Pack("authorizers", a)
	out, err := callRaw(ctx, r, book, data)
	if err != nil {
		return false, err
	}
	res, err := offerBookABI.Unpack("authorizers", out)
	if err != nil {
		return false, err
	}
	return res[0].(bool), nil
}

// ParseOfferCreated finds the offer seller created in a transaction's logs.
func ParseOfferCreated(logs []*types.Log, book, seller common.Address) (*big.Int, bool) {
	ev := offerBookABI.Events["OfferCreated"]
	for _, l := range logs {
		if l.Address == book && len(l.Topics) == 4 && l.Topics[0] == ev.ID && common.BytesToAddress(l.Topics[2].Bytes()) == seller {
			return new(big.Int).SetBytes(l.Topics[1].Bytes()), true
		}
	}
	return nil, false
}

// Fill is a purchase recorded in a transaction's logs.
type Fill struct {
	OfferID      *big.Int
	Taker        common.Address
	Recipient    common.Address
	PaymentToken common.Address
	Amount       *big.Int
	Payment      *big.Int
}

// ParseFills finds the purchases from book in a transaction's logs.
func ParseFills(logs []*types.Log, book common.Address) []Fill {
	ev := offerBookABI.Events["Filled"]
	var out []Fill
	for _, l := range logs {
		if l.Address != book || len(l.Topics) != 4 || l.Topics[0] != ev.ID {
			continue
		}
		vals, err := ev.Inputs.NonIndexed().Unpack(l.Data)
		if err != nil {
			continue
		}
		out = append(out, Fill{
			OfferID: new(big.Int).SetBytes(l.Topics[1].Bytes()), Taker: common.BytesToAddress(l.Topics[2].Bytes()), Recipient: common.BytesToAddress(l.Topics[3].Bytes()),
			PaymentToken: vals[0].(common.Address), Amount: vals[1].(*big.Int), Payment: vals[2].(*big.Int),
		})
	}
	return out
}

// ERC20Mint mints amount of an Ownable mintable token (TokenizedAsset,
// the internal balance token) to to; only the token's owner can.
func ERC20Mint(token, to common.Address, amount *big.Int) Call {
	data, _ := mintABI.Pack("mint", to, amount)
	return Call{To: token, Value: big.NewInt(0), Data: data}
}

var mintABI = mustABI(`[{"name":"mint","type":"function","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[]}]`)
