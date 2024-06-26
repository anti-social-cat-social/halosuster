package user

import (
	"errors"
	localError "halosuster/pkg/error"
	"halosuster/pkg/hasher"
	tokenizer "halosuster/pkg/jwt"
	"strconv"
	"time"
)

type IUserUsecase interface {
	NurseLogin(req NurseLoginDTO) (User, *localError.GlobalError)
	ITLogin(req ITLoginDTO) (*LoginResponse, *localError.GlobalError)
	ITRegister(req ITRegisterDTO) (*LoginResponse, *localError.GlobalError)
	FindByNIP(nip string) (*User, *localError.GlobalError)
	NurseRegister(req NurseRegisterDTO) (User, *localError.GlobalError)
	NurseAccess(req NurseAccessDTO, id string) *localError.GlobalError
	GetUsers(query UserQueryParams) ([]User, *localError.GlobalError)
	Delete(id string) *localError.GlobalError
	Update(id string, dto NurseUpdateDTO) *localError.GlobalError
}

type userUsecase struct {
	repo IUserRepository
}

func NewUserUsecase(repo IUserRepository) IUserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (a *userUsecase) ITLogin(req ITLoginDTO) (*LoginResponse, *localError.GlobalError) {
	// CHEK if NIP is valid
	if req.NIP[:3] != string(ITPrefix) {
		return nil, localError.ErrNotFound("Not authorized", errors.New("not an IT personel"))
	}

	// Searcd user by NIP
	user, err := a.repo.FindByNIP(req.NIP)
	if err != nil {
		return nil, localError.ErrNotFound("Account not found", err.Error)
	}

	// Check password
	passErr := hasher.CheckPassword(user.Password, req.Password)
	if passErr != nil {
		return nil, localError.ErrBadRequest(passErr.Error(), passErr)
	}

	// Generate Token
	tokenData := tokenizer.TokenData{
		ID:   user.ID,
		Name: user.Name,
		Role: string(user.Role),
	}

	token, tokenErr := tokenizer.GenerateToken(tokenData)
	if tokenErr != nil {
		return nil, localError.ErrInternalServer(tokenErr.Error(), tokenErr)
	}

	response := FormatLoginResponse(*user, token)

	return &response, nil
}

func (uc *userUsecase) ITRegister(req ITRegisterDTO) (*LoginResponse, *localError.GlobalError) {
	// FInd existing user data by NIP
	existingUser, _ := uc.repo.FindByNIP(req.NIP)
	if existingUser != nil {
		return nil, localError.ErrConflict("User already exists", errors.New("user already exists"))
	}

	// Generate Password
	password, errPass := hasher.HashPassword(req.Password)
	if errPass != nil {
		return nil, localError.ErrInternalServer(errPass.Error(), errPass)
	}

	// Create User
	user, err := uc.repo.Create(User{
		Role:                IT,
		NIP:                 req.NIP,
		Name:                req.Name,
		Password:            &password,
		IdentityCardScanImg: nil,
	})
	if err != nil {
		return nil, err
	}

	// Generate token
	tokenData := tokenizer.TokenData{
		ID:   user.ID,
		Name: user.Name,
		Role: string(user.Role),
	}

	token, errToken := tokenizer.GenerateToken(tokenData)
	if errToken != nil {
		return nil, localError.ErrInternalServer(errToken.Error(), errToken)
	}

	response := LoginResponse{
		UserId:      user.ID,
		NIP:         user.NIP,
		Name:        user.Name,
		AccessToken: token,
	}

	return &response, nil
}

// NurseLogin implements IUserUsecase.
func (a *userUsecase) NurseLogin(req NurseLoginDTO) (User, *localError.GlobalError) {
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

	if role == "it" && nip[:3] != "615" {
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
	if randomDigits < 0 || randomDigits > 99999 {
		return false
	}

	return true
}

// NurseRegister implements IUserUsecase.
func (a *userUsecase) NurseRegister(req NurseRegisterDTO) (User, *localError.GlobalError) {
	if !validateNIP(req.NIP, "nurse") {
		return User{}, localError.ErrNotFound("NIP not valid", nil)
	}

	// Search user by NIP
	existedNurse, _ := a.repo.FindByNIP(req.NIP)
	if existedNurse != nil {
		return User{}, localError.ErrConflict("Nurse already exists", nil)
	}

	nurse := User{
		NIP:                 req.NIP,
		Name:                req.Name,
		IdentityCardScanImg: &req.IdentityCardScanImg,
		Role:                UserRole("nurse"),
	}

	registeredNurse, err := a.repo.Create(nurse)
	if err != nil {
		return User{}, err
	}

	return *registeredNurse, nil
}

func (a *userUsecase) NurseAccess(req NurseAccessDTO, id string) *localError.GlobalError {
	// Search user by ID
	nurse, err := a.repo.FindById(id)
	if err != nil {
		return err
	}

	if nurse.Role != "nurse" {
		return localError.ErrNotFound("user not found", errors.New("user not found"))
	}

	// Generate user password
	password, errHash := hasher.HashPassword(req.Password)
	if errHash != nil {
		return localError.ErrInternalServer(errHash.Error(), errHash)
	}

	err = a.repo.UpdateById(id, "password", password)
	if err != nil {
		return err
	}

	return nil
}

func (a *userUsecase) GetUsers(query UserQueryParams) ([]User, *localError.GlobalError) {
	users, err := a.repo.FindAll(query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *userUsecase) Delete(id string) *localError.GlobalError {
	// Check if user not actually exists
	user, err := uc.repo.FindById(id)
	if err != nil {
		return err
	}

	if user.Role != Nurse {
		return localError.ErrNotFound("Non Nurse user can not be deleted", errors.New("non nurse accoun can't be deleted"))
	}

	err = uc.repo.Delete(user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *userUsecase) Update(id string, dto NurseUpdateDTO) *localError.GlobalError {
	// Get existing user
	nip := strconv.Itoa(dto.NIP)

	existingUser, errExists := uc.repo.FindByNIP(nip)

	// IF user not found
	if errExists != nil {
		return localError.ErrNotFound(errExists.Message.(string), errExists.Error)
	}

	// User not nurse
	if existingUser.Role != Nurse {
		return localError.ErrNotFound("NOT FOUND", errors.New("not found"))
	}

	// Return errror if user already exists and new NIP not same with request
	if existingUser.ID != id {
		return localError.ErrConflict("NIP already exists", errors.New("user already exists"))
	}

	// Update current data
	existingUser.NIP = nip
	existingUser.Name = dto.Name

	return uc.repo.Update(id, *existingUser)
}
