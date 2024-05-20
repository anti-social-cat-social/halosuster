package record

import (
	localError "halosuster/pkg/error"
	"strconv"
)

type IRecordUsecase interface {
	GetAll(params RecordQueryParam) ([]RecordResponse, *localError.GlobalError)
	Create(dto RecordDTO) *localError.GlobalError
}

type recordUsecase struct {
	repo IRecordRepository
}

func NewRecordUsecase(repo IRecordRepository) IRecordUsecase {
	return &recordUsecase{
		repo: repo,
	}
}

func (uc *recordUsecase) GetAll(params RecordQueryParam) ([]RecordResponse, *localError.GlobalError) {
	records := []RecordResponse{}

	result, err := uc.repo.GetAll(params)
	if err != nil {
		return []RecordResponse{}, err
	}

	for _, res := range result {
		response := FormatRecordResponse(res)

		records = append(records, response)
	}

	return records, err
}

func (uc *recordUsecase) Create(dto RecordDTO) *localError.GlobalError {
	// Change identity number
	identityNumber := strconv.Itoa(dto.IdentityNumber)

	// Create entity data
	record := Record{
		IdentityNumber: identityNumber,
		Symptomp:       dto.Symptomp,
		Medication:     dto.Medication,
	}

	// Utilize repo
	return uc.repo.Create(&record)
}
