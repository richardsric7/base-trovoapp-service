package users

import (
	"log"
	"strings"
	bc "trovo-wallet-api/internal/blockchainalgofuncs"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"
)

func SaveUserSecretQuestions(user *userModels.User, answer userModels.UserSecretAnswer, gc *sharedconfig.GlobalConfig) error {
	var existingAnswer userModels.UserSecretAnswer
	// check if  questions where repeated
	if answer.Q1 == answer.Q2 || answer.Q1 == answer.Q3 || answer.Q2 == answer.Q3 {
		return &tErrors.CustomError{Param: "username",
			Err:        "error-cannot-have-duplicate-question",
			ErrMessage: "You cannot have duplicate question.",
		}
	}
	e := gc.DB.Where("username = ?", user.Username).First(&existingAnswer).Error
	if e != nil {
		// possibly does not exist
		answer.Username = user.Username

		{
			answer.A1 = bc.EncodeSha256(strings.ToLower(answer.A1))
			answer.A2 = bc.EncodeSha256(strings.ToLower(answer.A2))
			answer.A3 = bc.EncodeSha256(strings.ToLower(answer.A3))
		}

		e := gc.DB.Save(&answer).Error
		if e != nil {
			log.Printf("[SaveSecretQuestions] error creating answers [%v]", e)
			return &tErrors.CustomError{Param: "id", Err: "error saving secret answers", ErrMessage: "Unable to save secret answers at this time"}
		}
		if user.HasSecretQuestions == 0 {
			user.HasSecretQuestions = 1
			gc.DB.Save(user)
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
	e = gc.DB.Save(&existingAnswer).Error
	if e != nil {
		log.Printf("[SaveSecretQuestions] error saving answers [%v]\n", e)
		return &tErrors.CustomError{Param: "id", Err: "error saving secret answers", ErrMessage: "Unable to save secret answers at this time"}
	}
	if user.HasSecretQuestions == 0 {
		user.HasSecretQuestions = 1
		gc.DB.Save(user)
	}
	return nil

}

func GetUserSecretAnswers(user *userModels.User, gc *sharedconfig.GlobalConfig) (answer userModels.UserSecretAnswer, err error) {
	if user.HasSecretQuestions == 0 {
		return answer, &tErrors.CustomError{Param: "username",
			Err:        "error-no-secret answers exist for user",
			ErrMessage: "No secret answers yet",
		}
	}
	if e := gc.DB.Where("username = ?", user.Username).First(&answer).Error; e != nil {
		return answer, &tErrors.CustomError{Param: "username",
			Err:        "error-no-secret answers exist for user",
			ErrMessage: "No secret answers yet",
		}
	}
	return answer, nil
}

func GetSecretQuestions(user *userModels.User, gc *sharedconfig.GlobalConfig) (questions []userModels.SecretQuestion) {

	if e := gc.DB.Order("question ASC").Find(&questions).Error; e != nil {
		return make([]userModels.SecretQuestion, 0)
	}
	return questions
}

func ValidateSecretAnswers(user *userModels.User, answer userModels.UserSecretAnswer, gc *sharedconfig.GlobalConfig) bool {

	storedAnswers, err := GetUserSecretAnswers(user, gc)
	if err != nil {
		return false
	}
	if storedAnswers.A1 != bc.EncodeSha256(strings.ToLower(answer.A1)) || storedAnswers.A2 != bc.EncodeSha256(strings.ToLower(answer.A2)) || storedAnswers.A3 != bc.EncodeSha256(strings.ToLower(answer.A3)) {
		return false
	}

	return true
}
