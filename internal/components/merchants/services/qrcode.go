package merchants

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	qrcode "github.com/yeqown/go-qrcode"
)

//GenerateQRCode generates QR Code in base64encoded string
func GenerateQRCode(dynamicLink string) (png string, err error) {
	if len(dynamicLink) == 0 {
		err = errors.New("no dynamic Link submitted for QRCode")
		return
	}
	qrc, err := qrcode.New(dynamicLink, qrcode.WithLogoImageFilePNG("xbnlogo.png"))
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return
	}

	// save file
	buf := new(bytes.Buffer)
	if err = qrc.SaveTo(buf); err != nil {
		fmt.Printf("could not save image: %v", err)
		return
	}
	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(buf.Bytes())

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Encoding, nil
}
