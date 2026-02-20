package handler

import (
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/repository"
    "github.com/myfarism/finance-tracker/internal/service"
    "github.com/myfarism/finance-tracker/pkg/response"
)

type TransactionHandler struct {
    txService service.TransactionService
}

func NewTransactionHandler(txService service.TransactionService) *TransactionHandler {
    return &TransactionHandler{txService}
}

// Helper ambil userID dari JWT context
func getUserID(c *gin.Context) uuid.UUID {
    return c.MustGet("userID").(uuid.UUID)
}

func (h *TransactionHandler) Create(c *gin.Context) {
    var input service.CreateTransactionInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    tx, err := h.txService.Create(getUserID(c), input)
    if err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    response.Created(c, "Transaction created", tx)
}

func (h *TransactionHandler) GetAll(c *gin.Context) {
    filter := repository.TransactionFilter{
        Type:       c.Query("type"),
        CategoryID: c.Query("category_id"),
        Search:     c.Query("search"),
    }

    // Parse tanggal jika ada
    if start := c.Query("start_date"); start != "" {
        t, err := time.Parse("2006-01-02", start)
        if err == nil {
            filter.StartDate = &t
        }
    }
    if end := c.Query("end_date"); end != "" {
        t, err := time.Parse("2006-01-02", end)
        if err == nil {
            filter.EndDate = &t
        }
    }

    transactions, err := h.txService.GetAll(getUserID(c), filter)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }

    response.OK(c, "Transactions fetched", transactions)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
    tx, err := h.txService.GetByID(c.Param("id"), getUserID(c))
    if err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    response.OK(c, "Transaction fetched", tx)
}

func (h *TransactionHandler) Update(c *gin.Context) {
    var input service.UpdateTransactionInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    tx, err := h.txService.Update(c.Param("id"), getUserID(c), input)
    if err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    response.OK(c, "Transaction updated", tx)
}

func (h *TransactionHandler) Delete(c *gin.Context) {
    if err := h.txService.Delete(c.Param("id"), getUserID(c)); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    response.OK(c, "Transaction deleted", nil)
}

func (h *TransactionHandler) GetSummary(c *gin.Context) {
    now := time.Now()
    month := now.Month()
    year := now.Year()

    if m := c.Query("month"); m != "" {
        if val, err := strconv.Atoi(m); err == nil {
            month = time.Month(val)
        }
    }
    if y := c.Query("year"); y != "" {
        if val, err := strconv.Atoi(y); err == nil {
            year = val
        }
    }

    summary, err := h.txService.GetSummary(getUserID(c), int(month), year)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }

    response.OK(c, "Summary fetched", summary)
}
