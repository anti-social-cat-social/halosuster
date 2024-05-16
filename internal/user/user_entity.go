package user

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"-" db:"createdat"`
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
