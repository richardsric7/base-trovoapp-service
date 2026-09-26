package usermetrics

import (
	"admin-panel-dashboard/internal/models"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

func GetUserMetrics(walletDB *gorm.DB) (*models.WalletUserMetricsResponse, error) {

	userCount, err := totalUserCount(walletDB)
	if err != nil {
		return nil, err
	}
	var suspendedCount int64
	if err := walletDB.Model(&models.User{}).Count(&suspendedCount).Where("suspended = 1").Error; err != nil {
		return nil, err
	}
	log.Println("checking numbers: ", suspendedCount, *userCount)
	// TODO Rethink this because ideally, activity of users should be tied to tracking user sessions/login frequency ie no of people who logged in today
	var dailyNewUsers int64
	if err := walletDB.Model(&models.User{}).
		Where("DATE(created_at) = DATE(NOW())").
		Count(&dailyNewUsers).Error; err != nil {
		return nil, err
	}
	var dailySuspendedUsers int64
	if err := walletDB.Model(&models.User{}).
		Where("DATE(created_at) = DATE(NOW())").
		Count(&dailySuspendedUsers).Where("suspended = 1").Error; err != nil {
		return nil, err
	}

	var weeklyNewUsers int64
	if err := walletDB.Model(&models.User{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Count(&weeklyNewUsers).Error; err != nil {
		return nil, err
	}

	var weeklySuspendedUsers int64
	if err := walletDB.Model(&models.User{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Count(&weeklySuspendedUsers).Where("suspended = 1").Error; err != nil {
		return nil, err
	}

	var monthlyNewUsers int64
	if err := walletDB.Model(&models.User{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Count(&monthlyNewUsers).Error; err != nil {
		return nil, err
	}

	var monthlySuspendedUsers int64
	if err := walletDB.Model(&models.User{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Count(&monthlySuspendedUsers).Where("suspended = 1").Error; err != nil {
		return nil, err
	}

	userMetrics := &models.WalletUserMetricsResponse{
		TotalUsers:            *userCount,
		TotalActiveUsers:      *userCount - suspendedCount,
		TotalSessions:         0,
		AverageTimePerSession: 0,
		DailyActiveAndNewUsers: models.DailyMetrics{
			ActiveUsers: dailyNewUsers - dailySuspendedUsers,
			NewUsers:    dailyNewUsers,
		},
		WeeklyActiveAndNewUsers: models.WeeklyMetrics{
			ActiveUsers: weeklyNewUsers - weeklySuspendedUsers,
			NewUsers:    weeklyNewUsers,
		},
		MonthlyActiveAndNewUsers: models.MonthlyMetrics{
			ActiveUsers: monthlyNewUsers - monthlySuspendedUsers,
			NewUsers:    monthlyNewUsers,
		},
		WeeklyPercentageChange: models.WeeklyChange{
			//ActiveUsers: weeklyActiveUsers,
			//NewUsers: dailyNewUsers,
		},
	}

	return userMetrics, nil
}
func totalUserCount(walletDB *gorm.DB) (*int64, error) {

	var userCount int64
	if err := walletDB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return nil, err
	}
	return &userCount, nil
}

// UNCOMPLETED SINCE I NEEED WALLET DISTRIBUTION DATA PERSISTED FROM MOBILE TEAM
func GetWalletVersionDistributionData1(walletDB *gorm.DB) ([]models.WalletVersionDistribution, error) {
	// Fetch wallet version data

	userCount, err := totalUserCount(walletDB)
	if err != nil {
		return nil, err
	}
	// Fetch wallet version data
	var walletVersions []models.WalletVersionDistribution
	walletDB.Find(&walletVersions)

	// Count the total number of wallet versions
	totalCount := len(walletVersions)

	// Create a slice to store distribution data
	var versionDistribution []models.WalletVersionDistribution

	// Calculate the distribution and percentage for each version
	for _, wv := range walletVersions {
		version := wv.Version
		var found bool

		// Check if the version is already in the slice
		for i, dist := range versionDistribution {
			if dist.Version == version {
				versionDistribution[i].Count++
				found = true
				break
			}
		}

		// If the version is not found, add it to the slice
		if !found {
			versionDistribution = append(versionDistribution, models.WalletVersionDistribution{
				Version:    version,
				Count:      1,
				Percentage: 0.0, // Initialize to 0
			})
		}
	}

	// Calculate the percentage for each version
	for i := range versionDistribution {
		versionDistribution[i].Percentage = float64(versionDistribution[i].Count) / float64(totalCount) * 100
		versionDistribution[i].Count = *userCount
	}

	return versionDistribution, nil

}

func GetWalletVersionDistributionData(walletDB *gorm.DB) ([]models.WalletVersionDistribution, error) {

	// userCount, err := totalUserCount()
	//if err != nil {
	//	return nil, err
	//}
	// Fetch wallet version data
	var walletVersions []models.WalletVersionDistribution
	walletDB.Find(&walletVersions)

	// Count the total number of wallet versions
	totalCount := len(walletVersions)

	// Create a slice to store distribution data
	var versionDistribution []models.WalletVersionDistribution

	// Calculate the distribution and percentage for each version
	for _, wv := range walletVersions {
		version := wv.Version
		var found bool

		// Check if the version is already in the slice
		for i, dist := range versionDistribution {
			if dist.Version == version {
				versionDistribution[i].Count++
				found = true
				break
			}
		}

		// If the version is not found, add it to the slice
		if !found {
			versionDistribution = append(versionDistribution, models.WalletVersionDistribution{
				Version:    version,
				Count:      1,
				Percentage: 0.0, // Initialize to 0
			})
		}
	}

	// Calculate the percentage for each version
	for i := range versionDistribution {
		if totalCount != 0 {
			versionDistribution[i].Percentage = float64(versionDistribution[i].Count) / float64(totalCount) * 100
		} else {
			versionDistribution[i].Percentage = 0 // Handle division by zero
		}
		// No need to set Count to userCount, keep it as it is
	}

	return versionDistribution, nil
}

type CountryCount struct {
	CountryCode string
	Count       int
}

type CountryUserCount struct {
	Country string
	Count   int
}

func CountUsersByCountry(walletDB *gorm.DB) ([]CountryUserCount, error) {

	// Query to count users by country_code
	var countryCounts []CountryCount
	err := walletDB.Model(&models.User{}).
		Select("country_code as country_code, COUNT(*) as count").
		Group("country_code").
		Scan(&countryCounts).Error
	if err != nil {
		return nil, err
	}

	// Iterate through the results and print the counts
	for _, c := range countryCounts {
		fmt.Printf("%s: %d\n", c.CountryCode, c.Count)
	}

	// Create a slice to store the results
	var results []CountryUserCount

	// Iterate through the results and populate the CountryUserCount struct
	for _, c := range countryCounts {
		countryUserCount := CountryUserCount{
			Country: c.CountryCode,
			Count:   c.Count,
		}
		results = append(results, countryUserCount)
	}

	return results, nil
}

// func CountUsersByCountry() ([]models.CountryUserCount, error) {
//	db, err := db.TrovoWalletDb()
//	if err != nil {
//		return nil, err
//	}
//
//	// Create a map to store the counts for each country
//	countryCounts := make(map[string]int)
//
//	// Query to count users by country_code
//	db.Table("users").
//		Select("country_code, COUNT(*) as count").
//		Group("country_code").
//		Scan(&countryCounts)
//
//	// Iterate through the results and print the counts
//	for countryCode, count := range countryCounts {
//		fmt.Printf("%s: %d\n", countryCode, count)
//	}
//
//	// Create a slice to store the results
//	var results []models.CountryUserCount
//
//	// Iterate through the results and populate the CountryUserCount struct
//	for countryCode, count := range countryCounts {
//		countryUserCount := models.CountryUserCount{
//			Country: countryCode,
//			Count:   count,
//		}
//		results = append(results, countryUserCount)
//	}
//
//	return results, nil
//}

func GetTopReferrers(walletDB *gorm.DB) ([]models.ReferrerCount, error) {

	var referrerCounts []models.ReferrerCount
	walletDB.Table("users").
		Select("COALESCE(NULLIF(referrer, ''), 'noreferrer') AS referrer, COUNT(*) AS count").
		Group("referrer").
		Order("count DESC").
		Limit(10).
		Scan(&referrerCounts)
	return referrerCounts, nil
}

func GetRecentRegistrations(request models.UserRequestDTO, walletDB *gorm.DB) ([]models.UserDTO, int, error) {
	var recentRegistrations []models.UserDTO
	var total int64

	// Base query with dynamic filtering
	query := walletDB.Model(&models.User{})
	if request.Username != "" {
		query = query.Where("LOWER(username) ILIKE LOWER(?)", request.Username)
	}
	if request.Email != "" {
		query = query.Where("LOWER(email) ILIKE LOWER(?)", request.Email)
	}
	if request.Phone != "" {
		query = query.Where("mobile ILIKE ?", request.Phone)
	}
	if request.FirstName != "" {
		query = query.Where("LOWER(first_name) ILIKE LOWER(?)", request.FirstName)
	}
	if request.LastName != "" {
		query = query.Where("LOWER(last_name) ILIKE LOWER(?)", request.LastName)
	}

	if request.Search != "" {
		searchTerm := "%" + request.Search + "%"
		query = query.Where(
			walletDB.Where("LOWER(username) ILIKE LOWER(?)", searchTerm).
				Or("LOWER(email) ILIKE LOWER(?)", searchTerm).
				Or("mobile ILIKE ?", searchTerm).
				Or("LOWER(first_name) ILIKE LOWER(?)", searchTerm).
				Or("LOWER(last_name) ILIKE LOWER(?)", searchTerm),
		)
	}

	// Count total entries for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	offset := (request.Page - 1) * request.PageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(request.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	for _, user := range users {
		recentRegistrations = append(recentRegistrations, models.UserDTO{
			Username:    user.Username,
			Address:     user.Address,
			FullName:    user.FirstName + " " + user.LastName,
			Email:       user.Email,
			PhoneNumber: user.Mobile,
			Date:        user.CreatedAt,
		})
	}

	return recentRegistrations, int(total), nil
}

//	func GetP2PUserStatistics(limit int) ([]UserDTO, error) {
//		var recentRegistrations []UserDTO
//		db, err := db.P2PDb()
//		if err != nil {
//			return nil, err
//		}
//		var users []models.User
//		db.Order("created_at desc").Limit(limit).Find(&users)
//
//		for _, user := range users {
//			recentRegistrations = append(recentRegistrations, UserDTO{
//				Username:    user.Username,
//				Address:     user.Address,
//				FullName:    user.FirstName + " " + user.LastName,
//				Email:       user.Email,
//				PhoneNumber: *user.Mobile,
//				Date:        user.CreatedAt,
//			})
//		}
//
//		return recentRegistrations, nil
//	}
// ---- P2P (new schema) ----
//
// The functions below replace the old P2P admin dashboard's queries
// against a dead legacy trading platform's schema (orders/offers/users
// with float64 amounts and offer_maker/offer_taker naming) with reads
// against app-backend's current P2P module, which now owns this data in
// the same physical database (see internal/db/main.go's AdminDB doc
// comment) - see internal/models/p2p.go for the read-only mirror structs.

// P2PStatistics is the marketplace-wide aggregate the "P2P statistics"
// admin dashboard card shows.
type P2PStatistics struct {
	TotalOffers      int64 `json:"totalOffers"`
	ActiveOffers     int64 `json:"activeOffers"`
	TotalOrders      int64 `json:"totalOrders"`
	CompletedOrders  int64 `json:"completedOrders"`
	OpenDisputes     int64 `json:"openDisputes"`
	ResolvedDisputes int64 `json:"resolvedDisputes"`
}

func GetP2PStatistics(p2pDB *gorm.DB) (*P2PStatistics, error) {
	var stats P2PStatistics
	if err := p2pDB.Model(&models.P2POffer{}).Count(&stats.TotalOffers).Error; err != nil {
		return nil, err
	}
	if err := p2pDB.Model(&models.P2POffer{}).Where("status = ?", "ACTIVE").Count(&stats.ActiveOffers).Error; err != nil {
		return nil, err
	}
	if err := p2pDB.Model(&models.P2POrder{}).Count(&stats.TotalOrders).Error; err != nil {
		return nil, err
	}
	if err := p2pDB.Model(&models.P2POrder{}).Where("order_status = ?", "COMPLETED").Count(&stats.CompletedOrders).Error; err != nil {
		return nil, err
	}
	if err := p2pDB.Model(&models.P2PDispute{}).Where("status = ?", "OPEN").Count(&stats.OpenDisputes).Error; err != nil {
		return nil, err
	}
	if err := p2pDB.Model(&models.P2PDispute{}).Where("status = ?", "RESOLVED").Count(&stats.ResolvedDisputes).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

// TradeStatistics backs the "Trading" tab's stat cards, top-traders table,
// and recent-trades preview.
type TradeStatistics struct {
	TotalTrades      int64             `json:"totalTrades"`
	CompletedTrades  int64             `json:"completedTrades"`
	CancelledTrades  int64             `json:"cancelledTrades"`
	ExpiredTrades    int64             `json:"expiredTrades"`
	TradeSuccessRate float64           `json:"tradeSuccessRate"`
	OpenDisputes     int64             `json:"openDisputes"`
	TopTraders       []TraderStats     `json:"topTraders"`
	RecentTrades     []models.P2POrder `json:"recentTrades"`
}

type TraderStats struct {
	MerchantUsername string `json:"merchantUsername"`
	CompletedTrades  int64  `json:"completedTrades"`
}

func GetTradeStatistics(p2pDB *gorm.DB) (*TradeStatistics, error) {
	var stats TradeStatistics

	if err := p2pDB.Model(&models.P2POrder{}).Count(&stats.TotalTrades).Error; err != nil {
		return nil, err
	}
	if err := p2pDB.Model(&models.P2POrder{}).Where("order_status = ?", "COMPLETED").Count(&stats.CompletedTrades).Error; err != nil {
		return nil, err
	}
	if err := p2pDB.Model(&models.P2POrder{}).Where("order_status = ?", "CANCELLED").Count(&stats.CancelledTrades).Error; err != nil {
		return nil, err
	}
	if err := p2pDB.Model(&models.P2POrder{}).Where("order_status = ?", "EXPIRED").Count(&stats.ExpiredTrades).Error; err != nil {
		return nil, err
	}
	if err := p2pDB.Model(&models.P2PDispute{}).Where("status = ?", "OPEN").Count(&stats.OpenDisputes).Error; err != nil {
		return nil, err
	}
	if stats.TotalTrades > 0 {
		stats.TradeSuccessRate = float64(stats.CompletedTrades) / float64(stats.TotalTrades) * 100
	}

	var topPerformers []models.MerchantPerformance
	if err := p2pDB.Order("completed_trades DESC").Limit(5).Find(&topPerformers).Error; err != nil {
		return nil, err
	}
	for _, p := range topPerformers {
		stats.TopTraders = append(stats.TopTraders, TraderStats{MerchantUsername: p.MerchantUsername, CompletedTrades: p.CompletedTrades})
	}

	if err := p2pDB.Order("created_at DESC").Limit(5).Find(&stats.RecentTrades).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// TradeListRequest filters the paginated "all trades" admin table.
type TradeListRequest struct {
	Page      int
	PageSize  int
	OfferType string
	Status    string
	Username  string
	CreatedAt string
}

func GetTradeList(p2pDB *gorm.DB, req TradeListRequest) ([]models.P2POrder, int64, error) {
	query := p2pDB.Model(&models.P2POrder{})
	if req.OfferType != "" {
		query = query.Where("offer_type = ?", req.OfferType)
	}
	if req.Status != "" {
		query = query.Where("order_status = ?", req.Status)
	}
	if req.Username != "" {
		query = query.Where("merchant_username = ? OR customer_username = ?", req.Username, req.Username)
	}
	if req.CreatedAt != "" {
		query = query.Where("DATE(created_at) = ?", req.CreatedAt)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var orders []models.P2POrder
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// P2PUserListRequest filters the "P2P users" admin table.
type P2PUserListRequest struct {
	Page     int
	PageSize int
	Username string
	Email    string
	Phone    string
	Search   string
}

// P2PUserDTO is a wallet user annotated with their P2P trading performance.
type P2PUserDTO struct {
	ID                      string    `json:"id"`
	Username                string    `json:"username"`
	Email                   string    `json:"email"`
	Phone                   string    `json:"phone"`
	RegistrationDate        time.Time `json:"registrationDate"`
	Suspended               bool      `json:"suspended"`
	IsMerchant              bool      `json:"isMerchant"`
	MerchantCompletedTrades int64     `json:"merchantCompletedTrades"`
	MerchantCompletionRate  string    `json:"merchantCompletionRate"`
	CustomerCompletedTrades int64     `json:"customerCompletedTrades"`
	CustomerCompletionRate  string    `json:"customerCompletionRate"`
}

// GetP2PUsers paginates wallet users (walletDB) and annotates each with
// their P2P performance (p2pDB) - walletDB and p2pDB are the same physical
// connection today (see package doc above) but kept as distinct
// parameters since they're conceptually different data sources.
func GetP2PUsers(walletDB, p2pDB *gorm.DB, req P2PUserListRequest) ([]P2PUserDTO, int64, error) {
	query := walletDB.Model(&models.User{})
	if req.Username != "" {
		query = query.Where("LOWER(username) ILIKE LOWER(?)", "%"+req.Username+"%")
	}
	if req.Email != "" {
		query = query.Where("LOWER(email) ILIKE LOWER(?)", "%"+req.Email+"%")
	}
	if req.Phone != "" {
		query = query.Where("mobile ILIKE ?", "%"+req.Phone+"%")
	}
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where(
			walletDB.Where("LOWER(username) ILIKE LOWER(?)", searchTerm).
				Or("LOWER(email) ILIKE LOWER(?)", searchTerm).
				Or("mobile ILIKE ?", searchTerm),
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []models.User
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	dtos := make([]P2PUserDTO, 0, len(users))
	for _, u := range users {
		dto := P2PUserDTO{
			ID:               u.ID,
			Username:         u.Username,
			Email:            u.Email,
			Phone:            u.Mobile,
			RegistrationDate: u.CreatedAt,
			Suspended:        u.Suspended == 1,
		}
		var merchantPerf models.MerchantPerformance
		if err := p2pDB.Where("merchant_id = ?", u.ID).First(&merchantPerf).Error; err == nil {
			dto.MerchantCompletedTrades = merchantPerf.CompletedTrades
			dto.MerchantCompletionRate = merchantPerf.CompletionRate
		}
		var offerCount int64
		p2pDB.Model(&models.P2POffer{}).Where("merchant_user_id = ?", u.ID).Count(&offerCount)
		dto.IsMerchant = offerCount > 0
		var customerPerf models.CustomerPerformance
		if err := p2pDB.Where("customer_id = ?", u.ID).First(&customerPerf).Error; err == nil {
			dto.CustomerCompletedTrades = customerPerf.CompletedTrades
			dto.CustomerCompletionRate = customerPerf.CompletionRate
		}
		dtos = append(dtos, dto)
	}
	return dtos, total, nil
}

func GetUserList(request models.UserRequestDTO, walletDB *gorm.DB) ([]models.UserDto, int64, error) {
	var userDtoList []models.UserDto
	var total int64

	// Base query with dynamic filtering
	query := walletDB.Model(&models.User{})
	if request.Username != "" {
		query = query.Where("LOWER(username) ILIKE LOWER(?)", request.Username)
	}
	if request.Email != "" {
		query = query.Where("LOWER(email) ILIKE LOWER(?)", request.Email)
	}
	if request.Phone != "" {
		query = query.Where("mobile ILIKE ?", request.Phone)
	}
	if request.FirstName != "" {
		query = query.Where("LOWER(first_name) ILIKE LOWER(?)", request.FirstName)
	}
	if request.LastName != "" {
		query = query.Where("LOWER(last_name) ILIKE LOWER(?)", request.LastName)
	}

	if request.City != "" {
		query = query.Where("LOWER(city) ILIKE LOWER(?)", request.City)
	}

	if request.Address != "" {
		query = query.Where("LOWER(address) ILIKE LOWER(?)", request.Address)
	}
	if request.Search != "" {
		searchTerm := "%" + strings.ToLower(request.Search) + "%"
		query = query.Where(
			walletDB.Where("LOWER(username) ILIKE ?", searchTerm).
				Or("LOWER(email) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(mobile, '')) ILIKE ?", searchTerm).
				Or("LOWER(first_name) ILIKE ?", searchTerm).
				Or("LOWER(last_name) ILIKE ?", searchTerm).
				// Or("LOWER(COALESCE(middle_name, '')) ILIKE ?", searchTerm).
				//Or("LOWER(COALESCE(contact_phone, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(city, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(region, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(region_name, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(time_zone, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(isp, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(public_ip, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(address, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(country_code, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(referrer, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(referral_link, '')) ILIKE ?", searchTerm).
				// Or("LOWER(COALESCE(bantu_talk, '')) ILIKE ?", searchTerm).
				Or("LOWER(COALESCE(suspension_reason, '')) ILIKE ?", searchTerm),
		)
	}

	// Count total entries for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	offset := (request.Page - 1) * request.PageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(request.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	for _, user := range users {
		userDto := models.UserDto{
			ID:        user.ID,
			UserName:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			//TelegramConnectedAt: user.TelegramConnectedAt,
			LastName:  user.LastName,
			FirstName: user.FirstName,
			//MiddleName:            user.MiddleName,
			Mobile: user.Mobile,
			//ContactPhone:          user.ContactPhone,
			//Telegram:              user.Telegram,
			Email: user.Email,
			//ImageThumbnail:        user.ImageThumbnail,
			CountryCode: user.CountryCode,
			Latitude:    user.Latitude,
			Longitude:   user.Longitude,
			City:        user.City,
			Region:      user.Region,
			RegionName:  user.RegionName,
			TimeZone:    user.TimeZone,
			ISP:         user.ISP,
			PublicIP:    user.PublicIP,
			Address:     user.Address,
			//BantuTalk:             user.BantuTalk,
			Suspended: uint(user.Suspended),
			//KYCLevel:              user.KYCLevel,
			//MaxAssetPerOffer:      user.MaxAssetPerOffer,
			//MaxAssetPerOrder:      user.MaxAssetPerOrder,
			SuspensionReason: user.SuspensionReason,
			//Offline:               user.Offline,
			//AdminLevel:            user.AdminLevel,
			//TelegramNotifications: user.TelegramNotifications,
		}
		userDtoList = append(userDtoList, userDto)
	}
	return userDtoList, total, nil
}

func GetUserProfileInfo(req models.UserProfile, walletDB *gorm.DB) (*models.User, error) {

	var user models.User
	query := walletDB.Model(&models.User{})

	if req.ID != "" {
		query = query.Where("id = ?", req.ID)
	} else if req.Username != "" {
		query = query.Where("username = ?", req.Username)
	} else if req.Email != "" {
		query = query.Where("email = ?", req.Email)
	} else {
		return nil, errors.New("no search parameter provided")
	}

	err := query.First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserWallets(req models.UserProfile, walletDB *gorm.DB) ([]models.UserWalletDTO, error) {
	var userWallets []models.UserWallet
	// Fetch wallets for the user
	if err := walletDB.Where("user_id = ?", req.UserID).Find(&userWallets).Error; err != nil {
		return nil, err
	}

	// Map the result to a slice of DTOs
	var walletDTOs []models.UserWalletDTO
	for _, wallet := range userWallets {
		walletDTOs = append(walletDTOs, models.UserWalletDTO{
			ID:            wallet.ID,
			Alias:         wallet.Alias,
			Signer:        wallet.Signer,
			UserID:        wallet.UserID,
			PrimaryWallet: wallet.PrimaryWallet,
		})
	}

	return walletDTOs, nil
}

func GetSuspensionReasonByID(reasonID uint, db *gorm.DB) (string, error) {
	var reason models.UserSuspensionReason
	// Retrieve the suspension reason by ID from the database
	if err := db.First(&reason, reasonID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("suspension reason not found")
		}
		return "", err
	}
	return reason.Reason, nil
}

func GetUserActivityHistory() {

}

func GetUserTransactionHistory() {

}

func GetUserRecoveryHistory() {

}

