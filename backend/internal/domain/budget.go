package domain

import (
    "github.com/google/uuid"
)

type Budget struct {
    ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    UserID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_budget_unique" json:"user_id"`
    CategoryID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_budget_unique" json:"category_id"`
    Category   Category  `json:"category"`
    Amount     float64   `gorm:"not null" json:"amount"`
    Month      int       `gorm:"not null;uniqueIndex:idx_budget_unique" json:"month"`
    Year       int       `gorm:"not null;uniqueIndex:idx_budget_unique" json:"year"`

    Spent     float64 `gorm:"-" json:"spent"`
    Remaining float64 `gorm:"-" json:"remaining"`
    IsOver    bool    `gorm:"-" json:"is_over"`
}
