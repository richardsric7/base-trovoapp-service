package callbacks

import (
	userModels "trovo-wallet-api/internal/components/users/models"

	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine, callBackRetryChan chan userModels.RetryCallbacks, gc *sharedconfig.GlobalConfig) {

	router.POST("/v1/callbacks/1l", postCallbacks1lHandler(gc))
	router.POST("/v1/callbacks/doja/webhook", postCallbacksDojaWebhookHandler(gc))
	router.POST("/v1/callbacks/flutterwave/webhook", postCallbacksFlutterwaveWebhookHandler(gc))
}
