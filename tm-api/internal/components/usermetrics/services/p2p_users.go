package usermetrics

import (
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get P2P User List
// @Description Retrieves a list of P2P users with pagination and optional filtering by various attributes.
// @ID GetP2PUserList
// @Tags P2P
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param page query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param username query string false "Filter by username"
// @Param email query string false "Filter by user email"
// @Param phone query string false "Filter by user phone number"
// @Param first_name query string false "Filter by user first_name"
// @Param last_name query string false "Filter by user last name"
// @Param city query string false "Filter by user city"
// @Param address query string false "Filter by address"
// @Param suspended query int false "Filter by suspension status (0 or 1)"
// @Param kyc_level query int false "Filter by KYC level"
// @Param admin_level query int false "Filter by admin level"
// @Param registration_date_from query string false "Filter by registration date (start) in YYYY-MM-DD format"
// @Param registration_date_to query string false "Filter by registration date (end) in YYYY-MM-DD format"
// @Param search query string false "General search across multiple fields"
// @Success 200 {object} response.Data
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /p2p/users [get]
func GetP2PUserList(walletDb, p2pdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
			return
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size"})
			return
		}

		filters := P2PUserRequestDTO{
			Page:                 page,
			PageSize:             pageSize,
			Username:             c.Query("username"),
			Email:                c.Query("email"),
			Phone:                c.Query("phone"),
			FirstName:            c.Query("first_name"),
			LastName:             c.Query("last_name"),
			City:                 c.Query("city"),
			Address:              c.Query("address"),
			Suspended:            c.DefaultQuery("suspended", ""),
			KYCLevel:             c.DefaultQuery("kyc_level", ""),
			AdminLevel:           c.DefaultQuery("admin_level", ""),
			Search:               c.Query("search"),
			RegistrationDateFrom: c.Query("registration_date_from"),
			RegistrationDateTo:   c.Query("registration_date_to"),
		}

		users, total, err := fetchP2PUsers(filters, p2pdb)
		if err != nil {
			log.Println("[P2P USERS] error fetching users:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		response := P2PUserResponse{
			Data:     users,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}
		serverResponse.JSON(c, http.StatusOK, "P2P user list fetched successfully", response, nil)
	}
}

func fetchP2PUsers(filters P2PUserRequestDTO, db *gorm.DB) ([]P2PUserDTO, int64, error) {
	var userDTOs []P2PUserDTO
	var total int64

	query := db.Model(&models.User{})

	if filters.Username != "" {
		query = query.Where("LOWER(username) ILIKE LOWER(?)", filters.Username)
	}
	if filters.Email != "" {
		query = query.Where("LOWER(email) ILIKE LOWER(?)", filters.Email)
	}
	if filters.Phone != "" {
		query = query.Where("mobile ILIKE ?", filters.Phone)
	}
	if filters.FirstName != "" {
		query = query.Where("LOWER(first_name) ILIKE LOWER(?)", filters.FirstName)
	}
	if filters.LastName != "" {
		query = query.Where("LOWER(last_name) ILIKE LOWER(?)", filters.LastName)
	}
	if filters.City != "" {
		query = query.Where("LOWER(city) ILIKE LOWER(?)", filters.City)
	}
	if filters.Address != "" {
		query = query.Where("LOWER(address) ILIKE LOWER(?)", filters.Address)
	}
	if filters.Suspended != "" {
		suspended, _ := strconv.Atoi(filters.Suspended)
		query = query.Where("suspended = ?", suspended)
	}
	if filters.KYCLevel != "" {
		kycLevel, _ := strconv.Atoi(filters.KYCLevel)
		query = query.Where("kyc_level = ?", kycLevel)
	}
	if filters.AdminLevel != "" {
		adminLevel, _ := strconv.Atoi(filters.AdminLevel)
		query = query.Where("admin_level = ?", adminLevel)
	}
	if filters.RegistrationDateFrom != "" {
		query = query.Where("created_at >= ?", filters.RegistrationDateFrom)
	}
	if filters.RegistrationDateTo != "" {
		query = query.Where("created_at <= ?", filters.RegistrationDateTo)
	}
	if filters.Search != "" {
		searchTerm := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where(
			db.Where("LOWER(username) ILIKE ?", searchTerm).
				Or("LOWER(email) ILIKE ?", searchTerm).
				Or("LOWER(first_name) ILIKE ?", searchTerm).
				Or("LOWER(last_name) ILIKE ?", searchTerm).
				Or("LOWER(city) ILIKE ?", searchTerm),
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.PageSize
	var users []user
	if err := query.Order("created_at desc").Offset(offset).Limit(filters.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	for _, usr := range users {
		userDTOs = append(userDTOs, P2PUserDTO{
			ID:               usr.ID,
			Username:         usr.Username,
			Email:            usr.Email,
			Phone:            usr.Mobile,
			RegistrationDate: usr.CreatedAt,
			AccountStatus:    usr.Suspended,
			//MerchantStatus:   user.AdminLevel,
			Violations: usr.SuspensionReason,
			City:       usr.City,
		})
	}

	return userDTOs, total, nil
}

// P2PUserRequestDTO defines the request structure for fetching P2P user list.
type P2PUserRequestDTO struct {
	Page                 int    `json:"page"`
	PageSize             int    `json:"pageSize"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	Phone                string `json:"phone"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	City                 string `json:"city"`
	Address              string `json:"address"`
	Suspended            string `json:"suspended"`
	KYCLevel             string `json:"kyc_level"`
	AdminLevel           string `json:"admin_level"`
	Search               string `json:"search"`
	RegistrationDateTo   string `json:"registration_date_to"`
	RegistrationDateFrom interface{}
}

// P2PUserResponse defines the response structure for fetching P2P user list.
type P2PUserResponse struct {
	Data     []P2PUserDTO `json:"data"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
}

// P2PUserDTO defines the data transfer object for a single P2P user.
type P2PUserDTO struct {
	ID               string     `json:"id"`
	Username         string     `json:"username"`
	Email            string     `json:"email"`
	Phone            *string    `json:"phone"`
	RegistrationDate *time.Time `json:"registration_date"`
	AccountStatus    int        `json:"account_status"`
	MerchantStatus   int64      `json:"merchant_status"`
	Violations       *string    `json:"violations"`
	City             *string    `json:"city"`
}

type user struct {
	ID                    string     `json:"id"`                     // character varying, not nullable
	Username              string     `json:"username"`               // character varying, not nullable
	CreatedAt             *time.Time `json:"created_at"`             // timestamp with time zone, nullable
	UpdatedAt             *time.Time `json:"updated_at"`             // timestamp with time zone, nullable
	TelegramConnectedAt   *time.Time `json:"telegram_connected_at"`  // timestamp with time zone, nullable
	LastName              string     `json:"last_name"`              // character varying, not nullable
	FirstName             string     `json:"first_name"`             // character varying, not nullable
	MiddleName            *string    `json:"middle_name"`            // character varying, nullable
	Mobile                *string    `json:"mobile"`                 // character varying, nullable
	ContactPhone          *string    `json:"contact_phone"`          // character varying, nullable
	Telegram              *int64     `json:"telegram"`               // bigint, nullable
	Email                 string     `json:"email"`                  // text, not nullable
	ImageThumbnail        *string    `json:"image_thumbnail"`        // text, nullable
	CountryCode           *string    `json:"country_code"`           // character varying, nullable
	Latitude              *float64   `json:"latitude"`               // numeric, nullable
	Longitude             *float64   `json:"longitude"`              // numeric, nullable
	City                  *string    `json:"city"`                   // character varying, nullable
	Region                *string    `json:"region"`                 // character varying, nullable
	RegionName            *string    `json:"region_name"`            // character varying, nullable
	TimeZone              *string    `json:"time_zone"`              // character varying, nullable
	ISP                   *string    `json:"isp"`                    // character varying, nullable
	PublicIP              *string    `json:"public_ip"`              // character varying, nullable
	Address               *string    `json:"address"`                // character varying, nullable
	BantuTalk             *string    `json:"bantu_talk"`             // character varying, nullable
	Suspended             int        `json:"suspended"`              // integer, not nullable, default 0
	KYCLevel              int        `json:"kyc_level"`              // integer, not nullable, default 0
	MaxAssetPerOffer      int        `json:"max_asset_per_offer"`    // integer, not nullable, default 0
	MaxAssetPerOrder      int        `json:"max_asset_per_order"`    // integer, not nullable, default 0
	SuspensionReason      *string    `json:"suspension_reason"`      // text, nullable
	Offline               int64      `json:"offline"`                // bigint, not nullable, default 0
	AdminLevel            int64      `json:"admin_level"`            // bigint, not nullable, default 0
	TelegramNotifications int64      `json:"telegram_notifications"` // bigint, not nullable, default 0
}

// @Summary Get Transaction History
// @Description Retrieves transaction history for a user based on email, including both offer_maker and offer_taker transactions.
// @ID GetTransactionHistory
// @Tags P2P
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param email query string true "User email for fetching transaction history"
// @Param page query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} response.Data
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /p2p/transaction/history [get]
func GetTransactionHistory(p2pdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
			return
		}

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
			return
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size"})
			return
		}

		transactions, total, err := fetchTransactionHistory(email, page, pageSize, p2pdb)
		if err != nil {
			log.Println("[TRANSACTION HISTORY] error fetching transactions:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		response := TransactionHistoryResponse{
			Data:     transactions,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}
		serverResponse.JSON(c, http.StatusOK, "P2P Transaction history fetched successfully", response, nil)
	}
}

func fetchTransactionHistory(email string, page, pageSize int, db *gorm.DB) ([]TransactionHistoryDTO, int64, error) {
	var transactions []TransactionHistoryDTO
	var total int64

	query := db.Model(&order{}).
		Where("LOWER(offer_maker) = LOWER(?) OR LOWER(offer_taker) = LOWER(?)", email, email)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var orders []order
	if err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	for _, ordr := range orders {
		transactions = append(transactions, TransactionHistoryDTO{
			ID:            ordr.ID,
			OrderType:     ordr.OrderType,
			CreatedAt:     ordr.CreatedAt,
			UpdatedAt:     ordr.UpdatedAt,
			ExpiresAt:     ordr.ExpiresAt,
			AcceptedAt:    ordr.AcceptedAt,
			OfferID:       ordr.OfferID,
			OfferType:     ordr.OfferType,
			OfferMaker:    ordr.OfferMaker,
			OfferTaker:    ordr.OfferTaker,
			OrderAmount:   ordr.OrderAmount,
			OrderMakerFee: ordr.OrderMakerFee,
			OrderTakerFee: ordr.OrderTakerFee,
			OrderStatusID: ordr.OrderStatusID,
		})
	}

	return transactions, total, nil
}

// TransactionHistoryDTO defines the data transfer object for a single transaction.
type TransactionHistoryDTO struct {
	ID            string     `json:"id"`
	OrderType     string     `json:"order_type"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	AcceptedAt    *time.Time `json:"accepted_at"`
	OfferID       string     `json:"offer_id"`
	OfferType     string     `json:"offer_type"`
	OfferMaker    string     `json:"offer_maker"`
	OfferTaker    string     `json:"offer_taker"`
	OrderAmount   float64    `json:"order_amount"`
	OrderMakerFee float64    `json:"order_maker_fee"`
	OrderTakerFee float64    `json:"order_taker_fee"`
	OrderStatusID int64      `json:"order_status_id"`
}

// TransactionHistoryResponse defines the response structure for fetching transaction history.
type TransactionHistoryResponse struct {
	Data     []TransactionHistoryDTO `json:"data"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"pageSize"`
}

type order struct {
	ID                                             string     `json:"id"`                                                   // character varying, not nullable
	OrderType                                      string     `json:"order_type"`                                           // text, not nullable
	CreatedAt                                      *time.Time `json:"created_at"`                                           // timestamp with time zone, nullable
	UpdatedAt                                      *time.Time `json:"updated_at"`                                           // timestamp with time zone, nullable
	ExpiresAt                                      *time.Time `json:"expires_at"`                                           // timestamp with time zone, nullable
	AcceptedAt                                     *time.Time `json:"accepted_at"`                                          // timestamp with time zone, nullable
	CancelAfter                                    *time.Time `json:"cancel_after"`                                         // timestamp with time zone, nullable
	OfferID                                        string     `json:"offer_id"`                                             // text, not nullable
	OfferType                                      string     `json:"offer_type"`                                           // text, not nullable
	OfferMaker                                     string     `json:"offer_maker"`                                          // text, not nullable
	OfferMakerPhone                                string     `json:"offer_maker_phone"`                                    // text, not nullable
	OfferMakerCountryCode                          *string    `json:"offer_maker_country_code"`                             // text, nullable
	OfferMaxTimePerTransaction                     int64      `json:"offer_max_time_per_transaction"`                       // bigint, not nullable, default 45
	OfferTaker                                     string     `json:"offer_taker"`                                          // text, not nullable
	OfferTakerPhone                                string     `json:"offer_taker_phone"`                                    // text, not nullable
	OfferPaymentMethodID                           *string    `json:"offer_payment_method_id"`                              // text, nullable
	OfferPaymentChannelID                          *string    `json:"offer_payment_channel_id"`                             // text, nullable
	OfferPaymentMethodName                         *string    `json:"offer_payment_method_name"`                            // text, nullable
	OfferPaymentMethodDestinationAccount           *string    `json:"offer_payment_method_destination_account"`             // text, nullable
	OfferPaymentMethodMemo                         *string    `json:"offer_payment_method_memo"`                            // text, nullable
	OfferPaymentMethodBankName                     *string    `json:"offer_payment_method_bank_name"`                       // text, nullable
	OfferPaymentMethodAccountOpeningBranch         *string    `json:"offer_payment_method_account_opening_branch"`          // text, nullable
	OfferCurrencyPaymentMethodID                   string     `json:"offer_currency_payment_method_id"`                     // text, not nullable
	OfferCurrencyPaymentChannelID                  string     `json:"offer_currency_payment_channel_id"`                    // text, not nullable
	OfferCurrencyPaymentMethodName                 *string    `json:"offer_currency_payment_method_name"`                   // text, nullable
	OfferCurrencyPaymentMethodDestinationAccount   string     `json:"offer_currency_payment_method_destination_account"`    // text, not nullable
	OfferCurrencyPaymentMethodMemo                 *string    `json:"offer_currency_payment_method_memo"`                   // text, nullable
	OfferCurrencyPaymentMethodBankName             *string    `json:"offer_currency_payment_method_bank_name"`              // text, nullable
	OfferCurrencyPaymentMethodAccountOpeningBranch *string    `json:"offer_currency_payment_method_account_opening_branch"` // text, nullable
	OfferCurrencyPaymentMethodCountryCode          *string    `json:"offer_currency_payment_method_country_code"`           // text, nullable
	OfferCurrencyPaymentMethodCurrencyID           *string    `json:"offer_currency_payment_method_currency_id"`            // text, nullable
	OfferCurrencyID                                string     `json:"offer_currency_id"`                                    // text, not nullable
	OfferAssetAmount                               float64    `json:"offer_asset_amount"`                                   // numeric, not nullable
	OfferAssetID                                   string     `json:"offer_asset_id"`                                       // character varying, not nullable
	OfferAssetPrice                                float64    `json:"offer_asset_price"`                                    // numeric, not nullable
	OfferMinTradeAmount                            float64    `json:"offer_min_trade_amount"`                               // numeric, not nullable
	OfferMaxTradeAmount                            float64    `json:"offer_max_trade_amount"`                               // numeric, not nullable
	OfferRemark                                    *string    `json:"offer_remark"`                                         // text, nullable
	TakerPaymentMethodID                           string     `json:"taker_payment_method_id"`                              // text, not nullable
	TakerPaymentChannelID                          string     `json:"taker_payment_channel_id"`                             // text, not nullable
	TakerPaymentMethodName                         *string    `json:"taker_payment_method_name"`                            // text, nullable
	TakerPaymentMethodDestinationAccount           string     `json:"taker_payment_method_destination_account"`             // text, not nullable
	TakerPaymentMethodMemo                         *string    `json:"taker_payment_method_memo"`                            // text, nullable
	TakerPaymentMethodBankName                     *string    `json:"taker_payment_method_bank_name"`                       // text, nullable
	TakerPaymentMethodAccountOpeningBranch         *string    `json:"taker_payment_method_account_opening_branch"`          // text, nullable
	TakerPaymentMethodCountryCode                  *string    `json:"taker_payment_method_country_code"`                    // character varying, nullable
	TakerPaymentMethodCurrencyID                   *string    `json:"taker_payment_method_currency_id"`                     // character varying, nullable
	OrderEscrowAddress                             string     `json:"order_escrow_address"`                                 // text, not nullable
	OrderPaymentMemo                               string     `json:"order_payment_memo"`                                   // character varying, not nullable
	OrderAmount                                    float64    `json:"order_amount"`                                         // numeric, not nullable
	OrderMakerFee                                  float64    `json:"order_maker_fee"`                                      // numeric, not nullable
	OrderTakerFee                                  float64    `json:"order_taker_fee"`                                      // numeric, not nullable
	OrderEscrowTransactionID                       *string    `json:"order_escrow_transaction_id"`                          // text, nullable
	OrderAssetReleaseTransactionID                 *string    `json:"order_asset_release_transaction_id"`                   // text, nullable
	OrderStatusID                                  int64      `json:"order_status_id"`                                      // bigint, not nullable, default 1
	DynamicLink                                    *string    `json:"dynamic_link"`                                         // text, nullable
	QRCode                                         *string    `json:"qr_code"`                                              // text, nullable
	FiatDepositTransactionID                       *string    `json:"fiat_deposit_transaction_id"`                          // text, nullable
}
