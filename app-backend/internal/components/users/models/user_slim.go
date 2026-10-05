package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
)

// A slim user is the users row alone: every column (so it can be saved back
// like any loaded user), but none of the associations GetFullUser preloads -
// wallets with their permissions, wallets shared with the user, patron
// membership, closed groups and fiat payment methods - which cost about six
// more queries per load. Use it wherever only the user's own fields are
// needed (authentication, username, push token, KYC, suspension, fee
// exemption, ...); code that reads UserWallets, WalletsSharedWithUser or the
// other associations must use GetFullUser / GetUserFromPrimarySigner.
//
// Slim users are cached under their own prefix (a slim object under
// "userObj" would look like a user without wallets to full readers) and are
// invalidated with the user's other cache entries (InvalidateUserCache).

const slimUserCachePrefix = "userLite "

const slimUserCacheSeconds = 2000

// slimUserCacheKeys are the keys a slim user is cached under.
func slimUserCacheKeys(u *User) []string {
	keys := make([]string, 0, 4)
	for _, v := range []string{u.Username, u.Email, u.PrimarySigner, u.ID} {
		if strings.TrimSpace(v) != "" {
			keys = append(keys, slimUserCachePrefix+v)
		}
	}
	return keys
}

func cachedSlimUser(key string, gc *sharedconfig.GlobalConfig) (User, bool) {
	var u User
	if gc == nil || gc.RedisCache == nil {
		return u, false
	}
	ok, raw := gc.RedisCache.GetCachedResultRaw(slimUserCachePrefix + key)
	if !ok || json.Unmarshal(raw, &u) != nil || u.ID == "" {
		return u, false
	}
	return u, true
}

func storeSlimUser(u *User, gc *sharedconfig.GlobalConfig) {
	if gc == nil || gc.RedisCache == nil {
		return
	}
	for _, k := range slimUserCacheKeys(u) {
		gc.RedisCache.StoreResultToCacheRaw(k, *u, slimUserCacheSeconds)
	}
}

func slimUserError(e error, param, what string) error {
	if errors.Is(e, gorm.ErrRecordNotFound) {
		return &tErrors.CustomError{Param: param, Err: "error-account-not-found", ErrMessage: what + " not found", Code: http.StatusNotFound}
	}
	return &tErrors.ErrorTemporaryServerError{}
}

// GetSlimUser loads a user by username or ID without associations.
func (u Username) GetSlimUser(db *gorm.DB, gc *sharedconfig.GlobalConfig) (owner User, err error) {
	key := strings.TrimSpace(string(u))
	if cached, ok := cachedSlimUser(key, gc); ok {
		return cached, nil
	}
	if e := db.Where("(username = ? OR id = ?)", key, key).First(&owner).Error; e != nil {
		return owner, slimUserError(e, "id", "Account")
	}
	storeSlimUser(&owner, gc)
	return owner, nil
}

// GetSlimUserBySigner loads the user whose primary signer is signer,
// without associations.
func GetSlimUserBySigner(signer string, db *gorm.DB, gc *sharedconfig.GlobalConfig) (owner User, err error) {
	key := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(signer), " ", ""))
	if cached, ok := cachedSlimUser(key, gc); ok {
		return cached, nil
	}
	if e := db.Where("primary_signer = ?", key).First(&owner).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return owner, &tErrors.CustomError{Param: "primarySigner", Err: "error primary signer does not exist", ErrMessage: "PrimarySigner does not exist", Code: http.StatusNotFound}
		}
		return owner, &tErrors.ErrorTemporaryServerError{}
	}
	storeSlimUser(&owner, gc)
	return owner, nil
}

// InvalidateSlimUserCache drops a user's slim cache entries.
func InvalidateSlimUserCache(u *User, gc *sharedconfig.GlobalConfig) {
	if u == nil || gc == nil || gc.RedisCache == nil {
		return
	}
	if keys := slimUserCacheKeys(u); len(keys) > 0 {
		gc.RedisCache.DeleteFromCache(keys...)
	}
}

// WalletsForInvalidation are the user's wallets as loaded, or read from
// the database when the user was loaded slim, so a slim user invalidates
// its wallets' cache entries too.
func (u *User) WalletsForInvalidation(gc *sharedconfig.GlobalConfig) []UserWallet {
	if len(u.UserWallets) > 0 || u.ID == "" || gc == nil || gc.DB == nil {
		return u.UserWallets
	}
	var wallets []UserWallet
	gc.DB.Select("id", "alias", "signer", "user_id").Where("user_id = ?", u.ID).Find(&wallets)
	return wallets
}

// SlimUserCacheKeysForTest exposes the slim cache keys to tests.
func SlimUserCacheKeysForTest(u *User) []string { return slimUserCacheKeys(u) }

// CachedSlimUser returns the slim user cached under key (a username, email,
// primary signer or user ID), if any.
func CachedSlimUser(key string, gc *sharedconfig.GlobalConfig) (User, bool) {
	return cachedSlimUser(strings.TrimSpace(key), gc)
}

// CacheSlimUser caches a user loaded without associations under its
// username, email, primary signer and ID.
func CacheSlimUser(u *User, gc *sharedconfig.GlobalConfig) { storeSlimUser(u, gc) }
