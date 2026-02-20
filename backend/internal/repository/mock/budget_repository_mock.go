package mock

import (
    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "github.com/stretchr/testify/mock"
)

type MockBudgetRepository struct {
    mock.Mock
}

func (m *MockBudgetRepository) Upsert(budget *domain.Budget) error {
    args := m.Called(budget)
    return args.Error(0)
}

func (m *MockBudgetRepository) FindByUserAndMonth(userID uuid.UUID, month, year int) ([]domain.Budget, error) {
    args := m.Called(userID, month, year)
    return args.Get(0).([]domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) FindByID(id uuid.UUID, userID uuid.UUID) (*domain.Budget, error) {
    args := m.Called(id, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
    args := m.Called(id, userID)
    return args.Error(0)
}
