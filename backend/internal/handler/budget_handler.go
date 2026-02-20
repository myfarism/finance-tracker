package handler

import (
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/myfarism/finance-tracker/internal/service"
    "github.com/myfarism/finance-tracker/pkg/response"
)

type BudgetHandler struct {
    budgetService service.BudgetService
}

func NewBudgetHandler(budgetService service.BudgetService) *BudgetHandler {
    return &BudgetHandler{budgetService}
}

func (h *BudgetHandler) Upsert(c *gin.Context) {
    var input service.UpsertBudgetInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    budget, err := h.budgetService.Upsert(getUserID(c), input)
    if err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    response.OK(c, "Budget disimpan", budget)
}

func (h *BudgetHandler) GetByMonth(c *gin.Context) {
    now := time.Now()
    month := int(now.Month())
    year := now.Year()

    if m := c.Query("month"); m != "" {
        if v, err := strconv.Atoi(m); err == nil {
            month = v
        }
    }
    if y := c.Query("year"); y != "" {
        if v, err := strconv.Atoi(y); err == nil {
            year = v
        }
    }

    budgets, err := h.budgetService.GetByMonth(getUserID(c), month, year)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }

    response.OK(c, "Budget fetched", budgets)
}

func (h *BudgetHandler) Delete(c *gin.Context) {
    if err := h.budgetService.Delete(c.Param("id"), getUserID(c)); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    response.OK(c, "Budget dihapus", nil)
}
