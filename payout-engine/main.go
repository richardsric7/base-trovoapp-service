// payout-engine pays tokenized-asset proceeds to token holders from a
// dedicated Safe. It has no HTTP interface: it works from app-backend's
// database (written by tm-api on admin actions) and is woken, and reports
// progress, over Redis. Run exactly one active instance; others wait as
// standbys (see the engine's heartbeat claim).
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"trovo-payout-engine/internal/bus"
	"trovo-payout-engine/internal/chain"
	"trovo-payout-engine/internal/config"
	"trovo-payout-engine/internal/engine"
	"trovo-payout-engine/internal/push"
	"trovo-payout-engine/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file; using the environment")
	}
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("configuration: %v", err)
	}
	signers, from, err := config.SignerSourceFromEnv()
	if err != nil {
		log.Fatalf("signers: %v", err)
	}
	db, err := store.Open()
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	client, err := chain.Dial(cfg.RPCURL)
	if err != nil {
		log.Fatalf("rpc: %v", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	chainID := cfg.ChainID
	if chainID == nil {
		if chainID, err = client.ChainID(ctx); err != nil {
			log.Fatalf("reading the chain id: %v", err)
		}
	}
	sender, err := push.NewFromEnv(ctx)
	if err != nil {
		log.Fatalf("push notifications: %v", err)
	}
	if sender == nil {
		log.Println("GC is not set: holders will not get push notifications")
	}
	b := bus.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisTLS)
	defer b.Close()
	log.Printf("payout-engine %s (%s) on chain %v: payout Safe %s, signers from %s", cfg.Version, cfg.Instance, chainID, cfg.PayoutSafe.Hex(), from)
	e := &engine.Engine{DB: db, Chain: client, ChainID: chainID, Cfg: cfg, Signers: signers, Push: sender, Bus: b}
	e.Run(ctx)
	log.Println("stopped")
}
