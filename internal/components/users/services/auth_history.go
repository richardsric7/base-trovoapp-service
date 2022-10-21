package users

import (
	"log"
	"strconv"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	db "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAuthorizationList(publicKeysSharedWithUser []string, gc *sharedconfig.GlobalConfig, c *gin.Context) (records userModels.PaginatedAuths) {
	var err error
	var authList []userModels.PendingAuth
	records.Records = make([]userModels.AuthJSON, 0)
	DB, _ := db.OpenDb()
	DBC, _ := db.OpenDb()

	var query *gorm.DB
	var countQuery *gorm.DB
	oD := "ASC"
	transactionType := strings.ToUpper(strings.TrimSpace(c.Query("transactionType")))
	transactionStatus := strings.ToUpper(strings.TrimSpace(c.Query("transactionStatus")))
	initiator := strings.TrimSpace(strings.ToLower(c.Query("initiator")))
	description := strings.TrimSpace(strings.ToLower(c.Query("description")))
	walletPublicKey := strings.TrimSpace(strings.ToUpper(c.Query("walletPublicKey")))
	walletAlias := strings.ToLower(strings.TrimSpace(c.Query("walletAlias")))
	limitU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("limit", "25")), 10, 64)
	limit := int(limitU)
	pageU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("page", "1")), 10, 64)
	page := int(pageU)

	// var walletOwnerPublicKey string
	if len(walletAlias) > 2 {
		wo, e := userModels.WalletAlias(walletAlias).GetWallet(gc.DB)
		if e == nil {
			walletPublicKey = wo.ID
		}
	}

	transactionID := strings.ToLower(strings.TrimSpace(c.Query("transactionID")))
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

	} else {
		query = query.Order("created_at DESC")
		countQuery = countQuery.Order("created_at DESC")
	}

	{
		query = query.Where("(wallet_public_key IN (?)", publicKeysSharedWithUser)
		countQuery = countQuery.Where("(wallet_public_key IN (?)", publicKeysSharedWithUser)

	}

	if len(walletPublicKey) == 56 {
		query = query.Where("(wallet_public_key = ?", walletPublicKey)
		countQuery = countQuery.Where("(wallet_public_key = ?", walletPublicKey)

	}

	if len(initiator) > 2 {
		query = query.Where("initiator = ?", initiator)
		countQuery = countQuery.Where("initiator = ?", initiator)

	}
	if len(transactionStatus) > 0 {
		query = query.Where("transaction_status = ?", transactionStatus)
		countQuery = countQuery.Where("transaction_status = ?", transactionStatus)

	}

	if len(description) > 2 {

		query = query.Where("lower(description) LIKE ?", "%"+strings.ToLower(description)+"%")
		countQuery = countQuery.Where("lower(description) LIKE ?", "%"+strings.ToLower(description)+"%")

	}
	if len(transactionType) > 3 {

		query = query.Where("transaction_type = ?", transactionType)
		countQuery = countQuery.Where("transaction_type = ?", transactionType)

	}
	if len(transactionID) > 40 {

		query = query.Where("transaction_id = ?", strings.TrimSpace(transactionID))
		countQuery = countQuery.Where("transaction_id = ?", strings.TrimSpace(transactionID))

	}

	if len(dateBetween) == 21 && strings.Contains(dateBetween, "|") {
		// 2020-01-01|2020-02-31 full range date
		dateRange := strings.Split(dateBetween, "|")
		query = query.Where("created_at::date BETWEEN ?::date AND ?::date", dateRange[0], dateRange[1])
		countQuery = countQuery.Where("created_at::date BETWEEN ?::date AND ?::date", dateRange[0], dateRange[1])

	}

	var countR int64

	errCount := countQuery.Find(&[]userModels.PendingAuth{}).Count(&countR).Error
	if errCount != nil {
		log.Println("[GetAuthorizationList]Count Error:", errCount)
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
	if err = query.Find(&authList).Error; err != nil {
		log.Println("[GetAuthorizationList] Query Error:", err)
		return
	}

	authListJSON := make([]userModels.AuthJSON, 0)

	for _, v := range authList {

		authListJSON = append(authListJSON, v.ToJSON(gc))

	}

	records = userModels.PaginatedAuths{CurrentPage: page, Pages: pages, TotalRecords: count, Limit: limit, Records: authListJSON}

	return records
}
