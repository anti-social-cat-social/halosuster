package record

import (
	"halosuster/internal/patient"
	"halosuster/internal/user"
	localError "halosuster/pkg/error"

	"github.com/jmoiron/sqlx"
)

type IRecordRepository interface {
	GetAll() ([]Record, *localError.GlobalError)
}

type recordRepo struct {
	db *sqlx.DB
}

func NewRecordRepo(db *sqlx.DB) IRecordRepository {
	return &recordRepo{db: db}
}

func (r *recordRepo) GetAll() ([]Record, *localError.GlobalError) {
	var records []Record

	q := "select r.*, u.nip, u.name as creator, u.id as user_id, p.* from records r join patients p on p.identity_number = r.identity_number 	join users u on u.id = r.created_by where 1=1 "

	rows, err := r.db.Queryx(q)
	defer rows.Close()

	if err != nil {
		return records, localError.ErrInternalServer(err.Error(), err)
	}

	for rows.Next() {
		var record Record
		var patient patient.Patient
		var creator user.User

		rows.Scan(&record.ID, &record.IdentityNumber, &record.Symptomp, &record.Medication, &record.CreatedAt, &creator.NIP, &creator.Name, &creator.ID, &patient.IdentityNumber, &patient.PhoneNumber, &patient.Name, &patient.BirthDate, &patient.Gender, &patient.IdentityCardScanImg)

		record.Patient = patient
		record.Creator = creator
		records = append(records, record)
	}

	return records, nil
}
