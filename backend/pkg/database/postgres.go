package database

import (
	"fmt"
	"log"
	"os"
    "strings"

	"github.com/google/uuid"
	"github.com/myfarism/finance-tracker/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    var dsn string

    if url := os.Getenv("DATABASE_URL"); url != "" {
        if !strings.Contains(url, "sslmode") {
            dsn = url + "?sslmode=require"
        } else {
            dsn = url
        }
    } else {
        host     := os.Getenv("DB_HOST")
        port     := os.Getenv("DB_PORT")
        user     := os.Getenv("DB_USER")
        password := os.Getenv("DB_PASSWORD")
        dbname   := os.Getenv("DB_NAME")

        if host == "" || port == "" || user == "" || dbname == "" {
            log.Fatal("Database configuration is incomplete")
        }

        dsn = fmt.Sprintf(
            "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
            host, port, user, password, dbname,
        )
    }

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
