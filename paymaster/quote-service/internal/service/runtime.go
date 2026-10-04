// Package service wires configuration, rate sources, the quoter and the
// deposit monitor together, and reloads them when the Vault secret changes.
package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"trovo-paymaster-quote-service/internal/chain"
	"trovo-paymaster-quote-service/internal/config"
	"trovo-paymaster-quote-service/internal/monitor"
	"trovo-paymaster-quote-service/internal/quote"
	"trovo-paymaster-quote-service/internal/rates"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Runtime is one consistent configuration and everything built from it.
type Runtime struct {
	Config    *config.Config
	Book      *rates.Book
	Quoter    *quote.Quoter
	Paymaster *chain.Paymaster
	Deposit   *monitor.Deposit
	client    *ethclient.Client
}

// Build connects to the chain, checks the paymaster matches the
// configuration, builds every rate source and fetches the first rates.
func Build(ctx context.Context, cfg *config.Config) (*Runtime, error) {
	client, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("connecting to RPC_URL: %w", err)
	}
	ok := false
	defer func() {
		if !ok {
			client.Close()
		}
	}()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading chain id: %w", err)
	}
	if chainID.Cmp(cfg.ChainID) != 0 {
		return nil, fmt.Errorf("RPC_URL is chain %s but CHAIN_ID is %s", chainID, cfg.ChainID)
	}

	pm := &chain.Paymaster{Caller: client, Address: cfg.Paymaster, EntryPoint: cfg.EntryPoint, CacheTTL: 15 * time.Second}
	ep, err := pm.ConfiguredEntryPoint(ctx)
	if err != nil {
		return nil, fmt.Errorf("PAYMASTER_ADDRESS %s does not look like TrovoTokenPaymaster: %w", cfg.Paymaster.Hex(), err)
	}
	if ep != cfg.EntryPoint {
		return nil, fmt.Errorf("paymaster uses EntryPoint %s but ENTRYPOINT_ADDRESS is %s", ep.Hex(), cfg.EntryPoint.Hex())
	}
	signer := crypto.PubkeyToAddress(cfg.SignerKey.PublicKey)
	if allowed, err := pm.IsQuoteSigner(ctx, signer); err != nil {
		return nil, err
	} else if !allowed {
		return nil, fmt.Errorf("QUOTE_SIGNER_PRIVATE_KEY's address %s is not a quote signer on the paymaster (owner must call setQuoteSigner)", signer.Hex())
	}

	book, err := rates.NewBook(cfg.RatePairs, rates.Env{
		Chain:  client,
		HTTP:   &http.Client{Timeout: cfg.SourceTimeout},
		Secret: cfg.Secret,
	}, cfg.RateMaxAge, cfg.SourceTimeout)
	if err != nil {
		return nil, fmt.Errorf("RATE_PAIRS: %w", err)
	}
	book.Refresh(ctx)

	dep := &monitor.Deposit{
		Read:       pm.Deposit,
		Low:        cfg.DepositLowWatermark,
		WebhookURL: cfg.AlertWebhookURL,
		Cooldown:   cfg.AlertCooldown,
		Label:      fmt.Sprintf("paymaster %s on chain %s", cfg.Paymaster.Hex(), cfg.ChainID),
	}
	dep.Check(ctx)

	q := &quote.Quoter{
		ChainID: cfg.ChainID, Paymaster: cfg.Paymaster, Key: cfg.SignerKey,
		Tokens: cfg.GasTokens, Spread: cfg.Spread, Rates: book, Chain: pm, Deposit: dep.Last,
		Validity: cfg.QuoteValidity, MaxValidity: cfg.QuoteMaxValidity, ClockSkew: cfg.ClockSkew,
		PaymasterVerificationGasLimit: cfg.PaymasterVerificationGasLimit,
		PaymasterPostOpGasLimit:       cfg.PaymasterPostOpGasLimit,
		Now:                           time.Now,
	}
	for _, t := range cfg.GasTokens {
		if r, err := q.Rate(t); err != nil {
			log.Printf("[startup] %s: no rate yet: %v", t.Symbol, err)
		} else {
			log.Printf("[startup] %s: 1 ETH = %s %s at market, %s quoted (spread %d bps)", t.Symbol, r.MarketRate, t.Symbol, r.QuotedRate, r.SpreadBps)
		}
	}
	ok = true
	return &Runtime{Config: cfg, Book: book, Quoter: q, Paymaster: pm, Deposit: dep, client: client}, nil
}

func (r *Runtime) Close() { r.client.Close() }

// Service holds the current Runtime and runs the background loops.
type Service struct {
	boot config.Bootstrap
	cur  atomic.Pointer[Runtime]
}

func New(boot config.Bootstrap, rt *Runtime) *Service {
	s := &Service{boot: boot}
	s.cur.Store(rt)
	return s
}

// Current is the runtime requests should use.
func (s *Service) Current() *Runtime { return s.cur.Load() }

// Run refreshes rates, checks the deposit and (if configured) reloads the
// configuration until ctx ends.
func (s *Service) Run(ctx context.Context) {
	rt := s.Current()
	rateTick := time.NewTicker(rt.Config.RateRefreshInterval)
	depTick := time.NewTicker(rt.Config.DepositCheckInterval)
	defer rateTick.Stop()
	defer depTick.Stop()
	var reload <-chan time.Time
	if rt.Config.ConfigReloadInterval > 0 {
		t := time.NewTicker(rt.Config.ConfigReloadInterval)
		defer t.Stop()
		reload = t.C
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-rateTick.C:
			s.Current().Book.Refresh(ctx)
		case <-depTick.C:
			s.Current().Deposit.Check(ctx)
		case <-reload:
			s.Reload(ctx)
		}
	}
}

// Reload re-reads the configuration and swaps in a new runtime if it is
// valid; on any error the current runtime keeps serving. Interval and
// listen-address changes need a restart.
func (s *Service) Reload(ctx context.Context) {
	values, err := s.boot.Load(ctx)
	if err != nil {
		log.Printf("[reload] %v (keeping the current configuration)", err)
		return
	}
	cfg, err := config.Parse(values)
	if err != nil {
		log.Printf("[reload] invalid configuration: %v (keeping the current configuration)", err)
		return
	}
	rt, err := Build(ctx, cfg)
	if err != nil {
		log.Printf("[reload] %v (keeping the current configuration)", err)
		return
	}
	old := s.cur.Swap(rt)
	if old.Config.HTTPAddr != cfg.HTTPAddr {
		log.Printf("[reload] HTTP_ADDR changed; restart the service to apply it")
	}
	log.Printf("[reload] configuration reloaded from %s", s.boot.Describe())
	// let in-flight requests on the old runtime finish before closing it
	time.AfterFunc(time.Minute, old.Close)
}
