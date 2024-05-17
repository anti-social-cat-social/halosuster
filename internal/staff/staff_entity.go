package staff

import "time"

type Staff struct {
	ID        	string     	`json:"id" db:"id"`
	Name      	string    	`json:"name"`
	PhoneNumber	string		`json:"phoneNumber" db:"phone_number"`
	Password 	string    	`json:"password"`
	CreatedAt	time.Time 	`json:"createdAt" db:"created_at"`
}

type StaffRegisterDTO struct {
	Name     	string `json:"name" validate:"required,min=5,max=50"`
	PhoneNumber	string `json:"phoneNumber" validate:"required,min=10,max=16"`
	Password	string `json:"password" validate:"required,min=5,max=15"`
}

type StaffLoginDTO struct {
	PhoneNumber	string `json:"phoneNumber" validate:"required,min=10,max=16"`
	Password	string `json:"password" validate:"required,min=5,max=15"`
}

type StaffRegAndLoginResponse struct {
	UserId		string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"` 
	Name 		string `json:"name"`
	AccessToken string `json:"accessToken"`
}