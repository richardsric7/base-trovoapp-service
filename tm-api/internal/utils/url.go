package utils

import "os"

func GetTrovomanagerFrontendBaseUrl() string {
	baseURL := os.Getenv("TROVO_MANAGER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://dashboard.dev.admin.trovo.app" // fallback to dev dashboard URL
	}
	return baseURL
}
