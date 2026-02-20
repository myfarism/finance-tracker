package otp

import (
    "crypto/rand"
    "fmt"
    "time"

    "github.com/patrickmn/go-cache"
)

// Cache dengan expiry 5 menit, cleanup tiap 10 menit
var otpCache = cache.New(5*time.Minute, 10*time.Minute)

// Generate OTP 6 digit
func Generate() string {
    b := make([]byte, 3)
    rand.Read(b)
    return fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2]))[:6]
}

// Simpan OTP ke cache dengan key = email
func Save(email, code string, expiry time.Duration) {
    otpCache.Set(email, code, expiry)
}

// Verifikasi OTP
func Verify(email, code string) bool {
    stored, found := otpCache.Get(email)
    if !found {
        return false
    }
    if stored.(string) != code {
        return false
    }
    // Hapus setelah berhasil diverifikasi (one-time use)
    otpCache.Delete(email)
    return true
}
