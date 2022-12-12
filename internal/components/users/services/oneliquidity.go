package users

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

func GetWithdrawalNetworks(currency string, gc *sharedconfig.GlobalConfig) (wdlNetworks []userModels.WithdrawalNetwork, err error) {
	var wdlNetworksResp userModels.CryptoWithdrawalNetworksResponse
	client := http.DefaultClient
	//get 'https://sandbox-api.oneliquidity.technology/wallets/v1/withdrawal/networks?currency=BTC'
	cacheKey := fmt.Sprintf("wallets/v1/withdrawal/networks?currency=%s", currency)
	// url := fmt.Sprintf("%s/%s?currency=%s", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/withdrawal/networks", currency)
	url := "https://sandbox-api.oneliquidity.technology/wallets/v1/withdrawal/networks?currency=BTC"
	{
		ok, rawData := gc.RedisCache.GetCachedResultRaw(cacheKey)
		if ok {
			json.Unmarshal(rawData, &wdlNetworksResp)
			return
		}
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetWithdrawalNetworks] error sending request:", err)
		return
	}
	log.Println("[GetWithdrawalNetworks] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&wdlNetworksResp); err != nil {
		log.Println("[GetWithdrawalNetworks] error decoding response for wdl networks id:", err)
		return
	}

	gc.RedisCache.StoreResultToCacheRaw(cacheKey, wdlNetworksResp, 1000)
	wdlNetworks = wdlNetworksResp.Data
	return wdlNetworks, nil

}

func ComplianceStartNewVerification(firstName, lastName string) (verificationID string, err error) {

	complianceAccountID := os.Getenv("COMPLIANCE_ACCOUNT_ID")
	vr := userModels.VerificationRequest{
		AccountID: complianceAccountID,
		FirstName: firstName,
		LastName:  lastName,
	}
	jbody, err := json.Marshal(vr)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	//get upload credentials
	url := fmt.Sprintf("%s/%s?", os.Getenv("ONELIQUIDITY_BASE_URL"), "compliance/v1/verification")
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[ComplianceStartNewVerification] error starting new verification:", err)
		return "", err
	}
	log.Println("[ComplianceStartNewVerification] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	var veriResp userModels.VerificationResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&veriResp); err != nil {
		log.Println("[ComplianceStartNewVerification] error decoding response for verification id:", err)
		return "", err
	}

	return veriResp.Data.VerificationID, nil

}

func GetProofOfresidencyCred(user *userModels.User) (fMCred userModels.ProofOfResidenceCred, verificationID string, err error) {
	if user.Corporate == 1 || user.LastName == nil {
		err = &tErrors.CustomError{
			Param:      "user",
			Err:        "error-not-available-for-corporate",
			ErrMessage: "Service not available for corporate accounts",
		}
		return
	}
	verificationID, err = ComplianceStartNewVerification(user.FirstName, *user.LastName)
	if err != nil {
		return fMCred, "", err
	}

	url := fmt.Sprintf("%s/compliance/v1/kyc/verification/signed-url?verificationId=%s&uploadDocType=proof_of_residency", os.Getenv("ONELIQUIDITY_BASE_URL"), verificationID)

	client := &http.Client{}
	//get upload credentials

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetProofOfresidencyCred] error starting new verification:", err)
		return fMCred, verificationID, err
	}
	log.Println("[GetProofOfresidencyCred] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	// var fmResp userModels.VerificationResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fMCred); err != nil {
		log.Println("[GetProofOfresidencyCred] error decoding proof of residency response:", err)
		return fMCred, verificationID, err
	}

	return fMCred, verificationID, nil

}
func GetFacematchPassportCred(user *userModels.User) (fMCred userModels.FaceMatchVerificationCred, verificationID string, err error) {
	if user.Corporate == 1 || user.LastName == nil {
		err = &tErrors.CustomError{
			Param:      "user",
			Err:        "error-not-available-for-corporate",
			ErrMessage: "Service not available for corporate accounts",
		}
		return
	}
	verificationID, err = ComplianceStartNewVerification(user.FirstName, *user.LastName)
	if err != nil {
		return fMCred, "", err
	}

	url := fmt.Sprintf("%s/compliance/v1/kyc/verification/signed-url?verificationId=%s&uploadDocType=facematch&sourceDocType=passport", os.Getenv("ONELIQUIDITY_BASE_URL"), verificationID)

	client := &http.Client{}
	//get upload credentials

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetFacematchPassportCred] error starting new verification:", err)
		return fMCred, verificationID, err
	}
	log.Println("[GetFacematchPassportCred] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	// var fmResp userModels.VerificationResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fMCred); err != nil {
		log.Println("[GetFacematchPassportCred] error decoding facematch response:", err)
		return fMCred, verificationID, err
	}

	return fMCred, verificationID, nil

}

func GetFacematchNationalIDCred(user *userModels.User) (fMCred userModels.FaceMatchVerificationCred, verificationID string, err error) {
	if user.Corporate == 1 || user.LastName == nil {
		err = &tErrors.CustomError{
			Param:      "user",
			Err:        "error-not-available-for-corporate",
			ErrMessage: "Service not available for corporate accounts",
		}
		return
	}
	verificationID, err = ComplianceStartNewVerification(user.FirstName, *user.LastName)
	if err != nil {
		return fMCred, "", err
	}

	url := fmt.Sprintf("%s/compliance/v1/kyc/verification/signed-url?verificationId=%s&uploadDocType=facematch&sourceDocType=national_id", os.Getenv("ONELIQUIDITY_BASE_URL"), verificationID)

	client := &http.Client{}
	//get upload credentials

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetFacematchPassportCred] error starting new verification:", err)
		return fMCred, verificationID, err
	}
	log.Println("[GetFacematchPassportCred] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	// var fmResp userModels.VerificationResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fMCred); err != nil {
		log.Println("[GetFacematchPassportCred] error decoding facematch response:", err)
		return fMCred, verificationID, err
	}

	return fMCred, verificationID, nil

}

func GetFacematchDrivingLicenseCred(user *userModels.User) (fMCred userModels.FaceMatchVerificationCred, verificationID string, err error) {
	if user.Corporate == 1 || user.LastName == nil {
		err = &tErrors.CustomError{
			Param:      "user",
			Err:        "error-not-available-for-corporate",
			ErrMessage: "Service not available for corporate accounts",
		}
		return
	}
	verificationID, err = ComplianceStartNewVerification(user.FirstName, *user.LastName)
	if err != nil {
		return fMCred, "", err
	}

	url := fmt.Sprintf("%s/compliance/v1/kyc/verification/signed-url?verificationId=%s&uploadDocType=facematch&sourceDocType=driving_license", os.Getenv("ONELIQUIDITY_BASE_URL"), verificationID)

	client := &http.Client{}
	//get upload credentials

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetFacematchPassportCred] error starting new verification:", err)
		return fMCred, verificationID, err
	}
	log.Println("[GetFacematchPassportCred] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	// var fmResp userModels.VerificationResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fMCred); err != nil {
		log.Println("[GetFacematchPassportCred] error decoding facematch response:", err)
		return fMCred, verificationID, err
	}

	return fMCred, verificationID, nil

}

func StartFacematchForPassport(user *userModels.User, selfieVideo, documentPicture *multipart.FileHeader, gc *sharedconfig.GlobalConfig) (err error) {
	// , gc *sharedconfig.GlobalConfig
	kycdData, err := user.GetKycData(gc)

	if err != nil {
		return errors.New("user country not set")
	}
	fmCred, verificationID, err := GetFacematchPassportCred(user)

	if err != nil {
		return err
	}

	sf := strings.Split(selfieVideo.Filename, ".")
	selfieExtension := fmt.Sprintf(".%s", sf[len(sf)-1])

	pf := strings.Split(documentPicture.Filename, ".")
	documentExtension := fmt.Sprintf(".%s", pf[len(pf)-1])

	// update credentials to db
	doc := userModels.FacematchPassport{
		ID:                uuid.NewString(),
		UserID:            user.ID,
		VerificationID:    verificationID,
		SelfieVideoFormat: mime.TypeByExtension(selfieExtension),
		PictureFormat:     mime.TypeByExtension(documentExtension),
		SelfieKey:         fmCred.Data.SelfieVideo.Fields.Key,
		DocumentKey:       fmCred.Data.Passport.Fields.Key,
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	e := dbTX.Create(&doc).Error
	if e != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	//uploadVideo
	err = UploadFileToS3(selfieVideo, &fmCred.Data.SelfieVideo)
	if err != nil {
		return err
	}
	//uploadPassport
	err = UploadFileToS3(documentPicture, &fmCred.Data.Passport)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/compliance/v1/kyc/facematch", os.Getenv("ONELIQUIDITY_BASE_URL"))

	client := &http.Client{}
	// get upload credentials
	fmrequest := userModels.StartDocumentRequest{
		VerificationID: verificationID,
	}
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "selfie_video",
		MimeType: mime.TypeByExtension(selfieExtension),
		Key:      fmCred.Data.SelfieVideo.Fields.Key,
	})
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "passport",
		MimeType: mime.TypeByExtension(documentExtension),
		Key:      fmCred.Data.Passport.Fields.Key,
		Side:     "front",
	})
	jbody, err := json.Marshal(fmrequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[StartFacematchForPassport] error starting new verification:", err)
		return err
	}
	log.Println("[StartFacematchForPassport] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	var fmResp userModels.OKResponse
	// Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fmResp); err != nil {
		log.Println("[StartFacematchForPassport] error decoding start facematch for passport response:", err)
		return err
	}

	dbTX.Commit()
	return nil

}

func StartFacematchForDrivingLicense(user *userModels.User, selfieVideo, documentPicture *multipart.FileHeader, gc *sharedconfig.GlobalConfig) (err error) {
	// , gc *sharedconfig.GlobalConfig
	kycdData, err := user.GetKycData(gc)

	if err != nil {
		return errors.New("user country not set")
	}

	fmCred, verificationID, err := GetFacematchDrivingLicenseCred(user)

	if err != nil {
		return err
	}

	sf := strings.Split(selfieVideo.Filename, ".")
	selfieExtension := fmt.Sprintf(".%s", sf[len(sf)-1])

	pf := strings.Split(documentPicture.Filename, ".")
	documentExtension := fmt.Sprintf(".%s", pf[len(pf)-1])

	// update credentials to db
	doc := userModels.FacematchDrivingLicense{
		ID:                uuid.NewString(),
		UserID:            user.ID,
		VerificationID:    verificationID,
		SelfieVideoFormat: mime.TypeByExtension(selfieExtension),
		PictureFormat:     mime.TypeByExtension(documentExtension),
		SelfieKey:         fmCred.Data.SelfieVideo.Fields.Key,
		DocumentKey:       fmCred.Data.Passport.Fields.Key,
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	e := dbTX.Create(&doc).Error
	if e != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	//uploadVideo
	err = UploadFileToS3(selfieVideo, &fmCred.Data.SelfieVideo)
	if err != nil {
		return err
	}
	//uploadPassport
	err = UploadFileToS3(documentPicture, &fmCred.Data.Passport)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/compliance/v1/kyc/facematch", os.Getenv("ONELIQUIDITY_BASE_URL"))

	client := &http.Client{}
	// get upload credentials
	fmrequest := userModels.StartDocumentRequest{
		VerificationID: verificationID,
	}
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "selfie_video",
		MimeType: mime.TypeByExtension(selfieExtension),
		Key:      fmCred.Data.SelfieVideo.Fields.Key,
	})
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "driving_license",
		MimeType: mime.TypeByExtension(documentExtension),
		Key:      fmCred.Data.Passport.Fields.Key,
		Side:     "front",
	})
	jbody, err := json.Marshal(fmrequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[StartFacematchForDrivingLicense] error starting new verification:", err)
		return err
	}
	log.Println("[StartFacematchForDrivingLicense] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	var fmResp userModels.OKResponse
	// Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fmResp); err != nil {
		log.Println("[StartFacematchForDrivingLicense] error decoding start facematch for passport response:", err)
		return err
	}
	dbTX.Commit()
	return nil

}

func StartFacematchForNationalID(user *userModels.User, selfieVideo, documentPicture *multipart.FileHeader, gc *sharedconfig.GlobalConfig) (err error) {
	// , gc *sharedconfig.GlobalConfig
	kycdData, err := user.GetKycData(gc)

	if err != nil {
		return errors.New("user country not set")
	}
	fmCred, verificationID, err := GetFacematchNationalIDCred(user)

	if err != nil {
		return err
	}

	sf := strings.Split(selfieVideo.Filename, ".")
	selfieExtension := fmt.Sprintf(".%s", sf[len(sf)-1])

	pf := strings.Split(documentPicture.Filename, ".")
	documentExtension := fmt.Sprintf(".%s", pf[len(pf)-1])

	// update credentials to db
	doc := userModels.FacematchNationalID{
		ID:                uuid.NewString(),
		UserID:            user.ID,
		VerificationID:    verificationID,
		SelfieVideoFormat: mime.TypeByExtension(selfieExtension),
		PictureFormat:     mime.TypeByExtension(documentExtension),
		SelfieKey:         fmCred.Data.SelfieVideo.Fields.Key,
		DocumentKey:       fmCred.Data.Passport.Fields.Key,
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	e := dbTX.Create(&doc).Error
	if e != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	//uploadVideo
	err = UploadFileToS3(selfieVideo, &fmCred.Data.SelfieVideo)
	if err != nil {
		return err
	}
	//uploadPassport
	err = UploadFileToS3(documentPicture, &fmCred.Data.Passport)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/compliance/v1/kyc/facematch", os.Getenv("ONELIQUIDITY_BASE_URL"))

	client := &http.Client{}
	// get upload credentials
	fmrequest := userModels.StartDocumentRequest{
		VerificationID: verificationID,
	}
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "selfie_video",
		MimeType: mime.TypeByExtension(selfieExtension),
		Key:      fmCred.Data.SelfieVideo.Fields.Key,
	})
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "national_id",
		MimeType: mime.TypeByExtension(documentExtension),
		Key:      fmCred.Data.Passport.Fields.Key,
		Side:     "front",
	})
	jbody, err := json.Marshal(fmrequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[StartFacematchForNationalID] error starting new verification:", err)
		return err
	}
	log.Println("[StartFacematchForNationalID] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	var fmResp userModels.OKResponse
	// Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fmResp); err != nil {
		log.Println("[StartFacematchForNationalID] error decoding start facematch for passport response:", err)
		return err
	}
	dbTX.Commit()
	return nil

}

func StartGovernmentIDCheckForProofOfResidency(user *userModels.User, documentPicture *multipart.FileHeader, gc *sharedconfig.GlobalConfig) (err error) {
	// , gc *sharedconfig.GlobalConfig
	kycdData, err := user.GetKycData(gc)

	if err != nil {
		return errors.New("user country not set")
	}
	fmCred, verificationID, err := GetProofOfresidencyCred(user)

	if err != nil {
		return err
	}

	pf := strings.Split(documentPicture.Filename, ".")
	pictureExtension := fmt.Sprintf(".%s", pf[len(pf)-1])

	// update credentials to db
	doc := userModels.GovermentIDProofOfResidency{
		ID:             uuid.NewString(),
		UserID:         user.ID,
		VerificationID: verificationID,
		PictureFormat:  mime.TypeByExtension(pictureExtension),
		DocumentKey:    fmCred.Data.PresignedURL.ProofOfResidency.Fields.Key,
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	e := dbTX.Create(&doc).Error
	if e != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}

	//uploadPassport
	err = UploadFileToS3(documentPicture, &fmCred.Data.PresignedURL.ProofOfResidency)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/compliance/v1/kyc/gov-id", os.Getenv("ONELIQUIDITY_BASE_URL"))

	client := &http.Client{}
	// get upload credentials
	fmrequest := userModels.StartDocumentRequest{
		VerificationID: verificationID,
	}

	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "proof_of_residency",
		MimeType: mime.TypeByExtension(pictureExtension),
		Key:      fmCred.Data.PresignedURL.ProofOfResidency.Fields.Key,
		Side:     "front",
	})
	jbody, err := json.Marshal(fmrequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[StartGovernmentIDCheckForProofOfResidency] error starting new verification:", err)
		return err
	}
	log.Println("[StartGovernmentIDCheckForProofOfResidency] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	var fmResp userModels.OKResponse
	// Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fmResp); err != nil {
		log.Println("[StartGovernmentIDCheckForProofOfResidency] error decoding start goverment id check response:", err)
		return err
	}
	dbTX.Commit()
	return nil

}
