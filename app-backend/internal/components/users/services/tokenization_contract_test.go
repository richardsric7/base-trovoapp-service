package users

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"

	tErrors "trovo-wallet-api/internal/errors"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// fakeToken answers the token contract calls verifyTokenizedAssetContract
// makes; a method missing from results reverts, like a contract that does
// not implement it.
type fakeToken struct {
	code    map[common.Address]bool
	results map[string][]interface{}
}

func (f *fakeToken) CodeAt(_ context.Context, a common.Address, _ *big.Int) ([]byte, error) {
	if f.code[a] {
		return []byte{0x60}, nil
	}
	return nil, nil
}

func (f *fakeToken) CallContract(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	for name, method := range tokenizedAssetTokenABI.Methods {
		if bytes.Equal(call.Data[:4], method.ID) {
			out, ok := f.results[name]
			if !ok {
				return nil, errors.New("execution reverted")
			}
			return method.Outputs.Pack(out...)
		}
	}
	return nil, errors.New("execution reverted")
}

var (
	testToken = common.HexToAddress("0x1000000000000000000000000000000000000001")
	testSafe  = common.HexToAddress("0x2000000000000000000000000000000000000002")
)

func validToken() *fakeToken {
	return &fakeToken{
		code: map[common.Address]bool{testToken: true, testSafe: true},
		results: map[string][]interface{}{
			"symbol":      {"REIT1"},
			"decimals":    {uint8(18)},
			"totalSupply": {big.NewInt(0)},
			"hasRole":     {true},
		},
	}
}

func TestVerifyTokenizedAssetContract(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name    string
		mutate  func(f *fakeToken)
		wantErr string
	}{
		{"minter role holder", func(f *fakeToken) {}, ""},
		{"ownable, owned by the Safe", func(f *fakeToken) {
			delete(f.results, "hasRole")
			f.results["owner"] = []interface{}{testSafe}
		}, ""},
		{"symbol case-insensitive", func(f *fakeToken) { f.results["symbol"] = []interface{}{"reit1"} }, ""},
		{"no code at token address", func(f *fakeToken) { delete(f.code, testToken) }, "No contract is deployed"},
		{"issuing wallet is not a Safe", func(f *fakeToken) { delete(f.code, testSafe) }, "not a deployed Safe"},
		{"symbol mismatch", func(f *fakeToken) { f.results["symbol"] = []interface{}{"OTHER"} }, "does not match the asset code"},
		{"no decimals", func(f *fakeToken) { delete(f.results, "decimals") }, "decimals()"},
		{"already has supply", func(f *fakeToken) { f.results["totalSupply"] = []interface{}{big.NewInt(1)} }, "already has tokens in circulation"},
		{"Safe cannot mint", func(f *fakeToken) {
			f.results["hasRole"] = []interface{}{false}
			f.results["owner"] = []interface{}{common.HexToAddress("0x3000000000000000000000000000000000000003")}
		}, "cannot mint"},
		{"neither AccessControl nor Ownable", func(f *fakeToken) { delete(f.results, "hasRole") }, "cannot mint"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := validToken()
			tc.mutate(f)
			err := verifyTokenizedAssetContract(ctx, f, testToken, testSafe, "REIT1")
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("expected the contract to be accepted, got %v", err)
				}
				return
			}
			ge, ok := err.(tErrors.GenericError)
			if !ok || ge.ErrorType() != "error-invalid-token-contract" || !strings.Contains(ge.Message(), tc.wantErr) {
				t.Fatalf("expected an error-invalid-token-contract error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}
