package users

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	userModels "trovo-wallet-api/internal/components/users/models"
)

func UploadFileToS3(uploadedFile *multipart.FileHeader, uploadcred *userModels.UploadData) error {

	request, err := newfileUploadRequest(uploadcred, uploadedFile)
	if err != nil {
		log.Println("[UploadFileToS3]error generating upload request to s3", err)
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[UploadFileToS3] error uploaing file to s3", err)
		return err
	}
	log.Println("[UploadFileToS3] succeeded with code: ", resp.StatusCode)
	var bodyContent []byte

	resp.Body.Read(bodyContent)
	resp.Body.Close()
	log.Println(bodyContent)
	return nil

}
func newfileUploadRequest(uploadcred *userModels.UploadData, uploadedFile *multipart.FileHeader) (request *http.Request, err error) {

	blobfile, err := uploadedFile.Open()
	if err != nil {
		return nil, err
	}

	bf, _ := io.ReadAll(blobfile)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", uploadedFile.Filename)
	if err != nil {
		return nil, err
	}

	part.Write(bf)
	_ = writer.WriteField("key", uploadcred.Fields.Key)
	_ = writer.WriteField("bucket", uploadcred.Fields.Bucket)
	_ = writer.WriteField("Policy", uploadcred.Fields.Policy)
	_ = writer.WriteField("X-Amz-Date", uploadcred.Fields.XAmzDate)
	_ = writer.WriteField("X-Amz-Algorithm", uploadcred.Fields.XAmzAlgorithm)
	_ = writer.WriteField("X-Amz-Signature", uploadcred.Fields.XAmzSignature)
	_ = writer.WriteField("X-Amz-Security-Token", uploadcred.Fields.XAmzSecurityToken)
	_ = writer.WriteField("X-Amz-Credential", uploadcred.Fields.XAmzCredential)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uploadcred.URL, body)
}
