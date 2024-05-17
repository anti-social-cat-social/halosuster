package customer

import (
	localError "eniqlo/pkg/error"
	"errors"

	"github.com/google/uuid"
)

type ICustomerUsecase interface {
	FindCustomers(query QueryParams) ([]Customer, *localError.GlobalError)
	Register(dto CustomerRegisterDTO) (*CustomerResponse, *localError.GlobalError)
}

type customerUsecase struct {
	repo ICustomerRepository
}

func NewCustomerUsecase(repo ICustomerRepository) ICustomerUsecase {
	return &customerUsecase{
		repo: repo,
	}
}

func (uc *customerUsecase) FindCustomers(query QueryParams) ([]Customer, *localError.GlobalError) {
	customers, err := uc.repo.FindAll(query)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

// Register implements ICustomerUsecase.
func (u *customerUsecase) Register(dto CustomerRegisterDTO) (*CustomerResponse, *localError.GlobalError) {
	// Validate customer request first

	// Check if customer with given phone is already exists
	existedCustomer, _ := u.repo.FindByPhone(dto.PhoneNumber)
	if existedCustomer != nil {
		return nil, localError.ErrConflict("Customer already exists", errors.New("customer already exists"))
	}

	// Map DTO to customer entity
	// This used for storing data to database
	customer := Customer{
		ID:          uuid.NewString(),
		Name:        dto.Name,
		PhoneNumber: dto.PhoneNumber,
	}

	registeredCustomer, err := u.repo.Create(customer)
	if err != nil {
		return nil, err
	}

	response := CustomerResponse{
		UserId:      registeredCustomer.ID,
		PhoneNumber: registeredCustomer.PhoneNumber,
		Name:        registeredCustomer.Name,
	}

	return &response, nil
}
