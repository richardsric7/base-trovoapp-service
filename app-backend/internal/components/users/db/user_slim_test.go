package users

import (
	"context"
	"os"
	"testing"

	"trovo-wallet-api/internal/cache"
	userModels "trovo-wallet-api/internal/components/users/models"
	appdb "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// slimTestGC is the app schema on SQLite with a real Redis cache
// (TEST_REDIS_ADDR, default 127.0.0.1:6379; skipped when none answers).
func slimTestGC(t *testing.T) *sharedconfig.GlobalConfig {
	t.Helper()
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	client := redis.NewClient(&redis.Options{Addr: addr})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skipf("no Redis at %s: %v", addr, err)
	}
	t.Cleanup(func() { client.Close() })
	// keys of this run only
	t.Setenv("CACHING_PARAMETER", "slimtest-"+t.Name())
	db, err := gorm.Open(sqlite.Open("file:"+t.TempDir()+"/slim.db"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	appdb.MigrateDB(db)
	return &sharedconfig.GlobalConfig{DB: db, RedisCache: &cache.RedisCache{Enabled: true, Client: client, Context: context.Background()}}
}

const (
	slimSigner = "0X00000000000000000000000000000000000000A1"
	slimWallet = "0X00000000000000000000000000000000000000B2"
)

func seedSlimUser(t *testing.T, gc *sharedconfig.GlobalConfig) {
	t.Helper()
	token := "token-1"
	mobile := "+2348000000001"
	u := userModels.User{ID: "7", Username: "alice", Email: "alice@example.com", PrimarySigner: slimSigner, Address: slimWallet, PushNotificationToken: &token, Mobile: &mobile}
	if err := gc.DB.Omit("UserWallets").Create(&u).Error; err != nil {
		t.Fatal(err)
	}
	w := userModels.UserWallet{ID: slimWallet, Alias: "alice_main", Signer: slimSigner, UserID: "7", PrimaryWallet: 1}
	if err := gc.DB.Create(&w).Error; err != nil {
		t.Fatal(err)
	}
}

func TestSlimUserLookups(t *testing.T) {
	gc := slimTestGC(t)
	seedSlimUser(t, gc)
	for _, id := range []string{"alice", "alice@example.com", "7", "+2348000000001", "alice_main", slimWallet} {
		u, err := GetSlimUser(id, gc.DB, gc)
		if err != nil || u.Username != "alice" {
			t.Fatalf("GetSlimUser(%q) = %q, %v", id, u.Username, err)
		}
		if len(u.UserWallets) != 0 {
			t.Fatalf("GetSlimUser(%q) loaded wallets", id)
		}
	}
	if u, err := GetSlimUserFromPrimarySigner(slimSigner, gc.DB, gc); err != nil || u.Username != "alice" {
		t.Fatalf("by signer: %q %v", u.Username, err)
	}
	if u, err := userModels.Username("7").GetSlimUser(gc.DB, gc); err != nil || u.Username != "alice" {
		t.Fatalf("by id: %q %v", u.Username, err)
	}
	if _, err := GetSlimUser("nobody", gc.DB, gc); err == nil {
		t.Fatal("unknown user found")
	}
	// the full loader keeps its own cache and still gets wallets
	full, err := userModels.Username("alice").GetFullUser(gc.DB, gc)
	if err != nil || len(full.UserWallets) != 1 {
		t.Fatalf("full user after slim loads: %d wallets, %v", len(full.UserWallets), err)
	}
	// GetUser by wallet address
	if byAddr, err := GetUser(slimWallet, gc.DB, gc); err != nil || byAddr.Username != "alice" {
		t.Fatalf("GetUser by address: %q %v", byAddr.Username, err)
	}
	// and a full object cached first is not served to slim readers
	if u, ok := userModels.CachedSlimUser("alice", gc); !ok || len(u.UserWallets) != 0 {
		t.Fatalf("slim cache holds wallets (%v)", ok)
	}
}

func TestSlimUserCacheInvalidation(t *testing.T) {
	gc := slimTestGC(t)
	seedSlimUser(t, gc)
	u, err := userModels.Username("alice").GetSlimUser(gc.DB, gc)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range userModels.SlimUserCacheKeysForTest(&u) {
		if ok, _ := gc.RedisCache.GetCachedResultRaw(k); !ok {
			t.Fatalf("%q not cached", k)
		}
	}
	// served from the cache: a direct DB change is not seen yet
	gc.DB.Model(&userModels.User{}).Where("id = ?", "7").Update("push_notification_token", "token-2")
	if c, _ := GetSlimUserFromPrimarySigner(slimSigner, gc.DB, gc); c.PushNotificationToken == nil || *c.PushNotificationToken != "token-1" {
		t.Fatal("expected the cached slim user")
	}
	// a wallet-keyed full entry, which a slim user has no wallets to name
	gc.RedisCache.StoreResultToCacheRaw("userObj alice_main", u, 60)

	// invalidating through the slim user clears its entries and its
	// wallets' entries
	u.InvalidateUserCache(gc)
	for _, k := range append(userModels.SlimUserCacheKeysForTest(&u), "userObj alice_main") {
		if ok, _ := gc.RedisCache.GetCachedResultRaw(k); ok {
			t.Fatalf("%q still cached", k)
		}
	}
	if c, _ := GetSlimUser("alice@example.com", gc.DB, gc); c.PushNotificationToken == nil || *c.PushNotificationToken != "token-2" {
		t.Fatal("expected the updated user after invalidation")
	}
	// the users/db invalidation path clears slim entries too
	InvalidateUserWalletCache(&u, gc)
	if _, ok := userModels.CachedSlimUser("alice", gc); ok {
		t.Fatal("slim entry survived InvalidateUserWalletCache")
	}
}
