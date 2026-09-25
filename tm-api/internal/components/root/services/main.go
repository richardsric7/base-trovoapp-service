package root

import (
	root "admin-panel-dashboard/internal/components/root/models"
)

// GetRootInfo returns root information on this bantupayapi instance.
// notably the organization and other relevant pieces of informations for general consumption
func GetRootInfo() root.RootInfo {
	var rootDefault root.RootInfo
	rootDefault.Service = "Peer To Peer Token Exchange Platform"

	return rootDefault
}
