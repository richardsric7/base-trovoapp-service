package users

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	db "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	pns "trovo-wallet-api/internal/pns"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetApprovalList(publicKeysSharedWithUser []string, gc *sharedconfig.GlobalConfig, c *gin.Context) (records userModels.PaginatedAuths) {
	var err error
	var authList []userModels.PendingAuth
	records.Records = make([]userModels.AuthJSON, 0)
	DB, _ := db.OpenDb()
	DBC, _ := db.OpenDb()

	if len(publicKeysSharedWithUser) == 0 {
		return
	}
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
		j := v.ToJSON(gc)
		j.Transaction = ""
		authListJSON = append(authListJSON, j)

	}

	records = userModels.PaginatedAuths{CurrentPage: page, Pages: pages, TotalRecords: count, Limit: limit, Records: authListJSON}

	return records
}

func GetApprovalRequest(id string, gc *sharedconfig.GlobalConfig) (approvalRequest userModels.PendingAuth, err error) {
	p, err := userModels.ApprovalID(id).GetSubmittedTransaction(gc.DB)
	if err != nil {
		return approvalRequest, err
	}

	return p, nil
}
func GetApprovalRequestJSON(id string, gc *sharedconfig.GlobalConfig) (approvalRequest userModels.AuthJSON, err error) {
	p, err := userModels.ApprovalID(id).GetSubmittedTransaction(gc.DB)
	if err != nil {
		return approvalRequest, err
	}

	return p.ToJSON(gc), nil
}

func ApproveTransaction(signerUser *userModels.User, p *userModels.PendingAuth, approvalInfo *userModels.ApprovalPayload, gc *sharedconfig.GlobalConfig) (err error) {
	approvalInfo.NetworkPassPhrase = gc.BantuNetworkPassphrase
	var pts userModels.PendingTransactionSignature
	var revokedList, modifiedList, addedList []userModels.WalletPermission

	initiatorUser, e := userModels.Username(p.Initiator).GetSimpleUser(gc.DB)
	if e != nil {
		log.Println("[ApproveTransaction] error getting initiator user object for modify shared access")
		return &tErrors.ErrorTemporaryServerError{}
	}

	walletOwner, e := userModels.UserWalletID(p.WalletPublicKey).GetWalletOwner(gc.DB)
	if e != nil {
		log.Println("[ApproveTransaction] error getting wallet owner user object for modify shared access")
		return &tErrors.ErrorTemporaryServerError{}
	}

	wallet, e := userModels.UserWalletID(p.WalletPublicKey).GetWallet(gc.DB)
	if e != nil {
		log.Println("[ApproveTransaction] error getting wallet object for modify shared access")
		return &tErrors.ErrorTemporaryServerError{}
	}
	if p.TransactionType == "MODIFY SHARED ACCESS" {
		//unmarshall trx
		tbyte := []byte(*p.TransactionInfoStr)
		var ts userModels.ModifySharedAccessInfo
		e = json.Unmarshal(tbyte, &ts)
		if e != nil {
			log.Println("[ApproveTransaction] error decoding json for modified shared access")
			return &tErrors.ErrorTemporaryServerError{}
		}

		ts.Commit = 0
		revokedList, modifiedList, addedList, e = ModifySharedWalletAccess(&initiatorUser, &walletOwner, &wallet, &ts, gc)
		if e != nil {
			log.Println("[ApproveTransaction] error dry running for modify shared access")
			return &tErrors.ErrorTemporaryServerError{}
		}
		p.TransactionXdr = ts.Transaction
		// e =gc.DB.Save(p).Error
		// if e != nil {
		// 	log.Println("[ApproveTransaction] error saving fresh transaction for modify shared access")
		// 	return &tErrors.ErrorTemporaryServerError{}
		// }

	}
	if len(approvalInfo.TransactionSignature) == 0 {

		approvalInfo.Transaction = p.TransactionXdr

		return
	}
	//signature exists
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()

	pts = userModels.PendingTransactionSignature{
		ID:                       uuid.NewString(),
		PendingAuthID:            p.ID,
		Approver:                 signerUser.Username,
		ApproverSignerPublicKey:  signerUser.PrimarySigner,
		TransactionWithSignature: approvalInfo.TransactionSignature,
	}

	e = dbTX.Create(&pts).Error
	if e != nil {
		log.Println("[ApproveTransaction]error saving transaction signature:", e)
		return &tErrors.ErrorTemporaryServerError{}
	}
	p.ApprovalsGotten++

	if p.ApprovalsGotten < p.ApprovalsNeeded {
		e = dbTX.Save(p).Error
		if e != nil {
			log.Println("[ApproveTransaction]error saving approval state:", e)
			return &tErrors.ErrorTemporaryServerError{}
		}
		dbTX.Commit()
		return nil
	}

	//total approval needed is complete. Process the transaction and submit to network.
	p.TransactionStatus = "COMPLETED"

	//if transaction fails on blockchain, then reverse all changes.
	{
		e = dbTX.Save(p).Error
		if e != nil {
			log.Println("[ApproveTransaction]error saving approval state:", e)
			return &tErrors.ErrorTemporaryServerError{}
		}
		//process submission routine here
		tHash, err := network.SubmitApprovalsXdrWithSignatures(gc.BantuExpansionClient, p.ID, dbTX)
		if err != nil {
			return err
		}
		//set transaction ID

		p.TransactionID = &tHash

	}

	//blockchain succeeded
	e = dbTX.Save(p).Error
	if e != nil {
		log.Println("[ApproveTransaction]error saving approval state:", e)
		return &tErrors.ErrorTemporaryServerError{}
	}
	{ //sub
		//process post blockchcain transaction
		if p.TransactionType == "DISABLE SHARED ACCESS" {

			accessList := wallet.Permissions
			// set shared access enabled to 0
			// delete access list
			// wallet.SharedAccessEnabled = 0
			// wallet.NumberOfApprovalsNeeded = 0
			wallet.SharedAccessEnabled = 0
			wallet.NumberOfApprovalsNeeded = 0
			wallet.Permissions = nil

			e = dbTX.Save(&wallet).Error
			if e != nil {
				log.Println("[ApproveTransaction] error saving wallet state:", e.Error())
			}
			e = dbTX.Delete(&accessList).Error
			if e != nil {
				log.Println("[ApproveTransaction] error deleting access list:", e.Error())
			}

			dbTX.Commit()
			for _, v := range accessList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil && v.Permission != "VIEW-ONLY" {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					pns.SendFirebaseMessage(*u.PushNotificationToken, fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)

				}
			}

			return nil
		} else if p.TransactionType == "MODIFY SHARED ACCESS" {

			//delete revoked access
			e = dbTX.Delete(&revokedList).Error
			if e != nil {
				log.Println("[ApproveTransaction] error deleting revoked list:", e.Error())
			}
			//save modified access
			e = dbTX.Save(&modifiedList).Error
			if e != nil {
				log.Println("[ApproveTransaction] error saving modified list:", e.Error())
			}
			//create added access
			e = dbTX.Create(&addedList).Error
			if e != nil {
				log.Println("[ApproveTransaction] error creating added list:", e.Error())
			}

			dbTX.Commit()

			for _, v := range revokedList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					pns.SendFirebaseMessage(*u.PushNotificationToken, fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)

				}
			}

			for _, v := range modifiedList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					pns.SendFirebaseMessage(*u.PushNotificationToken, fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)

				}
			}

			for _, v := range addedList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					pns.SendFirebaseMessage(*u.PushNotificationToken, fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)

				}
			}

			return nil

		} else if p.TransactionType == "PAYMENT" {
			tbyte := []byte(*p.TransactionInfoStr)
			var paymentInfo paymentModels.PaymentInfo
			e = json.Unmarshal(tbyte, &paymentInfo)
			if e != nil {
				log.Println("[ApproveTransaction] error decoding json for modified shared access")
				return &tErrors.ErrorTemporaryServerError{}
			}
			accessList := wallet.GetPermissionList(gc.DB)
			// send push notifications
			assetCode := paymentInfo.AssetCode
			if assetCode == "" {
				assetCode = "XBN"
			}
			dataPayload := make(map[string]string)
			dataPayload["route"] = "pendingAuth"
			for _, v := range accessList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					pns.SendFirebaseMessage(*u.PushNotificationToken, fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)

				}
			}

		}
	} //end sub

	dbTX.Commit()
	return nil
}
