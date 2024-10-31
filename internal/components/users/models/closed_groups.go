package users

type ClosedGroup struct {
	ID                      string `gorm:"primaryKey;size:100" json:"id"`
	GroupName               string `gorm:"size:100" json:"groupName"`
	GroupOwner              string `gorm:"size:45" json:"groupOwner"`
	GroupDescription        string `gorm:"size:300" json:"groupDescription"`
	RegisteredEntity        int    `gorm:"default:0" json:"registeredEntity"`
	RegistrationName        string `gorm:"size:200" json:"registrationName"`
	RegistrationNumber      string `gorm:"size:200" json:"registrationNumber"`
	RegistrationDocumentUrl string `json:"registrationDocumentUrl"`
}
type ClosedGroupJSONInput struct {
	ID                      string `gorm:"primaryKey;size:100" json:"id"`
	GroupName               string `gorm:"size:100" json:"groupName"`
	GroupDescription        string `gorm:"size:300" json:"groupDescription"`
	RegisteredEntity        int    `gorm:"default:0" json:"registeredEntity"`
	RegistrationName        string `gorm:"size:200" json:"registrationName"`
	RegistrationNumber      string `gorm:"size:200" json:"registrationNumber"`
	RegistrationDocumentUrl string `json:"registrationDocumentUrl"`
}

type UserClosedGroup struct {
	ID            string `gorm:"size:100" json:"id"`
	ClosedGroupID string `gorm:"size:100;index:idx_unique_closed_group,unique" json:"closedGroupId"`
	UserID        string `gorm:"size:100;index:idx_unique_closed_group,unique" json:"userId"`
}
