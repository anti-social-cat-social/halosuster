package user

import (
	"errors"
	localError "halosuster/pkg/error"
	"halosuster/pkg/hasher"
	tokenizer "halosuster/pkg/jwt"
)

type IUserUsecase interface {
	NurseLogin(req NurseLoginDTO) (User, *localError.GlobalError)
	ITLogin(req ITLoginDTO) (*LoginResponse, *localError.GlobalError)
	FindByNIP(nip string) (*User, *localError.GlobalError)
}

type userUsecase struct {
	repo IUserRepository
}

// FindByNIP implements IUserUsecase.
func (u *userUsecase) FindByNIP(email string) (*User, *localError.GlobalError) {
	return u.repo.FindByNIP(email)
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
	}

	token, tokenErr := tokenizer.GenerateToken(tokenData)
	if tokenErr != nil {
		return nil, localError.ErrInternalServer(tokenErr.Error(), tokenErr)
	}

	response := FormatLoginResponse(*user, token)

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
