package users

import (
	"fmt"
	"log"
	"mime/multipart"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUserClosedGroups(db *gorm.DB) (ucgs []userModels.UserClosedGroup) {
	ucgs = make([]userModels.UserClosedGroup, 0)
	db.Preload(clause.Associations).Order("closed_group_Id").Find(&ucgs)

	return
}

func GetClosedGroupByOwner(groupOwner string, db *gorm.DB) (cgs []userModels.ClosedGroup) {
	cgs = make([]userModels.ClosedGroup, 0)
	db.Preload(clause.Associations).Order("group_name").Where("group_owner = ?", groupOwner).Find(&cgs)

	return
}

func CreateClosedGroup(cgInput userModels.ClosedGroupJSONInput, groupUser *userModels.User, db *gorm.DB) (cg userModels.ClosedGroup, err error) {

	return
}

func UploadClosedGroupRegistrationDocument(groupOwner *userModels.User, file multipart.File, fileNameWithExt string, closedGroup *userModels.ClosedGroup, gc *sharedconfig.GlobalConfig) (string, error) {

	newThumbnail, err := gc.FirebaseStorageUploader.UploadFile(file, fileNameWithExt, "")
	if err != nil {
		return "", err
	}

	//update the thumbnail url
	url := fmt.Sprintf("https://storage.googleapis.com/%v/%v", gc.FirebaseStorageUploader.BucketName, newThumbnail)
	//check if document already saved and then retireve it:
	if len(closedGroup.RegistrationDocumentUrl) > 0 {
		//existing record match, update
		closedGroup.RegistrationDocumentUrl = url
		es := gc.DB.Save(closedGroup).Error
		if es != nil {

			log.Printf("[UploadClosedGroupRegistrationDocument]error saving existing document in database  [%v] for %v: %v\n", closedGroup, groupOwner.Username, es)
			return "", fmt.Errorf("error saving document %v", closedGroup.GroupName)

		}
	}

	groupOwner.InvalidateUserCache(gc)
	owner, err := userModels.Username(groupOwner.Username).GetFullUser(gc.DB, gc)
	if err == nil {
		if owner.Username == groupOwner.Username {
			groupOwner = &owner
		}

	}

	return url, nil
}
