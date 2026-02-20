package service

import (
    "errors"
    "os"
    "time"

    "github.com/google/uuid"
    "github.com/myfarism/finance-tracker/internal/domain"
    "github.com/myfarism/finance-tracker/internal/repository"
    jwtpkg "github.com/myfarism/finance-tracker/pkg/jwt"
    "github.com/myfarism/finance-tracker/pkg/mailer"
    "github.com/myfarism/finance-tracker/pkg/otp"
    "golang.org/x/crypto/bcrypt"
    "strconv"
)

type RegisterInput struct {
    Name     string `json:"name" binding:"required,min=2"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type VerifyOTPInput struct {
    Email string `json:"email" binding:"required,email"`
    Code  string `json:"code" binding:"required,len=6"`
}

type ResendOTPInput struct {
    Email string `json:"email" binding:"required,email"`
}

type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Token string      `json:"token"`
    User  domain.User `json:"user"`
}

type AuthService interface {
    Register(input RegisterInput) error              // kirim OTP, belum return token
    VerifyOTP(input VerifyOTPInput) (*AuthResponse, error) // verifikasi → return token
    ResendOTP(input ResendOTPInput) error
    Login(input LoginInput) (*AuthResponse, error)
}

type authService struct {
    userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
    return &authService{userRepo}
}

func (s *authService) Register(input RegisterInput) error {
    // Cek email sudah terdaftar
    existing, _ := s.userRepo.FindByEmail(input.Email)
    if existing != nil {
        if existing.IsVerified {
            return errors.New("email sudah terdaftar")
        }
        // Email ada tapi belum verified → kirim ulang OTP
        return s.sendOTP(existing.Name, existing.Email)
    }

    // Hash password
    hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := &domain.User{
        ID:         uuid.New(),
        Name:       input.Name,
        Email:      input.Email,
        Password:   string(hashed),
        IsVerified: false,
    }

    if err := s.userRepo.Create(user); err != nil {
        return err
    }

    return s.sendOTP(user.Name, user.Email)
}

func (s *authService) sendOTP(name, email string) error {
    expiryMinutes, _ := strconv.Atoi(os.Getenv("OTP_EXPIRY_MINUTES"))
    if expiryMinutes == 0 {
        expiryMinutes = 5
    }

    code := otp.Generate()
    otp.Save(email, code, time.Duration(expiryMinutes)*time.Minute)

    return mailer.SendOTP(email, name, code)
}

func (s *authService) VerifyOTP(input VerifyOTPInput) (*AuthResponse, error) {
    // Cek OTP valid
    if !otp.Verify(input.Email, input.Code) {
        return nil, errors.New("kode OTP tidak valid atau sudah kadaluarsa")
    }

    // Ambil user
    user, err := s.userRepo.FindByEmail(input.Email)
    if err != nil {
        return nil, errors.New("user tidak ditemukan")
    }

    // Update status verified
    if err := s.userRepo.UpdateVerified(user.ID, true); err != nil {
        return nil, err
    }
    user.IsVerified = true

    // Generate token
    token, err := jwtpkg.GenerateToken(user.ID, user.Email)
    if err != nil {
        return nil, err
    }

    return &AuthResponse{Token: token, User: *user}, nil
}

func (s *authService) ResendOTP(input ResendOTPInput) error {
    user, err := s.userRepo.FindByEmail(input.Email)
    if err != nil {
        return errors.New("email tidak ditemukan")
    }
    if user.IsVerified {
        return errors.New("akun sudah terverifikasi")
    }
    return s.sendOTP(user.Name, user.Email)
}

func (s *authService) Login(input LoginInput) (*AuthResponse, error) {
    user, err := s.userRepo.FindByEmail(input.Email)
    if err != nil {
        return nil, errors.New("email atau password salah")
    }

    // Cek sudah verified
    if !user.IsVerified {
        return nil, errors.New("akun belum diverifikasi, cek email kamu")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        return nil, errors.New("email atau password salah")
    }

    token, err := jwtpkg.GenerateToken(user.ID, user.Email)
    if err != nil {
        return nil, err
    }

    return &AuthResponse{Token: token, User: *user}, nil
}
