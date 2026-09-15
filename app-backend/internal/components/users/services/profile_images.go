package users

import (
	"fmt"
	"log"
	"mime/multipart"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm/clause"
)

func UploadProfilePicture(user *userModels.User, file multipart.File, fileNameWithExt string, gc *sharedconfig.GlobalConfig) (string, error) {
	var oldThumbnail string
	if user.ImageThumbnailURL != nil {
		oldThumbnail = *user.ImageThumbnailURL
	}
	newThumbnail, err := gc.FirebaseStorageUploader.UploadFile(file, fileNameWithExt, oldThumbnail)
	if err != nil {
		return "", err
	}

	//update the user thumbnail url
	url := fmt.Sprintf("https://storage.googleapis.com/%v/%v", gc.FirebaseStorageUploader.BucketName, newThumbnail)
	user.ImageThumbnailURL = &url
	e := gc.DB.Omit(clause.Associations).Save(user).Error
	if e != nil {
		log.Printf("[UploadProfilePicture]error saving profile picture for %v: %v\n", user.Username, e)
		return "", fmt.Errorf("error saving account profile picture url %v", url)
	}
	user.InvalidateUserCache(gc)
	owner, _ := userModels.Username(user.Username).GetFullUser(gc.DB, gc)
	// user = &owner
	if len(owner.ID) > 0 {
		user = &owner
	}
	return url, nil
}
