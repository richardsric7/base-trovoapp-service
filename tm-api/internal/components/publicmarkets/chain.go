package publicmarkets

import (
	"context"
	"errors"
	"math/big"
	"os"
	"strings"

	"admin-panel-dashboard/internal/network"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// TokenInfo is what Trovo Manager checks of an asset's token contract.
type TokenInfo struct {
	Owner       common.Address
	Decimals    int
	TotalSupply *big.Int
}

// ContractReader reads a token contract on Base.
type ContractReader interface {
	Token(ctx context.Context, contract common.Address) (TokenInfo, error)
}

// NewContractReader reads through the shared RPC client; nil when
// BASE_RPC_URL is not set (contracts then cannot be registered).
func NewContractReader() ContractReader {
	if strings.TrimSpace(os.Getenv("BASE_RPC_URL")) == "" {
		return nil
	}
	return rpcReader{}
}

type rpcReader struct{}

var (
	selOwner       = common.FromHex("0x8da5cb5b") // owner()
	selDecimals    = common.FromHex("0x313ce567") // decimals()
	selTotalSupply = common.FromHex("0x18160ddd") // totalSupply()
)

func (rpcReader) Token(ctx context.Context, contract common.Address) (TokenInfo, error) {
	client := network.GetBlockchainClient()
	call := func(sel []byte) ([]byte, error) {
		out, err := client.CallContract(ctx, ethereum.CallMsg{To: &contract, Data: sel}, nil)
		if err != nil {
			return nil, err
		}
		if len(out) < 32 {
			return nil, errors.New("not a token contract")
		}
		return out, nil
	}
	var info TokenInfo
	code, err := client.CodeAt(ctx, contract, nil)
	if err != nil {
		return info, err
	}
	if len(code) == 0 {
		return info, errors.New("no contract at this address")
	}
	out, err := call(selOwner)
	if err != nil {
		return info, err
	}
	info.Owner = common.BytesToAddress(out[12:32])
	if out, err = call(selDecimals); err != nil {
		return info, err
	}
	info.Decimals = int(new(big.Int).SetBytes(out[:32]).Int64())
	if out, err = call(selTotalSupply); err != nil {
		return info, err
	}
	info.TotalSupply = new(big.Int).SetBytes(out[:32])
	return info, nil
}
