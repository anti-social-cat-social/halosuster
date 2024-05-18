package user

import (
	localError "halosuster/pkg/error"
	"halosuster/pkg/hasher"
)

type IUserUsecase interface {
	NurseLogin(req NurseLoginDTO) (User, *localError.GlobalError)
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
