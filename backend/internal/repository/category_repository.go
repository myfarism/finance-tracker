package repository

import (
    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "gorm.io/gorm"
)

type CategoryRepository interface {
    FindAll() ([]domain.Category, error)
    FindByID(id uuid.UUID) (*domain.Category, error)
    Create(category *domain.Category) error
}

type categoryRepository struct {
    db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
    return &categoryRepository{db}
}

func (r *categoryRepository) FindAll() ([]domain.Category, error) {
    var categories []domain.Category
    err := r.db.Find(&categories).Error
    return categories, err
}

func (r *categoryRepository) FindByID(id uuid.UUID) (*domain.Category, error) {
    var category domain.Category
    err := r.db.First(&category, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &category, nil
}

func (r *categoryRepository) Create(category *domain.Category) error {
    return r.db.Create(category).Error
}
