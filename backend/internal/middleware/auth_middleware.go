package middleware

import (
    "strings"

    "github.com/gin-gonic/gin"
    jwtpkg "github.com/myfarism/finance-tracker/pkg/jwt"
    "github.com/myfarism/finance-tracker/pkg/response"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            response.Unauthorized(c, "Authorization header missing")
            c.Abort()
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := jwtpkg.ValidateToken(tokenStr)
        if err != nil {
            response.Unauthorized(c, err.Error())
            c.Abort()
            return
        }

        // Simpan data user ke context agar bisa diakses handler lain
        c.Set("userID", claims.UserID)
        c.Set("email", claims.Email)
        c.Next()
    }
}
