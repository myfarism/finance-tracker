package service_test

import (
    "testing"
    "time"
	"errors"

    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "github.com/myfarism/finance-tracker/internal/repository"
    repomock "github.com/myfarism/finance-tracker/internal/repository/mock"
    "github.com/myfarism/finance-tracker/internal/service"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestGetSummary_CalculatesBalanceCorrectly(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    userID := uuid.New()
    mockTxRepo.On("GetSummaryByUser", userID, 2, 2026).
        Return(5000000.0, 1500000.0, nil)

    summary, err := svc.GetSummary(userID, 2, 2026)

    assert.NoError(t, err)
    assert.Equal(t, 5000000.0, summary.Income)
    assert.Equal(t, 1500000.0, summary.Expense)
    assert.Equal(t, 3500000.0, summary.Balance) // 5jt - 1.5jt
    mockTxRepo.AssertExpectations(t)
}

func TestGetAll_WithFilter(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    userID := uuid.New()
    catID  := uuid.New()
    filter := repository.TransactionFilter{Type: "expense"}

    expectedTx := []domain.Transaction{
        {
            ID:     uuid.New(),
            UserID: userID,
            Type:   domain.Expense,
            Amount: 50000,
            Date:   time.Now(),
            Category: domain.Category{
                ID:   catID,
                Name: "Makanan",
                Icon: "🍜",
            },
        },
    }

    mockTxRepo.On("FindAllByUser", userID, filter).
        Return(expectedTx, nil)

    result, err := svc.GetAll(userID, filter)

    assert.NoError(t, err)
    assert.Len(t, result, 1)
    assert.Equal(t, domain.Expense, result[0].Type)
    mockTxRepo.AssertExpectations(t)
}

func TestDelete_TransactionNotFound(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    userID := uuid.New()
    randomID := uuid.New()

    mockTxRepo.On("FindByID", randomID, userID).
        Return(nil, assert.AnError)

    err := svc.Delete(randomID.String(), userID)

    assert.Error(t, err)
    assert.Equal(t, "transaction not found", err.Error())
    // Pastikan Delete di repo TIDAK dipanggil
    mockTxRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

// ──────────────────────────────────────────
// CREATE TRANSACTION TESTS
// ──────────────────────────────────────────

func TestCreate_Success(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    userID := uuid.New()
    catID  := uuid.New()
    cat    := &domain.Category{ID: catID, Name: "Makan", Icon: "🍜"}

    mockCatRepo.On("FindByID", catID).Return(cat, nil)
    mockTxRepo.On("Create", mock.AnythingOfType("*domain.Transaction")).Return(nil)
    mockTxRepo.On("FindByID", mock.AnythingOfType("uuid.UUID"), userID).
        Return(&domain.Transaction{
            ID:       uuid.New(),
            UserID:   userID,
            Category: *cat,
            Type:     domain.Expense,
            Amount:   50000,
        }, nil)

    result, err := svc.Create(userID, service.CreateTransactionInput{
        CategoryID:  catID.String(),
        Type:        "expense",
        Amount:      50000,
        Description: "Makan siang",
        Date:        "2026-02-20",
    })

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, domain.Expense, result.Type)
    assert.Equal(t, 50000.0, result.Amount)
    mockTxRepo.AssertExpectations(t)
    mockCatRepo.AssertExpectations(t)
}

func TestCreate_InvalidCategoryID(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    _, err := svc.Create(uuid.New(), service.CreateTransactionInput{
        CategoryID: "bukan-uuid-valid",
        Type:       "expense",
        Amount:     50000,
        Date:       "2026-02-20",
    })

    assert.Error(t, err)
    assert.Equal(t, "invalid category_id", err.Error())
    mockCatRepo.AssertNotCalled(t, "FindByID")
}

func TestCreate_CategoryNotFound(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    catID := uuid.New()
    mockCatRepo.On("FindByID", catID).Return(nil, errors.New("not found"))

    _, err := svc.Create(uuid.New(), service.CreateTransactionInput{
        CategoryID: catID.String(),
        Type:       "expense",
        Amount:     50000,
        Date:       "2026-02-20",
    })

    assert.Error(t, err)
    assert.Equal(t, "category not found", err.Error())
    mockTxRepo.AssertNotCalled(t, "Create")
}

func TestCreate_InvalidDateFormat(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    catID := uuid.New()
    mockCatRepo.On("FindByID", catID).
        Return(&domain.Category{ID: catID}, nil)

    _, err := svc.Create(uuid.New(), service.CreateTransactionInput{
        CategoryID: catID.String(),
        Type:       "expense",
        Amount:     50000,
        Date:       "20-02-2026", // ← format salah
    })

    assert.Error(t, err)
    assert.Equal(t, "invalid date format, use YYYY-MM-DD", err.Error())
}

// ──────────────────────────────────────────
// UPDATE TRANSACTION TESTS
// ──────────────────────────────────────────

func TestUpdate_Success(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    userID := uuid.New()
    txID   := uuid.New()
    catID  := uuid.New()

    existingTx := &domain.Transaction{
        ID:         txID,
        UserID:     userID,
        CategoryID: catID,
        Type:       domain.Expense,
        Amount:     50000,
        Date:       time.Now(),
    }

    updatedTx := &domain.Transaction{
        ID:     txID,
        UserID: userID,
        Amount: 75000, // amount berubah
    }

    mockTxRepo.On("FindByID", txID, userID).Return(existingTx, nil).Once()
    mockTxRepo.On("Update", mock.AnythingOfType("*domain.Transaction")).Return(nil)
    mockTxRepo.On("FindByID", txID, userID).Return(updatedTx, nil).Once()

    result, err := svc.Update(txID.String(), userID, service.UpdateTransactionInput{
        Amount: 75000,
    })

    assert.NoError(t, err)
    assert.Equal(t, 75000.0, result.Amount)
    mockTxRepo.AssertExpectations(t)
}

func TestUpdate_InvalidID(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    _, err := svc.Update("bukan-uuid", uuid.New(), service.UpdateTransactionInput{})

    assert.Error(t, err)
    assert.Equal(t, "invalid transaction id", err.Error())
}

// ──────────────────────────────────────────
// DELETE TRANSACTION TESTS
// ──────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    userID := uuid.New()
    txID   := uuid.New()

    mockTxRepo.On("FindByID", txID, userID).
        Return(&domain.Transaction{ID: txID, UserID: userID}, nil)
    mockTxRepo.On("Delete", txID, userID).Return(nil)

    err := svc.Delete(txID.String(), userID)

    assert.NoError(t, err)
    mockTxRepo.AssertExpectations(t)
}

func TestDelete_InvalidID(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    err := svc.Delete("bukan-uuid", uuid.New())

    assert.Error(t, err)
    assert.Equal(t, "invalid transaction id", err.Error())
}

// ──────────────────────────────────────────
// GET SUMMARY TESTS TAMBAHAN
// ──────────────────────────────────────────

func TestGetSummary_ZeroTransactions(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    userID := uuid.New()
    mockTxRepo.On("GetSummaryByUser", userID, 1, 2026).
        Return(0.0, 0.0, nil)

    summary, err := svc.GetSummary(userID, 1, 2026)

    assert.NoError(t, err)
    assert.Equal(t, 0.0, summary.Income)
    assert.Equal(t, 0.0, summary.Expense)
    assert.Equal(t, 0.0, summary.Balance)
}

func TestGetSummary_NegativeBalance(t *testing.T) {
    mockTxRepo  := new(repomock.MockTransactionRepository)
    mockCatRepo := new(repomock.MockCategoryRepository)
    svc := service.NewTransactionService(mockTxRepo, mockCatRepo)

    userID := uuid.New()
    // Pengeluaran lebih besar dari pemasukan
    mockTxRepo.On("GetSummaryByUser", userID, 2, 2026).
        Return(1000000.0, 3000000.0, nil)

    summary, err := svc.GetSummary(userID, 2, 2026)

    assert.NoError(t, err)
    assert.Equal(t, -2000000.0, summary.Balance)
}