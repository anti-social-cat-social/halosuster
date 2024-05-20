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

type RecordResponse struct {
	Patient    patient.Patient   `json:"identityDetail"`
	Symptomp   string            `json:"symptoms"`
	Medication string            `json:"medication"`
	CreatedAt  string            `json:"createdAt"`
	Creator    user.UserResponse `json:"createdBy"`
}

type RecordQueryParam struct {
	IdentityNumber int    `form:"identityNumber" validate:"numeric"`
	UserId         string `form:"userId" validate:"omitempty,uuid"`
	NIP            string `form:"nip" `
	Limit          int    `form:"limit" validate:"numeric"`
	Offset         int    `form:"offset" validate:"numeric"`
	CreatedAt      string `form:"createdAt"`
}

func FormatRecordResponse(rec Record) RecordResponse {
	return RecordResponse{
		Patient:    rec.Patient,
		Symptomp:   rec.Symptomp,
		Medication: rec.Medication,
		CreatedAt:  rec.CreatedAt.Format(time.RFC3339),
		Creator:    user.FormatUserResponse(rec.Creator),
	}
}
