package users

import (
	"testing"
	usermodels "trovo-wallet-api/internal/components/users/models"
	userservices "trovo-wallet-api/internal/components/users/services"
	errors "trovo-wallet-api/internal/errors"
)

type MissingParameterTestItem struct {
	expectedError string
	userInfo      usermodels.UserRegistrationInfo
}

func TestValidateUserRegistrationInfoErrorMissingParameter(t *testing.T) {

	var missingParameterTestItems []MissingParameterTestItem

	missingParameterTestItems = append(missingParameterTestItems, MissingParameterTestItem{expectedError: "username", userInfo: usermodels.UserRegistrationInfo{}})
	missingParameterTestItems = append(missingParameterTestItems, MissingParameterTestItem{expectedError: "publicKey", userInfo: usermodels.UserRegistrationInfo{Username: "bogus"}})
	missingParameterTestItems = append(missingParameterTestItems, MissingParameterTestItem{
		expectedError: "email",
		userInfo:      usermodels.UserRegistrationInfo{Username: "bogus", PublicKey: "GCNS"}})

	for idx, missingParameterTestItem := range missingParameterTestItems {

		err := userservices.ValidateUserRegistrationInfo(missingParameterTestItem.userInfo)

		if err == nil {
			t.Errorf("user registration should have thrown an error")
			continue
		}

		var ex, ok = err.(*errors.ErrorMissingParameter)

		if !ok {
			t.Fatalf("missing parameter exception was not thrown for item number %d", idx)
		}

		expectedMissingParameter := missingParameterTestItem.expectedError

		if ex.Parameter != expectedMissingParameter {
			t.Errorf("expected [%s] but got [%s]", expectedMissingParameter, ex.Parameter)
		}
	}

}
