package user

import "time"

type UserRole string

type NIPPrefix string

const (
	IT    UserRole = "it"
	Nurse UserRole = "nurse"

	ITPrefix    NIPPrefix = "615"
	NursePrefix NIPPrefix = "103"
)

type User struct {
	ID                  string    `json:"id" db:"id"`
	Role                UserRole  `json:"role" db:"role"`
	NIP                 string    `json:"nip" db:"nip"`
	Name                string    `json:"name" db:"name"`
	Password            *string   `json:"password" db:"password"`
	IdentityCardScanImg string    `json:"identityCardScanImg" db:"identity_card_scan_img"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
}

type ITRegisterDTO struct {
	NIP      string `json:"nip" binding:"required,numeric,min=13,max=15"`
	Name     string `json:"name" binding:"required,min=5,max=50"`
	Password string `json:"password" binding:"required,min=5,max=33"`
}

type NurseRegisterDTO struct {
	NIP                 string `json:"nip" binding:"required,numeric,min=13,max=15"`
	Name                string `json:"name" binding:"required,min=5,max=50"`
	IdentityCardScanImg string `json:"identityCardScanImg" binding:"required,url"`
}

type NurseAccessDTO struct {
	Password string `json:"password" binding:"required,min=5,max=33"`
}

type ITLoginDTO struct {
	NIP      string `json:"nip" validate:"required,min=13,valid_nip"`
	Password string `json:"password" binding:"required,min=5,max=33"`
}

type NurseLoginDTO struct {
	NIP      string `json:"nip" binding:"required,numeric,min=13,max=15"`
	Password string `json:"password" binding:"required,min=5,max=33"`
}

type LoginResponse struct {
	UserId      string `json:"userId"`
	NIP         string `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func FormatLoginResponse(nurse User, token string) LoginResponse {
	return LoginResponse{
		UserId:      nurse.ID,
		NIP:         nurse.NIP,
		Name:        nurse.Name,
		AccessToken: token,
	}
}

type NurseRegisterResponse struct {
	UserId string `json:"userId"`
	NIP    string `json:"nip"`
	Name   string `json:"name"`
}

func FormatNurseRegisterResponse(nurse User) NurseRegisterResponse {
	return NurseRegisterResponse{
		UserId: nurse.ID,
		NIP:    nurse.NIP,
		Name:   nurse.Name,
	}
}
