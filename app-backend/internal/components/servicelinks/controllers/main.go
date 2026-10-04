package servicelinks

import (
	"bytes"
	"os"
	"time"

	servicelinkModels "trovo-wallet-api/internal/components/servicelinks/models"
	"trovo-wallet-api/internal/sharedconfig"

	"log"
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

	// callbacks are delivered (and retried) through gc.SendCallback;
	// nothing is sent on this channel any more, so it holds no buffer
	callBackRetryChan := make(chan retryCallbacks)

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
	router.POST("/v1/servicelinks/login/request/:targetUser", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "servicelinks-login-request", 30, time.Minute), postServicelinksLoginRequestTargetUserHandler(gc))

	//user login approval url; uses signature algorithm bcos it is only called by trovoApp.
	router.POST("/v1/users/servicelinks/login/approval/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersServicelinksLoginApprovalTargetUserHandler(callBackRetryChan, gc))

	//service login verify url
	router.GET("/v1/servicelinks/login/verify/:ownerUsername/:targetUser/:loginID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "servicelinks-login-verify", 60, time.Minute), getServicelinksLoginVerifyOwnerUsernameTargetUserLoginIDHandler(gc))

	//service login refresh token url
	router.POST("/v1/servicelinks/token/refresh", postServicelinksTokenRefreshHandler(gc))

	//service token verify
	router.POST("/v1/servicelinks/token/verify", postServicelinksTokenVerifyHandler(gc))

	//service token verify
	router.DELETE("/v1/servicelinks/token", deleteServicelinksTokenHandler(gc))

	//service authorization request
	router.POST("/v1/servicelinks/authorize/request/:targetUser", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "servicelinks-authorize-request", 30, time.Minute), postServicelinksAuthorizeRequestTargetUserHandler(gc))

	//service authorization tokenizedAsset
	router.POST("/v1/servicelinks/authorize/tokenized-asset", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "servicelinks-authorize-tokenized-asset", 30, time.Minute), postServicelinksAuthorizeTokenizedAssetHandler(gc))

	//service event link request
	router.POST("/v1/servicelinks/events/request", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "servicelinks-events-request", 30, time.Minute), postServicelinksEventsRequestHandler(gc))

	//user authorization approval url
	router.POST("/v1/users/servicelinks/authorize/approval/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersServicelinksAuthorizeApprovalTargetUserHandler(callBackRetryChan, gc))

	//user events approval url
	router.POST("/v1/users/servicelinks/events/approval/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersServicelinksEventsApprovalTargetUserHandler(callBackRetryChan, gc))

	//service authorization verify url
	router.GET("/v1/servicelinks/authorize/verify/:ownerUsername/:targetUser/:authId", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "servicelinks-authorize-verify", 60, time.Minute), getServicelinksAuthorizeVerifyOwnerUsernameTargetUserAuthIdHandler(gc))
	//app authorization verify url
	router.GET("/v1/servicelinks/app/authorize/verify/:ownerUsername/:targetUser/:authId", middleware.AuthenticationMiddlewareUsingTimestamp(), getServicelinksAppAuthorizeVerifyOwnerUsernameTargetUserAuthIdHandler(gc))

	//service payment request
	router.GET("/v1/servicelinks/payment/request/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), getServicelinksPaymentRequestTargetUserHandler(gc))

	//api service payment request
	router.GET("/v1/trovo-api/payment/request/:targetUser", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-payment-request", 60, time.Minute), getTrovoApiPaymentRequestTargetUserHandler(gc))

	//service tokenized asset request
	router.GET("/v1/servicelinks/tokenized-asset/:assetCode", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "servicelinks-tokenized-asset", 60, time.Minute), getServicelinksTokenizedAssetAssetCodeHandler(gc))

	//SERVICELINK USER INFO request
	router.GET("/v1/servicelinks/:ownerUsername/:targetUser/userinfo", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "servicelinks-userinfo", 60, time.Minute), getServicelinksOwnerUsernameTargetUserUserinfoHandler(gc))

	//SERVICE push notification request
	router.POST("/v1/servicelinks/:ownerUsername/:targetUser/push", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "servicelinks-push", 30, time.Minute), postServicelinksOwnerUsernameTargetUserPushHandler(gc))

	//register user from service link
	router.POST("/v1/trovo-api/users/onboard", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-onboard", 20, time.Minute), postTrovoApiUsersOnboardHandler(gc))

	//update user kyc from service link
	router.POST("/v1/trovo-api/users/update-kyc", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-update-kyc", 20, time.Minute), postTrovoApiUsersUpdateKycHandler(gc))

	//mint token from service link
	router.POST("/v1/trovo-api/tokens/mint", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-mint", 10, time.Minute), postTrovoApiTokensMintHandler(gc))

	//get wallet balance from service link
	router.GET("/v1/trovo-api/users/balance/:walletAddress", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-balance", 60, time.Minute), getTrovoApiUsersBalanceWalletAddressHandler(gc))

	//get wallet payment history from service link
	router.GET("/v1/trovo-api/users/payment-history/:walletAddress", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-payment-history", 120, time.Minute), getTrovoApiUsersPaymentHistoryWalletAddressHandler(gc))

	//send payment from service link
	router.POST("/v1/trovo-api/users/payment", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-payment", 20, time.Minute), postTrovoApiUsersPaymentHandler(callBackRetryChan, gc))

	//create new subwallet from service link
	router.POST("/v1/trovo-api/users/subwallet", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-subwallet", 20, time.Minute), postTrovoApiUsersSubwalletHandler(gc))

	//Get tokenization parameters from service link
	router.GET("/v1/trovo-api/assets/parameters", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-parameters", 60, time.Minute), getTrovoApiAssetsParametersHandler(gc))

	//Get tokenization bank list from service link
	router.GET("/v1/trovo-api/assets/bank-list/:countryCode", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-bank-list", 60, time.Minute), getTrovoApiAssetsBankListCountryCodeHandler(gc))

	//Get tokenization list for Admin from service link
	router.GET("/v1/trovo-api/assets/admin/list", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-admin-list", 60, time.Minute), getTrovoApiAssetsAdminListHandler(gc))

	//Get tokenization list for market from service link
	router.GET("/v1/trovo-api/assets/marketplace/list", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-marketplace-list", 60, time.Minute), getTrovoApiAssetsMarketplaceListHandler(gc))

	//Apply for tokenization from service link
	router.POST("/v1/trovo-api/assets/apply", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-apply", 20, time.Minute), postTrovoApiAssetsApplyHandler(gc))

	//Upload logo for tokenization from service link
	router.PUT("/v1/trovo-api/assets/logo", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-logo", 20, time.Minute), putTrovoApiAssetsLogoHandler(gc))

	//Upload documents for tokenization from service link
	router.PUT("/v1/trovo-api/assets/documents", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-documents", 20, time.Minute), putTrovoApiAssetsDocumentsHandler(gc))

	// Store private stakeholder portal documents for Trovo Manager.
	router.POST("/v1/trovo-api/stakeholder-documents", middleware.AuthenticationMiddlewareUsingAPIKey(gc), requireActiveServiceLink(gc), middleware.RateLimitMiddleware(gc, "trovo-api-stakeholder-documents-post", 20, time.Minute), postStakeholderDocumentHandler(gc))
	router.GET("/v1/trovo-api/stakeholder-documents/:objectID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), requireActiveServiceLink(gc), middleware.RateLimitMiddleware(gc, "trovo-api-stakeholder-documents-get", 60, time.Minute), getStakeholderDocumentHandler(gc))
	router.DELETE("/v1/trovo-api/stakeholder-documents/:objectID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), requireActiveServiceLink(gc), middleware.RateLimitMiddleware(gc, "trovo-api-stakeholder-documents-delete", 20, time.Minute), deleteStakeholderDocumentHandler(gc))

	//Upload fee payment documents for tokenization from service link
	router.PUT("/v1/trovo-api/assets/fees/document", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-fees-document", 20, time.Minute), putTrovoApiAssetsFeesDocumentHandler(gc))

	//confirm fee payment for tokenization from service link
	router.POST("/v1/trovo-api/assets/fees/confirm/:tokenizationID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-fees-confirm", 20, time.Minute), postTrovoApiAssetsFeesConfirmTokenizationIDHandler(gc))

	//DELETE tokenization from service link
	router.DELETE("/v1/trovo-api/assets/:tokenizationID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-delete", 20, time.Minute), deleteTrovoApiAssetsTokenizationIDHandler(gc))

	//DELETE tokenization document from service link
	router.DELETE("/v1/trovo-api/assets/documents/:documentID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-documents-delete", 20, time.Minute), deleteTrovoApiAssetsDocumentsDocumentIDHandler(gc))

	//DELETE tokenization fee document from service link
	router.DELETE("/v1/trovo-api/assets/fees/documents/:documentID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-fees-documents-delete", 20, time.Minute), deleteTrovoApiAssetsFeesDocumentsDocumentIDHandler(gc))

	//Confirm Application for tokenization from service link
	router.POST("/v1/trovo-api/assets/confirm-application/:tokenizationID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-confirm-application", 20, time.Minute), postTrovoApiAssetsConfirmApplicationTokenizationIDHandler(gc))

	//Purchase primary sales tokenized asset from service link
	router.POST("/v1/trovo-api/assets/marketplace/primary", middleware.AuthenticationMiddlewareUsingAPIKey(gc), middleware.RateLimitMiddleware(gc, "trovo-api-assets-marketplace-primary", 10, time.Minute), postTrovoApiAssetsMarketplacePrimaryHandler(gc))

}
