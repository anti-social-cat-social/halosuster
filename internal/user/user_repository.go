package user

import (
	"database/sql"
	"errors"
	localError "halosuster/pkg/error"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	FindById(id string) (*User, *localError.GlobalError)
	Create(entity User) (*User, *localError.GlobalError)
	FindByNIP(nip string) (*User, *localError.GlobalError)
	UpdateById(id string, key string, value string) *localError.GlobalError
	Delete(id string) *localError.GlobalError
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

// Store new ueser to database
func (u *userRepository) Create(entity User) (*User, *localError.GlobalError) {
	// Generate user UUID
	userId := uuid.NewString()
	entity.ID = userId

	q := ""
	if string(entity.Role) == "nurse" {
		q = "INSERT INTO users (id, nip, role, name, identity_card_scan_img) values (:id, :nip, :role, :name, :identity_card_scan_img);"
	} else if string(entity.Role) == "it" {
		q = "INSERT INTO users (id, nip, role, name, password) values (:id, :nip, :role, :name, :password);"
	}

	// Insert into database
	_, err := u.db.NamedExec(q, &entity)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	return &entity, nil
}

// This can be use for authentication process
func (u *userRepository) FindById(id string) (*User, *localError.GlobalError) {
	user := User{}

	log.Println(id)

	if err := u.db.Get(&user, "SELECT * FROM users where id=$1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, localError.ErrNotFound("User data not found", err)
		}

		log.Println(err)

		return nil, &localError.GlobalError{
			Code:    400,
			Message: "Not found",
			Error:   err,
		}

	}

	return &user, nil
}

// Find user by NIP
// This can be use for authentication process
func (u *userRepository) FindByNIP(nip string) (*User, *localError.GlobalError) {
	var user User

	if err := u.db.Get(&user, "SELECT * FROM users where nip=$1", nip); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, localError.ErrNotFound("User data not found", err)
		}

		log.Println(err)
		return nil, &localError.GlobalError{
			Code:    400,
			Message: "Not found",
			Error:   err,
		}
	}

	return &user, nil
}

func (u *userRepository) UpdateById(id string, key string, value string) *localError.GlobalError {
	var err error

	query := "UPDATE users SET " + key + " = $1 WHERE id = $2;"
	_, err = u.db.Exec(query, value, id)
	if err != nil {
		return localError.ErrInternalServer(err.Error(), err)
	}

	return nil
}

func (u *userRepository) Delete(id string) *localError.GlobalError {
	var err error

	query := "DELETE from users where id = $1;"
	_, err = u.db.Exec(query, id)
	if err != nil {
		return localError.ErrInternalServer(err.Error(), err)
	}

	return nil
}
