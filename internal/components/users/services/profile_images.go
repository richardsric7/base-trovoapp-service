package users

import (
	"fmt"
	"log"
	"mime/multipart"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

func UploadProfilePicture(user *userModels.User, file multipart.File, fileNameWithExt string, gc *sharedconfig.GlobalConfig) (string,error) {

	err := gc.FirebaseStorageUploader.UploadFile(file, fileNameWithExt)
	if err != nil {
		return "",err
	}

	//update the user thumbnail url
	url := fmt.Sprintf("https://storage.googleapis.com/%v/%v/%v", gc.FirebaseStorageUploader.BucketName, gc.FirebaseStorageUploader.BucketName, fileNameWithExt)
	user.ImageThumbnailURL = &url
	e := gc.DB.Save(user).Error
	if e != nil {
		log.Printf("[UploadProfilePicture]error saving profile picture for %v: %v\n", user.Username, e)
		return "", fmt.Errorf("error saving account profile picture url %v", url)
	}
	return url, nil
}
