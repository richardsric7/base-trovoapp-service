package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	db "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/xdr"
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

func ApproveTransaction(signerUser *userModels.User, p *userModels.PendingAuth, approvalInfo *userModels.ApprovalPayload, retryCallbackChan chan userModels.RetryCallbacks, gc *sharedconfig.GlobalConfig) (err error) {
	approvalInfo.NetworkPassPhrase = gc.BantuNetworkPassphrase
	var pts userModels.PendingTransactionSignature
	var revokedList, modifiedList, addedList []userModels.WalletPermission
	var paymentInfo paymentModels.PaymentInfo
	var marketOffer userModels.MarketOffer
	// var swapInfo swapModels.SwapSendInfo
	// var pendingAssetClaim userModels.PendingAssetToClaim

	if p.TransactionStatus == "COMPLETED" {
		return &tErrors.ErrorCompletedRequest{ID: p.ID}
	}

	if p.TransactionStatus == "REJECTED" {
		return &tErrors.ErrorRejectedRequest{ID: p.ID}
	}

	initiatorUser, e := userModels.Username(p.Initiator).GetSimpleUser(gc.DB, gc)
	if e != nil {
		log.Println("[ApproveTransaction] error getting initiator user object for modify shared access")
		return &tErrors.ErrorTemporaryServerError{}
	}

	walletOwner, e := userModels.UserWalletID(p.WalletPublicKey).GetWalletOwner(gc.DB, gc)
	if e != nil {
		log.Println("[ApproveTransaction] error getting wallet owner user object for modify shared access")
		return &tErrors.ErrorTemporaryServerError{}
	}

	wallet, e := userModels.UserWalletID(p.WalletPublicKey).GetWallet(gc.DB)
	if e != nil {
		log.Println("[ApproveTransaction] error getting wallet object for modify shared access")
		return &tErrors.ErrorTemporaryServerError{}
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()

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
			return e
		}
		// p.TransactionXdr = ts.Transaction

	} else if p.TransactionType == "PAYMENT" {
		tbyte := []byte(*p.TransactionInfoStr)

		e = json.Unmarshal(tbyte, &paymentInfo)
		if e != nil {
			log.Println("[ApproveTransaction] error decoding json for modified shared access")
			return &tErrors.ErrorTemporaryServerError{}
		}
	} else if p.TransactionType == "MAKE MARKET OFFER" {
		tbyte := []byte(*p.TransactionInfoStr)

		e = json.Unmarshal(tbyte, &marketOffer)
		if e != nil {
			log.Println("[ApproveTransaction] error decoding json for modified shared access")
			return &tErrors.ErrorTemporaryServerError{}
		}
		e = dbTX.Create(&marketOffer).Error
		if e != nil {
			log.Printf("[ApproveTransaction]Error saving market offer: %+v\nError: %v\n", marketOffer, err)
			return &tErrors.ErrorTemporaryServerError{}
		}
	}

	if len(approvalInfo.TransactionSignature) == 0 {

		approvalInfo.Transaction = p.TransactionXdr

		return
	}
	//signature exists

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
	var txnResult horizon.Transaction
	{
		e = dbTX.Save(p).Error
		if e != nil {
			log.Println("[ApproveTransaction]error saving approval state:", e)
			return &tErrors.ErrorTemporaryServerError{}
		}
		//process submission routine here
		txnResult, err = network.SubmitApprovalsXdrWithSignaturesReturnsTrx(gc.BantuExpansionClient, p.ID, dbTX)
		if err != nil {
			return err
		}
		//set transaction ID

		p.TransactionID = &txnResult.Hash

	}

	//blockchain succeeded
	e = dbTX.Save(p).Error
	if e != nil {
		log.Println("[ApproveTransaction]error saving approval state:", e)
		// return &tErrors.ErrorTemporaryServerError{}
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
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil && v.Permission != "VIEW-ONLY" {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					u.SendPushMessage(fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc)
					u.InvalidateUserCache(gc)
				}
			}

			return nil
		} else if p.TransactionType == "MODIFY SHARED ACCESS" {

			if len(revokedList) > 0 {
				//delete revoked access
				e = dbTX.Delete(&revokedList).Error
				if e != nil {
					log.Println("[ApproveTransaction] error deleting revoked list:", e.Error())
				}
			}
			if len(modifiedList) > 0 {
				//save modified access
				e = dbTX.Save(&modifiedList).Error
				if e != nil {
					log.Println("[ApproveTransaction] error saving modified list:", e.Error())
				}
			}

			if len(revokedList) > 0 {
				//create added access
				e = dbTX.Create(&addedList).Error
				if e != nil {
					log.Println("[ApproveTransaction] error creating added list:", e.Error())
				}
			}

			dbTX.Commit()

			for _, v := range revokedList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					u.SendPushMessage(fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc)
					u.InvalidateUserCache(gc)
				}
			}

			for _, v := range modifiedList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					u.SendPushMessage(fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc)
					u.InvalidateUserCache(gc)
				}
			}

			for _, v := range addedList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					u.SendPushMessage(fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc)
					u.InvalidateUserCache(gc)
				}
			}

			return nil

		} else if p.TransactionType == "PAYMENT" {
			dbTX.Commit()
			accessList := wallet.GetPermissionList(gc.DB)
			// send push notifications
			assetCode := paymentInfo.AssetCode
			if assetCode == "" {
				assetCode = os.Getenv("NATIVE_ASSET_CODE")
			}
			dataPayload := make(map[string]string)
			dataPayload["route"] = "pendingAuth"
			for _, v := range accessList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					u.SendPushMessage(fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc)
					u.SendPushMessage("Trovo: Shared Wallet Debited!", fmt.Sprintf("Payment successfully sent %v %v from shared wallet with alias %v to %v", paymentInfo.Amount, assetCode, wallet.Alias, paymentInfo.Destination), "", dataPayload, gc)
					u.InvalidateUserCache(gc)
				}
			}
			{

				//start callback process here

				//TODO: make callback request if callback is available

				if d, ok := paymentInfo.CallbackURLS["orderPaymentCallbackUrl"]; ok && len(d) > 5 {
					log.Printf("[paymentNotification] found notification callbackUrl: [%v]\n\n", d)
					//make callback request
					// callbackResponse := new(map[string]interface{})
					type payload struct {
						Destination     string    `json:"destination"`
						Sender          string    `json:"sender"`
						Amount          string    `json:"amount"`
						AssetCode       string    `json:"assetCode"`
						AssetIssuer     string    `json:"assetIssuer"`
						TransactionID   string    `json:"transactionId"`
						TransactionMemo string    `json:"transactionMemo"`
						TransactionTime time.Time `json:"transactionTime"`
						DeviceID        string    `json:"deviceId"`
					}

					jsonPayload := payload{
						Destination:     paymentInfo.Destination,
						Sender:          wallet.Alias,
						Amount:          paymentInfo.Amount,
						AssetCode:       assetCode,
						AssetIssuer:     paymentInfo.AssetIssuer,
						TransactionID:   paymentInfo.TransactionID,
						TransactionMemo: paymentInfo.Memo,
						TransactionTime: time.Now(),
						DeviceID:        approvalInfo.DeviceID,
					}
					/////
					body, err := json.Marshal(jsonPayload)
					if err != nil {
						log.Printf("[paymentNotification] could not unmarshal callback message due to [%v]\n", err)

					}
					log.Printf("[paymentNotification] JSON STRING: [%v]\n", string(body))

					responseBody := bytes.NewBuffer(body)
					//Leverage Go's HTTP Post function to make request
					c := userModels.RetryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
					retryCallbackChan <- c
				}
			}

			{
				// send push notifications
				assetCode := paymentInfo.AssetCode
				if assetCode == "" {
					assetCode = os.Getenv("NATIVE_ASSET_CODE")
				}
				dataPayload := make(map[string]string)
				dataPayload["route"] = "basicTransactionHistory"
				if len(paymentInfo.Destination) < 31 {
					destWallet, e := userModels.WalletAlias(paymentInfo.Destination).GetWallet(gc.DB)
					if e == nil {
						dataPayload := make(map[string]string)
						dataPayload["none"] = ""
						if destWallet.SharedAccessEnabled == 1 {
							if destWallet.HasViewOnlyAccess(gc) {
								u, e := destWallet.GetWalletOwner(gc.DB, gc)
								if e == nil {
									if u.PushNotificationToken != nil {

										u.SendPushMessage("Trovo: Shared Wallet Credited!", fmt.Sprintf("You have received %v %v from %v to your shared wallet with alias %v", paymentInfo.Amount, assetCode, wallet.Alias, paymentInfo.Destination), "", dataPayload, gc)
										u.InvalidateUserCache(gc)
									}
								}

							}
							for _, v := range destWallet.Permissions {
								u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
								if e != nil {
									continue
								}
								if u.PushNotificationToken != nil {

									u.SendPushMessage("Trovo: Shared Wallet Credited!", fmt.Sprintf("You have received %v %v from %v to your shared wallet with alias %v", paymentInfo.Amount, assetCode, wallet.Alias, paymentInfo.Destination), "", dataPayload, gc)
									u.InvalidateUserCache(gc)
								}
							}
						} else {
							//shared access not enabled on destination wallet
							u, e := destWallet.GetWalletOwner(gc.DB, gc)
							if e == nil {
								if u.PushNotificationToken != nil {

									u.SendPushMessage("Trovo: Wallet Credited!", fmt.Sprintf("You have received %v %v from %v to your wallet with alias %v", paymentInfo.Amount, assetCode, wallet.Alias, paymentInfo.Destination), "", dataPayload, gc)
									u.InvalidateUserCache(gc)
								}
							}
						}

					}

				}

			}

		} else if p.TransactionType == "MAKE MARKET OFFER" {

			marketOffer.TransactionID = &txnResult.Hash
			//get and set the offerID
			{
				var re xdr.TransactionResult
				e := xdr.SafeUnmarshalBase64(txnResult.ResultXdr, &re)
				if e != nil {
					fmt.Println(e)
				}
				log.Println(re)
				or, _ := re.OperationResults()

				for _, r := range or {

					ms, ok := r.Tr.GetManageSellOfferResult()
					if !ok {
						continue
					}
					offerID := fmt.Sprintf("%v", ms.Success.Offer.Offer.OfferId)

					marketOffer.BlockchainOfferID = &offerID
				}
			}
			e = dbTX.Save(&marketOffer).Error
			if e != nil {
				log.Printf("[ApproveTransaction]Error saving transactionID on market offer: %+v\nError: %v\n", marketOffer, err)
				// return &tErrors.ErrorTemporaryServerError{}
			}
			dbTX.Commit()
			accessList := wallet.Permissions
			for _, v := range accessList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}
				// if u.PushNotificationToken != nil && v.Permission != "VIEW-ONLY" {
				dataPayload := make(map[string]string)
				dataPayload["none"] = ""
				u.SendPushMessage(fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc)
				u.InvalidateUserCache(gc)
				// }
			}

			return nil

		} else {

			accessList := wallet.GetPermissionList(gc.DB)

			dataPayload := make(map[string]string)
			dataPayload["route"] = "pendingAuth"
			for _, v := range accessList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}
				if u.PushNotificationToken != nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					u.SendPushMessage(fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias), fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description), "", dataPayload, gc)
					u.InvalidateUserCache(gc)
				}
			}

		}
	} //end sub

	dbTX.Commit()
	return nil
}

func RejectTransaction(signerUser *userModels.User, p *userModels.PendingAuth, rejectionInfo *userModels.RejectPayload, gc *sharedconfig.GlobalConfig) (err error) {

	if p.TransactionStatus == "COMPLETED" {
		return &tErrors.ErrorCompletedRequest{ID: p.ID}
	}
	if p.TransactionStatus == "REJECTED" {
		return &tErrors.ErrorRejectedRequest{ID: p.ID}
	}
	rejectionInfo.RejectionReason = strings.TrimSpace(rejectionInfo.RejectionReason)
	if len(rejectionInfo.RejectionReason) < 5 {
		return &tErrors.CustomError{
			Param:      "rejectionReason",
			Err:        "error-invalid-rejection-reason",
			ErrMessage: "Rejection reason is not valid enough. Must be elaborate and must contain at least 5 chracters",
		}
	}
	p.TransactionStatus = "REJECTED"
	p.RejectedBy = &signerUser.Username
	p.ReasonForRejection = &rejectionInfo.RejectionReason

	//signature exists
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	e := dbTX.Save(p).Error
	if e != nil {
		log.Println("[ApproveTransaction]error saving approval state:", e)
		return &tErrors.ErrorTemporaryServerError{}
	}
	dbTX.Commit()
	return nil

}

func CheckPendingSharedAccessApproval(walletPublicKey string, db *gorm.DB) (exists bool) {

	var pendingApproval userModels.PendingAuth
	e := db.Where("wallet_public_key = ? AND (transaction_type = ? OR transaction_type = ?) AND transaction_status = ?", walletPublicKey, "DISABLE SHARED ACCESS", "MODIFY SHARED ACCESS", "PENDING").First(&pendingApproval).Error
	return e == nil

}
