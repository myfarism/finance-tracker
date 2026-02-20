package repository

import (
    "time"

    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "gorm.io/gorm"
)

type TransactionFilter struct {
    Type       string
    CategoryID string
    StartDate  *time.Time
    EndDate    *time.Time
    Search     string
}

type TransactionRepository interface {
    Create(tx *domain.Transaction) error
    FindAllByUser(userID uuid.UUID, filter TransactionFilter) ([]domain.Transaction, error)
    FindByID(id uuid.UUID, userID uuid.UUID) (*domain.Transaction, error)
    Update(tx *domain.Transaction) error
    Delete(id uuid.UUID, userID uuid.UUID) error
    GetSummaryByUser(userID uuid.UUID, month, year int) (income, expense float64, err error)
}

type transactionRepository struct {
    db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
    return &transactionRepository{db}
}

func (r *transactionRepository) Create(tx *domain.Transaction) error {
    return r.db.Create(tx).Error
}

func (r *transactionRepository) FindAllByUser(userID uuid.UUID, filter TransactionFilter) ([]domain.Transaction, error) {
    var transactions []domain.Transaction

    query := r.db.Where("user_id = ?", userID).
        Preload("Category").
        Order("date DESC")

    if filter.Type != "" {
        query = query.Where("type = ?", filter.Type)
    }
    if filter.CategoryID != "" {
        query = query.Where("category_id = ?", filter.CategoryID)
    }
    if filter.StartDate != nil {
        query = query.Where("date >= ?", filter.StartDate)
    }
    if filter.EndDate != nil {
        query = query.Where("date <= ?", filter.EndDate)
    }
    if filter.Search != "" {
        query = query.Where("description ILIKE ?", "%"+filter.Search+"%")
    }

    err := query.Find(&transactions).Error
    return transactions, err
}

func (r *transactionRepository) FindByID(id uuid.UUID, userID uuid.UUID) (*domain.Transaction, error) {
    var tx domain.Transaction
    err := r.db.Where("id = ? AND user_id = ?", id, userID).
        Preload("Category").
        First(&tx).Error
    if err != nil {
        return nil, err
    }
    return &tx, nil
}

func (r *transactionRepository) Update(tx *domain.Transaction) error {
    return r.db.Save(tx).Error
}

func (r *transactionRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
    return r.db.Where("id = ? AND user_id = ?", id, userID).
        Delete(&domain.Transaction{}).Error
}

func (r *transactionRepository) GetSummaryByUser(userID uuid.UUID, month, year int) (float64, float64, error) {
    var income, expense float64

    r.db.Model(&domain.Transaction{}).
        Where("user_id = ? AND type = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
            userID, domain.Income, month, year).
        Select("COALESCE(SUM(amount), 0)").
        Scan(&income)

    r.db.Model(&domain.Transaction{}).
        Where("user_id = ? AND type = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
            userID, domain.Expense, month, year).
        Select("COALESCE(SUM(amount), 0)").
        Scan(&expense)

    return income, expense, nil
}
