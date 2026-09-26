package models

import "time"

type WalletUserMetricsResponse struct {
	TotalUsers               int64          `json:"total_users"`
	TotalActiveUsers         int64          `json:"total_active_users"`
	TotalSessions            int64          `json:"total_sessions"`
	AverageTimePerSession    int64          `json:"average_time_per_session"`
	DailyActiveAndNewUsers   DailyMetrics   `json:"daily_active_and_new_users"`
	WeeklyActiveAndNewUsers  WeeklyMetrics  `json:"weekly_active_and_new_users"`
	MonthlyActiveAndNewUsers MonthlyMetrics `json:"monthly_active_and_new_users"`
	WeeklyPercentageChange   WeeklyChange   `json:"weekly_percentage_change"`
}
type WalletVersionDistribution struct {
	TotalNumberOfUsers int64   `json:"total_number_of_users"`
	Version            string  `json:"version"`
	Count              int64   `json:"count"`
	Percentage         float64 `json:"percentage"`
}

type DailyMetrics struct {
	ActiveUsers int64 `json:"active_users"`
	NewUsers    int64 `json:"new_users"`
}

type WeeklyMetrics struct {
	ActiveUsers int64 `json:"active_users"`
	NewUsers    int64 `json:"new_users"`
}

type MonthlyMetrics struct {
	ActiveUsers int64 `json:"active_users"`
	NewUsers    int64 `json:"new_users"`
}

type WeeklyChange struct {
	NewUsers      float64 `json:"new_users_percentage_change"`
	ActiveUsers   float64 `json:"active_users_percentage_change"`
	TotalSessions float64 `json:"total_sessions_percentage_change"`
	AverageTime   float64 `json:"average_time_percentage_change"`
}

type CountryUserCount struct {
	Country string
	Count   int
}

type ReferrerCount struct {
	Referrer string `json:"referrer"`
	Count    int    `json:"count"`
}

type UserDTO struct {
	Username    string    `json:"username"`
	Address     string    `json:"address"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Date        time.Time `json:"date"`
}

type UserList struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Location     string `json:"location"`
	PhoneNumber  string `json:"phone_number"`
	PatronStatus string `json:"patron_status"`
}

type UserDto struct {
	ID                    string    `json:"id"`
	UserName              string    `json:"user_name"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	TelegramConnectedAt   time.Time `json:"telegram_connected_at"`
	LastName              string    `json:"last_name"`
	FirstName             string    `json:"first_name"`
	MiddleName            string    `json:"middle_name"`
	Mobile                string    `json:"mobile"`
	ContactPhone          string    `json:"contact_phone"`
	Telegram              int64     `json:"telegram"`
	Email                 string    `json:"email"`
	ImageThumbnail        string    `json:"image_thumbnail"`
	CountryCode           string    `json:"country_code"`
	Latitude              float64   `json:"latitude"`
	Longitude             float64   `json:"longitude"`
	City                  string    `json:"city"`
	Region                string    `json:"region"`
	RegionName            string    `json:"region_name"`
	TimeZone              string    `json:"time_zone"`
	ISP                   string    `json:"isp"`
	PublicIP              string    `json:"public_ip"`
	Address               string    `json:"address"`
	BantuTalk             string    `json:"bantu_talk"`
	Suspended             uint      `json:"suspended"`
	KYCLevel              uint      `json:"kyc_level"`
	MaxAssetPerOffer      float64   `json:"max_asset_per_offer"`
	MaxAssetPerOrder      float64   `json:"max_asset_per_order"`
	SuspensionReason      string    `json:"suspension_reason"`
	Offline               uint      `json:"offline"`
	AdminLevel            uint      `json:"admin_level"`
	TelegramNotifications uint      `json:"telegram_notifications"`
}

type UserResponse struct {
	Data     []UserDTO `json:"data"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

type UserResponseData struct {
	Data     []UserDto `json:"data"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserRequestDTO struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Search    string `json:"search"`
	City      string `json:"city"`
	Address   string `json:"address"`
}

type UserProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	UserID   string `json:"user_id"`
}

type UserWalletDTO struct {
	ID            string `json:"address"`
	Alias         string `json:"alias"`
	Signer        string `json:"signer"`
	UserID        string `json:"userId"`
	PrimaryWallet int    `json:"primaryWallet"`
}

type UserInfoResponse struct {
	Wallets  []UserWalletDTO `json:"wallets"`
	UserInfo *User           `json:"user_info"`
}

