package validators

import (
	"testing"
	"trovo-wallet-payment-history-engine/internal/errors"
	"trovo-wallet-payment-history-engine/internal/validators"
)

func TestInvalidUsernames(t *testing.T) {

	usernames := []string{"", "he", "123456789012345678901234567890", "ef ok", "e_restaur", "you.me", "héllo", "1ab"}

	for _, username := range usernames {
		err := validators.ValidateUsernameFormat(username)

		if err == nil {
			t.Errorf("should have thrown an error %v", username)
		}
	}

}

func TestValidUsernames(t *testing.T) {

	usernames := []string{"abcd", "a123bber", "youme"}

	for _, username := range usernames {
		err := validators.ValidateUsernameFormat(username)

		if err != nil {
			t.Errorf("should not have thrown an error %v", username)
		}
	}
}

func TestInvalidAddresses(t *testing.T) {

	var publicKeys []string = []string{"bogus", "fake", "SCBTEKU7J6BL3LXUWVZCPBQNVU3DCI4SQ7YFQN3VYMUFPEDB334JVAVT"}

	for idx, publicKey := range publicKeys {

		err := validators.ValidateAddressFormat(publicKey)

		if err == nil {
			t.Errorf("publicKey should have thrown an error %s", publicKey)
			continue
		}

		var ex, ok = err.(*errors.ErrorInvalidAddress)

		if !ok {
			t.Errorf("Invalid Public Key exception was not thrown for item number %d, [%v]", idx, err)
			continue
		}

		expectedInvalidAddress := publicKey

		if ex.Address != expectedInvalidAddress {
			t.Errorf("expected [%s] but got [%s]", expectedInvalidAddress, ex.Address)
		}
	}

}

func TestValidAddress(t *testing.T) {

	var publicKeys []string = []string{"0xc4FE8226634b79a06e49f5DFdB02475E2eA8AF28", "0xD59AC3c804E6A19d0418CDb2BA96F3F02A1166ab"}

	for idx, publicKey := range publicKeys {

		err := validators.ValidateAddressFormat(publicKey)

		if err != nil {
			t.Errorf("public key validation not have thrown an error [%d], %v", idx, err)
			continue
		}

	}

}

func TestInvalidEmails(t *testing.T) {

	var emails []string = []string{"bogus", "123@2.net@", "@metoyou@u.com", "you@me@t.com"}

	for _, email := range emails {

		valid := validators.ValidateEmailFormat(email)

		if valid {
			t.Errorf("email should be invalid %v", email)
			continue
		}
	}

}

func TestValidEmails(t *testing.T) {

	var emails []string = []string{"johndoe@gmail.com", "johndoe@g.mail.com"}

	for _, email := range emails {

		valid := validators.ValidateEmailFormat(email)

		if !valid {
			t.Errorf("email should be valid %v", email)
			continue
		}
	}

}
