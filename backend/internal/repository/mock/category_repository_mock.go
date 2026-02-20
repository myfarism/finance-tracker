// category_repository_mock.go
package mock

import (
    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
    mock.Mock
}

func (m *MockCategoryRepository) FindAll() ([]domain.Category, error) {
    args := m.Called()
    return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByID(id uuid.UUID) (*domain.Category, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Create(category *domain.Category) error {
    args := m.Called(category)
    return args.Error(0)
}
