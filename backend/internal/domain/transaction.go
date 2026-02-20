package domain

import (
    "time"
    "github.com/google/uuid"
)

type TransactionType string

const (
    Income  TransactionType = "income"
    Expense TransactionType = "expense"
)

type Category struct {
    ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    Name string    `gorm:"not null" json:"name"`
    Icon string    `json:"icon"`
}

type Transaction struct {
    ID          uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    UserID      uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
    User        User            `json:"-"`
    CategoryID  uuid.UUID       `gorm:"type:uuid;not null" json:"category_id"`
    Category    Category        `json:"category"`
    Type        TransactionType `gorm:"type:varchar(10);not null" json:"type"`
    Amount      float64         `gorm:"not null" json:"amount"`
    Description string          `json:"description"`
    Date        time.Time       `gorm:"not null;default:now()" json:"date"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
}
