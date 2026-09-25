package services

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	authServices "admin-panel-dashboard/internal/components/auth/services"
	p2pErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"
	"admin-panel-dashboard/internal/trovosdk"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LinkWalletOnly links Trovo Manager with Trovo Wallet in the DB only. Allowed once per user; rejects if already linked.
func LinkWalletOnly(server *serverModels.Server, c *gin.Context, walletUsername string) error {
	walletUsername = strings.ToLower(strings.TrimSpace(walletUsername))
	if !strings.HasSuffix(walletUsername, "@trovo") {
		return errors.New("invalid wallet username format (must contain @trovo)")
	}

	memberID := c.GetString("member_id")
	if memberID == "" {
		return errors.New("organization member context required; ensure you are authenticated as an organization member")
	}

	// Require wallet username to be an active Trovo user
	if !authServices.IsTrovoUser(server, walletUsername) {
		return errors.New("trovo wallet user not found or suspended")
	}

	return server.AdminDB.Transaction(func(tx *gorm.DB) error {
		var member models.OrganizationMember
		if err := tx.Where("id = ?", memberID).First(&member).Error; err != nil {
			return err
		}

		if member.IsWalletLinked && member.TrovoWalletUsername != nil && *member.TrovoWalletUsername != "" {
			return errors.New("wallet already linked; each user may link only once")
		}

		slog.Info("Linking Trovo wallet to member profile", "member_id", memberID, "wallet_username", walletUsername)
		member.TrovoWalletUsername = &walletUsername
		member.IsWalletLinked = true
		member.UpdatedAt = time.Now()

		if err := tx.Save(&member).Error; err != nil {
			return err
		}
		return nil
	})
}

// SendWalletLinkAuthorizationRequest starts wallet linking for an unlinked organization member.
// The link is persisted only after VerifyWalletLink confirms approval in Trovo Wallet.
func SendWalletLinkAuthorizationRequest(server *serverModels.Server, c *gin.Context, walletUsername string) (*trovosdk.TrovoWalletAuthorizationData, error) {
	memberID := c.GetString("member_id")
	if memberID == "" {
		return nil, errors.New("organization member context required")
	}

	var member models.OrganizationMember
	if err := server.AdminDB.Where("id = ?", memberID).First(&member).Error; err != nil {
		return nil, errors.New("member not found")
	}
	if member.IsWalletLinked {
		return nil, errors.New("wallet already linked; each user may link only once")
	}

	walletUsername = strings.ToLower(strings.TrimSpace(walletUsername))
	if !strings.HasSuffix(walletUsername, "@trovo") {
		return nil, errors.New("invalid wallet username format (must contain @trovo)")
	}
	if !authServices.IsTrovoUser(server, walletUsername) {
		return nil, errors.New("trovo wallet user not found or suspended")
	}
	cleanUsername := strings.TrimSuffix(walletUsername, "@trovo")
	deviceInfo := "Trovo Manager"
	callbackUrl := fmt.Sprintf("%v/%v", os.Getenv("LOGIN_CALLBACK_URL"), "trovo")

	authData, err := server.GC.ServiceLink.SendAuthorizationRequest(
		cleanUsername,
		"Authorize Trovo Manager",
		deviceInfo,
		callbackUrl,
		10,
	)
	if err != nil {
		log.Printf("[SendWalletLinkAuthorizationRequest] Error: %v", err)
		return nil, fmt.Errorf("failed to send authorization request: %v", err)
	}
	return authData, nil
}

// SendWalletLinkLoginRequest sends a login request (SendLoginRequest) for the already-linked member. Returns loginId, deeplink, QR (for use with VerifyWalletLinkLoginRequest).
func SendWalletLinkLoginRequest(server *serverModels.Server, c *gin.Context) (*trovosdk.LoginWithTrovoWalletData, error) {
	memberID := c.GetString("member_id")
	if memberID == "" {
		return nil, errors.New("organization member context required")
	}

	var member models.OrganizationMember
	if err := server.AdminDB.Where("id = ?", memberID).First(&member).Error; err != nil {
		return nil, errors.New("member not found")
	}
	if !member.IsWalletLinked || member.TrovoWalletUsername == nil || *member.TrovoWalletUsername == "" {
		return nil, errors.New("wallet not linked; complete link-wallet first")
	}

	cleanUsername := strings.Split(*member.TrovoWalletUsername, "@")[0]
	deviceInfo := "Trovo Manager"
	callbackUrl := fmt.Sprintf("%v/%v", os.Getenv("LOGIN_CALLBACK_URL"), "trovo")

	loginData, err := server.GC.ServiceLink.SendLoginRequest(cleanUsername, "", deviceInfo, callbackUrl)
	if err != nil {
		log.Printf("[SendWalletLinkLoginRequest] Error: %v", err)
		return nil, fmt.Errorf("failed to send login request: %v", err)
	}
	return loginData, nil
}

// VerifyWalletLinkLoginRequest verifies the login request and returns tokens (VerifyLoginRequest).
func VerifyWalletLinkLoginRequest(server *serverModels.Server, c *gin.Context, loginID string) (*trovosdk.LoginVerifyResponse, error) {
	if loginID == "" {
		return nil, errors.New("login_id is required")
	}

	memberID := c.GetString("member_id")
	if memberID == "" {
		return nil, errors.New("organization member context required")
	}

	var member models.OrganizationMember
	if err := server.AdminDB.Where("id = ?", memberID).First(&member).Error; err != nil {
		return nil, errors.New("member not found")
	}
	if !member.IsWalletLinked || member.TrovoWalletUsername == nil || *member.TrovoWalletUsername == "" {
		return nil, errors.New("wallet not linked")
	}

	cleanUsername := strings.Split(*member.TrovoWalletUsername, "@")[0]

	log.Println("[VerifyWalletLinkLoginRequest] Verifying login request for user", cleanUsername, "with login ID", loginID)
	resp, err := server.GC.ServiceLink.VerifyLoginRequest(cleanUsername, loginID)
	if err != nil {
		var ex p2pErrors.GenericError
		if ok := errors.As(err, &ex); ok {
			log.Printf("[VerifyWalletLinkLoginRequest] verification failed: %v", ex.JSONError())
		}
		return nil, fmt.Errorf("login verification failed: %v", err)
	}
	return resp, nil
}

// VerifyWalletLink verifies the authorization (authID) and updates the member's profile. Kept for any auth-based flow; primary token flow is VerifyWalletLinkLoginRequest.
func VerifyWalletLink(server *serverModels.Server, memberID string, walletUsername string, authID string) error {
	if authID == "" {
		return errors.New("authID is required")
	}
	walletUsername = strings.ToLower(strings.TrimSpace(walletUsername))
	if !strings.HasSuffix(walletUsername, "@trovo") {
		return errors.New("invalid wallet username format")
	}
	cleanUsername := strings.TrimSuffix(walletUsername, "@trovo")

	var existing models.OrganizationMember
	if err := server.AdminDB.Where("id = ?", memberID).First(&existing).Error; err != nil {
		return err
	}
	if existing.IsWalletLinked {
		return errors.New("wallet already linked; each user may link only once")
	}
	if !authServices.IsTrovoUser(server, walletUsername) {
		return errors.New("trovo wallet user not found or suspended")
	}

	authdata, err := server.GC.ServiceLink.VerifyAuthorizationRequest(cleanUsername, authID)
	if err != nil {
		log.Println("[VerifyWalletLink] Error verifying authorization for user", cleanUsername, "with authID", authID, "and wallet username", walletUsername, "AuthData:", authdata, "Error:", err)
		var ex p2pErrors.GenericError
		if ok := errors.As(err, &ex); ok {
			log.Printf("[VerifyWalletLink] verification failed with code %d: %v", ex.HTTPCode(), ex.JSONError())
		}
		return fmt.Errorf("authorization verification failed: %v", err)
	}
	log.Println("[VerifyWalletLink] Authorization verified for user", cleanUsername, "with authID", authID, "and wallet username", walletUsername, "AuthData:", authdata)
	return server.AdminDB.Transaction(func(tx *gorm.DB) error {
		var member models.OrganizationMember
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", memberID).First(&member).Error; err != nil {
			return err
		}
		if member.IsWalletLinked {
			return errors.New("wallet already linked; each user may link only once")
		}
		member.TrovoWalletUsername = &walletUsername
		member.IsWalletLinked = true
		member.UpdatedAt = time.Now()
		return tx.Save(&member).Error
	})
}
