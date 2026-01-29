package services

import (
	"kasir-api/models"
	"kasir-api/repositories"

	"kasir-api/pkg/validator"
	// "github.com/go-playground/validator/v10"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) Create(data *models.Product) error {
	if err := validator.Validate.Struct(data); err != nil {
		return err
	}
	return s.repo.Create(data)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(data *models.Product) error {
	if err := validator.Validate.Struct(data); err != nil {
		return err
	}
	return s.repo.Update(data)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
