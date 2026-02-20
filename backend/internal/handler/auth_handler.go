package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/myfarism/finance-tracker/internal/service"
    "github.com/myfarism/finance-tracker/pkg/response"
)

type AuthHandler struct {
    authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
    return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var input service.RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    if err := h.authService.Register(input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    response.OK(c, "Kode OTP telah dikirim ke email kamu", nil)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
    var input service.VerifyOTPInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    result, err := h.authService.VerifyOTP(input)
    if err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    response.OK(c, "Verifikasi berhasil", result)
}

func (h *AuthHandler) ResendOTP(c *gin.Context) {
    var input service.ResendOTPInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    if err := h.authService.ResendOTP(input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    response.OK(c, "Kode OTP baru telah dikirim", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
    var input service.LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    result, err := h.authService.Login(input)
    if err != nil {
        response.Unauthorized(c, err.Error())
        return
    }

    result.User.Password = ""
    response.OK(c, "Login berhasil", result)
}
