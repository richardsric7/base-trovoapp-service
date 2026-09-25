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
func GetP2PMetrics(p2pDB *gorm.DB) (*models.P2PUserMetricsResponse, error) {

	totalUserCountValue, err := totalUserCount(p2pDB)
	if err != nil {
		return nil, err
	}
	var suspendedCount int64
	if err := p2pDB.Model(&models.User{}).Count(&suspendedCount).Where("suspended = 1").Error; err != nil {
		return nil, err
	}

	// TODO Rethink this because ideally, activity of users should be tied to tracking user sessions/login frequency ie no of people who logged in today
	var dailyNewUsers int64
	if err := p2pDB.Model(&models.User{}).
		Where("DATE(created_at) = DATE(NOW())").
		Count(&dailyNewUsers).Error; err != nil {
		return nil, err
	}
	var dailySuspendedUsers int64
	if err := p2pDB.Model(&models.User{}).
		Where("DATE(created_at) = DATE(NOW())").
		Count(&dailySuspendedUsers).Where("suspended = 1").Error; err != nil {
		return nil, err
	}

	var weeklyNewUsers int64
	if err := p2pDB.Model(&models.User{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Count(&weeklyNewUsers).Error; err != nil {
		return nil, err
	}

	var weeklySuspendedUsers int64
	if err := p2pDB.Model(&models.User{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Count(&weeklySuspendedUsers).Where("suspended = 1").Error; err != nil {
		return nil, err
	}

	var monthlyNewUsers int64
	if err := p2pDB.Model(&models.User{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Count(&monthlyNewUsers).Error; err != nil {
		return nil, err
	}

	var monthlySuspendedUsers int64
	if err := p2pDB.Model(&models.User{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Count(&monthlySuspendedUsers).Where("suspended = 1").Error; err != nil {
		return nil, err
	}

	userMetrics := &models.P2PUserMetricsResponse{
		TotalUsers:            *totalUserCountValue,
		TotalActiveUsers:      *totalUserCountValue - suspendedCount,
		TotalSessions:         0,
		AverageTimePerSession: 0,
		DailyActiveAndNewUsers: models.DailyMetrics{
			ActiveUsers: *totalUserCountValue - dailySuspendedUsers,
			NewUsers:    dailyNewUsers,
		},
		WeeklyActiveAndNewUsers: models.WeeklyMetrics{
			ActiveUsers: *totalUserCountValue - weeklySuspendedUsers,
			NewUsers:    weeklyNewUsers,
		},
		MonthlyActiveAndNewUsers: models.MonthlyMetrics{
			ActiveUsers: *totalUserCountValue - monthlySuspendedUsers,
			NewUsers:    monthlyNewUsers,
		},
		WeeklyPercentageChange: models.WeeklyChange{
			//ActiveUsers: weeklyActiveUsers,
			//NewUsers: dailyNewUsers,
		},
	}
	// todo-1. get merchant value correctly
	return userMetrics, nil
}

// Main response struct representing the API's output
type Response struct {
	Filter string       `json:"filter"` // The applied filter (e.g., "today", "yesterday", "last_week", etc.)
	Data   []TimePeriod `json:"data"`   // The time divisions (e.g., days, hours, months)
	Meta   MetaData     `json:"meta"`   // Metadata containing aggregated totals
}

// TimePeriod struct for each time division
type TimePeriod struct {
	Label       string `json:"label"`        // Name of the time division (e.g., "Monday", "00:00", "January")
	ActiveUsers int    `json:"active_users"` // Count of active users in this time division
	NewUsers    int    `json:"new_users"`    // Count of new users in this time division
}

// MetaData struct for summary statistics
type MetaData struct {
	TotalActiveUsers int `json:"total_active_users"` // Total active users across all divisions
	TotalNewUsers    int `json:"total_new_users"`    // Total new users across all divisions
}

func GetP2PMetricsNew(p2pDB *gorm.DB, filter string, startDate, endDate *time.Time) (*Response, error) {
	var periods []TimePeriod
	var totalActiveUsers, totalNewUsers int64

	// Fetch total active users (not limited to date range)
	if err := p2pDB.Model(&models.User{}).
		Where("suspended = 0").
		Count(&totalActiveUsers).Error; err != nil {
		return nil, err
	}

	switch filter {
	case "custom":
		// Ensure startDate and endDate are provided
		if startDate == nil || endDate == nil {
			return nil, fmt.Errorf("start_date and end_date must be provided for custom filter")
		}

		// Fetch daily data within the custom range
		for day := *startDate; day.Before(*endDate) || day.Equal(*endDate); day = day.AddDate(0, 0, 1) {
			dayStart := day.Truncate(24 * time.Hour)
			dayEnd := dayStart.AddDate(0, 0, 1)

			_, newUsers, err := fetchMetricsForRange(p2pDB, dayStart, dayEnd)
			if err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       dayStart.Format("2006-01-02"), // Include exact date
				ActiveUsers: int(totalActiveUsers),         // Use global active users count
				NewUsers:    int(newUsers),
			})
			totalNewUsers += newUsers
		}

	case "today":
		// Fetch hourly data for today
		for i := 0; i < 24; i++ {
			start := time.Now().Truncate(24 * time.Hour).Add(time.Duration(i) * time.Hour)
			end := start.Add(time.Hour)

			_, newUsers, err := fetchMetricsForRange(p2pDB, start, end)
			if err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       fmt.Sprintf("%s %s", start.Format("15:00"), start.Format("2006-01-02")), // Include exact date
				ActiveUsers: int(totalActiveUsers),                                                   // Use global active users count
				NewUsers:    int(newUsers),
			})
			totalNewUsers += newUsers
		}

	case "yesterday":
		// Fetch hourly data for yesterday
		yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
		for i := 0; i < 24; i++ {
			start := yesterday.Add(time.Duration(i) * time.Hour)
			end := start.Add(time.Hour)

			_, newUsers, err := fetchMetricsForRange(p2pDB, start, end)
			if err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       fmt.Sprintf("%s %s", start.Format("15:00"), start.Format("2006-01-02")), // Include exact date
				ActiveUsers: int(totalActiveUsers),                                                   // Use global active users count
				NewUsers:    int(newUsers),
			})
			totalNewUsers += newUsers
		}

	case "last_week":
		// Fetch daily data for the last 7 days (Monday-Sunday)
		startOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
		for i := 0; i < 7; i++ {
			dayStart := startOfWeek.AddDate(0, 0, i)
			dayEnd := dayStart.AddDate(0, 0, 1)

			_, newUsers, err := fetchMetricsForRange(p2pDB, dayStart, dayEnd)
			if err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       dayStart.Weekday().String(), // Daily label
				ActiveUsers: int(totalActiveUsers),       // Use global active users count
				NewUsers:    int(newUsers),
			})
			totalNewUsers += newUsers
		}

	case "last_month":
		// Fetch weekly data for the last month
		startOfMonth := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
		for i := 0; i < 4; i++ {
			weekStart := startOfMonth.AddDate(0, 0, i*7)
			weekEnd := weekStart.AddDate(0, 0, 7)

			_, newUsers, err := fetchMetricsForRange(p2pDB, weekStart, weekEnd)
			if err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       fmt.Sprintf("Week %d (%s)", i+1, weekStart.Format("2006-01-02")), // Include start date of the week
				ActiveUsers: int(totalActiveUsers),                                            // Use global active users count
				NewUsers:    int(newUsers),
			})
			totalNewUsers += newUsers
		}

	case "last_3_months":
		// Fetch monthly data for the last 3 months
		for i := 3; i > 0; i-- {
			monthStart := time.Now().AddDate(0, -i, 0).Truncate(24 * time.Hour)
			monthEnd := monthStart.AddDate(0, 1, 0)

			_, newUsers, err := fetchMetricsForRange(p2pDB, monthStart, monthEnd)
			if err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       fmt.Sprintf("%s (%s)", monthStart.Format("January"), monthStart.Format("2006-01-02")), // Include exact date
				ActiveUsers: int(totalActiveUsers),                                                                 // Use global active users count
				NewUsers:    int(newUsers),
			})
			totalNewUsers += newUsers
		}

	case "last_year":
		// Fetch monthly data for the last year
		for i := 12; i > 0; i-- {
			monthStart := time.Now().AddDate(0, -i, 0).Truncate(24 * time.Hour)
			monthEnd := monthStart.AddDate(0, 1, 0)

			_, newUsers, err := fetchMetricsForRange(p2pDB, monthStart, monthEnd)
			if err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       fmt.Sprintf("%s (%s)", monthStart.Format("January"), monthStart.Format("2006-01-02")), // Include exact date
				ActiveUsers: int(totalActiveUsers),                                                                 // Use global active users count
				NewUsers:    int(newUsers),
			})
			totalNewUsers += newUsers
		}

	default:
		return nil, fmt.Errorf("invalid filter value")
	}

	// Return the response
	return &Response{
		Filter: filter,
		Data:   periods,
		Meta: MetaData{
			TotalActiveUsers: int(totalActiveUsers),
			TotalNewUsers:    int(totalNewUsers),
		},
	}, nil
}

func fetchMetricsForRange(p2pDB *gorm.DB, start, end time.Time) (int64, int64, error) {
	var newUsers int64

	if err := p2pDB.Model(&models.User{}).
		Where("created_at BETWEEN ? AND ? AND suspended = 0", start, end).
		Count(&newUsers).Error; err != nil {
		return 0, 0, err
	}

	return 0, newUsers, nil
}

func GetP2PMetricsNeww(p2pDB *gorm.DB, filter string, startDate, endDate *time.Time) (*Response, error) {
	var periods []TimePeriod
	var totalActiveUsers, totalNewUsers int64

	switch filter {
	case "custom":
		// Ensure startDate and endDate are provided
		if startDate == nil || endDate == nil {
			return nil, fmt.Errorf("start_date and end_date must be provided for custom filter")
		}

		// Fetch daily data within the custom range
		for day := *startDate; day.Before(*endDate) || day.Equal(*endDate); day = day.AddDate(0, 0, 1) {
			dayStart := day.Truncate(24 * time.Hour)
			dayEnd := dayStart.AddDate(0, 0, 1)

			var activeUsers, newUsers int64
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", dayStart, dayEnd).
				Count(&newUsers).Error; err != nil {
				return nil, err
			}
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", dayStart, dayEnd).
				Count(&activeUsers).Error; err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       dayStart.Format("2006-01-02"), // Format as YYYY-MM-DD
				ActiveUsers: int(activeUsers),
				NewUsers:    int(newUsers),
			})
			totalActiveUsers += activeUsers
			totalNewUsers += newUsers
		}

	case "today":
		// Fetch hourly data for today
		for i := 0; i < 24; i++ {
			start := time.Now().Truncate(24 * time.Hour).Add(time.Duration(i) * time.Hour)
			end := start.Add(time.Hour)

			var activeUsers, newUsers int64
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", start, end).
				Count(&newUsers).Error; err != nil {
				return nil, err
			}
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", start, end).
				Count(&activeUsers).Error; err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       start.Format("15:00"), // Hourly label
				ActiveUsers: int(activeUsers),
				NewUsers:    int(newUsers),
			})
			totalActiveUsers += activeUsers
			totalNewUsers += newUsers
		}

	case "yesterday":
		// Fetch hourly data for yesterday
		yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
		for i := 0; i < 24; i++ {
			start := yesterday.Add(time.Duration(i) * time.Hour)
			end := start.Add(time.Hour)

			var activeUsers, newUsers int64
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", start, end).
				Count(&newUsers).Error; err != nil {
				return nil, err
			}
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", start, end).
				Count(&activeUsers).Error; err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       start.Format("15:00"), // Hourly label
				ActiveUsers: int(activeUsers),
				NewUsers:    int(newUsers),
			})
			totalActiveUsers += activeUsers
			totalNewUsers += newUsers
		}

	case "last_week":
		// Fetch daily data for the last 7 days (Monday-Sunday)
		startOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
		for i := 0; i < 7; i++ {
			day := startOfWeek.AddDate(0, 0, i)

			var activeUsers, newUsers int64
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", day, day.AddDate(0, 0, 1)).
				Count(&newUsers).Error; err != nil {
				return nil, err
			}
			if err := p2pDB.Model(&models.User{}).
				Where("suspended = 0 AND created_at BETWEEN ? AND ?", day, day.AddDate(0, 0, 1)).
				Count(&activeUsers).Error; err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       day.Weekday().String(), // Daily label
				ActiveUsers: int(activeUsers),
				NewUsers:    int(newUsers),
			})
			totalActiveUsers += activeUsers
			totalNewUsers += newUsers
		}

	case "last_month":
		// Fetch weekly data for the last month
		startOfMonth := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
		for i := 0; i < 4; i++ {
			weekStart := startOfMonth.AddDate(0, 0, i*7)
			weekEnd := weekStart.AddDate(0, 0, 7)

			var activeUsers, newUsers int64
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", weekStart, weekEnd).
				Count(&newUsers).Error; err != nil {
				return nil, err
			}
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", weekStart, weekEnd).
				Count(&activeUsers).Error; err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       "Week " + string(rune(i+1)),
				ActiveUsers: int(activeUsers),
				NewUsers:    int(newUsers),
			})
			totalActiveUsers += activeUsers
			totalNewUsers += newUsers
		}

	case "last_3_months":
		// Fetch monthly data for the last 3 months
		for i := 3; i > 0; i-- {
			monthStart := time.Now().AddDate(0, -i, 0).Truncate(24 * time.Hour)
			monthEnd := monthStart.AddDate(0, 1, 0)

			var activeUsers, newUsers int64
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", monthStart, monthEnd).
				Count(&newUsers).Error; err != nil {
				return nil, err
			}
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", monthStart, monthEnd).
				Count(&activeUsers).Error; err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       monthStart.Format("January"), // Monthly label
				ActiveUsers: int(activeUsers),
				NewUsers:    int(newUsers),
			})
			totalActiveUsers += activeUsers
			totalNewUsers += newUsers
		}

	case "last_year":
		// Fetch monthly data for the last year
		for i := 12; i > 0; i-- {
			monthStart := time.Now().AddDate(0, -i, 0).Truncate(24 * time.Hour)
			monthEnd := monthStart.AddDate(0, 1, 0)

			var activeUsers, newUsers int64
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", monthStart, monthEnd).
				Count(&newUsers).Error; err != nil {
				return nil, err
			}
			if err := p2pDB.Model(&models.User{}).
				Where("created_at BETWEEN ? AND ? AND suspended = 0", monthStart, monthEnd).
				Count(&activeUsers).Error; err != nil {
				return nil, err
			}

			periods = append(periods, TimePeriod{
				Label:       monthStart.Format("January"), // Monthly label
				ActiveUsers: int(activeUsers),
				NewUsers:    int(newUsers),
			})
			totalActiveUsers += activeUsers
			totalNewUsers += newUsers
		}
	}

	// Return the response
	return &Response{
		Filter: filter,
		Data:   periods,
		Meta: MetaData{
			TotalActiveUsers: int(totalActiveUsers),
			TotalNewUsers:    int(totalNewUsers),
		},
	}, nil
}

// func GetTradesOnAppealList() ([]models.TradesOnAppeal, error) {
//	db, err := db.P2PDb()
//	if err != nil {
//		return nil, err
//	}
//	// Fetch data from the order_appeals table
//	var orderAppeals []models.OrderAppeal
//	//db.Table("order_appeals").Select("id, created_at, order_id, appealed_by, appeal_reason, appeal_detail").Scan(&orderAppeals)
//	db.Table("order_appeals").Scan(&orderAppeals)
//
//	var tradesOnAppeal []models.TradesOnAppeal
//	// Output the results
//	for _, orderAppeal := range orderAppeals {
//		tradeOnAppeal := models.TradesOnAppeal{
//			CustomerName: orderAppeal.AppealedBy,
//			Description:  orderAppeal.AppealDetail,
//			Status:       "PENDING",
//			UpdatedOn:    orderAppeal.CreatedAt, // Assuming you want to use CreatedAt here
//			AppealID:     orderAppeal.OrderID,
//			Attachment:   "Transaction Agreement.docx", // You might want to fetch this from somewhere
//		}
//		tradesOnAppeal = append(tradesOnAppeal, tradeOnAppeal)
//	}
//	return tradesOnAppeal, nil
//}

//	func GetAllTradesList() ([]models.P2POrder, error) {
//		db, err := db.P2PDb()
//		if err != nil {
//			return nil, err
//		}
//		// Fetch data from the order_appeals table
//		var orders []models.Order
//		db.Table("orders").Select("id, created_at, price, amount, total, traded_on, customer_name").Scan(&orders)
//		var p2POrders []models.P2POrder
//		// Output the results
//		for i := 0; i < len(orders); i++ {
//			for j := 0; j < len(p2POrders); j++ {
//				p2POrders[j].Price = orders[i].OfferAssetPrice
//				p2POrders[j].Amount = orders[i].OrderAmount
//				p2POrders[j].Total = 000
//				p2POrders[j].TradedOn = orders[i].CreatedAt
//				p2POrders[j].CustomerName = orders[i].OfferMaker
//			}
//		}
//		return p2POrders, nil
//
// }
//
//	func GetAllTradesList() ([]models.P2POrder, error) {
//		db, err := db.P2PDb()
//		if err != nil {
//			return nil, err
//		}
//		// Fetch data from the orders table
//		var orders []models.Order
//		//db.Table("orders").Select("id, created_at, price, amount, total, traded_on, customer_name").Scan(&orders)
//		db.Table("orders").Scan(&orders)
//
//		var p2POrders []models.P2POrder
//		// Log the data directly from the database
//
//		// Output the results
//		for _, order := range orders {
//			p2POrder := models.P2POrder{
//				Price:        order.OfferAssetPrice,
//				Amount:       order.OrderAmount,
//				Total:        order.OfferAssetPrice * order.OrderAmount,
//				TradedOn:     order.CreatedAt,
//				CustomerName: order.OfferMaker,
//				OrderType:    order.OrderType,
//			}
//			p2POrders = append(p2POrders, p2POrder)
//		}
//		return p2POrders, nil
//	}

func GetTradesOnAppealList(request models.AppealListRequest, p2pDB *gorm.DB) ([]models.TradesOnAppeal, int, error) {
	var tradesOnAppeal []models.TradesOnAppeal
	var total int64

	// Base query with dynamic filtering
	query := p2pDB.Model(&models.OrderAppeal{})
	if request.AppealID != "" {
		query = query.Where("order_id ILIKE ?", "%"+request.AppealID+"%")
	}
	if request.CustomerName != "" {
		query = query.Where("appealed_by ILIKE ?", "%"+request.CustomerName+"%")
	}
	if request.Description != "" {
		query = query.Where("appeal_detail ILIKE ?", "%"+request.Description+"%")
	}
	if request.UpdatedOn != "" {
		query = query.Where("DATE(created_at) = ?", request.UpdatedOn)
	}
	if request.Status != "" {
		query = query.Where("status ILIKE ?", "%"+request.Status+"%")
	}
	if request.Search != "" {
		searchTerm := "%" + request.Search + "%"
		query = query.Where(
			p2pDB.Where("order_id ILIKE ?", searchTerm).
				Or("appealed_by ILIKE ?", searchTerm).
				Or("appeal_detail ILIKE ?", searchTerm),
		)
	}

	// Count total entries for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orderAppeals []models.OrderAppeal
	offset := (request.Page - 1) * request.PageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(request.PageSize).Find(&orderAppeals).Error; err != nil {
		return nil, 0, err
	}

	// Convert OrderAppeal to TradesOnAppeal
	for _, orderAppeal := range orderAppeals {
		tradeOnAppeal := models.TradesOnAppeal{
			CustomerName: orderAppeal.AppealedBy,
			Description:  orderAppeal.AppealDetail,
			Status:       "PENDING",
			UpdatedOn:    orderAppeal.CreatedAt, // Assuming you want to use CreatedAt here
			AppealID:     orderAppeal.OrderID,
			Attachment:   "Transaction Agreement.docx", // You might want to fetch this from somewhere
		}
		tradesOnAppeal = append(tradesOnAppeal, tradeOnAppeal)
	}

	return tradesOnAppeal, int(total), nil
}

func GetAllTradesList(req models.TradeListRequest, p2pDB *gorm.DB) ([]models.Order, int64, error) {

	var total int64
	query := p2pDB.Table("orders")

	if req.ID != "" {
		query = query.Where("id = ?", req.ID)
	}
	if req.OrderType != "" {
		query = query.Where("order_type = ?", req.OrderType)
	}
	if !req.CreatedAt.IsZero() {
		query = query.Where("created_at = ?", req.CreatedAt)
	}
	if !req.UpdatedAt.IsZero() {
		query = query.Where("updated_at = ?", req.UpdatedAt)
	}
	if !req.ExpiresAt.IsZero() {
		query = query.Where("expires_at = ?", req.ExpiresAt)
	}
	if !req.AcceptedAt.IsZero() {
		query = query.Where("accepted_at = ?", req.AcceptedAt)
	}
	if !req.CancelAfter.IsZero() {
		query = query.Where("cancel_after = ?", req.CancelAfter)
	}
	if req.OfferID != "" {
		query = query.Where("offer_id = ?", req.OfferID)
	}
	if req.OfferType != "" {
		query = query.Where("offer_type = ?", req.OfferType)
	}
	if req.OfferMaker != "" {
		query = query.Where("offer_maker = ?", req.OfferMaker)
	}
	// req.Username matches either offer_maker or offer_taker
	if req.Username != "" { // Check if username is provided before modifying and applying filter
		usernameWithSuffix := req.Username + "@trovo"
		query = query.Where("offer_maker = ? OR offer_taker = ?", usernameWithSuffix, usernameWithSuffix)
	}
	if req.OfferMakerPhone != "" {
		query = query.Where("offer_maker_phone = ?", req.OfferMakerPhone)
	}
	if req.OfferMakerCountryCode != nil {
		query = query.Where("offer_maker_country_code = ?", *req.OfferMakerCountryCode)
	}
	if req.OfferMaxTimePerTransaction != 0 {
		query = query.Where("offer_max_time_per_transaction = ?", req.OfferMaxTimePerTransaction)
	}
	if req.OfferTaker != "" {
		query = query.Where("offer_taker = ?", req.OfferTaker)
	}
	if req.OfferTakerPhone != "" {
		query = query.Where("offer_taker_phone = ?", req.OfferTakerPhone)
	}
	if req.OfferPaymentMethodID != nil {
		query = query.Where("offer_payment_method_id = ?", *req.OfferPaymentMethodID)
	}
	if req.OfferPaymentChannelID != nil {
		query = query.Where("offer_payment_channel_id = ?", *req.OfferPaymentChannelID)
	}
	if req.OfferPaymentMethodName != nil {
		query = query.Where("offer_payment_method_name = ?", *req.OfferPaymentMethodName)
	}
	if req.OfferPaymentMethodDestinationAccount != nil {
		query = query.Where("offer_payment_method_destination_account = ?", *req.OfferPaymentMethodDestinationAccount)
	}
	if req.OfferPaymentMethodMemo != nil {
		query = query.Where("offer_payment_method_memo = ?", *req.OfferPaymentMethodMemo)
	}
	if req.OfferPaymentMethodBankName != nil {
		query = query.Where("offer_payment_method_bank_name = ?", *req.OfferPaymentMethodBankName)
	}
	if req.OfferPaymentMethodAccountOpeningBranch != nil {
		query = query.Where("offer_payment_method_account_opening_branch = ?", *req.OfferPaymentMethodAccountOpeningBranch)
	}
	if req.OfferCurrencyPaymentMethodID != "" {
		query = query.Where("offer_currency_payment_method_id = ?", req.OfferCurrencyPaymentMethodID)
	}
	if req.OfferCurrencyPaymentChannelID != "" {
		query = query.Where("offer_currency_payment_channel_id = ?", req.OfferCurrencyPaymentChannelID)
	}
	if req.OfferCurrencyPaymentMethodName != nil {
		query = query.Where("offer_currency_payment_method_name = ?", *req.OfferCurrencyPaymentMethodName)
	}
	if req.OfferCurrencyPaymentMethodDestinationAccount != "" {
		query = query.Where("offer_currency_payment_method_destination_account = ?", req.OfferCurrencyPaymentMethodDestinationAccount)
	}
	if req.OfferCurrencyPaymentMethodMemo != nil {
		query = query.Where("offer_currency_payment_method_memo = ?", *req.OfferCurrencyPaymentMethodMemo)
	}
	if req.OfferCurrencyPaymentMethodBankName != nil {
		query = query.Where("offer_currency_payment_method_bank_name = ?", *req.OfferCurrencyPaymentMethodBankName)
	}
	if req.OfferCurrencyPaymentMethodAccountOpeningBranch != nil {
		query = query.Where("offer_currency_payment_method_account_opening_branch = ?", *req.OfferCurrencyPaymentMethodAccountOpeningBranch)
	}
	if req.OfferCurrencyPaymentMethodCountryCode != nil {
		query = query.Where("offer_currency_payment_method_country_code = ?", *req.OfferCurrencyPaymentMethodCountryCode)
	}
	if req.OfferCurrencyPaymentMethodCurrencyID != nil {
		query = query.Where("offer_currency_payment_method_currency_id = ?", *req.OfferCurrencyPaymentMethodCurrencyID)
	}
	if req.OfferCurrencyID != "" {
		query = query.Where("offer_currency_id = ?", req.OfferCurrencyID)
	}
	if req.OfferAssetAmount != 0 {
		query = query.Where("offer_asset_amount = ?", req.OfferAssetAmount)
	}
	if req.OfferAssetID != "" {
		query = query.Where("offer_asset_id = ?", req.OfferAssetID)
	}
	if req.OfferAssetPrice != 0 {
		query = query.Where("offer_asset_price = ?", req.OfferAssetPrice)
	}
	if req.OfferMinTradeAmount != 0 {
		query = query.Where("offer_min_trade_amount = ?", req.OfferMinTradeAmount)
	}
	if req.OfferMaxTradeAmount != 0 {
		query = query.Where("offer_max_trade_amount = ?", req.OfferMaxTradeAmount)
	}
	if req.OfferRemark != nil {
		query = query.Where("offer_remark = ?", *req.OfferRemark)
	}
	if req.TakerPaymentMethodID != "" {
		query = query.Where("taker_payment_method_id = ?", req.TakerPaymentMethodID)
	}
	if req.TakerPaymentChannelID != "" {
		query = query.Where("taker_payment_channel_id = ?", req.TakerPaymentChannelID)
	}
	if req.TakerPaymentMethodName != nil {
		query = query.Where("taker_payment_method_name = ?", *req.TakerPaymentMethodName)
	}
	if req.TakerPaymentMethodDestinationAccount != "" {
		query = query.Where("taker_payment_method_destination_account = ?", req.TakerPaymentMethodDestinationAccount)
	}
	if req.TakerPaymentMethodMemo != nil {
		query = query.Where("taker_payment_method_memo = ?", *req.TakerPaymentMethodMemo)
	}
	if req.TakerPaymentMethodBankName != nil {
		query = query.Where("taker_payment_method_bank_name = ?", *req.TakerPaymentMethodBankName)
	}
	if req.TakerPaymentMethodAccountOpeningBranch != nil {
		query = query.Where("taker_payment_method_account_opening_branch = ?", *req.TakerPaymentMethodAccountOpeningBranch)
	}
	if req.TakerPaymentMethodCountryCode != nil {
		query = query.Where("taker_payment_method_country_code = ?", *req.TakerPaymentMethodCountryCode)
	}
	if req.TakerPaymentMethodCurrencyID != nil {
		query = query.Where("taker_payment_method_currency_id = ?", *req.TakerPaymentMethodCurrencyID)
	}
	if req.OrderEscrowAddress != "" {
		query = query.Where("order_escrow_address = ?", req.OrderEscrowAddress)
	}
	if req.OrderPaymentMemo != "" {
		query = query.Where("order_payment_memo = ?", req.OrderPaymentMemo)
	}
	if req.OrderAmount != 0 {
		query = query.Where("order_amount = ?", req.OrderAmount)
	}
	if req.OrderMakerFee != 0 {
		query = query.Where("order_maker_fee = ?", req.OrderMakerFee)
	}
	if req.OrderTakerFee != 0 {
		query = query.Where("order_taker_fee = ?", req.OrderTakerFee)
	}
	if req.OrderEscrowTransactionID != nil {
		query = query.Where("order_escrow_transaction_id = ?", *req.OrderEscrowTransactionID)
	}
	if req.OrderAssetReleaseTransactionID != nil {
		query = query.Where("order_asset_release_transaction_id = ?", *req.OrderAssetReleaseTransactionID)
	}
	if req.OrderStatusID != 0 {
		query = query.Where("order_status_id = ?", req.OrderStatusID)
	}
	// if req.OrderStatus != "" {
	//	query = query.Where("order_status = ?", req.OrderStatus)
	//}
	if req.DynamicLink != nil {
		query = query.Where("dynamic_link = ?", *req.DynamicLink)
	}
	if req.QRCode != nil {
		query = query.Where("qr_code = ?", *req.QRCode)
	}
	if req.FiatDepositTransactionID != nil {
		query = query.Where("fiat_deposit_transaction_id = ?", *req.FiatDepositTransactionID)
	}

	query.Count(&total)

	page := req.Page
	pageSize := req.PageSize
	offset := (page - 1) * pageSize

	query = query.Offset(offset).Limit(pageSize)

	var orders []models.Order
	if err := query.Table("orders").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

//	func GetUserList() ([]models.UserDto, error) {
//		db, err := db.TrovoWalletDb()
//		if err != nil {
//			return nil, err
//		}
//		// Fetch data from the users table
//		var userList []models.User
//		db.Table("users").Scan(&userList)
//		//db.Table("users").Select("id, first_name, last_name, username, email, contact_phone, country_code").Scan(&userList)
//		var userDtoList []models.UserDto
//		// Output the results
//		for _, user := range userList {
//			userDto := models.UserDto{
//				CustomerName: user.FirstName + " " + user.LastName,
//				UserName:     user.Username,
//				EmailAddress: user.Email,
//				Status:       strconv.Itoa(int(user.Suspended)), // You might want to set this based on some condition
//				PhoneNumber:  *user.Mobile,
//				Location:     *user.CountryCode,
//			}
//			userDtoList = append(userDtoList, userDto)
//		}
//		return userDtoList, nil
//	}
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

func GetTradeStatistics(p2pDB *gorm.DB) (*TradeStatistics, error) {
	var stats TradeStatistics
	var totalOfferAssetAmount float64
	var totalTrades int64

	err := p2pDB.Transaction(func(tx *gorm.DB) error {
		// Count total trades
		if err := tx.Model(&models.Order{}).Count(&totalTrades).Error; err != nil {
			return err
		}
		stats.TotalTrades = totalTrades

		// Count completed trades
		if err := tx.Model(&models.Order{}).Where("order_status_id = ?", 12).Count(&stats.CompletedTrades).Error; err != nil {
			return err
		}

		// Count cancelled trades
		if err := tx.Model(&models.Order{}).Where("order_status_id != 12").Count(&stats.CancelledTrades).Error; err != nil {
			return err
		}

		// Get total trade volume (sum of order amounts)
		if err := tx.Model(&models.Order{}).Select("COALESCE(SUM(offer_asset_amount), 0)").Scan(&totalOfferAssetAmount).Error; err != nil {
			return err
		}
		stats.TotalVolume = totalOfferAssetAmount

		// Calculate average trade size
		if totalTrades > 0 {
			stats.AverageTradeSize = totalOfferAssetAmount / float64(totalTrades)
		}

		// Calculate trade success rate
		if totalTrades > 0 {
			stats.TradeSuccessRate = (float64(stats.CompletedTrades) / float64(totalTrades)) * 100
		}

		// Count trades on appeal
		if err := tx.Model(&models.OrderAppeal{}).Count(&stats.TradesOnAppeal).Error; err != nil {
			return err
		}

		// Get top 5 traders by order count
		if err := tx.Model(&models.Order{}).
			Select("offer_maker, COUNT(*) as order_count").
			Group("offer_maker").
			Order("order_count DESC").
			Limit(5).
			Find(&stats.TopTraders).Error; err != nil {
			return err
		}

		// Get 5 most recent trades
		if err := tx.Model(&models.Order{}).
			Order("created_at DESC").
			Limit(5).
			Find(&stats.RecentTrades).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

type TradeStatistics struct {
	TotalTrades      int64          `json:"total_trades"`
	CompletedTrades  int64          `json:"completed_trades"`
	CancelledTrades  int64          `json:"cancelled_trades"`
	TotalVolume      float64        `json:"total_volume"`
	AverageTradeSize float64        `json:"average_trade_size"`
	TradeSuccessRate float64        `json:"trade_success_rate"`
	TradesOnAppeal   int64          `json:"trades_on_appeal"`
	TopTraders       []TraderStats  `json:"top_traders"`
	RecentTrades     []models.Order `json:"recent_trades"`
}

type TraderStats struct {
	OfferMaker string `json:"offer_maker"`
	OrderCount int    `json:"order_count"`
}
