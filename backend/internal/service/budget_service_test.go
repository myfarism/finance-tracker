package service_test

import (
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "github.com/myfarism/finance-tracker/internal/repository" // ← tambahkan kembali
    repomock "github.com/myfarism/finance-tracker/internal/repository/mock"
    "github.com/myfarism/finance-tracker/internal/service"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestUpsertBudget_Success(t *testing.T) {
    mockBudgetRepo := new(repomock.MockBudgetRepository)
    mockTxRepo     := new(repomock.MockTransactionRepository)
    svc := service.NewBudgetService(mockBudgetRepo, mockTxRepo)

    userID := uuid.New()
    catID  := uuid.New()

    mockBudgetRepo.On("Upsert", mock.AnythingOfType("*domain.Budget")).Return(nil)
    mockBudgetRepo.On("FindByUserAndMonth", userID, 2, 2026).
        Return([]domain.Budget{
            {
                ID:         uuid.New(),
                UserID:     userID,
                CategoryID: catID,
                Category:   domain.Category{ID: catID, Name: "Makan", Icon: "🍜"},
                Amount:     500000,
                Month:      2,
                Year:       2026,
            },
        }, nil)

    mockTxRepo.On("FindAllByUser", userID,
        mock.AnythingOfType("repository.TransactionFilter")). // ← fix
        Return([]domain.Transaction{}, nil)

    result, err := svc.Upsert(userID, service.UpsertBudgetInput{
        CategoryID: catID.String(),
        Amount:     500000,
        Month:      2,
        Year:       2026,
    })

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, 500000.0, result.Amount)
    mockBudgetRepo.AssertExpectations(t)
}

func TestUpsertBudget_InvalidCategoryID(t *testing.T) {
    mockBudgetRepo := new(repomock.MockBudgetRepository)
    mockTxRepo     := new(repomock.MockTransactionRepository)
    svc := service.NewBudgetService(mockBudgetRepo, mockTxRepo)

    _, err := svc.Upsert(uuid.New(), service.UpsertBudgetInput{
        CategoryID: "bukan-uuid",
        Amount:     500000,
        Month:      2,
        Year:       2026,
    })

    assert.Error(t, err)
    assert.Equal(t, "invalid category_id", err.Error())
    mockBudgetRepo.AssertNotCalled(t, "Upsert")
}

func TestUpsertBudget_DatabaseError(t *testing.T) {
    mockBudgetRepo := new(repomock.MockBudgetRepository)
    mockTxRepo     := new(repomock.MockTransactionRepository)
    svc := service.NewBudgetService(mockBudgetRepo, mockTxRepo)

    catID := uuid.New()

    mockBudgetRepo.On("Upsert", mock.AnythingOfType("*domain.Budget")).
        Return(assert.AnError)

    _, err := svc.Upsert(uuid.New(), service.UpsertBudgetInput{
        CategoryID: catID.String(),
        Amount:     500000,
        Month:      2,
        Year:       2026,
    })

    assert.Error(t, err)
    mockBudgetRepo.AssertNotCalled(t, "FindByUserAndMonth")
}

func TestGetBudgetByMonth_WithSpentCalculation(t *testing.T) {
    mockBudgetRepo := new(repomock.MockBudgetRepository)
    mockTxRepo     := new(repomock.MockTransactionRepository)
    svc := service.NewBudgetService(mockBudgetRepo, mockTxRepo)

    userID := uuid.New()
    catID  := uuid.New()

    mockBudgetRepo.On("FindByUserAndMonth", userID, 2, 2026).
        Return([]domain.Budget{
            {
                ID:         uuid.New(),
                UserID:     userID,
                CategoryID: catID,
                Category:   domain.Category{ID: catID, Name: "Makan"},
                Amount:     500000,
                Month:      2,
                Year:       2026,
            },
        }, nil)

    mockTxRepo.On("FindAllByUser", userID,
        mock.AnythingOfType("repository.TransactionFilter")). // ← fix
        Return([]domain.Transaction{
            {
                ID:         uuid.New(),
                UserID:     userID,
                CategoryID: catID,
                Type:       domain.Expense,
                Amount:     150000,
                Date:       time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC),
            },
        }, nil)

    budgets, err := svc.GetByMonth(userID, 2, 2026)

    assert.NoError(t, err)
    assert.Len(t, budgets, 1)
    assert.Equal(t, 150000.0, budgets[0].Spent)
    assert.Equal(t, 350000.0, budgets[0].Remaining)
    assert.False(t, budgets[0].IsOver)
}

func TestGetBudgetByMonth_OverBudget(t *testing.T) {
    mockBudgetRepo := new(repomock.MockBudgetRepository)
    mockTxRepo     := new(repomock.MockTransactionRepository)
    svc := service.NewBudgetService(mockBudgetRepo, mockTxRepo)

    userID := uuid.New()
    catID  := uuid.New()

    mockBudgetRepo.On("FindByUserAndMonth", userID, 2, 2026).
        Return([]domain.Budget{
            {
                ID:         uuid.New(),
                UserID:     userID,
                CategoryID: catID,
                Amount:     200000,
                Month:      2,
                Year:       2026,
            },
        }, nil)

    mockTxRepo.On("FindAllByUser", userID,
        mock.AnythingOfType("repository.TransactionFilter")). // ← fix
        Return([]domain.Transaction{
            {
                CategoryID: catID,
                Type:       domain.Expense,
                Amount:     350000,
                Date:       time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC),
            },
        }, nil)

    budgets, err := svc.GetByMonth(userID, 2, 2026)

    assert.NoError(t, err)
    assert.True(t, budgets[0].IsOver)
    assert.Equal(t, 350000.0, budgets[0].Spent)
    assert.Equal(t, -150000.0, budgets[0].Remaining)
}

func TestGetBudgetByMonth_EmptyBudgets(t *testing.T) {
    mockBudgetRepo := new(repomock.MockBudgetRepository)
    mockTxRepo     := new(repomock.MockTransactionRepository)
    svc := service.NewBudgetService(mockBudgetRepo, mockTxRepo)

    userID := uuid.New()

    mockBudgetRepo.On("FindByUserAndMonth", userID, 2, 2026).
        Return([]domain.Budget{}, nil)

    mockTxRepo.On("FindAllByUser", userID,
        mock.AnythingOfType("repository.TransactionFilter")). // ← fix
        Return([]domain.Transaction{}, nil)

    budgets, err := svc.GetByMonth(userID, 2, 2026)

    assert.NoError(t, err)
    assert.Empty(t, budgets)
}

func TestDeleteBudget_Success(t *testing.T) {
    mockBudgetRepo := new(repomock.MockBudgetRepository)
    mockTxRepo     := new(repomock.MockTransactionRepository)
    svc := service.NewBudgetService(mockBudgetRepo, mockTxRepo)

    budgetID := uuid.New()
    userID   := uuid.New()

    mockBudgetRepo.On("Delete", budgetID, userID).Return(nil)

    err := svc.Delete(budgetID.String(), userID)

    assert.NoError(t, err)
    mockBudgetRepo.AssertExpectations(t)
}

func TestDeleteBudget_InvalidID(t *testing.T) {
    mockBudgetRepo := new(repomock.MockBudgetRepository)
    mockTxRepo     := new(repomock.MockTransactionRepository)
    svc := service.NewBudgetService(mockBudgetRepo, mockTxRepo)

    err := svc.Delete("bukan-uuid", uuid.New())

    assert.Error(t, err)
    assert.Equal(t, "invalid budget id", err.Error())
    mockBudgetRepo.AssertNotCalled(t, "Delete")
}

// Pastikan import repository dipakai agar tidak error "imported and not used"
var _ = repository.TransactionFilter{}
