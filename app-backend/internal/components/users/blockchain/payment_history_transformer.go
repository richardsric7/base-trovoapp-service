package users

import (
	"sync"
	usermodels "trovo-wallet-api/internal/components/users/models"

	"gorm.io/gorm"
)

// v was a Horizon operations.Operation (Stellar); Base has no equivalent
// account-agnostic operation feed to source one from (see main.go's
// MonitorStream doc comment), so these are unreachable stubs kept for
// call-site compatibility.
func ProcessRetrievedPaymentOperation(publicKey string, i int, v interface{}, parsedRecIndexed IndexedBantuOperation, parsedRecChan chan IndexedBantuOperation, ownerData usermodels.User, wg *sync.WaitGroup, db *gorm.DB) {

	//end go routine here
}

func ProcessStreamPaymentOperation(publicKey string, v interface{}, ownerData usermodels.User, db *gorm.DB) (parsedRec PaymentItem) {

	return
}
