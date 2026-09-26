package models

import "time"

// RuleScore holds the rule score for the rule engine
type RuleScore struct {
	Rule        string `gorm:"size:20;not null"`
	Score       int    `gorm:"type:smallint;not null;default:0"`
	Description string
}

// HighRiskIP holds IPs and Subnets with highrisk for service
type HighRiskIP struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	IPAddress string    `gorm:"size:100;not null;index:highrisk_ip_index,unique"`
	Score     int       `gorm:"type:integer;not null;default:0" json:"score"`
}

// HighRiskCountry holds country codes with highrisk for service
type HighRiskCountry struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	CountryCode string    `gorm:"size:100;not null;index:highrisk_country_code_index,unique"`
	Score       int       `gorm:"type:integer;not null;default:0" json:"score"`
}

// HighRiskCity holds city with highrisk for service
type HighRiskCity struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	City      string    `gorm:"size:100;not null;index:highrisk_city_index,unique"`
	Score     int       `gorm:"type:integer;not null;default:0" json:"score"`
}

// MediumRiskCountry holds country codes with medium risk for service
type MediumRiskCountry struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	CountryCode string    `gorm:"size:100;not null;index:mediumrisk_country_code_index,unique"`
	Score       int       `gorm:"type:integer;not null;default:0" json:"score"`
}

// HighRiskISP holds ISPs with highrisk for service
type HighRiskISP struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	ISP       string    `gorm:"size:100;not null;index:highrisk_isp_index,unique"`
	Score     int       `gorm:"type:integer;not null;default:0" json:"score"`
}

// MediumRiskDomain holds domain with medium risk for service
type MediumRiskDomain struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	Domain    string    `gorm:"size:100;not null;index:mediumrisk_domain_index,unique"`
	Score     int       `gorm:"type:integer;not null;default:0" json:"score"`
}

// CountryPhoneCode holds country and their phone codes
type CountryPhoneCode struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	CountryCode string    `gorm:"size:17;not null;index:country_code_index,unique"`
	PhoneCode   string    `gorm:"size:100;not null;index:phone_code_index"`
}

// OffendingIP is model for storing offending IPs
type OffendingIP struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IPAddress string    `gorm:"size:100;not null;index:offending_ip_index,unique" json:"ip_address"`
}
