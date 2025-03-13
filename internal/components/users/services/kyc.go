package users

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"
	userDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetKYCLevels(gc *sharedconfig.GlobalConfig) (levels []userModels.KYCLevel) {
	levels = make([]userModels.KYCLevel, 0)
	gc.DB.Order("user_category ASC, ID ASC").Find(&levels)
	return
}

func GetUserKYCProgress(username string, gc *sharedconfig.GlobalConfig) (progress userModels.UserKYCProgress) {
	gc.DB.Where("username = ?", username).First(&progress)
	if len(progress.Username) == 0 {
		progress.Username = username
	}
	return
}

func GetKYCConfigs(gc *sharedconfig.GlobalConfig) (configs []userModels.KYCConfig) {
	configs = make([]userModels.KYCConfig, 0)
	gc.DB.Order("Service_Provider ASC").Find(&configs)
	return
}

func GetKYCConfigByServiceProvider(serviceProvider string, gc *sharedconfig.GlobalConfig) (config userModels.KYCConfig, err error) {
	err = gc.DB.Where("service_provider = ?", serviceProvider).First(&config).Error
	return
}

func SaveSumsubWebhook(sumSubwhInput *userModels.SumSubReviewResultInput, db *gorm.DB) (err error) {
	dbObj := userModels.SumSubReviewResult{
		ApplicantID:    sumSubwhInput.ApplicantID,
		InspectionID:   sumSubwhInput.InspectionID,
		CorrelationID:  sumSubwhInput.CorrelationID,
		ExternalUserID: sumSubwhInput.ExternalUserID,
		LevelName:      sumSubwhInput.LevelName,
		Type:           sumSubwhInput.Type,
		ReviewStatus:   sumSubwhInput.ReviewStatus,
		CreatedAtMs:    sumSubwhInput.CreatedAtMs,
	}
	bytes, err := json.Marshal(sumSubwhInput.ReviewResult)
	if err != nil {
		log.Printf("[SaveSumsubWebhook] failed to marshal data to bytes due to : %v\n", err)
		return
	}
	dbObj.ReviewResult = string(bytes)
	e := db.Omit(clause.Associations).Create(&dbObj).Error
	if e != nil {
		log.Printf("[SaveSumsubWebhook] error saving reviewresult to database  [%+v], %v\n", dbObj, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	return
}

// VerifySumSubWebhook verifies the webhook sender by comparing the provided x-payload-digest
// with the calculated HMAC-SHA256 digest of the payload.
func VerifySumSubWebhook(payload []byte, xPayloadDigestHeader, xPayloadDigestAlgHeader string, gc *sharedconfig.GlobalConfig) (bool, error) {
	// Ensure the algorithm is HMAC_SHA256_HEX
	if strings.ToUpper(xPayloadDigestAlgHeader) != "HMAC_SHA256_HEX" {
		log.Printf("[VerifySumSubWebhook] algo not supported: %v\n", xPayloadDigestAlgHeader)
		err := &tErrors.ErrorTemporaryServerError{}
		return false, err
	}

	// Create a new HMAC-SHA256 hasher with the secret key
	config, err := GetKYCConfigByServiceProvider("sumsub", gc)
	if err != nil {
		log.Printf("[VerifySumSubWebhook] error getting service config: %v\n", err)
		err = &tErrors.ErrorTemporaryServerError{}
		return false, err
	}
	hasher := hmac.New(sha256.New, []byte(config.SecretKey))

	// Write the raw payload to the hasher
	_, err = hasher.Write(payload)
	if err != nil {
		log.Printf("[VerifySumSubWebhook] error writing payload to hasher: %v\n", err)
		err = &tErrors.ErrorTemporaryServerError{}
		return false, err
	}

	// Calculate the HMAC-SHA256 digest
	calculatedDigest := hex.EncodeToString(hasher.Sum(nil))

	// Compare the calculated digest with the provided x-payload-digest header value
	return calculatedDigest == xPayloadDigestHeader, nil
}

func InitiateUserKYCProgressForSumsub(username, levelName string, gc *sharedconfig.GlobalConfig) (err error) {

	//get kycprogress
	var kycProgress userModels.UserKYCProgress

	gc.DB.Where("username = ?", username).First(&kycProgress)
	kycProgress.Username = username

	if strings.Contains(levelName, "level-1") {
		//level 1
		if kycProgress.KYCLevel1Done == 1 {
			return &tErrors.CustomError{
				Err:        "kyc-level-already-done",
				Param:      "levelName",
				ErrMessage: "KYC level already done",
			}
		}

		kycProgress.KYCLevel1Initiated = 1
	}

	if strings.Contains(levelName, "level-2") {
		//level 2
		if kycProgress.KYCLevel2Done == 1 {
			return &tErrors.CustomError{
				Err:        "kyc-level-already-done",
				Param:      "levelName",
				ErrMessage: "KYC level already done",
			}
		}
		if kycProgress.KYCLevel1Done == 0 {
			return &tErrors.CustomError{
				Err:        "kyc-level-1-not-done",
				Param:      "levelName",
				ErrMessage: "KYC level must be done in the order 1 to 3",
			}
		}
		kycProgress.KYCLevel2Initiated = 1
	}

	if strings.Contains(levelName, "level-3") {
		//level 3
		if kycProgress.KYCLevel3Done == 1 {
			return &tErrors.CustomError{
				Err:        "kyc-level-already-done",
				Param:      "levelName",
				ErrMessage: "KYC level already done",
			}
		}
		if kycProgress.KYCLevel1Done == 0 || kycProgress.KYCLevel2Done == 0 {
			return &tErrors.CustomError{
				Err:        "kyc-level-not-done",
				Param:      "levelName",
				ErrMessage: "KYC level must be done in the order 1 to 3",
			}
		}
		kycProgress.KYCLevel3Initiated = 1
	}

	e := gc.DB.Omit(clause.Associations).Save(&kycProgress).Error
	if e != nil {
		log.Printf("[UpdateUserKYCLevelFromSumsub] error saving user KYC progress to database  [%+v], %v\n", username, e)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	// user.InvalidateUserCache(gc)
	return nil
}

func ResetUserKYCProgressForSumsub(username, levelName string, db *gorm.DB) (err error) {

	//get kycprogress
	var kycProgress userModels.UserKYCProgress

	db.Where("username = ?", username).First(&kycProgress)
	kycProgress.Username = username

	if strings.Contains(levelName, "level-1") {
		//level 1

		kycProgress.KYCLevel1Initiated = 0
		kycProgress.KYCLevel1Done = 0

	}

	if strings.Contains(levelName, "level-2") {
		//level 2

		kycProgress.KYCLevel2Initiated = 0
		kycProgress.KYCLevel2Done = 0
	}

	if strings.Contains(levelName, "level-3") {
		//level 3

		kycProgress.KYCLevel3Initiated = 0
		kycProgress.KYCLevel3Done = 0
	}

	e := db.Omit(clause.Associations).Save(&kycProgress).Error
	if e != nil {
		log.Printf("[UpdateUserKYCLevelFromSumsub] error saving user KYC progress to database  [%+v], %v\n", username, e)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	// user.InvalidateUserCache(gc)
	return nil
}

func CompleteUserKYCProgressForSumsub(user *userModels.User, levelName string, gc *sharedconfig.GlobalConfig) (err error) {

	//get kycprogress
	var kycProgress userModels.UserKYCProgress

	gc.DB.Where("username = ?", user.Username).First(&kycProgress)
	kycProgress.Username = user.Username

	if strings.Contains(levelName, "level-1") {
		//level 1
		if kycProgress.KYCLevel1Done == 1 {
			return &tErrors.CustomError{
				Err:        "kyc-level-already-done",
				Param:      "levelName",
				ErrMessage: "KYC level already done",
			}
		}

		kycProgress.KYCLevel1Done = 1
	}

	if strings.Contains(levelName, "level-2") {
		//level 2
		if kycProgress.KYCLevel2Done == 1 {
			return &tErrors.CustomError{
				Err:        "kyc-level-already-done",
				Param:      "levelName",
				ErrMessage: "KYC level already done",
			}
		}
		if kycProgress.KYCLevel1Done == 0 {
			return &tErrors.CustomError{
				Err:        "kyc-level-1-not-done",
				Param:      "levelName",
				ErrMessage: "KYC level must be done in the order 1 to 3",
			}
		}
		kycProgress.KYCLevel2Done = 1
	}

	if strings.Contains(levelName, "level-3") {
		//level 3
		if kycProgress.KYCLevel3Done == 1 {
			return &tErrors.CustomError{
				Err:        "kyc-level-already-done",
				Param:      "levelName",
				ErrMessage: "KYC level already done",
			}
		}
		if kycProgress.KYCLevel1Done == 0 || kycProgress.KYCLevel2Done == 0 {
			return &tErrors.CustomError{
				Err:        "kyc-level-not-done",
				Param:      "levelName",
				ErrMessage: "KYC level must be done in the order 1 to 3",
			}
		}
		kycProgress.KYCLevel3Done = 1
	}

	e := gc.DB.Omit(clause.Associations).Save(&kycProgress).Error
	if e != nil {
		log.Printf("[CompleteUserKYCProgressForSumsub] error saving user KYC progress to database  [%+v], %v\n", user.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	// user.InvalidateUserCache(gc)
	return nil
}

func UpdateUserKYCLevelFromSumsub(user *userModels.User, sumSubwhInput *userModels.SumSubReviewResultInput, db *gorm.DB, gc *sharedconfig.GlobalConfig) (err error) {
	// get user
	// user, err := userDB.GetUser(sumSubwhInput.ExternalUserID, db, gc)
	// if err != nil {
	// 	log.Printf("[UpdateUserKYCLevelFromSumsub] error getting user: %v, %v\n", sumSubwhInput.ExternalUserID, err)
	// 	err = &tErrors.ErrorTemporaryServerError{}
	// 	return err
	// }
	// //get kycprogress
	// var kycProgress userModels.UserKYCProgress

	// e := db.Where("username = ?", user.Username).First(&kycProgress).Error
	// if e != nil {
	// 	log.Printf("[UpdateUserKYCLevelFromSumsub] error getting kyc progress: %v\n", sumSubwhInput.ExternalUserID)

	// 	err = &tErrors.ErrorTemporaryServerError{}
	// 	return err
	// }

	if strings.Contains(sumSubwhInput.LevelName, "level-1") {
		//level 1
		user.KYCVerified = 1
		// kycProgress.KYCLevel1Done = 1
	}

	if strings.Contains(sumSubwhInput.LevelName, "level-2") {
		//level 2
		user.KYCVerified = 2
		// kycProgress.KYCLevel2Done = 1
	}

	if strings.Contains(sumSubwhInput.LevelName, "level-3") {
		//level 3
		user.KYCVerified = 3
		// kycProgress.KYCLevel3Done = 1
	}

	// if strings.Contains(sumSubwhInput.LevelName, "level-4") {
	// 	//level 4
	// 	user.KYCVerified = 4
	// }

	e := gc.DB.Omit(clause.Associations).Save(user).Error
	if e != nil {
		log.Printf("[UpdateUserKYCLevelFromSumsub] error saving user KYC Level to database  [%+v], %v\n", user.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	// e = gc.DB.Omit(clause.Associations).Save(&kycProgress).Error
	// if e != nil {
	// 	log.Printf("[UpdateUserKYCLevelFromSumsub] error saving user KYC progress to database  [%+v], %v\n", user.Username, e)

	// 	err = &tErrors.ErrorTemporaryServerError{}
	// 	return
	// }
	user.InvalidateUserCache(gc)
	return nil
}

func ProcessSumsubwebhook(sumSubwhInput *userModels.SumSubReviewResultInput, gc *sharedconfig.GlobalConfig) (err error) {
	// get user
	user, err := userDB.GetUser(sumSubwhInput.ExternalUserID, gc.DB, gc)
	if err != nil {
		log.Printf("[ProcessSumsubwebhook] error getting user: %v, %v\n", sumSubwhInput.ExternalUserID, err)
		err = &tErrors.ErrorTemporaryServerError{}
		return err
	}

	dbTx := gc.DB.Begin()
	defer dbTx.Rollback()

	err = SaveSumsubWebhook(sumSubwhInput, dbTx)
	if err != nil {
		return
	}
	if sumSubwhInput.ReviewResult.ReviewAnswer == "GREEN" {
		err = UpdateUserKYCLevelFromSumsub(&user, sumSubwhInput, dbTx, gc)
		if err != nil {
			return
		}
	}
	if sumSubwhInput.ReviewResult.ReviewAnswer == "RED" {
		err = ResetUserKYCProgressForSumsub(user.Username, sumSubwhInput.LevelName, dbTx)
		if err != nil {
			return
		}
	}

	dbTx.Commit()
	return nil
}
