package users

import (
	"sync"
	usermodels "trovo-wallet-api/internal/components/users/models"

	"github.com/stellar/go/protocols/horizon/operations"
	"gorm.io/gorm"
)

func ProcessRetrievedPaymentOperation(publicKey string, i int, v operations.Operation, parsedRecIndexed IndexedBantuOperation, parsedRecChan chan IndexedBantuOperation, ownerData usermodels.User, wg *sync.WaitGroup, db *gorm.DB) {

	//end go routine here
}

func ProcessStreamPaymentOperation(publicKey string, v operations.Operation, ownerData usermodels.User, db *gorm.DB) (parsedRec PaymentItem) {

	return
}
