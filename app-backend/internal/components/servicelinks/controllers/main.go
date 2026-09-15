package servicelinks

import (
	"bytes"
	"os"
	"time"

	servicelinkModels "trovo-wallet-api/internal/components/servicelinks/models"
	"trovo-wallet-api/internal/sharedconfig"

	"log"
	"net/http"
	"trovo-wallet-api/internal/middleware"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// Init initializes /v1/services endpoint

// retryCallbacks stores failed callbacks
type retryCallbacks struct {
	Req         *bytes.Buffer
	CallbackURL string
	Count       int
}

func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {

	callBackRetryChan := make(chan retryCallbacks, 20000)
	go func(c chan retryCallbacks) {
		log.Println("#####@started Routine to retry failed auth/events/login callbacks....")
		//loop
		for {
			callbackObj := <-c
			if callbackObj.Count > 100 {
				//skip 100 retries
				continue
			}
			_, err := http.Post(callbackObj.CallbackURL, "application/json", callbackObj.Req)

			if err != nil {
				//send back into channel to retry later
				log.Printf("Callback retry failed: [%+v]\n", callbackObj)
				if callbackObj.Count <= 99 {
					callbackObj.Count++
					c <- callbackObj
				}

			}
			//wait 1 second
			time.Sleep(1 * time.Second)
		}

	}(callBackRetryChan)

	{
		//auto expire login sessions that where that are not within valid time.

		go func() {
			log.Println("@@@@@Started routine to Auto remove <SERVICELINK> LoginSessions")
			period := time.Duration(3)
			if os.Getenv("SERVICE_LINK_LOGIN_REQUEST_VALIDITY") != "" {
				m, e := decimal.NewFromString(os.Getenv("SERVICE_LINK_LOGIN_REQUEST_VALIDITY"))
				if e == nil {
					if m.IsPositive() {
						period = time.Duration(m.IntPart())
					}
				}
			}
			for {
				err := gc.DB.Where("created_at < ?", time.Now().Add(-1*period*time.Minute)).Delete(servicelinkModels.ServiceLinkLoginSession{}).Error
				if err != nil {
					log.Printf("[Expire Login Sessions Routine]unable to delete expired login sessions due to error [%v]\n", err)
				}
				time.Sleep(30 * time.Second)
			}
		}()
	}

	{
		//auto expire authorizations that are not within valid time.

		go func() {
			log.Println("@@@@@Started routine to Auto remove <service> Authorizations")

			for {
				err := gc.DB.Where("expires_at < ?", time.Now()).Delete(servicelinkModels.ServiceLinkAuthorization{}).Error
				if err != nil {
					log.Printf("[Expire Authorizations Routine]unable to delete expired authorizations requests due to error [%v]\n", err)
				}
				time.Sleep(30 * time.Second)
			}
		}()
	}

	{
		//auto expire events that are not within valid time.

		go func() {
			log.Println("@@@@@Started routine to Auto remove <service> events")

			for {
				err := gc.DB.Where("expires_at < ?", time.Now()).Delete(servicelinkModels.ServiceLinkEvent{}).Error
				if err != nil {
					log.Printf("[Expire Events Routine]unable to delete expired events requests due to error [%v]\n", err)
				}
				time.Sleep(30 * time.Second)
			}
		}()
	}

	//service login request
	router.POST("/v1/servicelinks/login/request/:targetUser", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postServicelinksLoginRequestTargetUserHandler(gc))

	//user login approval url; uses signature algorithm bcos it is only called by trovoApp.
	router.POST("/v1/users/servicelinks/login/approval/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersServicelinksLoginApprovalTargetUserHandler(callBackRetryChan, gc))

	//service login verify url
	router.GET("/v1/servicelinks/login/verify/:ownerUsername/:targetUser/:loginID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getServicelinksLoginVerifyOwnerUsernameTargetUserLoginIDHandler(gc))

	//service login refresh token url
	router.POST("/v1/servicelinks/token/refresh", postServicelinksTokenRefreshHandler(gc))

	//service token verify
	router.POST("/v1/servicelinks/token/verify", postServicelinksTokenVerifyHandler(gc))

	//service token verify
	router.DELETE("/v1/servicelinks/token", deleteServicelinksTokenHandler(gc))

	//service authorization request
	router.POST("/v1/servicelinks/authorize/request/:targetUser", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postServicelinksAuthorizeRequestTargetUserHandler(gc))

	//service authorization tokenizedAsset
	router.POST("/v1/servicelinks/authorize/tokenized-asset", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postServicelinksAuthorizeTokenizedAssetHandler(gc))

	//service event link request
	router.POST("/v1/servicelinks/events/request", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postServicelinksEventsRequestHandler(gc))

	//user authorization approval url
	router.POST("/v1/users/servicelinks/authorize/approval/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersServicelinksAuthorizeApprovalTargetUserHandler(callBackRetryChan, gc))

	//user events approval url
	router.POST("/v1/users/servicelinks/events/approval/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersServicelinksEventsApprovalTargetUserHandler(callBackRetryChan, gc))

	//service authorization verify url
	router.GET("/v1/servicelinks/authorize/verify/:ownerUsername/:targetUser/:authId", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getServicelinksAuthorizeVerifyOwnerUsernameTargetUserAuthIdHandler(gc))
	//app authorization verify url
	router.GET("/v1/servicelinks/app/authorize/verify/:ownerUsername/:targetUser/:authId", middleware.AuthenticationMiddlewareUsingTimestamp(), getServicelinksAppAuthorizeVerifyOwnerUsernameTargetUserAuthIdHandler(gc))

	//service payment request
	router.GET("/v1/servicelinks/payment/request/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), getServicelinksPaymentRequestTargetUserHandler(gc))

	//api service payment request
	router.GET("/v1/trovo-api/payment/request/:targetUser", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getTrovoApiPaymentRequestTargetUserHandler(gc))

	//service tokenized asset request
	router.GET("/v1/servicelinks/tokenized-asset/:assetCode", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getServicelinksTokenizedAssetAssetCodeHandler(gc))

	//SERVICELINK USER INFO request
	router.GET("/v1/servicelinks/:ownerUsername/:targetUser/userinfo", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getServicelinksOwnerUsernameTargetUserUserinfoHandler(gc))

	//SERVICE push notification request
	router.POST("/v1/servicelinks/:ownerUsername/:targetUser/push", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postServicelinksOwnerUsernameTargetUserPushHandler(gc))

	//register user from service link
	router.POST("/v1/trovo-api/users/onboard", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postTrovoApiUsersOnboardHandler(gc))

	//update user kyc from service link
	router.POST("/v1/trovo-api/users/update-kyc", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postTrovoApiUsersUpdateKycHandler(gc))

	//mint token from service link
	router.POST("/v1/trovo-api/tokens/mint", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postTrovoApiTokensMintHandler(gc))

	//get wallet balance from service link
	router.GET("/v1/trovo-api/users/balance/:walletPublicKey", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getTrovoApiUsersBalanceWalletPublicKeyHandler(gc))

	//get wallet payment history from service link
	router.GET("/v1/trovo-api/users/payment-history/:walletPublicKey", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getTrovoApiUsersPaymentHistoryWalletPublicKeyHandler(gc))

	//send payment from service link
	router.POST("/v1/trovo-api/users/payment", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postTrovoApiUsersPaymentHandler(callBackRetryChan, gc))

	//create new subwallet from service link
	router.POST("/v1/trovo-api/users/subwallet", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postTrovoApiUsersSubwalletHandler(gc))

	//Get tokenization parameters from service link
	router.GET("/v1/trovo-api/assets/parameters", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getTrovoApiAssetsParametersHandler(gc))

	//Get tokenization bank list from service link
	router.GET("/v1/trovo-api/assets/bank-list/:countryCode", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getTrovoApiAssetsBankListCountryCodeHandler(gc))

	//Get tokenization list for Admin from service link
	router.GET("/v1/trovo-api/assets/admin/list", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getTrovoApiAssetsAdminListHandler(gc))

	//Get tokenization list for market from service link
	router.GET("/v1/trovo-api/assets/marketplace/list", middleware.AuthenticationMiddlewareUsingAPIKey(gc), getTrovoApiAssetsMarketplaceListHandler(gc))

	//Apply for tokenization from service link
	router.POST("/v1/trovo-api/assets/apply", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postTrovoApiAssetsApplyHandler(gc))

	//Upload logo for tokenization from service link
	router.PUT("/v1/trovo-api/assets/logo", middleware.AuthenticationMiddlewareUsingAPIKey(gc), putTrovoApiAssetsLogoHandler(gc))

	//Upload documents for tokenization from service link
	router.PUT("/v1/trovo-api/assets/documents", middleware.AuthenticationMiddlewareUsingAPIKey(gc), putTrovoApiAssetsDocumentsHandler(gc))

	// Store private stakeholder portal documents for Trovo Manager.
	router.POST("/v1/trovo-api/stakeholder-documents", middleware.AuthenticationMiddlewareUsingAPIKey(gc), requireActiveServiceLink(gc), postStakeholderDocumentHandler(gc))
	router.GET("/v1/trovo-api/stakeholder-documents/:objectID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), requireActiveServiceLink(gc), getStakeholderDocumentHandler(gc))
	router.DELETE("/v1/trovo-api/stakeholder-documents/:objectID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), requireActiveServiceLink(gc), deleteStakeholderDocumentHandler(gc))

	//Upload fee payment documents for tokenization from service link
	router.PUT("/v1/trovo-api/assets/fees/document", middleware.AuthenticationMiddlewareUsingAPIKey(gc), putTrovoApiAssetsFeesDocumentHandler(gc))

	//confirm fee payment for tokenization from service link
	router.POST("/v1/trovo-api/assets/fees/confirm/:tokenizationID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postTrovoApiAssetsFeesConfirmTokenizationIDHandler(gc))

	//DELETE tokenization from service link
	router.DELETE("/v1/trovo-api/assets/:tokenizationID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), deleteTrovoApiAssetsTokenizationIDHandler(gc))

	//DELETE tokenization document from service link
	router.DELETE("/v1/trovo-api/assets/documents/:documentID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), deleteTrovoApiAssetsDocumentsDocumentIDHandler(gc))

	//DELETE tokenization fee document from service link
	router.DELETE("/v1/trovo-api/assets/fees/documents/:documentID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), deleteTrovoApiAssetsFeesDocumentsDocumentIDHandler(gc))

	//Confirm Application for tokenization from service link
	router.POST("/v1/trovo-api/assets/confirm-application/:tokenizationID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postTrovoApiAssetsConfirmApplicationTokenizationIDHandler(gc))

	//Purchase primary sales tokenized asset from service link
	router.POST("/v1/trovo-api/assets/marketplace/primary", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postTrovoApiAssetsMarketplacePrimaryHandler(gc))

}
