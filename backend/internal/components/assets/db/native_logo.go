package assets

import "os"

func NativeLogo() string {
	return os.Getenv("NATIVE_ASSET_IMAGE_URL")
}
