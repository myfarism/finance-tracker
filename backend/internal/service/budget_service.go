package service

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "github.com/myfarism/finance-tracker/internal/repository"
)

type UpsertBudgetInput struct {
    CategoryID string  `json:"category_id" binding:"required"`
    Amount     float64 `json:"amount" binding:"required,gt=0"`
    Month      int     `json:"month" binding:"required,min=1,max=12"`
    Year       int     `json:"year" binding:"required,min=2000"`
}

type BudgetService interface {
    Upsert(userID uuid.UUID, input UpsertBudgetInput) (*domain.Budget, error)
    GetByMonth(userID uuid.UUID, month, year int) ([]domain.Budget, error)
    Delete(id string, userID uuid.UUID) error
}

type budgetService struct {
    budgetRepo repository.BudgetRepository
    txRepo     repository.TransactionRepository
}

func NewBudgetService(
    budgetRepo repository.BudgetRepository,
    txRepo repository.TransactionRepository,
) BudgetService {
    return &budgetService{budgetRepo, txRepo}
}

// Di budget_service.go — ganti fungsi Upsert
func (s *budgetService) Upsert(userID uuid.UUID, input UpsertBudgetInput) (*domain.Budget, error) {
    catID, err := uuid.Parse(input.CategoryID)
    if err != nil {
        return nil, errors.New("invalid category_id")
    }

    budget := &domain.Budget{
        ID:         uuid.New(),
        UserID:     userID,
        CategoryID: catID,
        Amount:     input.Amount,
        Month:      input.Month,
        Year:       input.Year,
    }

    if err := s.budgetRepo.Upsert(budget); err != nil {
        return nil, err
    }

    // ✅ Langsung GetByMonth untuk dapat data lengkap + spent
    budgets, err := s.GetByMonth(userID, input.Month, input.Year)
    if err != nil {
        return nil, err
    }

    // Cari budget yang baru saja di-upsert by CategoryID
    for i := range budgets {
        if budgets[i].CategoryID == catID {
            return &budgets[i], nil
        }
    }

    return nil, errors.New("budget not found")
}


func (s *budgetService) GetByMonth(userID uuid.UUID, month, year int) ([]domain.Budget, error) {
    budgets, err := s.budgetRepo.FindByUserAndMonth(userID, month, year)
    if err != nil {
        return nil, err
    }

    // Ambil semua transaksi bulan ini untuk hitung spent
    now := time.Now()
    startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
    endDate := startDate.AddDate(0, 1, -1)

    filter := repository.TransactionFilter{
        StartDate: &startDate,
        EndDate:   &endDate,
    }

    transactions, err := s.txRepo.FindAllByUser(userID, filter)
    if err != nil {
        return nil, err
    }

    // Hitung spent per kategori dari transaksi
    spentMap := map[uuid.UUID]float64{}
    for _, tx := range transactions {
        if tx.Type == domain.Expense {
            spentMap[tx.CategoryID] += tx.Amount
        }
    }

    // Enrich budget dengan data spent
    for i := range budgets {
        spent := spentMap[budgets[i].CategoryID]
        budgets[i].Spent = spent
        budgets[i].Remaining = budgets[i].Amount - spent
        budgets[i].IsOver = spent > budgets[i].Amount
    }

    return budgets, nil
}

func (s *budgetService) Delete(id string, userID uuid.UUID) error {
    budgetID, err := uuid.Parse(id)
    if err != nil {
        return errors.New("invalid budget id")
    }
    return s.budgetRepo.Delete(budgetID, userID)
}

// Helper: cari satu budget dari slice dan return dengan data spent
func (s *budgetService) enrichBudget(userID uuid.UUID, budgets []domain.Budget, targetID uuid.UUID) (*domain.Budget, error) {
    for i := range budgets {
        if budgets[i].ID == targetID {
            return &budgets[i], nil
        }
    }
    return nil, errors.New("budget not found")
}
