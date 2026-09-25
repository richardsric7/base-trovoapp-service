package payments

import (
	"log"
	"strconv"
	"strings"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	db "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetPaymentHistory(targetAddress string, gc *sharedconfig.GlobalConfig, c *gin.Context) (records paymentModels.PaginatedPaymentHistory) {
	var err error
	var paymentHistories []paymentModels.PaymentHistory
	records.Records = make([]paymentModels.PaymentHistoryJSON, 0)
	// DB, _ := db.OpenDb()
	// DBC, _ := db.OpenDb()
	DB := gc.DB
	DBC := gc.DB

	// if err != nil {
	// 	log.Fatalf("[main]Error opening DB %s", err)
	// 	return
	// }
	// DB := gc.DB
	// DBC := gc.DB
	var query *gorm.DB
	var countQuery *gorm.DB
	oD := "ASC"
	transactionType := strings.TrimSpace(c.Query("transactionType"))
	fromAddress := strings.TrimSpace(strings.ToUpper(c.Query("fromAddress")))
	toAddress := strings.TrimSpace(strings.ToUpper(c.Query("toAddress")))
	name := strings.TrimSpace(c.Query("name"))
	memo := strings.TrimSpace(c.Query("memo"))
	limitU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("limit", "25")), 10, 64)
	limit := int(limitU)
	pageU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("page", "1")), 10, 64)
	page := int(pageU)
	contractAddress := strings.TrimSpace(strings.ToUpper(c.Query("contractAddress")))
	// var contractAddressVal *string
	assetCode := strings.TrimSpace(strings.ToUpper(c.Query("assetCode")))

	transactionID := strings.ToLower(strings.TrimSpace(c.Query("transactionID")))
	amountBetween := strings.TrimSpace(c.Query("amount"))
	dateBetween := strings.TrimSpace(c.Query("dateBetween"))

	orderBy := strings.TrimSpace(c.DefaultQuery("orderby", "transaction_date"))
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
		query = query.Order("transaction_date DESC")
		countQuery = countQuery.Order("transaction_date DESC")
	}

	{
		query = query.Where("(from_address = ? OR to_address = ?)", targetAddress, targetAddress)
		countQuery = countQuery.Where("(from_address = ? OR to_address = ?)", targetAddress, targetAddress)

	}

	if len(fromAddress) == 42 {
		query = query.Where("from_address = ?", fromAddress)
		countQuery = countQuery.Where("from_address = ?", fromAddress)

	}

	if len(toAddress) == 42 {
		query = query.Where("to_address = ?", toAddress)
		countQuery = countQuery.Where("to_address = ?", toAddress)

	}
	if len(assetCode) > 0 {
		query = query.Where("asset_code = ?", strings.ToUpper(assetCode))
		countQuery = countQuery.Where("asset_code = ?", strings.ToUpper(assetCode))

	}
	if len(contractAddress) == 42 {
		query = query.Where("contract_address = ?", contractAddress)
		countQuery = countQuery.Where("contract_address = ?", contractAddress)

	}
	if len(name) > 2 {

		query = query.Where(`(lower("from") LIKE ? OR lower("to") LIKE ?)`, "%"+strings.ToLower(name)+"%", "%"+strings.ToLower(name)+"%")
		countQuery = countQuery.Where(`(lower("from") LIKE ? OR lower("to") LIKE ?)`, "%"+strings.ToLower(name)+"%", "%"+strings.ToLower(name)+"%")

	}
	if len(memo) > 2 {

		query = query.Where("lower(memo) LIKE ?", "%"+strings.ToLower(memo)+"%")
		countQuery = countQuery.Where("lower(memo) LIKE ?", "%"+strings.ToLower(memo)+"%")

	}
	if len(transactionType) > 3 {
		if strings.EqualFold(transactionType, "payment") {
			query = query.Where("transaction_type = ?", strings.ToUpper(transactionType))
			countQuery = countQuery.Where("transaction_type = ?", strings.ToUpper(transactionType))

		}
		if strings.EqualFold(transactionType, "swap") {
			query = query.Where("transaction_type LIKE ?", strings.ToUpper(transactionType)+"%")
			countQuery = countQuery.Where("transaction_type LIKE ?", strings.ToUpper(transactionType)+"%")

		}

	}
	if len(transactionID) > 40 {

		query = query.Where("transaction_id = ?", strings.TrimSpace(transactionID))
		countQuery = countQuery.Where("transaction_id = ?", strings.TrimSpace(transactionID))

	}

	if len(dateBetween) == 21 && strings.Contains(dateBetween, "|") {
		// 2020-01-01|2020-02-31 full range date
		dateRange := strings.Split(dateBetween, "|")
		query = query.Where("DATE(transaction_date) BETWEEN DATE(?) AND DATE(?)", dateRange[0], dateRange[1])
		countQuery = countQuery.Where("DATE(transaction_date) BETWEEN DATE(?) AND DATE(?)", dateRange[0], dateRange[1])

	}
	if len(amountBetween) > 2 && strings.Contains(amountBetween, "|") {
		// 0|1
		amountRange := strings.Split(amountBetween, "|")
		query = query.Where("CAST(amount AS REAL) BETWEEN CAST(? AS REAL) AND CAST(? AS REAL)", amountRange[0], amountRange[1])
		countQuery = countQuery.Where("CAST(amount AS REAL) BETWEEN CAST(? AS REAL) AND CAST(? AS REAL)", amountRange[0], amountRange[1])

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

func GetCryptoDepositHistory(targetAddress, currency string, gc *sharedconfig.GlobalConfig, c *gin.Context) (records userModels.PaginatedCryptoDepositHistory) {
	var err error
	var depositHistory []userModels.CryptoDeposit
	records.Records = make([]userModels.CryptoDepositJSON, 0)
	recs := make([]userModels.CryptoDepositJSON, 0)
	// DB, _ := db.OpenDb()
	// DBC, _ := db.OpenDb()

	DB := gc.DB
	DBC := gc.DB

	var query *gorm.DB
	var countQuery *gorm.DB
	oD := "ASC"
	s := strings.TrimSpace(c.Query("s"))

	limitU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("limit", "25")), 10, 64)
	limit := int(limitU)
	pageU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("page", "1")), 10, 64)
	page := int(pageU)

	amountBetween := strings.TrimSpace(c.Query("amount"))
	dateBetween := strings.TrimSpace(c.Query("dateBetween"))

	orderBy := strings.TrimSpace(c.DefaultQuery("orderby", "id"))
	orderDirection := c.DefaultQuery("order", "DESC")

	query = DB.Preload(clause.Associations)
	countQuery = DBC.Group("id")

	if len(orderDirection) > 0 && strings.ToLower(orderDirection) == "desc" {
		oD = "DESC"
	}
	if len(orderBy) > 0 {
		query = query.Order(orderBy + " " + oD)
		countQuery = countQuery.Order(orderBy + " " + oD)

	}

	{
		query = query.Where("trovo_wallet_address = ?", targetAddress)
		countQuery = countQuery.Where("trovo_wallet_address = ?", targetAddress)

	}

	if len(s) >= 2 {
		query = query.Where("(from_address = ? OR to_address = ? OR upper(network) = upper(?) OR tx_id = ?)", s, s, s, s)
		countQuery = countQuery.Where("(from_address = ? OR to_address = ? OR upper(network) = upper(?) OR tx_id = ?)", s, s, s, s)

	}

	if len(currency) > 0 {
		query = query.Where("upper(currency) = upper(?)", currency)
		countQuery = countQuery.Where("upper(currency) = upper(?)", currency)

	}

	if len(dateBetween) == 21 && strings.Contains(dateBetween, "|") {
		// 2020-01-01|2020-02-31 full range date
		dateRange := strings.Split(dateBetween, "|")
		query = query.Where("DATE(created_at) BETWEEN DATE(?) AND DATE(?)", dateRange[0], dateRange[1])
		countQuery = countQuery.Where("DATE(created_at) BETWEEN DATE(?) AND DATE(?)", dateRange[0], dateRange[1])

	}
	if len(amountBetween) > 2 && strings.Contains(amountBetween, "|") {
		// 0|1
		amountRange := strings.Split(amountBetween, "|")
		query = query.Where("CAST(amount AS REAL) BETWEEN CAST(? AS REAL) AND CAST(? AS REAL)", amountRange[0], amountRange[1])
		countQuery = countQuery.Where("CAST(amount AS REAL) BETWEEN CAST(? AS REAL) AND CAST(? AS REAL)", amountRange[0], amountRange[1])

	}

	var countR int64

	errCount := countQuery.Find(&[]userModels.CryptoDeposit{}).Count(&countR).Error
	if errCount != nil {
		log.Println("[GetCryptoDepositHistory]Count Error:", errCount)
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
	if err = query.Find(&depositHistory).Error; err != nil {
		log.Println("[GetCryptoDepositHistory] Query Error:", err)
		return
	}

	for _, deposit := range depositHistory {
		recs = append(recs, deposit.ToJSON(gc))
	}

	records = userModels.PaginatedCryptoDepositHistory{CurrentPage: page, Pages: pages, TotalRecords: count, Limit: limit, Records: recs}

	return records
}

func GetCryptoWithdrawalHistory(targetAddress, currency string, gc *sharedconfig.GlobalConfig, c *gin.Context) (records userModels.PaginatedCryptoWithdrawalHistory) {
	var err error
	var wdlHistory []userModels.WithdrawalRequest
	records.Records = make([]userModels.WithdrawalRequest, 0)
	DB, _ := db.OpenDb()
	DBC, _ := db.OpenDb()

	var query *gorm.DB
	var countQuery *gorm.DB
	oD := "ASC"
	// s := strings.TrimSpace(c.Query("s"))
	withdrawalAddress := strings.ToUpper(strings.TrimSpace(c.Query("withdrawalAddress")))
	withdrawalStatus := strings.ToUpper(strings.TrimSpace(c.Query("withdrawalStatus")))

	limitU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("limit", "25")), 10, 64)
	limit := int(limitU)
	pageU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("page", "1")), 10, 64)
	page := int(pageU)

	amountBetween := strings.TrimSpace(c.Query("amount"))
	dateBetween := strings.TrimSpace(c.Query("dateBetween"))

	orderBy := strings.TrimSpace(c.DefaultQuery("orderby", "created_at"))
	orderDirection := c.DefaultQuery("order", "DESC")

	query = DB.Preload(clause.Associations)
	countQuery = DBC.Group("id")

	if len(orderDirection) > 0 && strings.ToLower(orderDirection) == "desc" {
		oD = "DESC"
	}
	if len(orderBy) > 0 {
		query = query.Order(orderBy + " " + oD)
		countQuery = countQuery.Order(orderBy + " " + oD)

	}

	{
		query = query.Where("(wallet_address = ?)", targetAddress)
		countQuery = countQuery.Where("(wallet_address = ?)", targetAddress)

	}

	if len(withdrawalAddress) >= 2 {
		query = query.Where("withdrawal_address = ?", withdrawalAddress)
		countQuery = countQuery.Where("withdrawal_address = ?", withdrawalAddress)

	}

	if len(withdrawalStatus) >= 2 {
		query = query.Where("withdrawal_status = ?)", withdrawalStatus)
		countQuery = countQuery.Where("withdrawal_status = ?)", withdrawalStatus)

	}

	// if len(s) >= 2 {
	// 	query = query.Where("(upper(withdrawal_network) = upper(?) OR withdrawal_status = ?)", s, s)
	// 	countQuery = countQuery.Where("(upper(withdrawal_network) = upper(?) OR withdrawal_status = ?)", s, s)

	// }

	// if len(s) >= 2 {
	// 	query = query.Where("(upper(withdrawal_network) = upper(?) OR withdrawal_status = ?)", s, s)
	// 	countQuery = countQuery.Where("(upper(withdrawal_network) = upper(?) OR withdrawal_status = ?)", s, s)

	// }

	if len(currency) > 0 {
		query = query.Where("upper(currency) = upper(?)", currency)
		countQuery = countQuery.Where("upper(currency) = upper(?)", currency)

	}

	if len(dateBetween) == 21 && strings.Contains(dateBetween, "|") {
		// 2020-01-01|2020-02-31 full range date
		dateRange := strings.Split(dateBetween, "|")
		query = query.Where("DATE(created_at) BETWEEN DATE(?) AND DATE(?)", dateRange[0], dateRange[1])
		countQuery = countQuery.Where("DATE(created_at) BETWEEN DATE(?) AND DATE(?)", dateRange[0], dateRange[1])

	}
	if len(amountBetween) > 2 && strings.Contains(amountBetween, "|") {
		// 0|1
		amountRange := strings.Split(amountBetween, "|")
		query = query.Where("CAST(amount_submitted AS REAL) BETWEEN CAST(? AS REAL) AND CAST(? AS REAL)", amountRange[0], amountRange[1])
		countQuery = countQuery.Where("CAST(amount_submitted AS REAL) BETWEEN CAST(? AS REAL) AND CAST(? AS REAL)", amountRange[0], amountRange[1])

	}

	var countR int64

	errCount := countQuery.Find(&[]userModels.WithdrawalRequest{}).Count(&countR).Error
	if errCount != nil {
		log.Println("[GetCryptoWithdrawalHistory]Count Error:", errCount)
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
	if err = query.Find(&wdlHistory).Error; err != nil {
		log.Println("[GetCryptoWithdrawalHistory] Query Error:", err)
		return
	}

	records = userModels.PaginatedCryptoWithdrawalHistory{CurrentPage: page, Pages: pages, TotalRecords: count, Limit: limit, Records: wdlHistory}

	return records
}
