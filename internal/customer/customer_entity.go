package customer

import (
	"time"
)

type Customer struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phoneNumber" db:"phone_number"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type CustomerRegisterDTO struct {
	Name        string `json:"name" validate:"required,min=5,max=50"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=16"`
}

type CustomerResponse struct {
	UserId      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

type QueryParams struct {
	Limit       int    `form:"limit"`
	Offset      int    `form:"offset"`
	PhoneNumber string `form:"phoneNumber"`
	Name        string `form:"name"`
}

func FormatCustomerResponse(customer Customer) CustomerResponse {
	return CustomerResponse{
		UserId:      customer.ID,
		Name:        customer.Name,
		PhoneNumber: customer.PhoneNumber,
	}
}

func FormatCustomersResponse(customers []Customer) []CustomerResponse {
	formattedCustomers := []CustomerResponse{}
	for _, customer := range customers {
		formattedCustomers = append(formattedCustomers, FormatCustomerResponse(customer))
	}

	return formattedCustomers
}
