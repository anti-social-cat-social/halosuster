package user

import "time"

type Role string

const (
	IT    Role = "it"
	Nurse Role = "nurse"
)

type User struct {
	ID           string `db:"id"`
	NIP          string `db:"nip"`
	Role         Role
	Name         string
	Password     string
	IdentityCard string    `db:"identity_card_scan_img"`
	CreatedAt    time.Time `db:"created_at"`
}

type UserDTO struct {
	Name     string `json:"name" validate:"required,min=5,max=50,valid_name"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,min=5,max=50,email"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}
