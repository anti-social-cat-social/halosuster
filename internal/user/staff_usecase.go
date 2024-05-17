package staff

import (
	localError "eniqlo/pkg/error"
	"eniqlo/pkg/hasher"
	localJwt "eniqlo/pkg/jwt"
	"errors"

	"github.com/google/uuid"
)

type IStaffUsecase interface {
	Login(dto StaffLoginDTO) (*StaffRegAndLoginResponse, *localError.GlobalError)
	Register(dto StaffRegisterDTO) (*StaffRegAndLoginResponse, *localError.GlobalError)
}

type staffUsecase struct {
	repo IStaffRepository
}

func NewStaffUsecase(repo IStaffRepository) IStaffUsecase {
	return &staffUsecase{
		repo: repo,
	}
}

// Register implements IStaffUsecase.
func (u *staffUsecase) Register(dto StaffRegisterDTO) (*StaffRegAndLoginResponse, *localError.GlobalError) {
	// Validate staff request first

	// Check if staff with given phone is already exists
	existedStaff, _ := u.repo.FindByPhone(dto.PhoneNumber)
	if existedStaff != nil {
		return nil, localError.ErrConflict("Staff already exists", errors.New("staff already exists"))
	}

	// Map DTO to staff entity
	// This used for storing data to database
	staff := Staff{
		Name:        dto.Name,
		PhoneNumber: dto.PhoneNumber,
	}

	// Generate staff password
	password, errHash := hasher.HashPassword(dto.Password)
	if errHash != nil {
		return nil, localError.ErrInternalServer(errHash.Error(), errHash)
	}
	// Assign staff password to struct if not error
	staff.Password = password

	staff.ID = uuid.NewString()

	tokenData := &localJwt.TokenData{
		ID:   staff.ID,
		Name: staff.Name,
	}

	// Generate token
	token, errToken := localJwt.GenerateToken(*tokenData)
	if errToken != nil {
		return nil, localError.ErrInternalServer(errToken.Error(), errToken)
	}

	registeredStaff, err := u.repo.Create(staff)
	if err != nil {
		return nil, err
	}

	response := StaffRegAndLoginResponse{
		UserId:      registeredStaff.ID,
		PhoneNumber: registeredStaff.PhoneNumber,
		Name:        registeredStaff.Name,
		AccessToken: token,
	}

	return &response, nil
}

// Login implements IStaffUsecase.
func (u *staffUsecase) Login(dto StaffLoginDTO) (*StaffRegAndLoginResponse, *localError.GlobalError) {
	// Search staff by phone
	result, errStaff := u.repo.FindByPhone(dto.PhoneNumber)
	if errStaff != nil {
		return nil, errStaff
	}

	// Compare staff password with stored password
	err := hasher.CheckPassword(result.Password, dto.Password)
	if err != nil {
		return nil, localError.ErrBadRequest("Credential not valid", err)
	}

	tokenData := &localJwt.TokenData{
		ID:   result.ID,
		Name: result.Name,
	}

	// Generate token if no error happened above
	// Token generated using JWT scheme
	token, err := localJwt.GenerateToken(*tokenData)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	response := StaffRegAndLoginResponse{
		UserId:      result.ID,
		PhoneNumber: result.PhoneNumber,
		Name:        result.Name,
		AccessToken: token,
	}

	return &response, nil
}
