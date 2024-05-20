package record

import (
	"fmt"
	"halosuster/internal/patient"
	"halosuster/internal/user"
	localError "halosuster/pkg/error"

	"github.com/jmoiron/sqlx"
)

type IRecordRepository interface {
	GetAll(params RecordQueryParam) ([]Record, *localError.GlobalError)
	Create(entity *Record) *localError.GlobalError
}

type recordRepo struct {
	db *sqlx.DB
}

func NewRecordRepo(db *sqlx.DB) IRecordRepository {
	return &recordRepo{db: db}
}

func (r *recordRepo) GetAll(param RecordQueryParam) ([]Record, *localError.GlobalError) {
	var records []Record

	q := "select r.*, u.nip, u.name as creator, u.id as user_id, p.* from records r join patients p on p.identity_number = r.identity_number 	join users u on u.id = r.created_by where 1=1 "

	// Handle the param
	if param.IdentityNumber != 0 {
		q += fmt.Sprintf("AND p.identity_number = '%d' ", param.IdentityNumber)
	}

	// BY NIP
	if param.NIP != "" {
		q += fmt.Sprintf("AND u.nip = '%s' ", param.NIP)
	}

	// BY user id
	if param.UserId != "" {
		q += fmt.Sprintf("AND u.id = '%s' ", param.UserId)
	}

	// Sorted
	if param.CreatedAt == "asc" || param.CreatedAt == "desc" {
		q += fmt.Sprintf("ORDER BY r.created_at %s ", param.CreatedAt)
	}

	// LIMIT and OFFSET
	limit := 5
	offset := 0
	if param.Limit != 0 {
		limit = param.Limit
	}

	if param.Offset != 0 {
		offset = param.Offset
	}

	q += fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.db.Queryx(q)
	if err != nil {
		return records, localError.ErrInternalServer(err.Error(), err)
	}

	defer rows.Close()

	for rows.Next() {
		var record Record
		var patient patient.Patient
		var creator user.User

		err := rows.Scan(
			&record.ID,
			&record.CreatedBy,
			&record.IdentityNumber,
			&record.Symptomp,
			&record.Medication,
			&record.CreatedAt,
			&creator.NIP,
			&creator.Name,
			&creator.ID,
			&patient.IdentityNumber,
			&patient.PhoneNumber,
			&patient.Name,
			&patient.BirthDate,
			&patient.Gender,
			&patient.IdentityCardScanImg,
			&patient.CreatedAt,
		)
		if err != nil {
			return records, localError.ErrInternalServer(err.Error(), err)
		}

		record.Patient = patient
		record.Creator = creator
		records = append(records, record)
	}

	return records, nil
}

func (repo *recordRepo) Create(entity *Record) *localError.GlobalError {
	// Creat equery
	q := "INSERT INTO records (symptoms, medications, identity_number, created_by) values (:symtoms, :medications, :identity, :creator)"

	// Execute Query
	_, err := repo.db.NamedExec(q, entity)
	if err != nil {
		return localError.ErrInternalServer(err.Error(), err)
	}

	return nil
}
