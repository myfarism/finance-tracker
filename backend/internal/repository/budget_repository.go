package repository

import (
    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

type BudgetRepository interface {
    Upsert(budget *domain.Budget) error
    FindByUserAndMonth(userID uuid.UUID, month, year int) ([]domain.Budget, error)
    FindByID(id uuid.UUID, userID uuid.UUID) (*domain.Budget, error)
    Delete(id uuid.UUID, userID uuid.UUID) error
}

type budgetRepository struct {
    db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) BudgetRepository {
    return &budgetRepository{db}
}

// Upsert: insert baru atau update jika kombinasi user+category+month+year sudah ada
func (r *budgetRepository) Upsert(budget *domain.Budget) error {
    return r.db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "user_id"}, {Name: "category_id"}, {Name: "month"}, {Name: "year"}},
        DoUpdates: clause.AssignmentColumns([]string{"amount"}),
    }).Create(budget).Error
}

func (r *budgetRepository) FindByUserAndMonth(userID uuid.UUID, month, year int) ([]domain.Budget, error) {
    var budgets []domain.Budget
    err := r.db.
        Where("user_id = ? AND month = ? AND year = ?", userID, month, year).
        Preload("Category").
        Find(&budgets).Error
    return budgets, err
}

func (r *budgetRepository) FindByID(id uuid.UUID, userID uuid.UUID) (*domain.Budget, error) {
    var budget domain.Budget
    err := r.db.
        Where("id = ? AND user_id = ?", id, userID).
        Preload("Category").
        First(&budget).Error
    if err != nil {
        return nil, err
    }
    return &budget, nil
}

func (r *budgetRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
    return r.db.
        Where("id = ? AND user_id = ?", id, userID).
        Delete(&domain.Budget{}).Error
}
