package users

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	// "github.com/k0kubun/pp"

	"github.com/pkg/errors"
)

const URL = "https://api.sumsub.com"

// const SumsubAppToken = "sbx:6L6rqHEtRVvBKKt7P1A03k2x.h6OsEOXWpyaXAjvBVNnx3ccXNGTBLHkw" // Example: sbx:uY0CgwELmgUAEyl4hNWxLngb.0WSeQeiYny4WEqmAALEAiK2qTC96fBad
// const SumsubSecretKey = "EraepapR4Grr2vI1eZWtTkFDhbhsC5EI"                             // Example: Hej2ch71kG2kTd1iIUDZFNsO5C1lh5Gq
//Please don't forget to change token and secret key values to production ones when switching to production

func GetSumsubIndividualApplicantKYC(user *userModels.User, levelName string, gc *sharedconfig.GlobalConfig) (token string, applicant userModels.SumsubApplicant, err error) {

	if user.Corporate != 0 {
		return token, applicant, &tErrors.CustomError{
			Param: "countryCode",
			Err:   "only individuals are allowed to use this service endpoint",
		}
	}
	var externalUserId = user.Username

	//  applicant = userModels.SumsubApplicant{}
	var fixedInfo = userModels.SumsubInfo{}
	if user.CountryCode != nil {
		fixedInfo.Country = *user.CountryCode
	}

	fixedInfo.FirstName = user.FirstName
	fixedInfo.LastName = *user.LastName

	// applicant.ID = user.Username
	applicant.FixedInfo = fixedInfo
	applicant.ExternalUserID = externalUserId

	// https://docs.sumsub.com/reference/create-applicant
	applicant, err = CreateApplicant(applicant, levelName, gc)
	if err != nil {
		return
	}
	// https://docs.sumsub.com/reference/add-id-documents
	// idDoc := AddDocument(applicant.ID)
	// fmt.Println(idDoc)

	// https://docs.sumsub.com/reference/get-applicant-data
	applicant, err = GetApplicantInfo(applicant, gc)
	if err != nil {
		return
	}
	// https://docs.sumsub.com/reference/generate-access-token-query
	// accessToken := GenerateAccessToken(applicant, levelName,gc)
	accessToken, err := GenerateAccessToken(applicant, levelName, gc)
	if err != nil {
		return
	}

	return accessToken.Token, applicant, nil
	// fmt.Println(accessToken.Token)
}

func GenerateAccessToken(applicant userModels.SumsubApplicant, levelName string, gc *sharedconfig.GlobalConfig) (userModels.SumsubAccessToken, error) {
	var token userModels.SumsubAccessToken
	b, err := _makeSumsubRequest("/resources/accessTokens?userId="+applicant.ExternalUserID+"&levelName="+levelName,
		"POST",
		"application/json",
		[]byte(""), gc)
	if err != nil {
		log.Printf("[GenerateAccessToken] error generating acess token: %v\n", err)
		return token, err
	}
	err = os.WriteFile("generateAccessToken.json", b, 0777)
	if err != nil {
		log.Printf("[GenerateAccessToken] error writing  acess token json file: %v\n", err)
		return token, &tErrors.ErrorTemporaryServerError{}
	}
	err = json.Unmarshal(b, &token)
	if err != nil {
		log.Printf("[GenerateAccessToken] error unmarshaling acess token: %v\n", err)
		return token, &tErrors.ErrorTemporaryServerError{}
	}

	return token, err
}

func CreateApplicant(applicant userModels.SumsubApplicant, levelName string, gc *sharedconfig.GlobalConfig) (userModels.SumsubApplicant, error) {
	postBody, _ := json.Marshal(applicant)

	var ac userModels.SumsubApplicant
	log.Printf("[SumsubCreateApplicant] Applicant to be created: %+v\n", applicant)
	b, err := _makeSumsubRequest(
		"/resources/applicants?levelName="+levelName,
		"POST",
		"application/json",
		postBody, gc)

	if err != nil {
		log.Printf("[CreateApplicant] error creating applicant: %v\n", err)
		return ac, err
	}
	// pp.Println(string(b))
	os.WriteFile("createApplicant.json", b, 0777)

	err = json.Unmarshal(b, &ac)
	if err != nil {
		log.Printf("[CreateApplicant] error unmarshaling applicant: %v\n", err)
		return ac, &tErrors.ErrorTemporaryServerError{}
	}
	log.Printf("[SumsubCreateApplicant] Created Applicant: %+v\n", applicant)

	return ac, nil
}

func GetApplicantInfo(applicant userModels.SumsubApplicant, gc *sharedconfig.GlobalConfig) (userModels.SumsubApplicant, error) {
	p := fmt.Sprintf("/resources/applicants/%s/one", applicant.ID)
	b, err := _makeSumsubRequest(
		p,
		"GET",
		"application/json",
		nil,
		gc)

	if err != nil {
		log.Printf("[CreateApplicant] error creating applicant: %v\n", err)
		return applicant, err
	}
	os.WriteFile("getApplicant.json", b, 0777)

	var r userModels.SumsubApplicant
	err = json.Unmarshal(b, &r)
	if err != nil {
		log.Printf("[CreateApplicant] error unmarshaling applicantINFO: %v\n", err)
		return applicant, &tErrors.ErrorTemporaryServerError{}
	}
	// pp.Println(r)

	return r, nil
}

func AddDocument(applicantId string, gc *sharedconfig.GlobalConfig) userModels.SumsubIdDoc {
	file, err := os.Open("resources/images/sumsub-logo.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	meta, err := json.Marshal(userModels.SumsubIdDoc{
		IdDocType: "PASSPORT",
		Country:   "GBR",
	})

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	var fw io.Writer
	if fw, err = w.CreateFormFile("content", file.Name()); err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		log.Fatal(err)
	}

	if fw, err = w.CreateFormField("metadata"); err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(fw, strings.NewReader(string(meta))); err != nil {
		log.Fatal(err)
	}
	w.Close()

	resp, err := _makeSumsubRequest(
		"/resources/applicants/"+applicantId+"/info/idDoc",
		"POST",
		w.FormDataContentType(),
		b.Bytes(),
		gc,
	)

	var doc userModels.SumsubIdDoc
	err = json.Unmarshal(resp, &doc)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

// X-App-Token - an App Token that you generate in our dashboard
// X-App-Access-Sig - signature of the request in the hex format (see below)
// X-App-Access-Ts - number of seconds since Unix Epoch in UTC
func _makeSumsubRequest(path, method, contentType string, body []byte, gc *sharedconfig.GlobalConfig) ([]byte, error) {

	config, err := GetKYCConfigByServiceProvider("sumsub", gc)
	if err != nil {
		log.Printf("[VerifySumSubWebhook] error getting service config: %v\n", err)
		err = &tErrors.ErrorTemporaryServerError{}
		return nil, err

	}

	request, err := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	ts := fmt.Sprintf("%d", time.Now().Unix())

	request.Header.Add("X-App-Token", config.Token)

	request.Header.Add("X-App-Access-Sig", _sign(ts, config.SecretKey, method, path, &body))
	request.Header.Add("X-App-Access-Ts", ts)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", contentType)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return b, nil
}

func _sign(ts string, secret string, method string, path string, body *[]byte) string {
	hash := hmac.New(sha256.New, []byte(secret))
	data := []byte(ts + method + path)

	if body != nil {
		data = append(data, *body...)
	}

	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
