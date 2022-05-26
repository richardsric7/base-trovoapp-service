package currency

//AppVersion model
type AppVersion struct {
	ID          uint64
	Version     string `gorm:"size:50;null" json:"version"`
	MinVersion  string `gorm:"size:50;null" json:"minVersion"`
	IosUrL      string `gorm:"null" json:"iosUrl"`
	AndroidUrL  string `gorm:"null" json:"androidUrl"`
	ForceUpdate uint64 `gorm:"type:integer;not null;default:0" json:"forceUpdate"`
}
