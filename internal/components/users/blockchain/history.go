package users

import (
	"log"
	"sort"
	"strings"
	"sync"
	"time"
	dbassets "trovo-wallet-api/internal/components/assets/db"
	assetsmodels "trovo-wallet-api/internal/components/assets/models"
	usermodels "trovo-wallet-api/internal/components/users/models"
	bantudb "trovo-wallet-api/internal/db"
	bantupayerrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"github.com/stellar/go/clients/horizonclient"
	"gorm.io/gorm"
)

//PaymentItem holds blockchain payment item
type PaymentItem struct {
	DestinationUsername       string    `json:"destinationUsername"`
	DestinationFirstName      string    `json:"destinationFirstName"`
	DestinationLastName       string    `json:"destinationLastName"`
	DestinationImageThumbnail string    `json:"-"`
	DestinationVerified       int       `json:"destinationVerified"`
	FunderUsername            string    `json:"senderUsername"`
	FunderFirstName           string    `json:"senderFirstName"`
	FunderLastName            string    `json:"senderLastName"`
	FunderImageThumbnail      string    `json:"-"`
	FunderVerified            int       `json:"funderVerified"`
	Amount                    string    `json:"amount"`
	AssetCode                 string    `json:"assetCode"`
	AssetIssuer               string    `json:"assetIssuer"`
	SourceAssetCode           string    `json:"sourceAssetCode"`        //for path payment
	SourceAssetIssuer         string    `json:"sourceAssetIssuer"`      //for path payment
	DestinationAssetCode      string    `json:"destinationAssetCode"`   //for path payment
	DestinationAssetIssuer    string    `json:"destinationAssetIssuer"` //for path payment
	TransactionID             string    `json:"transactionID"`
	FeeCharged                int64     `json:"-"`
	Ledger                    int32     `json:"-"`
	TransactionMemo           string    `json:"transactionMemo"`
	TransactionTime           time.Time `json:"transactionTime"`
	Cursor                    string    `json:"cursor"` //paging token returned as cursor to resume streaming from this state
}

//PaymentEntities holds entities in the payment history batch
type PaymentEntities struct {
	Username       string `json:"username"`
	ImageThumbnail string `json:"thumbnail"`
}

//PaymentHistory holds blockchain payment history
type PaymentHistory struct {
	PageCursor string                     `json:"pageCursor"`
	Payments   []PaymentItem              `json:"payments"`
	Entities   map[string]PaymentEntities `json:"entities"`
}

//IndexedBantuOperation for holding the index of the operation for sorting later
type IndexedBantuOperation struct {
	Index     int
	Operation PaymentItem
}

//GetPaymentHistory fetches the bantu account information using public key
func GetPaymentHistory(ownerData usermodels.User, accountKey string, limit uint, order horizonclient.Order, cursor, forTransactionHash string, includeHash bool, db *gorm.DB) (history PaymentHistory, err error) {

	bantudb.PrintDBStats("GetPaymentHistory", db)
	client := network.GetBlockchainClient()
	publicKey := accountKey
	accountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	_, err = client.AccountDetail(accountRequest)
	if err != nil {
		log.Println("[GetBlockchainAccountDetail]: ", err)
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "dial") {
			log.Printf("[client.AccountDetail Network Failure]: %s", "Error Connecting to Expansion Service")
			return history, &bantupayerrors.ErrorTemporaryServerError{}
		} else if strings.Contains(strings.ToLower(err.Error()), "missing") {

			err = &bantupayerrors.ErrorBlockchainAccountNotActivated{}
		} else {

			err = &bantupayerrors.ErrorTemporaryServerError{}
		}
		return
	}
	if len(string(order)) == 0 {
		order = horizonclient.OrderDesc
	}

	// var includeFailed bool
	if limit <= 0 {
		limit = 25
	}
	// fmt.Printf("%+v\n", input)
	var paymentRequest horizonclient.OperationRequest
	if cursor != "" {
		paymentRequest = horizonclient.OperationRequest{
			ForAccount: publicKey,
			Order:      order,
			Cursor:     cursor,
			Limit:      limit,
			Join:       "transactions",
		}
	} else {
		paymentRequest = horizonclient.OperationRequest{
			ForAccount: publicKey,
			Order:      order,
			Limit:      limit,
			Join:       "transactions",
		}
	}

	ops, err := client.Payments(paymentRequest)
	if err != nil {
		log.Println("[GetPayments] error: ", err)
		if strings.Contains(err.Error(), "tls") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {

		} else {
			if hError, ok := err.(*horizonclient.Error); ok {

				//something went wrong, verify stage and check approprate action
				rCode, _ := hError.ResultCodes()
				rS, _ := hError.ResultString()
				log.Println("[GetPayments] error: Problem in Transaction:", hError.Problem, ", result code: ", rCode, ", result string: ", rS, "public key:", publicKey)
				log.Printf("[GetPayments] Problem in Transaction - RESPONSE: %+v\n", hError.Response)
			}
		}

	}

	parsedOps := ops
	unparsedRecs := parsedOps.Embedded.Records
	parsedRecChan := make(chan IndexedBantuOperation, len(unparsedRecs))
	// paymentEntityChan := make(chan PaymentEntities, len(unparsedRecs)*2)
	paymentEntityMap := make(map[string]PaymentEntities)

	var parsedRecIndexed IndexedBantuOperation
	var payments []PaymentItem
	var pageCursor string
	// var parsedRec BantuOperation
	mapCuratedAssets := make(map[string]assetsmodels.CuratedAsset)
	curatedAssets, _ := dbassets.GetCuratedAssets(false, db)
	for _, av := range curatedAssets {
		mapCuratedAssets[av.AssetCode+":"+av.AssetIssuer] = av
	}
	var wg sync.WaitGroup
	if forTransactionHash != "" && !includeHash {
		for i, v := range unparsedRecs {
			pageCursor = v.PagingToken()
			if forTransactionHash == v.GetTransactionHash() {
				break
			}
			wg.Add(1)
			//start go routine here

			go processRetrievedPaymentOperation(publicKey, i, v, parsedRecIndexed, parsedRecChan, ownerData, mapCuratedAssets, &wg, db)

		}
	} else if forTransactionHash != "" && includeHash {
		for i, v := range unparsedRecs {
			pageCursor = v.PagingToken()
			wg.Add(1)
			//start go routine here
			go processRetrievedPaymentOperation(publicKey, i, v, parsedRecIndexed, parsedRecChan, ownerData, mapCuratedAssets, &wg, db)
			if forTransactionHash == v.GetTransactionHash() {
				break
			}
		}
	} else {
		for i, v := range unparsedRecs {
			pageCursor = v.PagingToken()
			wg.Add(1)
			//start go routine here

			go processRetrievedPaymentOperation(publicKey, i, v, parsedRecIndexed, parsedRecChan, ownerData, mapCuratedAssets, &wg, db)

		}
	}
	wg.Wait()
	close(parsedRecChan)
	var parsedIndexOps = make(map[int]IndexedBantuOperation)
	for parsedIndexOp := range parsedRecChan {
		parsedIndexOps[parsedIndexOp.Index] = parsedIndexOp
		if _, ok := paymentEntityMap[parsedIndexOp.Operation.DestinationUsername]; !ok {
			paymentEntityMap[parsedIndexOp.Operation.DestinationUsername] = PaymentEntities{
				Username:       parsedIndexOp.Operation.DestinationUsername,
				ImageThumbnail: parsedIndexOp.Operation.DestinationImageThumbnail}
		}
		if _, ok := paymentEntityMap[parsedIndexOp.Operation.FunderUsername]; !ok {
			paymentEntityMap[parsedIndexOp.Operation.FunderUsername] = PaymentEntities{
				Username:       parsedIndexOp.Operation.FunderUsername,
				ImageThumbnail: parsedIndexOp.Operation.FunderImageThumbnail}
		}

	}
	// To store the keys in slice in sorted order
	// var entities []PaymentEntities
	// for _, v := range paymentEntityMap {
	// 	entities = append(entities, v)
	// }
	keys := make([]int, len(parsedIndexOps))
	i := 0
	for k := range parsedIndexOps {
		keys[i] = k
		i++
	}
	sort.Ints(keys)

	// To perform the opertion you want
	for _, k := range keys {
		// fmt.Println("Key:", k, "Value:", m[k])
		parsedRec := parsedIndexOps[k].Operation
		parsedRec.DestinationImageThumbnail = ""
		parsedRec.FunderImageThumbnail = ""
		payments = append(payments, parsedRec)
	}
	history.Payments = payments
	history.PageCursor = pageCursor
	history.Entities = paymentEntityMap

	return
}
