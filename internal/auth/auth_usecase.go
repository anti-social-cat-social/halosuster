package auth

import (
	"1-cat-social/internal/user"
	localError "1-cat-social/pkg/error"
	"1-cat-social/pkg/hasher"
	"strconv"
	"time"
)

type IAuthUsecase interface {
	NurseLogin(req NurseLoginRequest) (user.User, *localError.GlobalError)
}

type authUsecase struct {
	userUc user.IUserUsecase
}

func NewAuthUsecase(userUc user.IUserUsecase) IAuthUsecase {
	return &authUsecase{
		userUc: userUc,
	}
}

// NurseLogin implements IAuthUsecase.
func (a *authUsecase) NurseLogin(req NurseLoginRequest) (user.User, *localError.GlobalError) {
	if !validateNIP(req.NIP, "nurse") {
		return user.User{}, localError.ErrNotFound("NIP not valid", nil)
	}

	// Search user by NIP
	nurse, err := a.userUc.FindByNIP(req.NIP)
	if err != nil {
		return *nurse, err
	}

	// Compare user password with stored password
	er := hasher.CheckPassword(nurse.Password, req.Password)
	if er != nil {
		return user.User{}, localError.ErrBadRequest("Password not match", er)
	}

	return *nurse, nil
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
