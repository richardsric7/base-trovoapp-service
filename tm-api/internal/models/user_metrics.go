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
type P2PUserMetricsResponse struct {
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
	PublicKey   string    `json:"public_key"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Date        time.Time `json:"date"`
}

type TradesOnAppeal struct {
	AppealID     string    `json:"appeal_id"`
	CustomerName string    `json:"customer_name"`
	Description  string    `json:"description"`
	Attachment   string    `json:"attachment"`
	UpdatedOn    time.Time `json:"updated_on"`
	Status       string    `json:"status"`
}
type UserList struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Location     string `json:"location"`
	PhoneNumber  string `json:"phone_number"`
	PatronStatus string `json:"patron_status"`
}

type P2POrder struct {
	ID                                             string    `json:"id"`
	OrderType                                      string    `json:"order_type"`
	CreatedAt                                      time.Time `json:"created_at"`
	UpdatedAt                                      time.Time `json:"updated_at"`
	ExpiresAt                                      time.Time `json:"expires_at"`
	AcceptedAt                                     time.Time `json:"accepted_at"`
	CancelAfter                                    time.Time `json:"cancel_after"`
	OfferID                                        string    `json:"offer_id"`
	OfferType                                      string    `json:"offer_type"`
	OfferMaker                                     string    `json:"offer_maker"`
	OfferMakerPhone                                string    `json:"offer_maker_phone"`
	OfferMakerCountryCode                          *string   `json:"offer_maker_country_code"`
	OfferMaxTimePerTransaction                     uint      `json:"offer_max_time_per_transaction"`
	OfferTaker                                     string    `json:"offer_taker"`
	OfferTakerPhone                                string    `json:"offer_taker_phone"`
	OfferPaymentMethodID                           *string   `json:"offer_payment_method_id"`
	OfferPaymentChannelID                          *string   `json:"offer_payment_channel_id"`
	OfferPaymentMethodName                         *string   `json:"offer_payment_method_name"`
	OfferPaymentMethodDestinationAccount           *string   `json:"offer_payment_method_destination_account"`
	OfferPaymentMethodMemo                         *string   `json:"offer_payment_method_memo"`
	OfferPaymentMethodBankName                     *string   `json:"offer_payment_method_bank_name"`
	OfferPaymentMethodAccountOpeningBranch         *string   `json:"offer_payment_method_account_opening_branch"`
	OfferCurrencyPaymentMethodID                   string    `json:"offer_currency_payment_method_id"`
	OfferCurrencyPaymentChannelID                  string    `json:"offer_currency_payment_channel_id"`
	OfferCurrencyPaymentMethodName                 *string   `json:"offer_currency_payment_method_name"`
	OfferCurrencyPaymentMethodDestinationAccount   string    `json:"offer_currency_payment_method_destination_account"`
	OfferCurrencyPaymentMethodMemo                 *string   `json:"offer_currency_payment_method_memo"`
	OfferCurrencyPaymentMethodBankName             *string   `json:"offer_currency_payment_method_bank_name"`
	OfferCurrencyPaymentMethodAccountOpeningBranch *string   `json:"offer_currency_payment_method_account_opening_branch"`
	OfferCurrencyPaymentMethodCountryCode          *string   `json:"offer_currency_payment_method_country_code"`
	OfferCurrencyPaymentMethodCurrencyID           *string   `json:"offer_currency_payment_method_currency_id"`
	OfferCurrencyID                                string    `json:"offer_currency_id"`
	OfferAssetAmount                               float64   `json:"offer_asset_amount"`
	OfferAssetID                                   string    `json:"offer_asset_id"`
	OfferAssetPrice                                float64   `json:"offer_asset_price"`
	OfferMinTradeAmount                            float64   `json:"offer_min_trade_amount"`
	OfferMaxTradeAmount                            float64   `json:"offer_max_trade_amount"`
	OfferRemark                                    *string   `json:"offer_remark"`
	TakerPaymentMethodID                           string    `json:"taker_payment_method_id"`
	TakerPaymentChannelID                          string    `json:"taker_payment_channel_id"`
	TakerPaymentMethodName                         *string   `json:"taker_payment_method_name"`
	TakerPaymentMethodDestinationAccount           string    `json:"taker_payment_method_destination_account"`
	TakerPaymentMethodMemo                         *string   `json:"taker_payment_method_memo"`
	TakerPaymentMethodBankName                     *string   `json:"taker_payment_method_bank_name"`
	TakerPaymentMethodAccountOpeningBranch         *string   `json:"taker_payment_method_account_opening_branch"`
	TakerPaymentMethodCountryCode                  *string   `json:"taker_payment_method_country_code"`
	TakerPaymentMethodCurrencyID                   *string   `json:"taker_payment_method_currency_id"`
	OrderEscrowAddress                             string    `json:"order_escrow_address"`
	OrderPaymentMemo                               string    `json:"order_payment_memo"`
	OrderAmount                                    float64   `json:"order_amount"`
	OrderMakerFee                                  float64   `json:"order_maker_fee"`
	OrderTakerFee                                  float64   `json:"order_taker_fee"`
	OrderEscrowTransactionID                       *string   `json:"order_escrow_transaction_id"`
	OrderAssetReleaseTransactionID                 *string   `json:"order_asset_release_transaction_id"`
	OrderStatusID                                  uint      `json:"order_status_id"`
	// OrderStatus                                    string    `json:"order_status"`
	DynamicLink              *string `json:"dynamic_link"`
	QRCode                   *string `json:"qr_code"`
	FiatDepositTransactionID *string `json:"fiat_deposit_transaction_id"`
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
	PublicKey             string    `json:"public_key"`
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

type AppealResponse struct {
	Data     []TradesOnAppeal `json:"data"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
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
	PublicKey string `json:"public_key"`
}

type UserProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	UserID   string `json:"user_id"`
}

type UserWalletDTO struct {
	ID            string `json:"publicKey"`
	Alias         string `json:"alias"`
	Signer        string `json:"signer"`
	UserID        string `json:"userId"`
	PrimaryWallet int    `json:"primaryWallet"`
}

type UserInfoResponse struct {
	Wallets  []UserWalletDTO `json:"wallets"`
	UserInfo *User           `json:"user_info"`
}

type AppealListRequest struct {
	Page         int
	PageSize     int
	AppealID     string
	CustomerName string
	Description  string
	UpdatedOn    string
	Status       string
	Search       string
}

type TradeListRequest struct {
	ID                                             string    `form:"id"`
	Username                                       string    `form:"username"`
	OrderType                                      string    `form:"order_type"`
	CreatedAt                                      time.Time `form:"created_at"`
	UpdatedAt                                      time.Time `form:"updated_at"`
	ExpiresAt                                      time.Time `form:"expires_at"`
	AcceptedAt                                     time.Time `form:"accepted_at"`
	CancelAfter                                    time.Time `form:"cancel_after"`
	OfferID                                        string    `form:"offer_id"`
	OfferType                                      string    `form:"offer_type"`
	OfferMaker                                     string    `form:"offer_maker"`
	OfferMakerPhone                                string    `form:"offer_maker_phone"`
	OfferMakerCountryCode                          *string   `form:"offer_maker_country_code"`
	OfferMaxTimePerTransaction                     uint      `form:"offer_max_time_per_transaction"`
	OfferTaker                                     string    `form:"offer_taker"`
	OfferTakerPhone                                string    `form:"offer_taker_phone"`
	OfferPaymentMethodID                           *string   `form:"offer_payment_method_id"`
	OfferPaymentChannelID                          *string   `form:"offer_payment_channel_id"`
	OfferPaymentMethodName                         *string   `form:"offer_payment_method_name"`
	OfferPaymentMethodDestinationAccount           *string   `form:"offer_payment_method_destination_account"`
	OfferPaymentMethodMemo                         *string   `form:"offer_payment_method_memo"`
	OfferPaymentMethodBankName                     *string   `form:"offer_payment_method_bank_name"`
	OfferPaymentMethodAccountOpeningBranch         *string   `form:"offer_payment_method_account_opening_branch"`
	OfferCurrencyPaymentMethodID                   string    `form:"offer_currency_payment_method_id"`
	OfferCurrencyPaymentChannelID                  string    `form:"offer_currency_payment_channel_id"`
	OfferCurrencyPaymentMethodName                 *string   `form:"offer_currency_payment_method_name"`
	OfferCurrencyPaymentMethodDestinationAccount   string    `form:"offer_currency_payment_method_destination_account"`
	OfferCurrencyPaymentMethodMemo                 *string   `form:"offer_currency_payment_method_memo"`
	OfferCurrencyPaymentMethodBankName             *string   `form:"offer_currency_payment_method_bank_name"`
	OfferCurrencyPaymentMethodAccountOpeningBranch *string   `form:"offer_currency_payment_method_account_opening_branch"`
	OfferCurrencyPaymentMethodCountryCode          *string   `form:"offer_currency_payment_method_country_code"`
	OfferCurrencyPaymentMethodCurrencyID           *string   `form:"offer_currency_payment_method_currency_id"`
	OfferCurrencyID                                string    `form:"offer_currency_id"`
	OfferAssetAmount                               float64   `form:"offer_asset_amount"`
	OfferAssetID                                   string    `form:"offer_asset_id"`
	OfferAssetPrice                                float64   `form:"offer_asset_price"`
	OfferMinTradeAmount                            float64   `form:"offer_min_trade_amount"`
	OfferMaxTradeAmount                            float64   `form:"offer_max_trade_amount"`
	OfferRemark                                    *string   `form:"offer_remark"`
	TakerPaymentMethodID                           string    `form:"taker_payment_method_id"`
	TakerPaymentChannelID                          string    `form:"taker_payment_channel_id"`
	TakerPaymentMethodName                         *string   `form:"taker_payment_method_name"`
	TakerPaymentMethodDestinationAccount           string    `form:"taker_payment_method_destination_account"`
	TakerPaymentMethodMemo                         *string   `form:"taker_payment_method_memo"`
	TakerPaymentMethodBankName                     *string   `form:"taker_payment_method_bank_name"`
	TakerPaymentMethodAccountOpeningBranch         *string   `form:"taker_payment_method_account_opening_branch"`
	TakerPaymentMethodCountryCode                  *string   `form:"taker_payment_method_country_code"`
	TakerPaymentMethodCurrencyID                   *string   `form:"taker_payment_method_currency_id"`
	OrderEscrowAddress                             string    `form:"order_escrow_address"`
	OrderPaymentMemo                               string    `form:"order_payment_memo"`
	OrderAmount                                    float64   `form:"order_amount"`
	OrderMakerFee                                  float64   `form:"order_maker_fee"`
	OrderTakerFee                                  float64   `form:"order_taker_fee"`
	OrderEscrowTransactionID                       *string   `form:"order_escrow_transaction_id"`
	OrderAssetReleaseTransactionID                 *string   `form:"order_asset_release_transaction_id"`
	OrderStatusID                                  uint      `form:"order_status_id"`
	// OrderStatus                                    string    `form:"order_status"`
	DynamicLink              *string `form:"dynamic_link"`
	QRCode                   *string `form:"qr_code"`
	FiatDepositTransactionID *string `form:"fiat_deposit_transaction_id"`
	Page                     int     `form:"page" binding:"min=1"`
	PageSize                 int     `form:"page_size" binding:"min=1,max=100"`
}
