// Command quote-service prices and signs gas quotes for TrovoTokenPaymaster.
//
// @title                       Trovo Paymaster Quote Service
// @version                     1.0
// @description                 Signs ERC-4337 paymaster quotes so Trovo wallets can pay gas in stablecoins (USDC, USDT, cNGN, ...) at market rates plus Trovo's spread.
// @BasePath                    /
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        X-API-Key
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "trovo-paymaster-quote-service/docs"
	"trovo-paymaster-quote-service/internal/api"
	"trovo-paymaster-quote-service/internal/config"
	"trovo-paymaster-quote-service/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	boot, err := config.BootstrapFromEnv()
	if err != nil {
		log.Fatalf("bootstrap: %v", err)
	}
	values, err := boot.Load(ctx)
	if err != nil {
		log.Fatalf("loading configuration: %v", err)
	}
	cfg, err := config.Parse(values)
	if err != nil {
		log.Fatalf("configuration from %s: %v", boot.Describe(), err)
	}
	rt, err := service.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("startup: %v", err)
	}
	svc := service.New(boot, rt)
	log.Printf("quote service: config from %s, chain %s, paymaster %s, signer %s", boot.Describe(), cfg.ChainID, cfg.Paymaster.Hex(), rt.Quoter.Signer().Hex())

	go svc.Run(ctx)

	// SIGHUP reloads the configuration (e.g. after changing a spread in Vault)
	hup := make(chan os.Signal, 1)
	signal.Notify(hup, syscall.SIGHUP)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-hup:
				svc.Reload(ctx)
			}
		}
	}()

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: api.Router(svc), ReadHeaderTimeout: 10 * time.Second}
	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdown)
	svc.Current().Close()
}
