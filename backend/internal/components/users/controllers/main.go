package users

import (
	"os"
	userModels "trovo-wallet-api/internal/components/users/models"

	"log"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// Init initializes /v1/users endpoint
func Init(router *gin.Engine, callBackRetryChan chan userModels.RetryCallbacks, gc *sharedconfig.GlobalConfig) {
	//websocket stream
	router.GET("/v1/stream/ws/:targetUser", getStreamWsTargetUserHandler(callBackRetryChan, gc))

	//websocket stream
	router.GET("/v1/stream/orderbook", getStreamOrderbookHandler(callBackRetryChan, gc))
	//websocket stream
	router.GET("/v1/stream/tradechart", getStreamTradechartHandler(callBackRetryChan, gc))
	//shortlink
	router.GET("/v1/shortlinks/:linkID", getShortlinksLinkIDHandler(callBackRetryChan, gc))

	router.GET("/v1/users/payments/:targetPublicKeyForHistory", middleware.AuthenticationMiddlewareUsingTimestamp(), getUsersPaymentsTargetPublicKeyForHistoryHandler(callBackRetryChan, gc))

	router.GET("/v1/curated-assets/users", middleware.AuthenticationMiddlewareUsingTimestamp(), getCuratedAssetsUsersHandler(callBackRetryChan, gc))

	router.GET("/v1/users/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), getUsersTargetUserHandler(callBackRetryChan, gc))

	router.POST("/v1/users", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersHandler(callBackRetryChan, gc))

	router.DELETE("/v1/users", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteUsersHandler(callBackRetryChan, gc))

	router.POST("/v1/webhook/sumsub/kyc/individual", postWebhookSumsubKycIndividualHandler(callBackRetryChan, gc))

	router.POST("/v1/users/kyc/sumsub/initiate/:levelName", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersKycSumsubInitiateLevelNameHandler(callBackRetryChan, gc))

	router.GET("/v1/users/kyc/sumsub/progress", middleware.AuthenticationMiddlewareUsingTimestamp(), getUsersKycSumsubProgressHandler(callBackRetryChan, gc))

	router.GET("/v1/users/kyc/doja/progress", middleware.AuthenticationMiddlewareUsingTimestamp(), getUsersKycDojaProgressHandler(callBackRetryChan, gc))

	router.GET("/v1/users/activate/fiat", middleware.AuthenticationMiddlewareUsingTimestamp(), getUsersActivateFiatHandler(callBackRetryChan, gc))

	router.GET("/v1/users/fiat/payments", middleware.AuthenticationMiddlewareUsingTimestamp(), getUsersFiatPaymentsHandler(callBackRetryChan, gc))

	router.POST("/v1/users/fiat/flutterwave", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersFiatFlutterwaveHandler(callBackRetryChan, gc))

	router.POST("/v1/users/kyc/sumsub/complete/:levelName", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersKycSumsubCompleteLevelNameHandler(callBackRetryChan, gc))

	router.GET("/v1/users/kyc/sumsub/configs", middleware.AuthenticationMiddlewareUsingTimestamp(), getUsersKycSumsubConfigsHandler(callBackRetryChan, gc))

	router.GET("/v1/users/kyc/doja/configs", middleware.AuthenticationMiddlewareUsingTimestamp(), getUsersKycDojaConfigsHandler(callBackRetryChan, gc))

	router.POST("/v1/users/subwallet", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersSubwalletHandler(callBackRetryChan, gc))

	router.PUT("/v1/users/upload-picture", middleware.AuthenticationMiddlewareUsingTimestamp(), putUsersUploadPictureHandler(callBackRetryChan, gc))

	// get config
	var config userModels.StablerailConfig
	gc.DB.First(&config)

	if len(config.ApiKey) == 1 && config.EnableStablerail == 1 {
		//STABLERAIL ENDPOINTS
		log.Println("<<<<<<<< STABLERAIL ENDPOINTS ACTIVATED >>>>>>>>>")
		router.GET("/v1/users/stablerail/banks", middleware.AuthenticationMiddlewareUsingTimestamp(), getUsersStablerailBanksHandler(callBackRetryChan, gc))

		router.POST("/v1/users/stablerail/onboarduser/:bvn", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersStablerailOnboarduserBvnHandler(callBackRetryChan, gc))

		router.POST("/v1/users/stablerail/onrampcngn/:amount", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersStablerailOnrampcngnAmountHandler(callBackRetryChan, gc))

	}
	router.POST("/v1/users/asset/opt-in", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersAssetOptInHandler(callBackRetryChan, gc))

	router.POST("/v1/shared-access/users/asset/opt-in", middleware.AuthenticationMiddlewareUsingTimestamp(), postSharedAccessUsersAssetOptInHandler(callBackRetryChan, gc))

	router.DELETE("/v1/users/asset/opt-out", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteUsersAssetOptOutHandler(callBackRetryChan, gc))

	router.DELETE("/v1/shared-access/users/asset/opt-out", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteSharedAccessUsersAssetOptOutHandler(callBackRetryChan, gc))

	router.PUT("/v1/users/actions/claim-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), putUsersActionsClaimAssetHandler(callBackRetryChan, gc))

	router.DELETE("/v1/users/actions/reject-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteUsersActionsRejectAssetHandler(callBackRetryChan, gc))

	router.PUT("/v1/shared-access/users/actions/claim-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), putSharedAccessUsersActionsClaimAssetHandler(callBackRetryChan, gc))

	router.DELETE("/v1/shared-access/users/actions/reject-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteSharedAccessUsersActionsRejectAssetHandler(callBackRetryChan, gc))

	router.GET("/v1/users/payment/generate/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), getUsersPaymentGenerateTargetUserHandler(callBackRetryChan, gc))

	router.POST("/v1/security-questions", middleware.AuthenticationMiddlewareUsingTimestamp(), postSecurityQuestionsHandler(callBackRetryChan, gc))

	router.GET("/v1/security-questions/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), getSecurityQuestionsTargetUserHandler(callBackRetryChan, gc))

	router.POST("/v1/verify-answers/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), postVerifyAnswersTargetUserHandler(callBackRetryChan, gc))

	router.POST("/v1/account/recovery/request-email-otp/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), postAccountRecoveryRequestEmailOtpTargetUserHandler(callBackRetryChan, gc))

	router.POST("/v1/account/recovery/verify-email-otp/:targetUser/:otp", middleware.AuthenticationMiddlewareUsingTimestamp(), postAccountRecoveryVerifyEmailOtpTargetUserOtpHandler(callBackRetryChan, gc))

	router.POST("/v1/verify-email-otp/:targetUser/:otp", middleware.AuthenticationMiddlewareUsingTimestamp(), postVerifyEmailOtpTargetUserOtpHandler(callBackRetryChan, gc))

	router.POST("/v1/users/account/recovery", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersAccountRecoveryHandler(callBackRetryChan, gc))

	router.DELETE("/v1/users/account/recovery", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteUsersAccountRecoveryHandler(callBackRetryChan, gc))

	router.POST("/v1/users/account/recover", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersAccountRecoverHandler(callBackRetryChan, gc))

	router.POST("/v1/users/inactive-account/recover", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersInactiveAccountRecoverHandler(callBackRetryChan, gc))

	router.POST("/v1/shared-access/users/account", middleware.AuthenticationMiddlewareUsingTimestamp(), postSharedAccessUsersAccountHandler(callBackRetryChan, gc))

	//modify shared access
	router.PUT("/v1/shared-access/users/account", middleware.AuthenticationMiddlewareUsingTimestamp(), putSharedAccessUsersAccountHandler(callBackRetryChan, gc))

	router.DELETE("/v1/shared-access/users/account", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteSharedAccessUsersAccountHandler(callBackRetryChan, gc))

	//get transaction list
	router.GET("/v1/shared-access/approvals", middleware.AuthenticationMiddlewareUsingTimestamp(), getSharedAccessApprovalsHandler(callBackRetryChan, gc))

	//get specific transaction
	router.GET("/v1/shared-access/approval/:ID", middleware.AuthenticationMiddlewareUsingTimestamp(), getSharedAccessApprovalIDHandler(callBackRetryChan, gc))

	//submit specific approval signature
	router.POST("/v1/shared-access/approval/:ID", middleware.AuthenticationMiddlewareUsingTimestamp(), postSharedAccessApprovalIDHandler(callBackRetryChan, gc))

	//reject specific approval signature
	router.DELETE("/v1/shared-access/approval/:ID", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteSharedAccessApprovalIDHandler(callBackRetryChan, gc))

	//get specific  wallet balance
	router.GET("/v1/shared-access/wallet-balances", middleware.AuthenticationMiddlewareUsingTimestamp(), getSharedAccessWalletBalancesHandler(callBackRetryChan, gc))

	//get specific  wallet balance
	router.GET("/v1/trovo-manager/wallet-balances/:walletPublicKey", middleware.JwtTokenAuthMiddleware(), getTrovoManagerWalletBalancesWalletPublicKeyHandler(callBackRetryChan, gc))

	// market making
	{
		router.POST("/v1/users/trades", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersTradesHandler(callBackRetryChan, gc))

	}
	//CRYPTO
	{
		router.GET("/v1/crypto/withdrawal-history/:currency/:targetPublicKeyForHistory", middleware.AuthenticationMiddlewareUsingTimestamp(), getCryptoWithdrawalHistoryCurrencyTargetPublicKeyForHistoryHandler(callBackRetryChan, gc))

		router.GET("/v1/crypto/deposit-history/:currency/:targetPublicKeyForHistory", middleware.AuthenticationMiddlewareUsingTimestamp(), getCryptoDepositHistoryCurrencyTargetPublicKeyForHistoryHandler(callBackRetryChan, gc))

		//get specific  wallet balance, middleware.AuthenticationMiddlewareUsingTimestamp()
		router.GET("/v1/crypto/withdrawal-networks/:currency", middleware.AuthenticationMiddlewareUsingTimestamp(), getCryptoWithdrawalNetworksCurrencyHandler(callBackRetryChan, gc))

		router.POST("/v1/crypto/withdrawals", middleware.AuthenticationMiddlewareUsingTimestamp(), postCryptoWithdrawalsHandler(callBackRetryChan, gc))

		router.POST("/v1/shared-access/crypto/withdrawals", middleware.AuthenticationMiddlewareUsingTimestamp(), postSharedAccessCryptoWithdrawalsHandler(callBackRetryChan, gc))

		router.POST("/v1/crypto/generate-addresses/:currency", middleware.AuthenticationMiddlewareUsingTimestamp(), postCryptoGenerateAddressesCurrencyHandler(callBackRetryChan, gc))

	}

	//PATRON
	if os.Getenv("ENABLE_PATRON") == "1" {
		router.GET("/v1/patron", middleware.AuthenticationMiddlewareUsingTimestamp(), getPatronHandler(callBackRetryChan, gc))

		router.POST("/v1/patron", middleware.AuthenticationMiddlewareUsingTimestamp(), postPatronHandler(callBackRetryChan, gc))

	}

	if os.Getenv("ENABLE_ASSET_TOKENIZATION") == "1" {
		log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>ASSET TOKENIZATION is enabled!")
		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) != 56 {
			log.Fatalln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ENV TOKENIZATION_ISSUING_PROFILE_WALLET is missing!")

		}
		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE")) == 0 {
			log.Fatalln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ENV TOKENIZATION_ISSUING_PROFILE is missing!")

		}
		if len(os.Getenv("TOKENIZATION_FEE_WALLET")) != 56 {
			log.Fatalln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ENV TOKENIZATION_FEE_WALLET is missing!")

		}

		if len(os.Getenv("TOKENIZATION_APPLICATION_FEE_ASSET")) < 56 {
			log.Fatalln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ENV TOKENIZATION_APPLICATION_FEE_ASSET is invaid!")

		}

		if len(os.Getenv("TOKENIZATION_APPLICATION_FEE_AMOUNT")) == 0 {
			log.Fatalln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ENV TOKENIZATION_APPLICATION_FEE_AMOUNT is invaid!")

		}

		if len(os.Getenv("TOKENIZATION_APPLICATION_FEE_WALLET")) == 0 {
			log.Fatalln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ENV TOKENIZATION_APPLICATION_FEE_WALLET is missing!")

		}

		router.GET("/v1/closed-groups", middleware.AuthenticationMiddlewareUsingTimestamp(), getClosedGroupsHandler(callBackRetryChan, gc))

		router.GET("/v1/banks/:countryCode", middleware.AuthenticationMiddlewareUsingTimestamp(), getBanksCountryCodeHandler(callBackRetryChan, gc))

		router.GET("/v1/forms/:formId", middleware.AuthenticationMiddlewareUsingTimestamp(), getFormsFormIdHandler(callBackRetryChan, gc))

		router.GET("/v1/trovo-manager/banks/:countryCode", middleware.JwtTokenAuthMiddleware(), getTrovoManagerBanksCountryCodeHandler(callBackRetryChan, gc))

		router.GET("/v1/public/tokenization", getPublicTokenizationHandler(callBackRetryChan, gc))

		router.GET("/v1/tokenization", middleware.AuthenticationMiddlewareUsingTimestamp(), getTokenizationHandler(callBackRetryChan, gc))

		router.GET("/v1/trovo-manager/tokenization", middleware.JwtTokenAuthMiddleware(), getTrovoManagerTokenizationHandler(callBackRetryChan, gc))

		router.GET("/v1/tokenization/detail/:tid", middleware.AuthenticationMiddlewareUsingTimestamp(), getTokenizationDetailTidHandler(callBackRetryChan, gc))

		router.GET("/v1/trovo-manager/tokenization/detail/:tid", middleware.JwtTokenAuthMiddleware(), getTrovoManagerTokenizationDetailTidHandler(callBackRetryChan, gc))

		router.GET("/v1/tokenization/list", middleware.AuthenticationMiddlewareUsingTimestamp(), getTokenizationListHandler(callBackRetryChan, gc))

		router.GET("/v1/trovo-manager/tokenization/list", middleware.JwtTokenAuthMiddleware(), getTrovoManagerTokenizationListHandler(callBackRetryChan, gc))

		router.GET("/v1/trovo-manager/tokenization/stat", middleware.JwtTokenAuthMiddleware(), getTrovoManagerTokenizationStatHandler(callBackRetryChan, gc))

		router.POST("/v1/tokenization/expressed-interests/:tokenizedAssetID", middleware.AuthenticationMiddlewareUsingTimestamp(), postTokenizationExpressedInterestsTokenizedAssetIDHandler(callBackRetryChan, gc))

		router.POST("/v1/tokenization/subscriptions/:tokenizedAssetID", middleware.AuthenticationMiddlewareUsingTimestamp(), postTokenizationSubscriptionsTokenizedAssetIDHandler(callBackRetryChan, gc))

		router.POST("/v1/shared-access/tokenization/subscriptions/:tokenizedAssetID", middleware.AuthenticationMiddlewareUsingTimestamp(), postSharedAccessTokenizationSubscriptionsTokenizedAssetIDHandler(callBackRetryChan, gc))

		router.GET("/v1/tokenization/expressed-interests", middleware.AuthenticationMiddlewareUsingTimestamp(), getTokenizationExpressedInterestsHandler(callBackRetryChan, gc))

		router.GET("/v1/tokenization/subscriptions", middleware.AuthenticationMiddlewareUsingTimestamp(), getTokenizationSubscriptionsHandler(callBackRetryChan, gc))

		router.POST("/v1/tokenization", middleware.AuthenticationMiddlewareUsingTimestamp(), postTokenizationHandler(callBackRetryChan, gc))

		router.PUT("/v1/trovo-manager/tokenization/update/:tid", middleware.JwtTokenAuthMiddleware(), putTrovoManagerTokenizationUpdateTidHandler(callBackRetryChan, gc))

		router.PUT("/v1/trovo-manager/tokenization/salesdate/:tid", middleware.JwtTokenAuthMiddleware(), putTrovoManagerTokenizationSalesdateTidHandler(callBackRetryChan, gc))

		router.PUT("/v1/trovo-manager/tokenization/vet/:tid", middleware.JwtTokenAuthMiddleware(), putTrovoManagerTokenizationVetTidHandler(callBackRetryChan, gc))

		router.POST("/v1/trovo-manager/tokenization/faildd/:tid", middleware.JwtTokenAuthMiddleware(), postTrovoManagerTokenizationFailddTidHandler(callBackRetryChan, gc))

		router.POST("/v1/trovo-manager/tokenization/fee/:tid", middleware.JwtTokenAuthMiddleware(), postTrovoManagerTokenizationFeeTidHandler(callBackRetryChan, gc))

		router.POST("/v1/trovo-manager/tokenization/mint/:tid", middleware.JwtTokenAuthMiddleware(), postTrovoManagerTokenizationMintTidHandler(callBackRetryChan, gc))

		router.PUT("/v1/tokenization/confirm/:tokenizationID", middleware.AuthenticationMiddlewareUsingTimestamp(), putTokenizationConfirmTokenizationIDHandler(callBackRetryChan, gc))

		router.DELETE("/v1/tokenization/:tokenizationID", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteTokenizationTokenizationIDHandler(callBackRetryChan, gc))

		router.DELETE("/v1/trovo-manager/tokenization/:tokenizationID", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteTrovoManagerTokenizationTokenizationIDHandler(callBackRetryChan, gc))

		router.PUT("/v1/tokenization/document", middleware.AuthenticationMiddlewareUsingTimestamp(), putTokenizationDocumentHandler(callBackRetryChan, gc))

		router.PUT("/v1/tokenization/logo", middleware.AuthenticationMiddlewareUsingTimestamp(), putTokenizationLogoHandler(callBackRetryChan, gc))

		router.PUT("/v1/trovo-manager/tokenization/logo/:tid", middleware.JwtTokenAuthMiddleware(), putTrovoManagerTokenizationLogoTidHandler(callBackRetryChan, gc))

		router.PUT("/v1/trovo-manager/tokenization/document", middleware.JwtTokenAuthMiddleware(), putTrovoManagerTokenizationDocumentHandler(callBackRetryChan, gc))

		router.POST("/v1/tokenization/fee/:tokenizationID", middleware.AuthenticationMiddlewareUsingTimestamp(), postTokenizationFeeTokenizationIDHandler(callBackRetryChan, gc))

		router.PUT("/v1/tokenization/fee/:tokenizedAssetID", middleware.AuthenticationMiddlewareUsingTimestamp(), putTokenizationFeeTokenizedAssetIDHandler(callBackRetryChan, gc))

		router.DELETE("/v1/tokenization/document/:documentID", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteTokenizationDocumentDocumentIDHandler(callBackRetryChan, gc))

		router.DELETE("/v1/trovo-manager/tokenization/document/:documentID", middleware.JwtTokenAuthMiddleware(), deleteTrovoManagerTokenizationDocumentDocumentIDHandler(callBackRetryChan, gc))

		router.DELETE("/v1/tokenization/fee/:documentID", middleware.AuthenticationMiddlewareUsingTimestamp(), deleteTokenizationFeeDocumentIDHandler(callBackRetryChan, gc))

	}

}
