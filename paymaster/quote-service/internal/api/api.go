// Package api is the quote service's HTTP API.
package api

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"trovo-paymaster-quote-service/internal/monitor"
	"trovo-paymaster-quote-service/internal/quote"
	"trovo-paymaster-quote-service/internal/rates"
	"trovo-paymaster-quote-service/internal/service"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Quantity is a JSON number, decimal string or 0x-hex string.
type Quantity struct{ *big.Int }

func (q *Quantity) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "null" {
		return nil
	}
	s = strings.Trim(s, `"`)
	v := new(big.Int)
	var ok bool
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		v, ok = v.SetString(s[2:], 16)
		if s == "0x" || s == "0X" {
			v, ok = big.NewInt(0), true
		}
	} else {
		v, ok = v.SetString(s, 10)
	}
	if !ok || v.Sign() < 0 {
		return fmt.Errorf("invalid quantity %s", string(b))
	}
	q.Int = v
	return nil
}

// UserOperation is an ERC-4337 v0.7 UserOperation in its RPC form, without
// the paymaster fields and signature (those come from the quote).
type UserOperation struct {
	Sender                        common.Address  `json:"sender" swaggertype:"string" example:"0x5a6B0c1C6f1C0F4b1d7aA2d1C3e0bEe8a9F2d4c1"`
	Nonce                         Quantity        `json:"nonce" swaggertype:"string" example:"0x0"`
	Factory                       *common.Address `json:"factory,omitempty" swaggertype:"string" example:"0x4e1DCf7AD4e460CfD30791CCC4F9c8a4f820ec67"`
	FactoryData                   hexutil.Bytes   `json:"factoryData,omitempty" swaggertype:"string" example:"0x1688f0b9..."`
	InitCode                      hexutil.Bytes   `json:"initCode,omitempty" swaggertype:"string" example:"0x"`
	CallData                      hexutil.Bytes   `json:"callData" swaggertype:"string" example:"0x7bb37428..."`
	CallGasLimit                  Quantity        `json:"callGasLimit" swaggertype:"string" example:"0x30d40"`
	VerificationGasLimit          Quantity        `json:"verificationGasLimit" swaggertype:"string" example:"0x493e0"`
	PreVerificationGas            Quantity        `json:"preVerificationGas" swaggertype:"string" example:"0xea60"`
	MaxFeePerGas                  Quantity        `json:"maxFeePerGas" swaggertype:"string" example:"0x3b9aca00"`
	MaxPriorityFeePerGas          Quantity        `json:"maxPriorityFeePerGas" swaggertype:"string" example:"0xf4240"`
	PaymasterVerificationGasLimit Quantity        `json:"paymasterVerificationGasLimit,omitempty" swaggertype:"string" example:"0x249f0"`
	PaymasterPostOpGasLimit       Quantity        `json:"paymasterPostOpGasLimit,omitempty" swaggertype:"string" example:"0x13880"`
}

// QuoteRequest asks for a signed gas quote.
type QuoteRequest struct {
	// Gas token symbol ("USDC") or address.
	Token string `json:"token" example:"cNGN"`
	// How long the quote stays valid; 0 = the service default. Shared
	// wallets waiting for approvals ask for longer (up to QUOTE_MAX_VALIDITY).
	ValiditySeconds int64         `json:"validitySeconds" example:"600"`
	UserOp          UserOperation `json:"userOp"`
}

// ErrorResponse is every error body.
type ErrorResponse struct {
	Code    string `json:"code" example:"rate-unavailable"`
	Message string `json:"message" example:"rate USD/NGN is stale"`
}

// StatusResponse describes the service.
type StatusResponse struct {
	ChainID    string         `json:"chainId"`
	Paymaster  common.Address `json:"paymaster" swaggertype:"string"`
	EntryPoint common.Address `json:"entryPoint" swaggertype:"string"`
	Signer     common.Address `json:"signer" swaggertype:"string"`
	Deposit    monitor.Status `json:"deposit"`
}

type handlers struct{ svc *service.Service }

// Router builds the HTTP API.
func Router(svc *service.Service) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	h := handlers{svc}

	r.GET("/healthz", h.healthz)
	r.GET("/readyz", h.readyz)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/v1", h.auth)
	v1.POST("/quote", h.quote)
	v1.GET("/tokens", h.tokens)
	v1.GET("/rates", h.rates)
	v1.GET("/status", h.status)
	return r
}

func (h handlers) auth(c *gin.Context) {
	key := c.GetHeader("X-API-Key")
	if key == "" {
		key = strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	}
	for _, k := range h.svc.Current().Config.APIKeys {
		if key != "" && subtle.ConstantTimeCompare([]byte(key), []byte(k)) == 1 {
			c.Next()
			return
		}
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Code: "unauthorized", Message: "missing or invalid API key"})
}

// healthz godoc
// @Summary      Liveness
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /healthz [get]
func (h handlers) healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// readyz godoc
// @Summary      Readiness: every gas token has a fresh rate and the deposit is known
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      503  {object}  map[string]interface{}
// @Router       /readyz [get]
func (h handlers) readyz(c *gin.Context) {
	rt := h.svc.Current()
	problems := map[string]string{}
	for _, t := range rt.Config.GasTokens {
		if _, err := rt.Quoter.Rate(t); err != nil {
			problems[t.Symbol] = err.Error()
		}
	}
	if _, ok := rt.Deposit.Last(); !ok {
		problems["deposit"] = "not read yet"
	}
	if len(problems) > 0 {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "problems": problems})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

// quote godoc
// @Summary      Sign a gas quote for a UserOperation
// @Description  Prices the operation's gas in the chosen stablecoin (market rate plus Trovo's spread), signs it, and returns the paymaster fields to put in the UserOperation before the wallet owner signs it. The quote is bound to every field of the operation; changing any of them after quoting makes the paymaster reject it.
// @Tags         quotes
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      QuoteRequest  true  "Operation and gas token"
// @Success      200      {object}  quote.Result
// @Failure      400      {object}  ErrorResponse  "invalid request or unknown token"
// @Failure      401      {object}  ErrorResponse
// @Failure      409      {object}  ErrorResponse  "the wallet owes the paymaster"
// @Failure      503      {object}  ErrorResponse  "rates, chain or deposit unavailable"
// @Router       /v1/quote [post]
func (h handlers) quote(c *gin.Context) {
	var req QuoteRequest
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "invalid-request", Message: err.Error()})
		return
	}
	op, err := req.UserOp.toOp()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "invalid-user-operation", Message: err.Error()})
		return
	}
	res, err := h.svc.Current().Quoter.Quote(c.Request.Context(), quote.Request{
		Op: op, Token: req.Token, Validity: time.Duration(req.ValiditySeconds) * time.Second,
	})
	if err != nil {
		var qe *quote.Error
		if errors.As(err, &qe) {
			c.JSON(qe.Status, ErrorResponse{Code: qe.Code, Message: qe.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Code: "internal", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (u UserOperation) toOp() (quote.UserOp, error) {
	initCode := []byte(u.InitCode)
	if u.Factory != nil {
		if len(initCode) != 0 {
			return quote.UserOp{}, fmt.Errorf("give either initCode or factory/factoryData, not both")
		}
		initCode = append(u.Factory.Bytes(), u.FactoryData...)
	}
	return quote.UserOp{
		Sender: u.Sender, Nonce: u.Nonce.Int, InitCode: initCode, CallData: u.CallData,
		VerificationGasLimit: u.VerificationGasLimit.Int, CallGasLimit: u.CallGasLimit.Int,
		PreVerificationGas: u.PreVerificationGas.Int, MaxFeePerGas: u.MaxFeePerGas.Int,
		MaxPriorityFeePerGas:          u.MaxPriorityFeePerGas.Int,
		PaymasterVerificationGasLimit: u.PaymasterVerificationGasLimit.Int,
		PaymasterPostOpGasLimit:       u.PaymasterPostOpGasLimit.Int,
	}, nil
}

// TokenInfo is a gas token and its current rate, if available.
type TokenInfo struct {
	quote.TokenRate
	Error string `json:"error,omitempty"`
}

// tokens godoc
// @Summary      Gas tokens and the current price of 1 ETH in each
// @Description  For showing users an estimated fee before they send: fee in token = gas (wei) × exchangeRate / 1e18 base units.
// @Tags         quotes
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   TokenInfo
// @Failure      401  {object}  ErrorResponse
// @Router       /v1/tokens [get]
func (h handlers) tokens(c *gin.Context) {
	rt := h.svc.Current()
	out := make([]TokenInfo, 0, len(rt.Config.GasTokens))
	for _, t := range rt.Config.GasTokens {
		r, err := rt.Quoter.Rate(t)
		info := TokenInfo{TokenRate: r}
		if err != nil {
			info.TokenRate = quote.TokenRate{Token: t.Address, Symbol: t.Symbol, Decimals: t.Decimals, SpreadBps: rt.Config.Spread(t)}
			info.Error = err.Error()
		}
		out = append(out, info)
	}
	c.JSON(http.StatusOK, out)
}

// rates godoc
// @Summary      Every rate pair with each source's latest answer
// @Tags         quotes
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   rates.Price
// @Failure      401  {object}  ErrorResponse
// @Router       /v1/rates [get]
func (h handlers) rates(c *gin.Context) {
	c.JSON(http.StatusOK, h.svc.Current().Book.Snapshot())
}

var _ = rates.Price{} // referenced by the swagger annotations

// status godoc
// @Summary      Paymaster, signer and deposit status
// @Tags         health
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  StatusResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /v1/status [get]
func (h handlers) status(c *gin.Context) {
	rt := h.svc.Current()
	c.JSON(http.StatusOK, StatusResponse{
		ChainID: rt.Config.ChainID.String(), Paymaster: rt.Config.Paymaster, EntryPoint: rt.Config.EntryPoint,
		Signer: rt.Quoter.Signer(), Deposit: rt.Deposit.Status(),
	})
}
