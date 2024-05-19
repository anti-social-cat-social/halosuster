package patient

import "time"

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type Patient struct {
	IdentityNumber      int64     `json:"identityNumber"`
	PhoneNumber         string    `json:"phoneNumber"`
	Name                string    `json:"name"`
	BirthDate           time.Time `json:"birthDate"`
	Gender              Gender    `json:"gender"`
	IdentityCardScanImg string    `json:"identityCardScanImg"`
}
