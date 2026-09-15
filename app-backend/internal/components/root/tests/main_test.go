package root

import (
	"log"
	"testing"
	rootservices "trovo-wallet-api/internal/components/root/services"

	"github.com/joho/godotenv"
)

func setupTest() {
	err := godotenv.Load(".env_test")
	if err != nil {
		log.Fatal("Error loading root .env file")
	}
}

func TestOrganisationName(t *testing.T) {
	setupTest()
	organisationName := rootservices.GetRootInfo()

	if len(organisationName.Organisation) == 0 {
		t.Errorf("organisation name is empty")
	}

}
