package users

import (
	"log"
	"strings"
	bc "trovo-wallet-api/internal/blockchainalgofuncs"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveUserSecurityQuestions(user *userModels.User, answer userModels.UserSecurityAnswer, db *gorm.DB) error {
	var existingAnswer userModels.UserSecurityAnswer
	// check if  questions where repeated
	if answer.Q1 == answer.Q2 || answer.Q1 == answer.Q3 || answer.Q2 == answer.Q3 {
		return &tErrors.CustomError{Param: "username",
			Err:        "error-cannot-have-duplicate-question",
			ErrMessage: "You cannot have duplicate question.",
		}
	}
	e := db.Where("username = ?", user.Username).First(&existingAnswer).Error
	if e != nil {
		// possibly does not exist
		answer.Username = user.Username

		{
			answer.A1 = bc.EncodeSha256(strings.ToLower(answer.A1))
			answer.A2 = bc.EncodeSha256(strings.ToLower(answer.A2))
			answer.A3 = bc.EncodeSha256(strings.ToLower(answer.A3))
		}

		e := db.Save(&answer).Error
		if e != nil {
			log.Printf("[SaveUserSecurityQuestions] error creating answers [%v]", e)
			return &tErrors.CustomError{Param: "id", Err: "error saving security answers", ErrMessage: "Unable to save security answers at this time"}
		}
		if user.HasSecurityQuestions == 0 {
			user.HasSecurityQuestions = 1
			db.Omit(clause.Associations).Save(user)
		}

		return nil

	}

	//existing answer
	// existingAnswer.A1 = answer.A1
	{
		existingAnswer.A1 = bc.EncodeSha256(strings.ToLower(answer.A1))
		existingAnswer.A2 = bc.EncodeSha256(strings.ToLower(answer.A2))
		existingAnswer.A3 = bc.EncodeSha256(strings.ToLower(answer.A3))
	}
	existingAnswer.Q1 = answer.Q1
	// existingAnswer.A2 = answer.A2
	existingAnswer.Q2 = answer.Q2
	// existingAnswer.A3 = answer.A3
	existingAnswer.Q3 = answer.Q3
	e = db.Save(&existingAnswer).Error
	if e != nil {
		log.Printf("[SaveUserSecurityQuestions] error saving answers [%v]\n", e)
		return &tErrors.CustomError{Param: "id", Err: "error saving security answers", ErrMessage: "Unable to save security answers at this time"}
	}
	if user.HasSecurityQuestions == 0 {
		user.HasSecurityQuestions = 1
		db.Omit(clause.Associations).Save(user)
	}

	return nil

}

func GetUserSecurityAnswers(user *userModels.User, gc *sharedconfig.GlobalConfig) (answer userModels.UserSecurityAnswer, err error) {
	if user.HasSecurityQuestions == 0 {
		return answer, &tErrors.CustomError{Param: "username",
			Err:        "error-no-security-answers-exist-for-user",
			ErrMessage: "No security answers yet",
		}
	}
	if e := gc.DB.Where("username = ?", user.Username).First(&answer).Error; e != nil {
		return answer, &tErrors.CustomError{Param: "username",
			Err:        "error-no-security answers exist for user",
			ErrMessage: "No security answers yet",
		}
	}
	return answer, nil
}

func GetSecurityQuestions(user *userModels.User, gc *sharedconfig.GlobalConfig) (questions []userModels.SecurityQuestion) {

	if e := gc.DB.Order("question ASC").Find(&questions).Error; e != nil {
		return make([]userModels.SecurityQuestion, 0)
	}
	return questions
}

func ValidateSecurityAnswers(user *userModels.User, answer userModels.UserSecurityAnswer, gc *sharedconfig.GlobalConfig) bool {

	storedAnswers, err := GetUserSecurityAnswers(user, gc)
	if err != nil {
		return false
	}
	if storedAnswers.A1 != bc.EncodeSha256(strings.ToLower(answer.A1)) || storedAnswers.A2 != bc.EncodeSha256(strings.ToLower(answer.A2)) || storedAnswers.A3 != bc.EncodeSha256(strings.ToLower(answer.A3)) {
		return false
	}

	return true
}
