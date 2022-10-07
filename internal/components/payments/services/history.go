package payments

import (
	"log"
	"strconv"
	"strings"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	db "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetPaymentHistory(targetPublicKey string, gc *sharedconfig.GlobalConfig, c *gin.Context) (records paymentModels.PaginatedPaymentHistory) {
	var err error
	var paymentHistories []paymentModels.PaymentHistory
	records.Records = make([]paymentModels.PaymentHistoryJSON, 0)
	DB, _ := db.OpenDb()
	DBC, _ := db.OpenDb()

	// if err != nil {
	// 	log.Fatalf("[main]Error opening DB %s", err)
	// 	return
	// }
	// DB := gc.DB
	// DBC := gc.DB
	var query *gorm.DB
	var countQuery *gorm.DB
	oD := "ASC"
	transactionType := c.Query("transactionType")
	fromPublicKey := strings.ToUpper(c.Query("fromPublicKey"))
	toPublicKey := strings.ToUpper(c.Query("toPublicKey"))
	name := c.Query("name")
	memo := c.Query("memo")
	limitU, _ := strconv.ParseUint(c.DefaultQuery("limit", "25"), 10, 64)
	limit := int(limitU)
	pageU, _ := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	page := int(pageU)
	assetIssuer := strings.ToUpper(c.Query("assetIssuer"))
	// var assetIssuerVal *string
	assetCode := strings.ToUpper(c.Query("assetCode"))
	// if strings.EqualFold(assetCode, "XBN") {
	// 	assetIssuerVal = nil
	// }
	// if len(assetCode) > 1 && len(assetIssuer) == 56 {
	// 	assetIssuerVal = &assetIssuer
	// }
	transactionID := c.Query("transactionID")
	amountBetween := c.Query("amount")
	dateBetween := c.Query("dateBetween")

	orderBy := c.DefaultQuery("orderby", "transaction_date")
	orderDirection := c.DefaultQuery("order", "DESC")

	query = DB.Preload(clause.Associations)
	countQuery = DBC.Group("id")

	if len(orderDirection) > 0 && strings.ToLower(orderDirection) == "desc" {
		oD = "DESC"
	}
	if len(orderBy) > 0 {
		query = query.Order(orderBy + " " + oD)
		countQuery = countQuery.Order(orderBy + " " + oD)

	} else {
		query = query.Order("transactionDate DESC")
		countQuery = countQuery.Order("transactionDate DESC")
	}

	{
		query = query.Where("(from_public_key = ? OR to_public_key = ?)", targetPublicKey, targetPublicKey)
		countQuery = countQuery.Where("(from_public_key = ? OR to_public_key = ?)", targetPublicKey, targetPublicKey)

	}

	if len(fromPublicKey) == 56 {
		query = query.Where("from_public_key = ?", fromPublicKey)
		countQuery = countQuery.Where("from_public_key = ?", fromPublicKey)

	}

	if len(toPublicKey) == 56 {
		query = query.Where("to_public_key = ?", toPublicKey)
		countQuery = countQuery.Where("to_public_key = ?", toPublicKey)

	}
	if len(assetCode) > 0 {
		query = query.Where("asset_code = ?", strings.ToUpper(assetCode))
		countQuery = countQuery.Where("asset_code = ?", strings.ToUpper(assetCode))

	}
	if len(assetIssuer) == 56 {
		query = query.Where("asset_issuer = ?", assetIssuer)
		countQuery = countQuery.Where("asset_issuer = ?", assetIssuer)

	}
	if len(name) > 2 {

		query = query.Where(`(lower("from") LIKE ? OR lower("to") LIKE ?)`, "%"+strings.ToLower(name)+"%", "%"+strings.ToLower(name)+"%")
		countQuery = countQuery.Where(`(lower("from") LIKE ? OR lower("to") LIKE ?)`, "%"+strings.ToLower(name)+"%", "%"+strings.ToLower(name)+"%")

	}
	if len(memo) > 2 {

		query = query.Where("lower(memo) LIKE ?", strings.ToLower(memo)+"%")
		countQuery = countQuery.Where("lower(memo) LIKE ?", strings.ToLower(memo)+"%")

	}
	if len(transactionType) > 4 {

		query = query.Where("lower(transaction_type) LIKE ?", strings.ToUpper(transactionType)+"%")
		countQuery = countQuery.Where("lower(transaction_type) LIKE ?", strings.ToUpper(transactionType)+"%")

	}
	if len(transactionID) > 4 {

		query = query.Where("transaction_id = ?", strings.TrimSpace(transactionID))
		countQuery = countQuery.Where("transaction_id = ?", strings.TrimSpace(transactionID))

	}

	if len(dateBetween) == 21 && strings.Contains(dateBetween, ";") {
		// 2020-01-01:2020-02-31 full range date
		dateRange := strings.Split(dateBetween, ";")
		query = query.Where("transaction_date::date BETWEEN ?::date AND ?::date", dateRange[0], dateRange[1])
		countQuery = countQuery.Where("transaction_date::date BETWEEN ?::date AND ?::date", dateRange[0], dateRange[1])

	}
	if len(amountBetween) > 2 && strings.Contains(amountBetween, ";") {
		// 0;1
		amountRange := strings.Split(amountBetween, ";")
		query = query.Where("amount::numeric BETWEEN ?::numeric AND ?::numeric", amountRange[0], amountRange[1])
		countQuery = countQuery.Where("amount::numeric BETWEEN ?::numeric AND ?::numeric", amountRange[0], amountRange[1])

	}

	var countR int64

	errCount := countQuery.Find(&[]paymentModels.PaymentHistory{}).Count(&countR).Error
	if errCount != nil {
		log.Println("[GetPaymentHistory]Count Error:", errCount)
		return records
	}

	if limit > 0 {
		query.Limit(limit)
	}
	count := int(countR)
	pages := 1
	if count > limit {
		// fmt.Println("count / limit = ", count/limit, "count%limit = ", count%limit)
		pages = count / limit
		if count%limit > 0 {
			pages = pages + 1
		}
	}
	if page > pages {
		page = pages
	}
	if page > 1 {
		// fmt.Println("Offset = ", (page-1)*limit)
		query.Offset(((page - 1) * limit))
	}
	if err = query.Find(&paymentHistories).Error; err != nil {
		log.Println("[GetPaymentHistory] Query Error:", err)
		return
	}

	historiesJSON := make([]paymentModels.PaymentHistoryJSON, 0)

	for _, v := range paymentHistories {

		historiesJSON = append(historiesJSON, v.ToJSON())

	}

	records = paymentModels.PaginatedPaymentHistory{CurrentPage: page, Pages: pages, TotalRecords: count, Limit: limit, Records: historiesJSON}

	return records
}
