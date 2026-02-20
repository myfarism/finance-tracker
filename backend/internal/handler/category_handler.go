package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "github.com/myfarism/finance-tracker/internal/service"
    "github.com/myfarism/finance-tracker/pkg/response"
)

type CategoryHandler struct {
    catService service.CategoryService
}

func NewCategoryHandler(catService service.CategoryService) *CategoryHandler {
    return &CategoryHandler{catService}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
    categories, err := h.catService.GetAll()
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.OK(c, "Categories fetched", categories)
}

func (h *CategoryHandler) Create(c *gin.Context) {
    var input struct {
        Name string `json:"name" binding:"required"`
        Icon string `json:"icon"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    cat := &domain.Category{
        ID:   uuid.New(),
        Name: input.Name,
        Icon: input.Icon,
    }

    if err := h.catService.Create(cat); err != nil {
        response.InternalError(c, err.Error())
        return
    }

    response.Created(c, "Category created", cat)
}
