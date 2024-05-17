package product

import (
	localError "eniqlo/pkg/error"
	"eniqlo/pkg/validation"
)

type IProductUsecase interface {
	CreateProduct(req CreateProductRequest) (*Product, *localError.GlobalError)
	FindProducts(query QueryParams) ([]Product, *localError.GlobalError)
	DeleteProduct(id string, userId string) *localError.GlobalError
	GetPublicProducts(filter ProductFilter) ([]Product, *localError.GlobalError)
}

type productUsecase struct {
	repo IProductRepository
}

func NewProductUsecase(repo IProductRepository) IProductUsecase {
	return &productUsecase{
		repo: repo,
	}
}

func (uc *productUsecase) FindProducts(query QueryParams) ([]Product, *localError.GlobalError) {
	products, err := uc.repo.FindAllProduct(query)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (uc *productUsecase) CreateProduct(req CreateProductRequest) (*Product, *localError.GlobalError) {
	if !validation.IsValidURL(req.ImageURL) {
		return nil, localError.ErrBadRequest("Invalid image URL", nil)
	}

	product, err := uc.repo.FindBySku(req.Sku)
	if err != nil {
		return nil, err
	}

	if product.ID != "" {
		return nil, localError.ErrBadRequest("SKU already exists", nil)
	}

	err = uc.repo.CreateProduct(req)
	if err != nil {
		return nil, err
	}

	product, err = uc.repo.FindBySku(req.Sku)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (uc *productUsecase) DeleteProduct(id string, userId string) *localError.GlobalError {
	if userId == "" {
		return localError.ErrUnauthorized("Unauthorized", nil)
	}

	product, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}

	if product.ID == "" {
		return localError.ErrNotFound("Product not found", nil)
	}

	err = uc.repo.DeleteProduct(id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *productUsecase) GetPublicProducts(filter ProductFilter) ([]Product, *localError.GlobalError) {
	products, err := uc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	return products, nil
}
