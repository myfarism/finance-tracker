package service

import (
    "github.com/myfarism/finance-tracker/internal/domain"
    "github.com/myfarism/finance-tracker/internal/repository"
)

type CategoryService interface {
    GetAll() ([]domain.Category, error)
    Create(cat *domain.Category) error
}

type categoryService struct {
    catRepo repository.CategoryRepository
}

func NewCategoryService(catRepo repository.CategoryRepository) CategoryService {
    return &categoryService{catRepo}
}

func (s *categoryService) GetAll() ([]domain.Category, error) {
    return s.catRepo.FindAll()
}

func (s *categoryService) Create(cat *domain.Category) error {
    return s.catRepo.Create(cat)
}
