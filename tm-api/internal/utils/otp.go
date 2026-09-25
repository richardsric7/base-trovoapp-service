package utils

import (
	"admin-panel-dashboard/internal/mail"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
)

// GenerateOTP generates a random OTP of specified length
func GenerateOTP(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("invalid OTP length: %d", n)
	}

	minn := int64(1)
	for i := 1; i < n; i++ {
		minn *= 10
	}
	maxx := minn*10 - 1

	rangeVal := big.NewInt(maxx - minn + 1)
	randNum, err := rand.Int(rand.Reader, rangeVal)
	if err != nil {
		slog.Error("Failed to generate OTP", "error", err)
		return "", err
	}

	otp := fmt.Sprintf("%d", randNum.Int64()+minn)
	return otp, nil
}

// SendOTPEmail sends an OTP verification email using the organization_otp_email template
func SendOTPEmail(server interface{}, invite *models.OrganizationInvite) error {
	// Get the server instance
	s, ok := server.(*serverModels.Server)
	if !ok {
		return fmt.Errorf("invalid server type")
	}
	// Send email using the dedicated function
	_, _, err := mail.SendOrganizationOTPEmail(invite.Email, invite.VerificationOTP, invite.OTPExpiresAt, s)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %v", err)
	}

	return nil
}
