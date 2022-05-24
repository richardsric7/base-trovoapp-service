package users

import (
	"errors"
	"sync"
	assetsmodels "trovo-wallet-api/internal/components/assets/models"
	users "trovo-wallet-api/internal/components/users/db"
	usermodels "trovo-wallet-api/internal/components/users/models"

	"github.com/stellar/go/protocols/horizon/operations"
	"gorm.io/gorm"
)

func processRetrievedPaymentOperation(publicKey string, i int, v operations.Operation, parsedRecIndexed IndexedBantuOperation, parsedRecChan chan IndexedBantuOperation, ownerData usermodels.User, curatedAssets map[string]assetsmodels.CuratedAsset, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()
	parsedRec := parsedRecIndexed.Operation

	if v.GetType() == "payment" {
		//payment transaction.
		pmt := interface{}(v).(operations.Payment)
		// fmt.Println("Payment made to: ", pmt.To)
		// pv := interface{}(pmt).(operations.Operation)
		if pmt.Issuer != "" {
			if _, ok := curatedAssets[pmt.Code+":"+pmt.Issuer]; !ok {
				return
			}
		}
		if publicKey == pmt.To {
			// same as owner of account. transaction is made to owner.

			if ownerData.Username == "" {
				parsedRec.DestinationUsername = pmt.To

			} else {
				parsedRec.DestinationUsername = ownerData.Username
				parsedRec.DestinationVerified = ownerData.Verified
			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.DestinationFirstName = ownerData.FirstName

			parsedRec.DestinationLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.To, db)
			if dberr != nil {
				if dberr.Error() == "error-username-does-not-exist" {
					parsedRec.DestinationUsername = pmt.To

				}
			} else {
				parsedRec.DestinationUsername = funder.Username
				parsedRec.DestinationVerified = funder.Verified

			}

			if funder.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.DestinationFirstName = funder.FirstName

			parsedRec.DestinationLastName = funder.LastName

		}
		if publicKey == pmt.From {
			// same as owner of account. sent from owner account
			if ownerData.Username == "" {
				parsedRec.FunderUsername = pmt.From
			} else {
				parsedRec.FunderUsername = ownerData.Username
				parsedRec.FunderVerified = ownerData.Verified

			}
			if ownerData.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.FunderFirstName = ownerData.FirstName

			parsedRec.FunderLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.From, db)
			if dberr != nil {
				if dberr.Error() == "error-username-does-not-exist" {
					parsedRec.FunderUsername = pmt.From
				}
			} else {
				parsedRec.FunderUsername = funder.Username
				parsedRec.FunderVerified = funder.Verified

			}
			if funder.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.FunderFirstName = funder.FirstName

			parsedRec.FunderLastName = funder.LastName

		}
		parsedRec.Amount = pmt.Amount

		if pmt.Asset.Type != "native" {
			parsedRec.AssetIssuer = pmt.Asset.Issuer
			parsedRec.AssetCode = pmt.Asset.Code
		}
		parsedRec.TransactionID = pmt.TransactionHash
		parsedRec.TransactionTime = pmt.Transaction.LedgerCloseTime
		parsedRec.TransactionMemo = pmt.Transaction.Memo

		//send through channel
		parsedRecIndexed = IndexedBantuOperation{Operation: parsedRec, Index: i}

		// fmt.Printf("Indexed Payment: %+v", parsedRecIndexed)
		parsedRecChan <- parsedRecIndexed

	} else if v.GetType() == "create_account" {
		//payment transaction.
		pmt := interface{}(v).(operations.CreateAccount)
		// fmt.Println("Account Created: ", pmt.Account)

		if publicKey == pmt.Account {
			// same as owner of account. owner account created.
			if ownerData.Username == "" {
				parsedRec.DestinationUsername = pmt.Account
			} else {
				parsedRec.DestinationUsername = ownerData.Username
				parsedRec.DestinationVerified = ownerData.Verified

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.DestinationFirstName = ownerData.FirstName

			parsedRec.DestinationLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.Account, db)
			if dberr != nil {
				if dberr.Error() == "error-username-does-not-exist" {
					parsedRec.DestinationUsername = pmt.Account
				}
			} else {
				parsedRec.DestinationUsername = funder.Username
				parsedRec.DestinationVerified = funder.Verified
			}

			if funder.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.DestinationFirstName = funder.FirstName

			parsedRec.DestinationLastName = funder.LastName

		}

		if publicKey == pmt.Funder {
			// same as owner of account. owner activated account
			if ownerData.Username == "" {
				parsedRec.FunderUsername = pmt.Funder
			} else {
				parsedRec.FunderUsername = ownerData.Username
				parsedRec.FunderVerified = ownerData.Verified

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.FunderFirstName = ownerData.FirstName

			parsedRec.FunderLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.Funder, db)
			if dberr != nil {
				if dberr.Error() == "error-username-does-not-exist" {
					parsedRec.FunderUsername = pmt.Funder
				}
			} else {
				parsedRec.FunderUsername = funder.Username
				parsedRec.FunderVerified = funder.Verified
			}

			if funder.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.FunderFirstName = funder.FirstName

			parsedRec.FunderLastName = funder.LastName

		}
		parsedRec.Amount = pmt.StartingBalance
		parsedRec.AssetIssuer = ""
		parsedRec.AssetCode = ""

		parsedRec.TransactionID = pmt.TransactionHash
		parsedRec.TransactionTime = pmt.Transaction.LedgerCloseTime
		parsedRec.TransactionMemo = pmt.Transaction.Memo
		//send through channel
		parsedRecIndexed = IndexedBantuOperation{Operation: parsedRec, Index: i}
		parsedRecChan <- parsedRecIndexed

	} else if v.GetType() == "account_merge" {
		//payment transaction.
		pmt := interface{}(v).(operations.AccountMerge)
		// fmt.Println("Account Created: ", pmt.Account)

		if publicKey == pmt.Into {
			// same as owner of account. owner account received merge tokens.
			if ownerData.Username == "" {
				parsedRec.DestinationUsername = pmt.Into
			} else {
				parsedRec.DestinationUsername = ownerData.Username
				parsedRec.DestinationVerified = ownerData.Verified

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.DestinationFirstName = ownerData.FirstName

			parsedRec.DestinationLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.Account, db)
			if dberr != nil {
				if dberr.Error() == "error-username-does-not-exist" {
					parsedRec.DestinationUsername = pmt.Account
				}
			} else {
				parsedRec.DestinationUsername = funder.Username
				parsedRec.DestinationVerified = funder.Verified
			}

			if funder.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.DestinationFirstName = funder.FirstName

			parsedRec.DestinationLastName = funder.LastName

		}

		if publicKey == pmt.Account {
			// same as owner of account. owner activated account
			if ownerData.Username == "" {
				parsedRec.FunderUsername = pmt.Account
			} else {
				parsedRec.FunderUsername = ownerData.Username
				parsedRec.FunderVerified = ownerData.Verified

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.FunderFirstName = ownerData.FirstName

			parsedRec.FunderLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.Account, db)

			if dberr != nil {

				if dberr.Error() == "error-username-does-not-exist" {
					parsedRec.FunderUsername = pmt.Account
				}
			} else {
				parsedRec.FunderUsername = funder.Username
				parsedRec.FunderVerified = funder.Verified
			}

			if funder.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.FunderFirstName = funder.FirstName

			parsedRec.FunderLastName = funder.LastName

		}
		parsedRec.Amount = "0"
		parsedRec.AssetIssuer = ""
		parsedRec.AssetCode = ""

		parsedRec.TransactionID = pmt.TransactionHash
		parsedRec.TransactionTime = pmt.Transaction.LedgerCloseTime
		parsedRec.TransactionMemo = pmt.Transaction.Memo
		//send through channel
		parsedRecIndexed = IndexedBantuOperation{Operation: parsedRec, Index: i}
		parsedRecChan <- parsedRecIndexed

	} else if v.GetType() == "path_payment_strict_send" {
		//payment transaction.
		pmt := interface{}(v).(operations.PathPaymentStrictSend)
		// fmt.Println("Payment made to: ", pmt.To)

		if pmt.Issuer != "" {
			if _, ok := curatedAssets[pmt.Code+":"+pmt.Issuer]; !ok {
				return
			}
		}
		if pmt.SourceAssetIssuer != "" {
			if _, ok := curatedAssets[pmt.SourceAssetCode+":"+pmt.SourceAssetIssuer]; !ok {
				return
			}
		}
		if publicKey == pmt.To {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.DestinationUsername = pmt.To
			} else {
				parsedRec.DestinationUsername = ownerData.Username
				parsedRec.DestinationVerified = ownerData.Verified

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *ownerData.ImageThumbnail

			}
			parsedRec.DestinationFirstName = ownerData.FirstName

			parsedRec.DestinationLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.To, db)

			if dberr != nil {
				if dberr.Error() == "error-username-does-not-exist" {
					parsedRec.DestinationUsername = pmt.To
				}
			} else {
				parsedRec.DestinationUsername = funder.Username
				parsedRec.DestinationVerified = funder.Verified

			}

			if funder.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.DestinationFirstName = funder.FirstName

			parsedRec.DestinationLastName = funder.LastName

		}
		if publicKey == pmt.From {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.FunderUsername = pmt.From
			} else {
				parsedRec.FunderUsername = ownerData.Username
				parsedRec.FunderVerified = ownerData.Verified

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.FunderFirstName = ownerData.FirstName

			parsedRec.FunderLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.From, db)

			if dberr != nil {
				if dberr.Error() == "error-username-does-not-exist" {
					parsedRec.FunderUsername = pmt.From
				}
			} else {
				parsedRec.FunderUsername = funder.Username
				parsedRec.FunderVerified = funder.Verified

			}

			if funder.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.FunderFirstName = funder.FirstName

			parsedRec.FunderLastName = funder.LastName

		}
		parsedRec.Amount = pmt.Amount
		if pmt.Asset.Type != "native" {
			parsedRec.AssetCode = pmt.Asset.Code
			parsedRec.AssetIssuer = pmt.Asset.Issuer
			parsedRec.DestinationAssetCode = pmt.Asset.Code
			parsedRec.DestinationAssetIssuer = pmt.Asset.Issuer

		}
		if pmt.SourceAssetType != "native" {
			parsedRec.SourceAssetCode = pmt.SourceAssetCode
			parsedRec.SourceAssetIssuer = pmt.SourceAssetIssuer

		}
		parsedRec.TransactionID = pmt.TransactionHash
		parsedRec.TransactionTime = pmt.LedgerCloseTime
		parsedRec.TransactionMemo = pmt.Transaction.Memo
		//send through channel
		parsedRecIndexed = IndexedBantuOperation{Operation: parsedRec, Index: i}
		parsedRecChan <- parsedRecIndexed

	} else if v.GetType() == "path_payment" {
		//payment transaction.
		pmt := interface{}(v).(operations.PathPayment)
		// fmt.Println("Payment made to: ", pmt.To)
		if pmt.Issuer != "" {
			if _, ok := curatedAssets[pmt.Code+":"+pmt.Issuer]; !ok {
				return
			}
		}
		if pmt.SourceAssetIssuer != "" {
			if _, ok := curatedAssets[pmt.SourceAssetCode+":"+pmt.SourceAssetIssuer]; !ok {
				return
			}
		}
		if publicKey == pmt.To {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.DestinationUsername = pmt.To
			} else {
				parsedRec.DestinationUsername = ownerData.Username
				parsedRec.DestinationVerified = ownerData.Verified

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *ownerData.ImageThumbnail

			}
			parsedRec.DestinationFirstName = ownerData.FirstName

			parsedRec.DestinationLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.To, db)

			if dberr != nil {
				if dberr.Error() == "error-username-does-not-exist" {
					parsedRec.DestinationUsername = pmt.To
				}
			} else {
				parsedRec.DestinationUsername = funder.Username
				parsedRec.DestinationVerified = funder.Verified

			}

			if funder.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.DestinationFirstName = funder.FirstName

			parsedRec.DestinationLastName = funder.LastName

		}
		if publicKey == pmt.From {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.FunderUsername = pmt.From
			} else {
				parsedRec.FunderUsername = ownerData.Username
				parsedRec.FunderVerified = ownerData.Verified

			}
			if ownerData.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *ownerData.ImageThumbnail

			}
			parsedRec.FunderFirstName = ownerData.FirstName
			parsedRec.FunderLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.From, db)

			if dberr != nil {
				if dberr.Error() == "error-username-does-not-exist" {
					parsedRec.FunderUsername = pmt.From
				}
			} else {
				parsedRec.FunderUsername = funder.Username
				parsedRec.FunderVerified = funder.Verified

			}

			if funder.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *funder.ImageThumbnail

			}
			parsedRec.FunderFirstName = funder.FirstName
			parsedRec.FunderLastName = funder.LastName

		}
		parsedRec.Amount = pmt.Amount

		if pmt.Asset.Type != "native" {
			parsedRec.AssetCode = pmt.Asset.Code
			parsedRec.AssetIssuer = pmt.Asset.Issuer
			parsedRec.DestinationAssetCode = pmt.Asset.Code
			parsedRec.DestinationAssetIssuer = pmt.Asset.Issuer

		}
		if pmt.SourceAssetType != "native" {
			parsedRec.SourceAssetCode = pmt.SourceAssetCode
			parsedRec.SourceAssetIssuer = pmt.SourceAssetIssuer

		}

		parsedRec.TransactionID = pmt.Transaction.ID
		parsedRec.TransactionTime = pmt.Transaction.LedgerCloseTime
		parsedRec.TransactionMemo = pmt.Transaction.Memo
		//send through channel
		parsedRecIndexed = IndexedBantuOperation{Operation: parsedRec, Index: i}
		parsedRecChan <- parsedRecIndexed
	}

	//end go routine here
}

func ProcessStreamPaymentOperation(publicKey string, v operations.Operation, ownerData usermodels.User, db *gorm.DB) (parsedRec PaymentItem) {
	//asssign the cursor to use to resume streaming payment streaming from here.
	parsedRec.Cursor = v.PagingToken()
	if v.GetType() == "payment" {
		//payment transaction.
		pmt := interface{}(v).(operations.Payment)
		// fmt.Println("Payment made to: ", pmt.To)
		// pv := interface{}(pmt).(operations.Operation)

		if publicKey == pmt.To {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.DestinationUsername = pmt.To
			} else {
				parsedRec.DestinationUsername = ownerData.Username

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.DestinationFirstName = ownerData.FirstName

			parsedRec.DestinationLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.To, db)
			if errors.Is(dberr, gorm.ErrRecordNotFound) {
				parsedRec.DestinationUsername = pmt.To

			} else {
				parsedRec.DestinationUsername = funder.Username

			}
			if funder.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.DestinationFirstName = funder.FirstName

			parsedRec.DestinationLastName = funder.LastName

		}
		if publicKey == pmt.From {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.FunderUsername = pmt.From
			} else {
				parsedRec.FunderUsername = ownerData.Username

			}
			if ownerData.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.FunderFirstName = ownerData.FirstName

			parsedRec.FunderLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.From, db)
			if errors.Is(dberr, gorm.ErrRecordNotFound) {
				parsedRec.FunderUsername = pmt.From
			} else {
				parsedRec.FunderUsername = funder.Username

			}
			if funder.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.FunderFirstName = funder.FirstName

			parsedRec.FunderLastName = funder.LastName

		}
		parsedRec.Amount = pmt.Amount

		if pmt.Asset.Type != "native" {
			parsedRec.AssetIssuer = pmt.Asset.Issuer
			parsedRec.AssetCode = pmt.Asset.Code
		}
		parsedRec.TransactionID = pmt.TransactionHash
		parsedRec.TransactionTime = pmt.Transaction.LedgerCloseTime
		parsedRec.TransactionMemo = pmt.Transaction.Memo

	} else if v.GetType() == "create_account" {
		//payment transaction.
		pmt := interface{}(v).(operations.CreateAccount)
		// fmt.Println("Account Created: ", pmt.Account)
		if publicKey == pmt.Account {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.DestinationUsername = pmt.Account
			} else {
				parsedRec.DestinationUsername = ownerData.Username

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.DestinationFirstName = ownerData.FirstName

			parsedRec.DestinationLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.Account, db)
			if errors.Is(dberr, gorm.ErrRecordNotFound) {
				parsedRec.DestinationUsername = pmt.Account
			} else {
				parsedRec.DestinationUsername = funder.Username
			}
			if funder.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.DestinationFirstName = funder.FirstName

			parsedRec.DestinationLastName = funder.LastName

		}

		if publicKey == pmt.Funder {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.FunderUsername = pmt.Funder
			} else {
				parsedRec.FunderUsername = ownerData.Username

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.FunderFirstName = ownerData.FirstName

			parsedRec.FunderLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.Funder, db)
			if errors.Is(dberr, gorm.ErrRecordNotFound) {
				parsedRec.FunderUsername = pmt.Funder
			} else {
				parsedRec.FunderUsername = funder.Username
			}

			if funder.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.FunderFirstName = funder.FirstName

			parsedRec.FunderLastName = funder.LastName

		}
		parsedRec.Amount = pmt.StartingBalance
		parsedRec.AssetIssuer = ""
		parsedRec.AssetCode = ""

		parsedRec.TransactionID = pmt.TransactionHash
		parsedRec.TransactionTime = pmt.Transaction.LedgerCloseTime
		parsedRec.TransactionMemo = pmt.Transaction.Memo

	} else if v.GetType() == "path_payment_strict_send" {
		//payment transaction.
		pmt := interface{}(v).(operations.PathPaymentStrictSend)
		// fmt.Println("Payment made to: ", pmt.To)

		if publicKey == pmt.To {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.DestinationUsername = pmt.To
			} else {
				parsedRec.DestinationUsername = ownerData.Username

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *ownerData.ImageThumbnail

			}
			parsedRec.DestinationFirstName = ownerData.FirstName

			parsedRec.DestinationLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.To, db)
			if errors.Is(dberr, gorm.ErrRecordNotFound) {
				parsedRec.DestinationUsername = pmt.To
			} else {
				parsedRec.DestinationUsername = funder.Username

			}

			if funder.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.DestinationFirstName = funder.FirstName

			parsedRec.DestinationLastName = funder.LastName

		}
		if publicKey == pmt.From {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.FunderUsername = pmt.From
			} else {
				parsedRec.FunderUsername = ownerData.Username

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *ownerData.ImageThumbnail

			}

			parsedRec.FunderFirstName = ownerData.FirstName

			parsedRec.FunderLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.From, db)
			if errors.Is(dberr, gorm.ErrRecordNotFound) {
				parsedRec.FunderUsername = pmt.From
			} else {
				parsedRec.FunderUsername = funder.Username

			}

			if funder.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.FunderFirstName = funder.FirstName

			parsedRec.FunderLastName = funder.LastName

		}
		parsedRec.Amount = pmt.Amount
		if pmt.Asset.Type != "native" {
			parsedRec.AssetCode = pmt.Asset.Code
			parsedRec.AssetIssuer = pmt.Asset.Issuer
			parsedRec.DestinationAssetCode = pmt.Asset.Code
			parsedRec.DestinationAssetIssuer = pmt.Asset.Issuer

		}
		if pmt.SourceAssetType != "native" {
			parsedRec.SourceAssetCode = pmt.SourceAssetCode
			parsedRec.SourceAssetIssuer = pmt.SourceAssetIssuer

		}
		parsedRec.TransactionID = pmt.TransactionHash
		parsedRec.TransactionTime = pmt.Transaction.LedgerCloseTime
		parsedRec.TransactionMemo = pmt.Transaction.Memo

	} else if v.GetType() == "path_payment" {
		//payment transaction.
		pmt := interface{}(v).(operations.PathPayment)
		// fmt.Println("Payment made to: ", pmt.To)

		if publicKey == pmt.To {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.DestinationUsername = pmt.To
			} else {
				parsedRec.DestinationUsername = ownerData.Username

			}

			if ownerData.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *ownerData.ImageThumbnail

			}
			parsedRec.DestinationFirstName = ownerData.FirstName

			parsedRec.DestinationLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.To, db)
			if errors.Is(dberr, gorm.ErrRecordNotFound) {
				parsedRec.DestinationUsername = pmt.To
			} else {
				parsedRec.DestinationUsername = funder.Username

			}

			if funder.ImageThumbnail != nil {
				parsedRec.DestinationImageThumbnail = *funder.ImageThumbnail

			}

			parsedRec.DestinationFirstName = funder.FirstName

			parsedRec.DestinationLastName = funder.LastName

		}
		if publicKey == pmt.From {
			// same as owner of account
			if ownerData.Username == "" {
				parsedRec.FunderUsername = pmt.From
			} else {
				parsedRec.FunderUsername = ownerData.Username

			}
			if ownerData.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *ownerData.ImageThumbnail

			}
			parsedRec.FunderFirstName = ownerData.FirstName
			parsedRec.FunderLastName = ownerData.LastName

		} else {
			funder, dberr := users.GetUserInfo(pmt.From, db)
			if errors.Is(dberr, gorm.ErrRecordNotFound) {
				parsedRec.FunderUsername = pmt.From
			} else {
				parsedRec.FunderUsername = funder.Username

			}
			if funder.ImageThumbnail != nil {
				parsedRec.FunderImageThumbnail = *funder.ImageThumbnail

			}
			parsedRec.FunderFirstName = funder.FirstName
			parsedRec.FunderLastName = funder.LastName

		}
		parsedRec.Amount = pmt.Amount

		if pmt.Asset.Type != "native" {
			parsedRec.AssetCode = pmt.Asset.Code
			parsedRec.AssetIssuer = pmt.Asset.Issuer
			parsedRec.DestinationAssetCode = pmt.Asset.Code
			parsedRec.DestinationAssetIssuer = pmt.Asset.Issuer

		}
		if pmt.SourceAssetType != "native" {
			parsedRec.SourceAssetCode = pmt.SourceAssetCode
			parsedRec.SourceAssetIssuer = pmt.SourceAssetIssuer

		}

		parsedRec.TransactionID = pmt.TransactionHash
		parsedRec.TransactionTime = pmt.Transaction.LedgerCloseTime
		parsedRec.TransactionMemo = pmt.Transaction.Memo

	}

	return
}
