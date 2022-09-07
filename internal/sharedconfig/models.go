package sharedconfig

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"
	"trovo-wallet-api/internal/cache"

	cs "cloud.google.com/go/storage"
	"firebase.google.com/go/messaging"
	"firebase.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/stellar/go/clients/horizonclient"
	"gorm.io/gorm"
)

type GlobalConfig struct {
	DynamicLinkServiceURLChan chan string
	PushNotificationClient    *messaging.Client
	FirebaseStorageUploader   *ClientUploader
	PNSContext                context.Context
	RedisCache                *cache.RedisCache
	DB                        *gorm.DB
	RoachDB                   *gorm.DB
	BantuExpansionClient      *horizonclient.Client
	BantuNetworkPassphrase    string
}

type ClientUploader struct {
	Client     *storage.Client
	ProjectID  string
	BucketName string
	UploadPath string
}

func (c *ClientUploader) UploadFile(fileInput multipart.File, fileName string) error {

	// create an id
	id := uuid.New()
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	// wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	sh, err := c.Client.Bucket(c.BucketName)
	if err != nil {
		//no bucket with that name exists
		log.Printf("[UploadFile] error getting bucket handle %v: %v\n", c.BucketName, err)
		return fmt.Errorf("error getting bucket handle %v: %v", c.BucketName, err)
	}

	_, err = sh.Attrs(ctx)

	if err != nil {
		//no bucket with that name exists, create it
		err := sh.Create(ctx, c.ProjectID, nil)
		if err != nil {
			log.Printf("[UploadFile] error creating bucket handle %v: %v\n", c.BucketName, err)
			return fmt.Errorf("error creating bucket handle %v: %v", c.BucketName, err)
		}
	}
	object := sh.Object(c.UploadPath + fileName)

	{
		//check if object already exists and delete it.
		if _, err := object.Attrs(ctx); err == nil {
			object.Delete(ctx)
			//set the object again
			object = sh.Object(c.UploadPath + "/" + fileName)
		}
	}
	writer := object.NewWriter(ctx)

	//Set the attribute
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	defer writer.Close()

	if _, err := io.Copy(writer, fileInput); err != nil {
		log.Printf("[UploadFile] error uploading file %v: %v\n", fileName, err)
		return fmt.Errorf("error uploading file %v: %v", fileName, err)
	}

	if err := object.ACL().Set(context.Background(), cs.AllUsers, cs.RoleReader); err != nil {
		log.Printf("[UploadFile] error setting file permission %v: %v\n", fileName, err)
		return fmt.Errorf("error setting file permission %v: %v", fileName, err)
	}
	return nil
}
