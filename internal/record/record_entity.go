package record

import (
	"halosuster/internal/patient"
	"halosuster/internal/user"
	"time"
)

type Record struct {
	ID             string          `json:"-"`
	Patient        patient.Patient `json:"identityDetail"`
	IdentityNumber string          `json:"-"`
	Symptomp       string          `json:"symptoms"`
	Medication     string          `json:"medication"`
	CreatedAt      time.Time       `json:"createdAt"`
	Creator        user.User       `json:"createdBy"`
	CreatedBy      string          `json:"-"`
}

type RecordQueryParam struct {
	IdentityNumber string `form:"identityNumber"`
	UserId         string `json:"userId"`
	NIP            string `json:"nip"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	CreatedAt      string `json:"createdAt"`
}
