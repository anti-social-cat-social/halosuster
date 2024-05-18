package user

import (
	localError "halosuster/pkg/error"
	"halosuster/pkg/hasher"
	"strconv"
	"time"
)

type IUserUsecase interface {
	NurseLogin(req NurseLoginDTO) (User, *localError.GlobalError)
	FindByNIP(nip string) (*User, *localError.GlobalError)
	// Create(dto UserDTO) (*User, *localError.GlobalError)
}

type userUsecase struct {
	repo IUserRepository
}

func NewUserUsecase(repo IUserRepository) IUserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

// NurseLogin implements IUserUsecase.
func (a *userUsecase) NurseLogin(req NurseLoginDTO) (User, *localError.GlobalError) {
	if !validateNIP(req.NIP, "nurse") {
		return User{}, localError.ErrNotFound("NIP not valid", nil)
	}

	// Search user by NIP
	nurse, err := a.repo.FindByNIP(req.NIP)
	if err != nil {
		return *nurse, err
	}

	// Compare user password with stored password
	er := hasher.CheckPassword(nurse.Password, req.Password)
	if er != nil {
		return User{}, localError.ErrBadRequest("Password not match", er)
	}

	return *nurse, nil
}

// FindByNIP implements IUserUsecase.
func (u *userUsecase) FindByNIP(nip string) (*User, *localError.GlobalError) {
	return u.repo.FindByNIP(nip)
}


func validateNIP(nip string, role string) bool {
	// Check if first three digits are '303'
	if role == "nurse" && nip[:3] != "303" {
		return false
	}

	// Check the fourth digit based on gender
	genderDigit, _ := strconv.Atoi(nip[3:4])
	if genderDigit != 1 && genderDigit != 2 {
		return false
	}

	// Check if the fifth and eighth digit represent a valid year
	year, _ := strconv.Atoi(nip[4:8])
	currentYear := time.Now().Year()
	if year < 2000 || year > currentYear {
		return false
	}

	// Check if the ninth and tenth digit represent a valid month
	month, _ := strconv.Atoi(nip[8:10])
	if month < 1 || month > 12 {
		return false
	}

	// Check if the eleventh and thirteenth digits are within range
	randomDigits, _ := strconv.Atoi(nip[10:])
	if randomDigits < 0 || randomDigits > 999 {
		return false
	}

	return true
}
