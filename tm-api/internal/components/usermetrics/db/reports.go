package usermetrics

import (
	"admin-panel-dashboard/internal/models"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// This file backs the "P2P Market Reports" admin dashboard - the reports
// a P2P marketplace operator typically needs (trading volume, order-status
// mix, most-traded assets/currencies, geographic spread, dispute health,
// fee revenue, merchant leaderboard, marketplace growth). Every query below
// pulls raw rows within a bounded date window and aggregates them in Go
// rather than relying on DB-specific date/grouping SQL (date_trunc is
// Postgres-only; AdminDB also runs on SQLite in dev/test - see
// internal/db/main.go's SQLite support) - this keeps every report portable
// across both without duplicating query logic per driver.

// ResolveReportRange turns a "range" query param into a [from, to] window
// ending now (UTC). Defaults to 30 days for an unrecognized/empty value.
func ResolveReportRange(rangeParam string) (from, to time.Time) {
	to = time.Now().UTC()
	days := 30
	switch rangeParam {
	case "7d":
		days = 7
	case "30d":
		days = 30
	case "90d":
		days = 90
	case "1y":
		days = 365
	}
	from = to.AddDate(0, 0, -days)
	return from, to
}

func dayKey(t time.Time) string { return t.UTC().Format("2006-01-02") }

// dayRange lists every calendar day in [from, to] so a report's series has
// no gaps on days with zero activity.
func dayRange(from, to time.Time) []string {
	days := make([]string, 0)
	start := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		days = append(days, dayKey(d))
	}
	return days
}

func addDecimalStr(a, b string) string {
	x, _ := decimal.NewFromString(a)
	y, err := decimal.NewFromString(b)
	if err != nil {
		return x.String()
	}
	return x.Add(y).String()
}

// --- Volume report: trading activity over time ---

type VolumeReportPoint struct {
	Date            string `json:"date"`
	TotalOrders     int64  `json:"totalOrders"`
	CompletedOrders int64  `json:"completedOrders"`
	CompletedVolume string `json:"completedVolume"`
}

type VolumeReport struct {
	Series          []VolumeReportPoint `json:"series"`
	TotalOrders     int64                `json:"totalOrders"`
	CompletedOrders int64                `json:"completedOrders"`
	CompletedVolume string               `json:"completedVolume"`
}

func GetP2PVolumeReport(p2pDB *gorm.DB, from, to time.Time) (*VolumeReport, error) {
	var orders []models.P2POrder
	if err := p2pDB.Select("id", "created_at", "order_status", "payment_amount").
		Where("created_at >= ? AND created_at <= ?", from, to).Find(&orders).Error; err != nil {
		return nil, err
	}

	buckets := map[string]*VolumeReportPoint{}
	for _, day := range dayRange(from, to) {
		buckets[day] = &VolumeReportPoint{Date: day, CompletedVolume: "0"}
	}

	report := &VolumeReport{CompletedVolume: "0"}
	for _, o := range orders {
		day := dayKey(o.CreatedAt)
		b, ok := buckets[day]
		if !ok {
			b = &VolumeReportPoint{Date: day, CompletedVolume: "0"}
			buckets[day] = b
		}
		b.TotalOrders++
		report.TotalOrders++
		if o.OrderStatus == "COMPLETED" {
			b.CompletedOrders++
			report.CompletedOrders++
			b.CompletedVolume = addDecimalStr(b.CompletedVolume, o.PaymentAmount)
			report.CompletedVolume = addDecimalStr(report.CompletedVolume, o.PaymentAmount)
		}
	}

	days := make([]string, 0, len(buckets))
	for d := range buckets {
		days = append(days, d)
	}
	sort.Strings(days)
	for _, d := range days {
		report.Series = append(report.Series, *buckets[d])
	}
	return report, nil
}

// --- Distribution report: what's trading, where ---

// CountItem is a plain key/count breakdown (order-status mix, dispute
// subjects/resolutions).
type CountItem struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

// VolumeItem additionally carries the summed payment volume for that key
// (assets/currencies/countries, where "how much traded" matters as much
// as "how many trades").
type VolumeItem struct {
	Key    string `json:"key"`
	Count  int64  `json:"count"`
	Volume string `json:"volume"`
}

type DistributionReport struct {
	ByAsset    []VolumeItem `json:"byAsset"`
	ByCurrency []VolumeItem `json:"byCurrency"`
	ByCountry  []VolumeItem `json:"byCountry"`
	ByStatus   []CountItem  `json:"byStatus"`
}

type volumeAcc struct {
	count  int64
	volume decimal.Decimal
}

func GetP2PDistributionReport(p2pDB *gorm.DB, from, to time.Time) (*DistributionReport, error) {
	var orders []models.P2POrder
	if err := p2pDB.Select("id", "asset", "currency", "country_code", "order_status", "payment_amount", "created_at").
		Where("created_at >= ? AND created_at <= ?", from, to).Find(&orders).Error; err != nil {
		return nil, err
	}

	byAsset := map[string]*volumeAcc{}
	byCurrency := map[string]*volumeAcc{}
	byCountry := map[string]*volumeAcc{}
	byStatus := map[string]int64{}

	accumulate := func(m map[string]*volumeAcc, key, amountStr string) {
		if key == "" {
			key = "UNKNOWN"
		}
		acc, ok := m[key]
		if !ok {
			acc = &volumeAcc{volume: decimal.Zero}
			m[key] = acc
		}
		acc.count++
		if amt, err := decimal.NewFromString(amountStr); err == nil {
			acc.volume = acc.volume.Add(amt)
		}
	}

	for _, o := range orders {
		accumulate(byAsset, o.Asset, o.PaymentAmount)
		accumulate(byCurrency, o.Currency, o.PaymentAmount)
		accumulate(byCountry, o.CountryCode, o.PaymentAmount)
		status := o.OrderStatus
		if status == "" {
			status = "UNKNOWN"
		}
		byStatus[status]++
	}

	toVolumeItems := func(m map[string]*volumeAcc) []VolumeItem {
		items := make([]VolumeItem, 0, len(m))
		for k, v := range m {
			items = append(items, VolumeItem{Key: k, Count: v.count, Volume: v.volume.String()})
		}
		sort.Slice(items, func(i, j int) bool { return items[i].Count > items[j].Count })
		return items
	}
	toCountItems := func(m map[string]int64) []CountItem {
		items := make([]CountItem, 0, len(m))
		for k, v := range m {
			items = append(items, CountItem{Key: k, Count: v})
		}
		sort.Slice(items, func(i, j int) bool { return items[i].Count > items[j].Count })
		return items
	}

	return &DistributionReport{
		ByAsset:    toVolumeItems(byAsset),
		ByCurrency: toVolumeItems(byCurrency),
		ByCountry:  toVolumeItems(byCountry),
		ByStatus:   toCountItems(byStatus),
	}, nil
}

// --- Dispute report: marketplace health/trust signal ---

type DisputeReportPoint struct {
	Date     string `json:"date"`
	Opened   int64  `json:"opened"`
	Resolved int64  `json:"resolved"`
}

type DisputeReport struct {
	Series                       []DisputeReportPoint `json:"series"`
	TotalOpened                  int64                `json:"totalOpened"`
	TotalResolved                int64                `json:"totalResolved"`
	StillOpen                    int64                `json:"stillOpen"`
	ResolutionRate                float64              `json:"resolutionRate"` // percent of opened disputes resolved
	AverageResolutionTimeSeconds int64                `json:"averageResolutionTimeSeconds"`
	BySubject                    []CountItem          `json:"bySubject"`
	ByResolution                 []CountItem          `json:"byResolution"`
}

// GetP2PDisputeReport windows on OpenedAt (a dispute counts toward the
// report it was opened in) - ResolvedAt may fall after `to` and still
// contributes to TotalResolved/resolution-time stats, but only feeds the
// day-by-day series when its resolution day is inside the window.
func GetP2PDisputeReport(p2pDB *gorm.DB, from, to time.Time) (*DisputeReport, error) {
	var disputes []models.P2PDispute
	if err := p2pDB.Where("opened_at >= ? AND opened_at <= ?", from, to).Find(&disputes).Error; err != nil {
		return nil, err
	}

	buckets := map[string]*DisputeReportPoint{}
	for _, day := range dayRange(from, to) {
		buckets[day] = &DisputeReportPoint{Date: day}
	}

	report := &DisputeReport{}
	subjectCounts := map[string]int64{}
	resolutionCounts := map[string]int64{}
	var totalResolutionSeconds int64
	var resolvedWithTiming int64

	for _, d := range disputes {
		day := dayKey(d.OpenedAt)
		b, ok := buckets[day]
		if !ok {
			b = &DisputeReportPoint{Date: day}
			buckets[day] = b
		}
		b.Opened++
		report.TotalOpened++
		subject := d.Subject
		if subject == "" {
			subject = "UNKNOWN"
		}
		subjectCounts[subject]++

		if d.Status == "RESOLVED" {
			report.TotalResolved++
			resolution := d.Resolution
			if resolution == "" {
				resolution = "UNKNOWN"
			}
			resolutionCounts[resolution]++
			if d.ResolvedAt != nil {
				if rb, ok := buckets[dayKey(*d.ResolvedAt)]; ok {
					rb.Resolved++
				}
			}
			if d.ResolutionTimeSeconds > 0 {
				totalResolutionSeconds += d.ResolutionTimeSeconds
				resolvedWithTiming++
			}
		} else {
			report.StillOpen++
		}
	}

	if report.TotalOpened > 0 {
		report.ResolutionRate = float64(report.TotalResolved) / float64(report.TotalOpened) * 100
	}
	if resolvedWithTiming > 0 {
		report.AverageResolutionTimeSeconds = totalResolutionSeconds / resolvedWithTiming
	}

	days := make([]string, 0, len(buckets))
	for d := range buckets {
		days = append(days, d)
	}
	sort.Strings(days)
	for _, d := range days {
		report.Series = append(report.Series, *buckets[d])
	}
	for k, v := range subjectCounts {
		report.BySubject = append(report.BySubject, CountItem{Key: k, Count: v})
	}
	sort.Slice(report.BySubject, func(i, j int) bool { return report.BySubject[i].Count > report.BySubject[j].Count })
	for k, v := range resolutionCounts {
		report.ByResolution = append(report.ByResolution, CountItem{Key: k, Count: v})
	}
	sort.Slice(report.ByResolution, func(i, j int) bool { return report.ByResolution[i].Count > report.ByResolution[j].Count })

	return report, nil
}

// --- Revenue report: fees actually collected (completed orders only) ---

type RevenueReportPoint struct {
	Date          string `json:"date"`
	PlatformFee   string `json:"platformFee"`
	RegulatoryFee string `json:"regulatoryFee"`
	Vat           string `json:"vat"`
	Total         string `json:"total"`
}

type RevenueReport struct {
	Series             []RevenueReportPoint `json:"series"`
	TotalPlatformFee   string               `json:"totalPlatformFee"`
	TotalRegulatoryFee string               `json:"totalRegulatoryFee"`
	TotalVat           string               `json:"totalVat"`
	TotalRevenue       string               `json:"totalRevenue"`
}

// GetP2PRevenueReport only counts COMPLETED orders - a fee is calculated
// and snapshotted on every order at creation time, but nothing is actually
// collected until settlement completes.
func GetP2PRevenueReport(p2pDB *gorm.DB, from, to time.Time) (*RevenueReport, error) {
	var orders []models.P2POrder
	if err := p2pDB.Select("id", "created_at", "order_status", "combined_platform_fee", "combined_regulatory_fee", "combined_vat").
		Where("order_status = ? AND created_at >= ? AND created_at <= ?", "COMPLETED", from, to).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	buckets := map[string]*RevenueReportPoint{}
	for _, day := range dayRange(from, to) {
		buckets[day] = &RevenueReportPoint{Date: day, PlatformFee: "0", RegulatoryFee: "0", Vat: "0", Total: "0"}
	}

	report := &RevenueReport{TotalPlatformFee: "0", TotalRegulatoryFee: "0", TotalVat: "0", TotalRevenue: "0"}
	for _, o := range orders {
		day := dayKey(o.CreatedAt)
		b, ok := buckets[day]
		if !ok {
			b = &RevenueReportPoint{Date: day, PlatformFee: "0", RegulatoryFee: "0", Vat: "0", Total: "0"}
			buckets[day] = b
		}
		b.PlatformFee = addDecimalStr(b.PlatformFee, o.CombinedPlatformFee)
		b.RegulatoryFee = addDecimalStr(b.RegulatoryFee, o.CombinedRegulatoryFee)
		b.Vat = addDecimalStr(b.Vat, o.CombinedVat)
		orderTotal := addDecimalStr(addDecimalStr(o.CombinedPlatformFee, o.CombinedRegulatoryFee), o.CombinedVat)
		b.Total = addDecimalStr(b.Total, orderTotal)

		report.TotalPlatformFee = addDecimalStr(report.TotalPlatformFee, o.CombinedPlatformFee)
		report.TotalRegulatoryFee = addDecimalStr(report.TotalRegulatoryFee, o.CombinedRegulatoryFee)
		report.TotalVat = addDecimalStr(report.TotalVat, o.CombinedVat)
		report.TotalRevenue = addDecimalStr(report.TotalRevenue, orderTotal)
	}

	days := make([]string, 0, len(buckets))
	for d := range buckets {
		days = append(days, d)
	}
	sort.Strings(days)
	for _, d := range days {
		report.Series = append(report.Series, *buckets[d])
	}
	return report, nil
}

// --- Growth report: new offers and first-time merchants over time ---

type GrowthReportPoint struct {
	Date         string `json:"date"`
	NewOffers    int64  `json:"newOffers"`
	NewMerchants int64  `json:"newMerchants"`
}

type GrowthReport struct {
	Series          []GrowthReportPoint `json:"series"`
	TotalNewOffers  int64                `json:"totalNewOffers"`
	TotalNewMerchants int64              `json:"totalNewMerchants"`
}

// firstOfferByMerchant maps merchant_user_id -> their earliest offer ever
// (across the whole table, not just the report window) - needed to tell
// "a merchant's first offer" apart from "any offer by an existing
// merchant". Computed in Go over every offer's (merchant_user_id,
// created_at) rather than a raw MIN()/GROUP BY query: scanning an
// aggregate expression's result column into time.Time is driver-specific
// (SQLite hands back a plain string, breaking the scan) - fetching through
// the typed model sidesteps that entirely, at the cost of one full table
// read.
func firstOfferByMerchant(p2pDB *gorm.DB) (map[string]time.Time, error) {
	var offers []models.P2POffer
	if err := p2pDB.Select("merchant_user_id", "created_at").Find(&offers).Error; err != nil {
		return nil, err
	}
	out := make(map[string]time.Time, len(offers))
	for _, o := range offers {
		if first, ok := out[o.MerchantUserID]; !ok || o.CreatedAt.Before(first) {
			out[o.MerchantUserID] = o.CreatedAt
		}
	}
	return out, nil
}

func GetP2PGrowthReport(p2pDB *gorm.DB, from, to time.Time) (*GrowthReport, error) {
	var offers []models.P2POffer
	if err := p2pDB.Select("id", "merchant_user_id", "created_at").
		Where("created_at >= ? AND created_at <= ?", from, to).Find(&offers).Error; err != nil {
		return nil, err
	}
	firstOffers, err := firstOfferByMerchant(p2pDB)
	if err != nil {
		return nil, err
	}

	buckets := map[string]*GrowthReportPoint{}
	for _, day := range dayRange(from, to) {
		buckets[day] = &GrowthReportPoint{Date: day}
	}

	report := &GrowthReport{}
	for _, o := range offers {
		day := dayKey(o.CreatedAt)
		b, ok := buckets[day]
		if !ok {
			b = &GrowthReportPoint{Date: day}
			buckets[day] = b
		}
		b.NewOffers++
		report.TotalNewOffers++
		if first, ok := firstOffers[o.MerchantUserID]; ok && dayKey(first) == day {
			b.NewMerchants++
			report.TotalNewMerchants++
		}
	}

	days := make([]string, 0, len(buckets))
	for d := range buckets {
		days = append(days, d)
	}
	sort.Strings(days)
	for _, d := range days {
		report.Series = append(report.Series, *buckets[d])
	}
	return report, nil
}

// --- Hourly activity report: what time of day the marketplace is busiest ---

type HourlyActivityPoint struct {
	Hour  int   `json:"hour"` // 0-23, UTC
	Count int64 `json:"count"`
}

type HourlyActivityReport struct {
	Series   []HourlyActivityPoint `json:"series"`
	PeakHour int                   `json:"peakHour"`
}

// GetP2PHourlyActivityReport counts orders by hour-of-day (UTC, summed
// across every day in the window) - an operational report for staffing/
// support-hours planning, not a per-day trend.
func GetP2PHourlyActivityReport(p2pDB *gorm.DB, from, to time.Time) (*HourlyActivityReport, error) {
	var orders []models.P2POrder
	if err := p2pDB.Select("id", "created_at").
		Where("created_at >= ? AND created_at <= ?", from, to).Find(&orders).Error; err != nil {
		return nil, err
	}

	counts := make([]int64, 24)
	for _, o := range orders {
		counts[o.CreatedAt.UTC().Hour()]++
	}

	report := &HourlyActivityReport{Series: make([]HourlyActivityPoint, 24)}
	peakCount := int64(-1)
	for h := 0; h < 24; h++ {
		report.Series[h] = HourlyActivityPoint{Hour: h, Count: counts[h]}
		if counts[h] > peakCount {
			peakCount = counts[h]
			report.PeakHour = h
		}
	}
	return report, nil
}

// --- Merchant leaderboard: full performance report, sortable ---

type MerchantLeaderboardRequest struct {
	Page     int
	PageSize int
	SortBy   string // completedTrades (default) | completedVolume | completionRate | disputesOpened
}

// GetP2PMerchantLeaderboard paginates MerchantPerformance. completedTrades
// and disputesOpened are real integer columns, sorted/paginated in SQL;
// completedVolume/completionRate are stored as decimal strings (so they
// sort correctly for values like "9.50" vs "10.00"), which a plain SQL
// ORDER BY would get wrong as a string comparison - those two sort orders
// are resolved in Go instead, over the full (realistically bounded -
// one row per merchant) result set.
func GetP2PMerchantLeaderboard(p2pDB *gorm.DB, req MerchantLeaderboardRequest) ([]models.MerchantPerformance, int64, error) {
	var total int64
	if err := p2pDB.Model(&models.MerchantPerformance{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize

	if req.SortBy == "completedVolume" || req.SortBy == "completionRate" {
		var all []models.MerchantPerformance
		if err := p2pDB.Find(&all).Error; err != nil {
			return nil, 0, err
		}
		sort.Slice(all, func(i, j int) bool {
			var a, b decimal.Decimal
			if req.SortBy == "completedVolume" {
				a, _ = decimal.NewFromString(all[i].CompletedTradeVolume)
				b, _ = decimal.NewFromString(all[j].CompletedTradeVolume)
			} else {
				a, _ = decimal.NewFromString(all[i].CompletionRate)
				b, _ = decimal.NewFromString(all[j].CompletionRate)
			}
			return a.GreaterThan(b)
		})
		end := offset + req.PageSize
		if offset >= len(all) {
			return []models.MerchantPerformance{}, total, nil
		}
		if end > len(all) {
			end = len(all)
		}
		return all[offset:end], total, nil
	}

	orderCol := "completed_trades DESC"
	if req.SortBy == "disputesOpened" {
		orderCol = "disputes_opened DESC"
	}
	var page []models.MerchantPerformance
	if err := p2pDB.Order(orderCol).Offset(offset).Limit(req.PageSize).Find(&page).Error; err != nil {
		return nil, 0, err
	}
	return page, total, nil
}
