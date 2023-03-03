package sharedconfig

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"sync"
	"time"
	"trovo-wallet-api/internal/cache"

	cs "cloud.google.com/go/storage"
	"firebase.google.com/go/messaging"
	"firebase.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
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
	ChannelAccounts           chan *keypair.Full
	InUseChannelAccounts      map[string]*keypair.Full
	Mutex                     sync.Mutex
}

type ClientUploader struct {
	Client     *storage.Client
	ProjectID  string
	BucketName string
	UploadPath string
}

func (c *ClientUploader) UploadFile(fileInput multipart.File, fileName, imageThumbnailURL string) (string, error) {

	// create an id
	id := uuid.New()
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	sh, err := c.Client.Bucket(c.BucketName)
	if err != nil {
		//no bucket with that name exists
		log.Printf("[UploadFile] error getting bucket handle %v: %v\n", c.BucketName, err)
		return "", fmt.Errorf("error getting bucket handle %v: %v", c.BucketName, err)
	}

	_, err = sh.Attrs(ctx)

	if err != nil {
		//no bucket with that name exists, create it
		rules := make([]cs.ACLRule, 0)
		rules = append(rules, cs.ACLRule{Entity: "allUsers", Role: "READER"})
		err := sh.Create(ctx, c.ProjectID, &cs.BucketAttrs{ACL: rules})
		if err != nil {
			log.Printf("[UploadFile] error creating bucket handle %v: %v\n", c.BucketName, err)
			return "", fmt.Errorf("error creating bucket handle %v: %v", c.BucketName, err)
		}
	}
	newImageThumbnailName := c.UploadPath + "/" + id.String() + fileName
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
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	defer writer.Close()

	if _, err := io.Copy(writer, fileInput); err != nil {
		log.Printf("[UploadFile] error uploading file %v: %v\n", newImageThumbnailName, err)
		return "", fmt.Errorf("error uploading file %v: %v", newImageThumbnailName, err)
	}

	return newImageThumbnailName, nil
}

func (c *ClientUploader) SaveQrCodeAsFileToCloud(fileInput *os.File, fileName, imageThumbnailURL string) (string, error) {

	// create an id
	id := uuid.New()
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	sh, err := c.Client.Bucket(c.BucketName)
	if err != nil {
		//no bucket with that name exists
		log.Printf("[SaveQrCodeAsFileToCloud] error getting bucket handle %v: %v\n", c.BucketName, err)
		return "", fmt.Errorf("error getting bucket handle %v: %v", c.BucketName, err)
	}

	_, err = sh.Attrs(ctx)

	if err != nil {
		//no bucket with that name exists, create it
		rules := make([]cs.ACLRule, 0)
		rules = append(rules, cs.ACLRule{Entity: "allUsers", Role: "READER"})
		err := sh.Create(ctx, c.ProjectID, &cs.BucketAttrs{ACL: rules})
		if err != nil {
			log.Printf("[SaveQrCodeAsFileToCloud] error creating bucket handle %v: %v\n", c.BucketName, err)
			return "", fmt.Errorf("error creating bucket handle %v: %v", c.BucketName, err)
		}
	}
	newImageThumbnailName := c.UploadPath + "/" + id.String() + fileName
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
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	defer writer.Close()

	if _, err := io.Copy(writer, fileInput); err != nil {
		log.Printf("[SaveQrCodeAsFileToCloud] error uploading file %v: %v\n", newImageThumbnailName, err)
		return "", fmt.Errorf("error uploading file %v: %v", newImageThumbnailName, err)
	}

	return newImageThumbnailName, nil
}

func (gc *GlobalConfig) ReleaseInUseChannelAccount(pk string) {
	if len(pk) == 0 {
		return
	}
	gc.Mutex.Lock()
	defer gc.Mutex.Unlock()
	ca, ok := gc.InUseChannelAccounts[pk]
	if ok {
		gc.ChannelAccounts <- ca
	}
	delete(gc.InUseChannelAccounts, pk)
}

func (gc *GlobalConfig) StoreInUseChannelAccount(kp *keypair.Full) {
	if kp == nil {
		return
	}
	gc.Mutex.Lock()
	defer gc.Mutex.Unlock()
	gc.InUseChannelAccounts[kp.Address()] = kp
}
