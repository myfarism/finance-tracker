package service

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "github.com/myfarism/finance-tracker/internal/repository"
)

type CreateTransactionInput struct {
    CategoryID  string  `json:"category_id" binding:"required"`
    Type        string  `json:"type" binding:"required,oneof=income expense"`
    Amount      float64 `json:"amount" binding:"required,gt=0"`
    Description string  `json:"description"`
    Date        string  `json:"date" binding:"required"` // format: "2006-01-02"
}

type UpdateTransactionInput struct {
    CategoryID  string  `json:"category_id"`
    Type        string  `json:"type" binding:"omitempty,oneof=income expense"`
    Amount      float64 `json:"amount" binding:"omitempty,gt=0"`
    Description string  `json:"description"`
    Date        string  `json:"date"`
}

type SummaryResponse struct {
    Income  float64 `json:"income"`
    Expense float64 `json:"expense"`
    Balance float64 `json:"balance"`
    Month   int     `json:"month"`
    Year    int     `json:"year"`
}

type TransactionService interface {
    Create(userID uuid.UUID, input CreateTransactionInput) (*domain.Transaction, error)
    GetAll(userID uuid.UUID, filter repository.TransactionFilter) ([]domain.Transaction, error)
    GetByID(id string, userID uuid.UUID) (*domain.Transaction, error)
    Update(id string, userID uuid.UUID, input UpdateTransactionInput) (*domain.Transaction, error)
    Delete(id string, userID uuid.UUID) error
    GetSummary(userID uuid.UUID, month, year int) (*SummaryResponse, error)
}

type transactionService struct {
    txRepo  repository.TransactionRepository
    catRepo repository.CategoryRepository
}

func NewTransactionService(txRepo repository.TransactionRepository, catRepo repository.CategoryRepository) TransactionService {
    return &transactionService{txRepo, catRepo}
}

func (s *transactionService) Create(userID uuid.UUID, input CreateTransactionInput) (*domain.Transaction, error) {
    catID, err := uuid.Parse(input.CategoryID)
    if err != nil {
        return nil, errors.New("invalid category_id")
    }

    // Validasi kategori ada
    _, err = s.catRepo.FindByID(catID)
    if err != nil {
        return nil, errors.New("category not found")
    }

    date, err := time.Parse("2006-01-02", input.Date)
    if err != nil {
        return nil, errors.New("invalid date format, use YYYY-MM-DD")
    }

    tx := &domain.Transaction{
        ID:          uuid.New(),
        UserID:      userID,
        CategoryID:  catID,
        Type:        domain.TransactionType(input.Type),
        Amount:      input.Amount,
        Description: input.Description,
        Date:        date,
    }

    if err := s.txRepo.Create(tx); err != nil {
        return nil, err
    }

    // Reload dengan relasi Category
    return s.txRepo.FindByID(tx.ID, userID)
}

func (s *transactionService) GetAll(userID uuid.UUID, filter repository.TransactionFilter) ([]domain.Transaction, error) {
    return s.txRepo.FindAllByUser(userID, filter)
}

func (s *transactionService) GetByID(id string, userID uuid.UUID) (*domain.Transaction, error) {
    txID, err := uuid.Parse(id)
    if err != nil {
        return nil, errors.New("invalid transaction id")
    }
    return s.txRepo.FindByID(txID, userID)
}

func (s *transactionService) Update(id string, userID uuid.UUID, input UpdateTransactionInput) (*domain.Transaction, error) {
    txID, err := uuid.Parse(id)
    if err != nil {
        return nil, errors.New("invalid transaction id")
    }

    tx, err := s.txRepo.FindByID(txID, userID)
    if err != nil {
        return nil, errors.New("transaction not found")
    }

    if input.CategoryID != "" {
        catID, err := uuid.Parse(input.CategoryID)
        if err != nil {
            return nil, errors.New("invalid category_id")
        }
        tx.CategoryID = catID
    }
    if input.Type != "" {
        tx.Type = domain.TransactionType(input.Type)
    }
    if input.Amount > 0 {
        tx.Amount = input.Amount
    }
    if input.Description != "" {
        tx.Description = input.Description
    }
    if input.Date != "" {
        date, err := time.Parse("2006-01-02", input.Date)
        if err != nil {
            return nil, errors.New("invalid date format, use YYYY-MM-DD")
        }
        tx.Date = date
    }

    if err := s.txRepo.Update(tx); err != nil {
        return nil, err
    }

    return s.txRepo.FindByID(txID, userID)
}

func (s *transactionService) Delete(id string, userID uuid.UUID) error {
    txID, err := uuid.Parse(id)
    if err != nil {
        return errors.New("invalid transaction id")
    }

    _, err = s.txRepo.FindByID(txID, userID)
    if err != nil {
        return errors.New("transaction not found")
    }

    return s.txRepo.Delete(txID, userID)
}

func (s *transactionService) GetSummary(userID uuid.UUID, month, year int) (*SummaryResponse, error) {
    income, expense, err := s.txRepo.GetSummaryByUser(userID, month, year)
    if err != nil {
        return nil, err
    }

    return &SummaryResponse{
        Income:  income,
        Expense: expense,
        Balance: income - expense,
        Month:   month,
        Year:    year,
    }, nil
}
