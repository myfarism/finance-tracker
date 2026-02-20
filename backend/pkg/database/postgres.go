package database

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/myfarism/finance-tracker/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate semua tabel
    db.AutoMigrate(
        &domain.User{},
        &domain.Category{},
        &domain.Transaction{},
        &domain.Budget{},
    )
    seedCategories(db)

    DB = db
    log.Println("✅ Database connected successfully")
}

func seedCategories(db *gorm.DB) {
    var count int64
    db.Model(&domain.Category{}).Count(&count)
    if count > 0 {
        return // sudah ada data, skip
    }

    categories := []domain.Category{
        {ID: uuid.New(), Name: "Gaji", Icon: "💼"},
        {ID: uuid.New(), Name: "Freelance", Icon: "💻"},
        {ID: uuid.New(), Name: "Investasi", Icon: "📈"},
        {ID: uuid.New(), Name: "Makanan", Icon: "🍜"},
        {ID: uuid.New(), Name: "Transportasi", Icon: "🚗"},
        {ID: uuid.New(), Name: "Belanja", Icon: "🛍️"},
        {ID: uuid.New(), Name: "Kesehatan", Icon: "🏥"},
        {ID: uuid.New(), Name: "Hiburan", Icon: "🎮"},
        {ID: uuid.New(), Name: "Tagihan", Icon: "📄"},
        {ID: uuid.New(), Name: "Lainnya", Icon: "📦"},
    }

    db.Create(&categories)
    log.Println("✅ Categories seeded")
}
