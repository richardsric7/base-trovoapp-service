package users

type KYCLevel struct {
	ID           string `json:"kycLevel"`
	UserCategory string `json:"userCategory"`
}

type KYCConfig struct {
	ID              uint64 `json:"-"`
	ServiceProvider string `json:"serviceProvider"`
	Token           string `json:"token"`
	SecretKey       string `json:"secretKey"`
}

type SumSubReviewResultInput struct {
	ApplicantID    string       `json:"applicantId"`
	InspectionID   string       `json:"inspectionId"`
	CorrelationID  string       `json:"correlationId"`
	ExternalUserID string       `json:"externalUserId"`
	LevelName      string       `json:"levelName"`
	Type           string       `json:"type"`
	ReviewResult   ReviewResult `json:"reviewResult"`
	ReviewStatus   string       `json:"reviewStatus"`
	CreatedAtMs    string       `json:"createdAtMs"`
}

type ReviewResult struct {
	ModerationComment string   `json:"moderationComment"`
	ClientComment     string   `json:"clientComment"`
	ReviewAnswer      string   `json:"reviewAnswer"`
	RejectLabels      []string `json:"rejectLabels"`
	ReviewRejectType  string   `json:"reviewRejectType"`
	ButtonIds         []string `json:"buttonIds"`
}

type SumSubReviewResult struct {
	ID             uint64
	ApplicantID    string `json:"applicantId"`
	InspectionID   string `json:"inspectionId"`
	CorrelationID  string `json:"correlationId"`
	ExternalUserID string `json:"externalUserId"`
	LevelName      string `json:"levelName"`
	Type           string `json:"type"`
	ReviewResult   string `json:"reviewResult"`
	ReviewStatus   string `json:"reviewStatus"`
	CreatedAtMs    string `json:"createdAtMs"`
}

type UserKYCProgress struct {
	ID                 uint64 `json:"-"`
	Username           string `gorm:"not null;index:,unique;" json:"Username"`
	KYCLevel1Initiated int    `gorm:"type:integer;not null; default:0" json:"kycLevel1Initiated"`
	KYCLevel1Done      int    `gorm:"type:integer;not null; default:0" json:"kycLevel1Done"`
	KYCLevel2Initiated int    `gorm:"type:integer;not null; default:0" json:"kycLevel2Initiated"`
	KYCLevel2Done      int    `gorm:"type:integer;not null; default:0" json:"kycLevel2Done"`
	KYCLevel3Initiated int    `gorm:"type:integer;not null; default:0" json:"kycLevel3Initiated"`
	KYCLevel3Done      int    `gorm:"type:integer;not null; default:0" json:"kycLevel3Done"`
}

type SumsubInfo struct {
	FirstName    string        `json:"firstName,omitempty"`
	FirstNameEn  string        `json:"firstNameEn,omitempty"`
	MiddleName   string        `json:"middleName,omitempty"`
	MiddleNameEn string        `json:"middleNameEn,omitempty"`
	LastName     string        `json:"lastName,omitempty"`
	LastNameEn   string        `json:"lastNameEn,omitempty"`
	Dob          string        `json:"dob,omitempty"` //yyyy-mm-dd format
	Gender       string        `json:"gender,omitempty"`
	Country      string        `json:"country,omitempty"`
	Phone        string        `json:"phone,omitempty"`
	IdDocs       []SumsubIdDoc `json:"idDocs,omitempty"`
}

type SumsubIdDoc struct {
	IdDocType    string `json:"idDocType,omitempty"`
	Country      string `json:"country,omitempty"`
	FirstName    string `json:"firstName,omitempty"`
	FirstNameEn  string `json:"firstNameEn,omitempty"`
	MiddleName   string `json:"middleName,omitempty"`
	MiddleNameEn string `json:"middleNameEn,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	LastNameEn   string `json:"lastNameEn,omitempty"`
	DateOfBirth  string `json:"dob,omitempty"` // yyyy-mm-dd format
}
type SumsubApplicant struct {
	ID             string     `json:"id,omitempty"`
	CreatedAt      string     `json:"createdAt,omitempty"`
	Key            string     `json:"key,omitempty"`
	ClientID       string     `json:"clientId,omitempty"`
	InspectionID   string     `json:"inspectionId,omitempty"`
	ExternalUserID string     `json:"externalUserId,omitempty"`
	Info           SumsubInfo `json:"info,omitempty"`
	FixedInfo      SumsubInfo `json:"fixedInfo,omitempty"`
	Review         struct {
		ElapsedSincePendingMs int    `json:"elapsedSincePendingMs,omitempty"`
		ElapsedSinceQueuedMs  int    `json:"elapsedSinceQueuedMs,omitempty"`
		Reprocessing          bool   `json:"reprocessing,omitempty"`
		CreateDate            string `json:"createDate,omitempty"`
		ReviewDate            string `json:"reviewDate,omitempty"`
		StartDate             string `json:"startDate,omitempty"`
		ReviewResult          struct {
			ReviewAnswer string `json:"reviewAnswer,omitempty"`
		} `json:"reviewResult,omitempty"`
		ReviewStatus           string `json:"reviewStatus,omitempty"`
		NotificationFailureCnt int    `json:"notificationFailureCnt,omitempty"`
		Priority               int    `json:"priority,omitempty"`
	} `json:"review,omitempty"`
	Lang string `json:"lang,omitempty"`
	Type string `json:"type,omitempty"`
}

type SumsubAccessToken struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}
