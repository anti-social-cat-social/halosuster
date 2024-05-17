package staff

import (
	localError "eniqlo/pkg/error"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type IStaffRepository interface {
	FindByPhone(phone string) (*Staff, *localError.GlobalError)
	Create(entity Staff) (*Staff, *localError.GlobalError)
}

type staffRepository struct {
	db *sqlx.DB
}

func NewStaffRepository(db *sqlx.DB) IStaffRepository {
	return &staffRepository{
		db: db,
	}
}

// Store new staff to database
func (u *staffRepository) Create(entity Staff) (*Staff, *localError.GlobalError) {
	// entity.ID = uuid.NewString()

	// Insert into database
	query := "INSERT INTO staffs (id, phone_number, name, password) values (:id, :phone_number, :name, :password);"
	_, err := u.db.NamedExec(query, &entity)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	return &entity, nil
}

// Find staff by phone
// This can be use for authentication process
func (u *staffRepository) FindByPhone(phone string) (*Staff, *localError.GlobalError) {
	var staff Staff

	if err := u.db.Get(&staff, "SELECT * FROM staffs where phone_number=$1", phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, localError.ErrNotFound("Staff data not found", err)
		}

		return nil, &localError.GlobalError{
			Code:    400,
			Message: "Not found",
			Error:   err,
		}

	}

	return &staff, nil
}
