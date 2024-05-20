package record

import (
	localError "halosuster/pkg/error"
)

type IRecordUsecase interface {
	GetAll(params RecordQueryParam) ([]RecordResponse, *localError.GlobalError)
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
