package dynamiclinks

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"trovo-wallet-api/internal/cache"

	qrv2 "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// GenerateQRCode generates QR Code in base64encoded string
func GenerateQRCode(dynamicLink string, redisCache *cache.RedisCache) (png string, err error) {
	if len(dynamicLink) == 0 {
		err = errors.New("no dynamic Link submitted for QRCode")
		return
	}
	cacheKey := dynamicLink + "_qrcodev2"
	{

		// search cache for link

		ok, response := redisCache.GetCachedResult(cacheKey)

		if ok {
			log.Printf("[GenerateQRCode][%v], served from cache\n", cacheKey)
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

	f, _ := os.CreateTemp("", "*.png")
	fileName := f.Name()

	defer os.Remove(f.Name())
	w := standard.NewWithWriter(f,
		standard.WithCircleShape(),
		standard.WithFgColorRGBHex("#2c2c32"),
		standard.WithBgColorRGBHex("#ffffff"),
		standard.WithQRWidth(20),
		standard.WithBorderWidth(20),
		standard.WithHalftone("ht2.png"),
	)

	err = qrc.Save(w)
	if err != nil {
		fmt.Printf("[GenerateQRCode]could not save QRCode: %v", err)
		return
	}
	var fileContents []byte
	fo, e := os.Open(fileName)
	if e == nil {
		f = fo
	}

	bufio.NewReader(f).Read(fileContents)

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(fileContents)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += base64.StdEncoding.EncodeToString(fileContents)
	//store to cache
	redisCache.StoreResultToCache(cacheKey, base64Encoding, (525960 * 3 * 60))
	f.Close()
	return base64Encoding, nil
}
