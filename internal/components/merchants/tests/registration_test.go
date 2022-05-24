package users

import (
	"testing"
	usermodels "trovo-wallet-api/internal/components/users/models"
	userservices "trovo-wallet-api/internal/components/users/services"

	"github.com/asaskevich/govalidator"
)

// func setupTest() (*gorm.DB, error) {

// 	err := godotenv.Load(".env_test")
// 	if err != nil {
// 		log.Fatal("Error loading users .env file")
// 		return nil, err
// 	}

// 	var DB *gorm.DB
// 	DB, err = db.OpenDb(os.Getenv("DB_TYPE"), os.Getenv("DB_CONNECTION_STRING"), 10, 10)

// 	if err != nil {
// 		log.Fatalf("Error opening DB %s", err)
// 		return nil, err
// 	}

// 	return DB, nil

// }

func TestVerificationCode(t *testing.T) {

	var user usermodels.UserRegistrationInfo

	user.Email = "bogus@example.com"
	user.PublicIP = "1.2.3.4"
	user.PublicKey = "GCNS"
	user.Username = "userName"

	code1 := userservices.GenerateEmailVerificationCode(user, "mysalt")

	if !govalidator.IsInt(code1) {
		t.Fatalf("expected an int, but go %v", code1)
	}

	if len(code1) != userservices.NumberOfCharactersInVerificationCode {
		t.Fatalf("expected %d characters, but got %v characters", userservices.NumberOfCharactersInVerificationCode, len(code1))
	}

	code2 := userservices.GenerateEmailVerificationCode(user, "mysalt2")

	if !govalidator.IsInt(code2) {
		t.Fatalf("expected an int, but go %v", code2)
	}

	if len(code2) != userservices.NumberOfCharactersInVerificationCode {
		t.Fatalf("expected %d characters, but got %v characters", userservices.NumberOfCharactersInVerificationCode, len(code2))
	}

	//the two codes should not be the same

	if code1 == code2 {
		t.Fatalf("expected codes to be different but got the same code %s", code1)
	}

}
