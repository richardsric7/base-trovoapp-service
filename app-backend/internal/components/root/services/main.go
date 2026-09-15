package root

import (
	"fmt"
	"os"
	root "trovo-wallet-api/internal/components/root/models"
)

// GetRootInfo returns root information on this bantupayapi instance.
// notably the organization and other relevant pieces of informations for general consumption
func GetRootInfo() root.RootInfo {
	var rootDefault root.RootInfo
	rootDefault.Organisation = os.Getenv("ORGANISATION")
	fmt.Println(rootDefault.Organisation)
	return rootDefault
}
