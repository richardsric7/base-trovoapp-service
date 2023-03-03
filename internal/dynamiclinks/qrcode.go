package dynamiclinks

import (
	"errors"
	"fmt"
	"log"
	"os"
	"trovo-wallet-api/internal/sharedconfig"

	qrv2 "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// GenerateQRCode generates QR Code in base64encoded string
func GenerateQRCode(dynamicLink string, gc *sharedconfig.GlobalConfig) (png string, err error) {
	if len(dynamicLink) == 0 {
		err = errors.New("no dynamic Link submitted for QRCode")
		return
	}
	cacheKey := dynamicLink + "qrcode2Link"
	{

		// search cache for link

		ok, response := gc.RedisCache.GetCachedResult(cacheKey)

		if ok {
			// log.Printf("[GenerateQRCode][%v], served from cache\n", cacheKey)
			png = response.(string)
			return
		}

	}
	qrc, err := qrv2.NewWith(dynamicLink,
		qrv2.WithErrorCorrectionLevel(qrv2.ErrorCorrectionHighest),
	)
	if err != nil {
		fmt.Printf("[GenerateQRCode]could not generate QRCode: %v", err)
		return
	}
	// buf :=new(bytes.Buffer)

	f, err := os.CreateTemp("", "*.png")
	if err != nil {
		log.Printf("[GenerateQRCode]could not generate QRCode: %v\n", err)

		return
	}

	fileName := f.Name()
	defer os.Remove(f.Name())

	w, err := standard.New(fileName,
		standard.WithCircleShape(),
		// standard.WithFgColorRGBHex("#2c2c32"),
		standard.WithBgColorRGBHex("#ffffff"),
		standard.WithQRWidth(20),
		standard.WithBorderWidth(20),
		// standard.WithHalftone("ht2.png"),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
	)
	if err != nil {
		log.Printf("[GenerateQRCode]could not save QRCode: %v\n", err)
		return
	}
	err = qrc.Save(w)
	if err != nil {
		log.Printf("[GenerateQRCode]could not save QRCode: %v\n", err)
		return
	}

	// fileContents, err := os.ReadFile(fileName)
	// if err != nil {
	// 	log.Printf("[GenerateQRCode]could read QRCode: %v\n", err)
	// 	return
	// }

	// log.Println(fileContents)
	// var fileNameWithExt string

	// Determine the content type of the image file
	// mimeType := http.DetectContentType(fileContents)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	// switch mimeType {
	// case "image/jpeg":
	// 	fileNameWithExt = fileName + ".jpeg"
	// case "image/png":
	// 	fileNameWithExt = fileName + ".png"
	// }

	newThumbnail, err := gc.FirebaseStorageUploader.SaveQrCodeAsFileToCloud(f, fileName, "")
	if err != nil {
		log.Printf("[GenerateQRCode]could upload QRCode: %v\n", err)

		return
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%v/%v", gc.FirebaseStorageUploader.BucketName, newThumbnail)

	//store to cache
	gc.RedisCache.StoreResultToCache(cacheKey, url, (525960 * 3 * 60))
	f.Close()
	return url, nil
}
