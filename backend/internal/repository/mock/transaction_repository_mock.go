package mock

import (
    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "github.com/myfarism/finance-tracker/internal/repository"
    "github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
    mock.Mock
}

func (m *MockTransactionRepository) Create(tx *domain.Transaction) error {
    args := m.Called(tx)
    return args.Error(0)
}

func (m *MockTransactionRepository) FindAllByUser(userID uuid.UUID, filter repository.TransactionFilter) ([]domain.Transaction, error) {
    args := m.Called(userID, filter)
    return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByID(id uuid.UUID, userID uuid.UUID) (*domain.Transaction, error) {
    args := m.Called(id, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(tx *domain.Transaction) error {
    args := m.Called(tx)
    return args.Error(0)
}

func (m *MockTransactionRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
    args := m.Called(id, userID)
    return args.Error(0)
}

func (m *MockTransactionRepository) GetSummaryByUser(userID uuid.UUID, month, year int) (float64, float64, error) {
    args := m.Called(userID, month, year)
    return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}
