package users

import (
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

func WalletCountViewOnlyAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "VIEW-ONLY" {
			accessCount++
		}
	}

	return
}
func WalletHasViewOnlyAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	viewOnly = true
	if wallet.ManagedAccessEnabled == 0 {
		return false
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel != "VIEW-ONLY" {
			return false
		}
	}

	return
}
func WalletCountAuthorizerAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "AUTHORIZER" {
			accessCount++
		}
	}

	return
}
func WalletCountInitiatorAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "INITIATOR" {
			accessCount++
		}
	}

	return
}

func PublicKeyCountViewOnlyAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}
	wallet, err := userModels.UserWalletID(publicKey).GetWallet(gc.DB)
	if err != nil {
		return 0
	}

	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "VIEW-ONLY" {
			accessCount++
		}
	}

	return
}

func PublicKeyHasViewOnlyAccess(publicKey string, gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	viewOnly = true
	if publicKey == "" {
		return false
	}
	wallet, err := userModels.UserWalletID(publicKey).GetWallet(gc.DB)
	if err != nil {
		return false
	}

	if wallet.ManagedAccessEnabled == 0 {
		return false
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel != "VIEW-ONLY" {
			return false
		}
	}

	return
}

func PublicKeyCountAuthorizerAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}
	wallet, err := userModels.UserWalletID(publicKey).GetWallet(gc.DB)
	if err != nil {
		return 0
	}

	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "AUTHORIZER" {
			accessCount++
		}
	}

	return
}

func PublicKeyCountInitiatorAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}
	wallet, err := userModels.UserWalletID(publicKey).GetWallet(gc.DB)
	if err != nil {
		return 0
	}

	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "INITIATOR" {
			accessCount++
		}
	}

	return
}
