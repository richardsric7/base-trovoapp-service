package users

import (
	"time"

	"gorm.io/gorm"
)

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

/*
*
These are possible verification status : Ongoing, Completed, Pending, Failed
*
*/
type UserDojaKYCProgress struct {
	ID                 uint64 `json:"-"`
	Username           string `gorm:"not null;index:,unique;" json:"Username"`
	KYCLevel1Submitted int    `gorm:"type:integer;not null; default:0" json:"kycLevel1Submitted"`
	KYCLevel1Completed int    `gorm:"type:integer;not null; default:0" json:"kycLevel1Completed"`
	KYCLevel2Submitted int    `gorm:"type:integer;not null; default:0" json:"kycLevel2Submitted"`
	KYCLevel2Completed int    `gorm:"type:integer;not null; default:0" json:"kycLevel2Completed"`
	KYCLevel3Submitted int    `gorm:"type:integer;not null; default:0" json:"kycLevel3Submitted"`
	KYCLevel3Completed int    `gorm:"type:integer;not null; default:0" json:"kycLevel3Completed"`
	KYCLevel4Submitted int    `gorm:"type:integer;not null; default:0" json:"kycLevel4Submitted"`
	KYCLevel4Completed int    `gorm:"type:integer;not null; default:0" json:"kycLevel4Completed"`
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

type KycWebhookRequest struct {
	CreatedAt       time.Time
	ID              uint64
	ServiceProvider string
	Data            string
}

/**
These are possible verification status : Ongoing, Completed, Pending, Failed
**/
//Obj.Data.GovernmentData.Data.Bvn.Entity.Bvn
//Obj.IDType="BVN", Obj.Value="BVN VALUE"
//Obj.Status=
type DojaKYCResponse struct {
	Metadata struct {
		Ipinfo struct {
			Status     string  `json:"status"`
			Country    string  `json:"country"`
			City       string  `json:"city"`
			District   string  `json:"district"`
			Zip        string  `json:"zip"`
			Lat        float64 `json:"lat"`
			Lon        float64 `json:"lon"`
			Timezone   string  `json:"timezone"`
			Isp        string  `json:"isp"`
			Org        string  `json:"org"`
			As         string  `json:"as"`
			Mobile     bool    `json:"mobile"`
			Proxy      bool    `json:"proxy"`
			Hosting    bool    `json:"hosting"`
			Query      string  `json:"query"`
			RegionName string  `json:"region_name"`
		} `json:"ipinfo"`
		DeviceInfo string `json:"device_info"`
		UserID     string `json:"user_id"`
	} `json:"metadata"`
	Data struct {
		Index struct {
			Data struct {
			} `json:"data"`
			Message string `json:"message"`
			Status  bool   `json:"status"`
		} `json:"index"`
		Email struct {
			Data struct {
				Email string `json:"email"`
			} `json:"data"`
			Status  bool   `json:"status"`
			Message string `json:"message"`
		} `json:"email"`
		UserData struct {
			Data struct {
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Dob       string `json:"dob"`
				Email     string `json:"email"`
			} `json:"data"`
			Message string `json:"message"`
			Status  bool   `json:"status"`
		} `json:"user_data"`
		Countries struct {
			Data struct {
				Country string `json:"country"`
			} `json:"data"`
			Message string `json:"message"`
			Status  bool   `json:"status"`
		} `json:"countries"`
		GovernmentData struct {
			Data struct {
				Bvn struct {
					Entity struct {
						Customer           string    `json:"customer"`
						AppID              any       `json:"app_id"`
						Bvn                string    `json:"bvn"`
						FirstName          string    `json:"first_name"`
						LastName           string    `json:"last_name"`
						MiddleName         string    `json:"middle_name"`
						Gender             string    `json:"gender"`
						DateOfBirth        string    `json:"date_of_birth"`
						PhoneNumber1       string    `json:"phone_number1"`
						PhoneNumber2       string    `json:"phone_number2"`
						ImageURL           any       `json:"image_url"`
						Email              string    `json:"email"`
						EnrollmentBank     string    `json:"enrollment_bank"`
						EnrollmentBranch   string    `json:"enrollment_branch"`
						LevelOfAccount     string    `json:"level_of_account"`
						LgaOfOrigin        string    `json:"lga_of_origin"`
						LgaOfResidence     string    `json:"lga_of_residence"`
						MaritalStatus      string    `json:"marital_status"`
						NameOnCard         string    `json:"name_on_card"`
						Nationality        string    `json:"nationality"`
						Nin                string    `json:"nin"`
						RegistrationDate   string    `json:"registration_date"`
						ResidentialAddress string    `json:"residential_address"`
						StateOfOrigin      string    `json:"state_of_origin"`
						StateOfResidence   string    `json:"state_of_residence"`
						Title              string    `json:"title"`
						Type               string    `json:"type"`
						Xc                 any       `json:"xc"`
						Sc                 bool      `json:"sc"`
						WatchListed        string    `json:"watch_listed"`
						CreatedAt          time.Time `json:"createdAt"`
						UpdatedAt          time.Time `json:"updatedAt"`
					} `json:"entity"`
				} `json:"bvn"`
			} `json:"data"`
			Message string `json:"message"`
			Status  bool   `json:"status"`
		} `json:"government_data"`
		PhoneNumber struct {
			Data struct {
				Phone string `json:"phone"`
			} `json:"data"`
			Message string `json:"message"`
			Status  bool   `json:"status"`
		} `json:"phone_number"`
		Address struct {
			Message string `json:"message"`
			Status  bool   `json:"status"`
			Data    struct {
				Location struct {
					UserLocation struct {
						Latitude  string `json:"latitude"`
						Longitude string `json:"longitude"`
						Name      string `json:"name"`
					} `json:"user_location"`
					AddressLocation struct {
						Latitude  string `json:"latitude"`
						Longitude string `json:"longitude"`
					} `json:"address_location"`
				} `json:"location"`
			} `json:"data"`
		} `json:"address"`
		AdditionalDocument []struct {
			DocumentType string `json:"document_type"`
			DocumentURL  string `json:"document_url"`
		} `json:"additional_document"`
	} `json:"data"`
	IDType            string `json:"id_type"`
	Value             string `json:"value"`
	Message           string `json:"message"`
	ReferenceID       string `json:"reference_id"`
	WidgetID          string `json:"widget_id"`
	VerificationMode  string `json:"verification_mode"`
	VerificationType  string `json:"verification_type"`
	VerificationValue string `json:"verification_value"`
	VerificationURL   string `json:"verification_url"`
	SelfieURL         string `json:"selfie_url"`
	Status            bool   `json:"status"`
	Aml               struct {
		Status bool `json:"status"`
	} `json:"aml"`
	VerificationStatus string `json:"verification_status"` //These are possible verification status : Ongoing, Completed, Pending, Failed
}

/*
*
the following are widget IDs that will be used for different levels of kyc verification via dojah
- level-1-individual: 67e6968cc5f45aec8ae35e0a
- level-2-individual: 6841be0ae4b315e6b0e1f6c0
- level-3-individual: 6841c03de4b315e6b0e279d5
- level-4-individual: 6841c137e4b315e6b0e27c2a
- level-1-corporate: 6841c26b1b455accb4e74525
- level-2-corporate: 6841c2bf4cb77a65a43ea95e
- level-3-corporate: 6841c34ac100ce46ac98a6b5
- level-4-corporate: 6841c424c100ce46ac98aa98
*
*/
type DojaWidget struct {
	ID        string `json:"id"`
	Level     int    `json:"level"`
	Corporate int    `json:"corporate"` //0 = individual, 1 = corporate
}

type DojaWidgetID string

func (d DojaWidgetID) GetByID(db *gorm.DB) (w DojaWidget) {
	db.Where("id = ?", string(d)).First(&w)
	return
}

func (u Username) GetUserDojaKYCProgress(db *gorm.DB) (progress UserDojaKYCProgress) {
	db.Where("username = ?", string(u)).First(&progress)
	if len(progress.Username) == 0 {
		progress.Username = string(u)
	}
	return
}
