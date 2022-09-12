package users

import (
	"log"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"
)

func SaveSecretQuestions(user *userModels.User, answer userModels.UserSecretAnswer, gc *sharedconfig.GlobalConfig) error {
	var existingAnswer userModels.UserSecretAnswer
	e := gc.DB.Where("username = ?", user.Username).First(&existingAnswer).Error
	if e != nil {
		// possibly does not exist
		answer.Username = user.Username
		e := gc.DB.Save(&answer).Error
		if e != nil {
			log.Printf("[SaveSecretQuestions] error creating answers [%v]", e)
			return &tErrors.CustomError{Param: "id", Err: "error saving secret answers", ErrMessage: "Unable to save secret answers at this time"}
		}
		user.HasSecretQuestions = 1
		gc.DB.Save(user)
		return nil

	}

	//existing answer
	existingAnswer.A1 = answer.A1
	existingAnswer.Q1 = answer.Q1
	existingAnswer.A2 = answer.A2
	existingAnswer.Q2 = answer.Q2
	existingAnswer.A3 = answer.A3
	existingAnswer.Q3 = answer.Q3
	e = gc.DB.Save(&existingAnswer).Error
	if e != nil {
		log.Printf("[SaveSecretQuestions] error saving answers [%v]", e)
		return &tErrors.CustomError{Param: "id", Err: "error saving secret answers", ErrMessage: "Unable to save secret answers at this time"}
	}
	user.HasSecretQuestions = 1
	gc.DB.Save(user)
	return nil

}
