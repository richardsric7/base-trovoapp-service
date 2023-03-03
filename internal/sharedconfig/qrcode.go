package sharedconfig

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	cs "cloud.google.com/go/storage"
	"github.com/google/uuid"
)

func (c *ClientUploader) UploadQrCode(fileInput *os.File, fileName, imageThumbnailURL string) (string, error) {

	// create an id
	id := uuid.NewString()
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	sh, err := c.Client.Bucket(c.BucketName)
	if err != nil {
		//no bucket with that name exists
		log.Printf("[UploadQrCode] error getting bucket handle %v: %v\n", c.BucketName, err)
		return "", fmt.Errorf("error getting bucket handle %v: %v", c.BucketName, err)
	}

	if _, err := sh.Attrs(ctx); err != nil {
		//no bucket with that name exists, create it
		rules := make([]cs.ACLRule, 0)
		rules = append(rules, cs.ACLRule{Entity: "allUsers", Role: "READER"})
		err := sh.Create(ctx, c.ProjectID, &cs.BucketAttrs{ACL: rules})
		if err != nil {
			log.Printf("[UploadQrCode] error creating bucket handle %v: %v\n", c.BucketName, err)
			return "", fmt.Errorf("error creating bucket handle %v: %v", c.BucketName, err)
		}
	}
	newImageThumbnailName := c.UploadPath + "/" + id + fileName
	object := sh.Object(newImageThumbnailName)

	if len(imageThumbnailURL) > 3 {
		// ImageThumbnailURL is full https url. strip the unnecessary portion
		oldName := strings.ReplaceAll(imageThumbnailURL, fmt.Sprintf("https://storage.googleapis.com/%v/", c.BucketName), "")
		oldObject := sh.Object(oldName)
		//check if object already exists and delete it.
		if _, err := oldObject.Attrs(ctx); err == nil {
			oldObject.Delete(ctx)
		}
	}
	writer := object.NewWriter(ctx)

	//Set the attribute
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}
	defer writer.Close()

	if _, err := io.Copy(writer, fileInput); err != nil {
		log.Printf("[UploadQrCode] error uploading file %v: %v\n", newImageThumbnailName, err)
		return "", fmt.Errorf("error uploading file %v: %v", newImageThumbnailName, err)
	}

	return newImageThumbnailName, nil
}
