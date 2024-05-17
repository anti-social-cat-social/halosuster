package customer

import (
	"database/sql"
	localError "eniqlo/pkg/error"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ICustomerRepository interface {
	FindAll(params QueryParams) ([]Customer, *localError.GlobalError)
	FindByPhone(phone string) (*Customer, *localError.GlobalError)
	Create(entity Customer) (*Customer, *localError.GlobalError)
}

type customerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) ICustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (r *customerRepository) FindAll(params QueryParams) ([]Customer, *localError.GlobalError) {
	customers := []Customer{}

	query := "SELECT * FROM customers"
	if params.Name != "" {
		query += fmt.Sprintf(" WHERE name ILIKE '%%%s%%'", params.Name)
	}
	if params.PhoneNumber != "" {
		if params.Name != "" {
			query += fmt.Sprintf(" AND phone_number ILIKE '%%%s%%'", params.PhoneNumber)
		} else {
			query += fmt.Sprintf(" WHERE phone_number ILIKE '%%%s%%'", params.PhoneNumber)
		}
	}
	if params.Limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", params.Limit)
	} else {
		query += " LIMIT 10"
	}
	if params.Offset != 0 {
		query += fmt.Sprintf(" OFFSET %d", params.Offset)
	} else {
		query += " OFFSET 0"
	}

	err := r.db.Select(&customers, query)
	if err != nil {
		return customers, localError.ErrInternalServer("Failed to find customers", err)
	}

	return customers, nil
}

// Store new customer to database
func (u *customerRepository) Create(entity Customer) (*Customer, *localError.GlobalError) {
	// entity.ID = uuid.NewString()

	// Insert into database
	query := "INSERT INTO customers (id, phone_number, name) values (:id, :phone_number, :name);"
	_, err := u.db.NamedExec(query, &entity)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	return &entity, nil
}

// Find customer by phone
// This can be use for authentication process
func (u *customerRepository) FindByPhone(phone string) (*Customer, *localError.GlobalError) {
	var customer Customer

	if err := u.db.Get(&customer, "SELECT * FROM customers where phone_number=$1", phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, localError.ErrNotFound("Customer data not found", err)
		}

		return nil, &localError.GlobalError{
			Code:    400,
			Message: "Not found",
			Error:   err,
		}

	}

	return &customer, nil
}
