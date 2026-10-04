package aa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

// QuoteClient calls the paymaster quote service
// (paymaster/quote-service, PAYMASTER_QUOTE_SERVICE_URL).
type QuoteClient struct {
	URL    string
	APIKey string
	HTTP   *http.Client
}

func NewQuoteClient(url, apiKey string) *QuoteClient {
	return &QuoteClient{URL: url, APIKey: apiKey, HTTP: &http.Client{Timeout: 15 * time.Second}}
}

// Quote is the quote service's answer (see paymaster/INTEGRATION.md).
type Quote struct {
	Paymaster                     common.Address  `json:"paymaster"`
	PaymasterVerificationGasLimit *hexutil.Big    `json:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       *hexutil.Big    `json:"paymasterPostOpGasLimit"`
	PaymasterAndData              hexutil.Bytes   `json:"paymasterAndData"`
	Token                         common.Address  `json:"token"`
	Symbol                        string          `json:"symbol"`
	ExchangeRate                  *hexutil.Big    `json:"exchangeRate"`
	MarketRate                    decimal.Decimal `json:"marketRate"`
	QuotedRate                    decimal.Decimal `json:"quotedRate"`
	SpreadBps                     int64           `json:"spreadBps"`
	ValidAfter                    uint64          `json:"validAfter"`
	ValidUntil                    uint64          `json:"validUntil"`
	MaxTokenCost                  *hexutil.Big    `json:"maxTokenCost"`
	MaxTokenCostFormatted         decimal.Decimal `json:"maxTokenCostFormatted"`
}

// QuoteError is a refusal from the quote service.
type QuoteError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *QuoteError) Error() string {
	return fmt.Sprintf("paymaster quote %d %s: %s", e.Status, e.Code, e.Message)
}

// Temporary reports whether paying in the token may work later (as opposed
// to a malformed request).
func (e *QuoteError) Temporary() bool { return e.Status == http.StatusServiceUnavailable }

type quoteUserOp struct {
	Sender                        common.Address `json:"sender"`
	Nonce                         *hexutil.Big   `json:"nonce"`
	InitCode                      hexutil.Bytes  `json:"initCode,omitempty"`
	CallData                      hexutil.Bytes  `json:"callData"`
	CallGasLimit                  *hexutil.Big   `json:"callGasLimit"`
	VerificationGasLimit          *hexutil.Big   `json:"verificationGasLimit"`
	PreVerificationGas            *hexutil.Big   `json:"preVerificationGas"`
	MaxFeePerGas                  *hexutil.Big   `json:"maxFeePerGas"`
	MaxPriorityFeePerGas          *hexutil.Big   `json:"maxPriorityFeePerGas"`
	PaymasterVerificationGasLimit *hexutil.Big   `json:"paymasterVerificationGasLimit,omitempty"`
	PaymasterPostOpGasLimit       *hexutil.Big   `json:"paymasterPostOpGasLimit,omitempty"`
}

// Quote asks for a signed quote for op paid in token. pmVerificationGas /
// pmPostOpGas may be nil (the service's defaults).
func (q *QuoteClient) Quote(ctx context.Context, op UserOperation, token common.Address, validity time.Duration, pmVerificationGas, pmPostOpGas *big.Int) (*Quote, error) {
	req := map[string]interface{}{
		"token":           token.Hex(),
		"validitySeconds": int64(validity / time.Second),
		"userOp": quoteUserOp{
			Sender: op.Sender, Nonce: hb(op.Nonce), InitCode: op.InitCode, CallData: op.CallData,
			CallGasLimit: hb(op.CallGasLimit), VerificationGasLimit: hb(op.VerificationGasLimit),
			PreVerificationGas: hb(op.PreVerificationGas), MaxFeePerGas: hb(op.MaxFeePerGas),
			MaxPriorityFeePerGas:          hb(op.MaxPriorityFeePerGas),
			PaymasterVerificationGasLimit: (*hexutil.Big)(pmVerificationGas),
			PaymasterPostOpGasLimit:       (*hexutil.Big)(pmPostOpGas),
		},
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, q.URL+"/v1/quote", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", q.APIKey)
	res, err := q.HTTP.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("paymaster quote service: %w", err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		qe := &QuoteError{Status: res.StatusCode}
		_ = json.Unmarshal(raw, qe)
		return nil, qe
	}
	var out Quote
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("paymaster quote service: invalid response: %w", err)
	}
	return &out, nil
}
