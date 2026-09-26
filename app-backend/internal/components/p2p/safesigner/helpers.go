package safesigner

import (
	"io"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func stringsReader(s string) io.Reader {
	return strings.NewReader(s)
}

func mustBytes32Type() abi.Type {
	t, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		panic("safesigner: invalid bytes32 type: " + err.Error())
	}
	return t
}

func toBytes32(b []byte) [32]byte {
	var out [32]byte
	copy(out[:], b)
	return out
}

func callMsg(to common.Address, data []byte) ethereum.CallMsg {
	return ethereum.CallMsg{To: &to, Data: data}
}
