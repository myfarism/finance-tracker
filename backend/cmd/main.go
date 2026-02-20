package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/myfarism/finance-tracker/internal/handler"
    "github.com/myfarism/finance-tracker/internal/middleware"
    "github.com/myfarism/finance-tracker/internal/repository"
    "github.com/myfarism/finance-tracker/internal/service"
    "github.com/myfarism/finance-tracker/pkg/database"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    database.Connect()

    // Repositories
    userRepo := repository.NewUserRepository(database.DB)
    catRepo  := repository.NewCategoryRepository(database.DB)
    txRepo   := repository.NewTransactionRepository(database.DB)
    budgetRepo := repository.NewBudgetRepository(database.DB)


    // Services
    authSvc := service.NewAuthService(userRepo)
    catSvc  := service.NewCategoryService(catRepo)
    txSvc   := service.NewTransactionService(txRepo, catRepo)
    budgetSvc     := service.NewBudgetService(budgetRepo, txRepo)

    // Handlers
    authHandler := handler.NewAuthHandler(authSvc)
    catHandler  := handler.NewCategoryHandler(catSvc)
    txHandler   := handler.NewTransactionHandler(txSvc)
    budgetHandler := handler.NewBudgetHandler(budgetSvc)

    r := gin.Default()

    // CORS
    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    api := r.Group("/api/v1")
    {
        // Public routes
        auth := api.Group("/auth")
        auth.POST("/register", authHandler.Register)
        auth.POST("/verify-otp", authHandler.VerifyOTP)
        auth.POST("/resend-otp", authHandler.ResendOTP)
        auth.POST("/login", authHandler.Login)

        // Protected routes
        protected := api.Group("/")
        protected.Use(middleware.AuthMiddleware())
        {
            // Categories
            protected.GET("/categories", catHandler.GetAll)
            protected.POST("/categories", catHandler.Create)

            // Transactions
            protected.POST("/transactions", txHandler.Create)
            protected.GET("/transactions", txHandler.GetAll)
            protected.GET("/transactions/:id", txHandler.GetByID)
            protected.PUT("/transactions/:id", txHandler.Update)
            protected.DELETE("/transactions/:id", txHandler.Delete)

            // Summary (untuk dashboard chart)
            protected.GET("/transactions/summary", txHandler.GetSummary)

            // Budgets
            protected.GET("/budgets", budgetHandler.GetByMonth)
            protected.POST("/budgets", budgetHandler.Upsert)
            protected.DELETE("/budgets/:id", budgetHandler.Delete)
        }
    }

    port := os.Getenv("PORT")
    log.Printf("🚀 Server running on port %s", port)
    r.Run(":" + port)
}
