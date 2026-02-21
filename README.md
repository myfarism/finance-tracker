# 💰 Finance Tracker

Aplikasi manajemen keuangan pribadi berbasis web yang dibangun dengan arsitektur fullstack modern. Memungkinkan pengguna mencatat pemasukan dan pengeluaran, memantau saldo bulanan, mengatur budget per kategori, serta mengekspor data transaksi.

🔗 **Live Demo:** *(coming soon)*  
📦 **API Base URL:** *(coming soon)*

---

## ✨ Features

- 🔐 **Authentication** — Register, Login dengan JWT + verifikasi OTP via Email
- 📊 **Dashboard** — Ringkasan pemasukan, pengeluaran, dan saldo bulan ini
- 💸 **Transaksi** — CRUD transaksi dengan filter tipe, pencarian, dan kategori
- 📈 **Visualisasi** — Bar chart arus kas mingguan dan pie chart pengeluaran per kategori
- 🎯 **Budget** — Atur batas pengeluaran per kategori dengan progress bar real-time
- 📥 **Export CSV** — Unduh riwayat transaksi dalam format CSV
- 📱 **Responsif** — Mobile-first design, optimal di semua ukuran layar
- 🧪 **Unit Tested** — 33 test cases, coverage 70.9% pada service layer

---

## 🛠️ Tech Stack

| Layer | Teknologi |
|---|---|
| **Frontend** | React 18, TypeScript, Vite, TailwindCSS, Zustand, Recharts |
| **Backend** | Golang 1.26, Gin, GORM |
| **Database** | PostgreSQL |
| **Auth** | JWT + OTP (SMTP Gmail) |
| **Testing** | Testify, Mock |
| **Deploy** | Vercel (FE) · Railway (BE) |

---

## 📁 Project Structure

```
finance-tracker/
├── backend/
│   ├── cmd/                  # Entry point
│   ├── internal/
│   │   ├── domain/           # Struct & entity
│   │   ├── repository/       # Database layer
│   │   │   └── mock/         # Mock untuk testing
│   │   ├── service/          # Business logic + unit tests
│   │   └── handler/          # HTTP handler (Gin)
│   └── pkg/
│       ├── jwt/              # JWT helper
│       ├── otp/              # OTP cache
│       ├── mailer/           # SMTP email sender
│       └── database/         # PostgreSQL connection
└── frontend/
    └── src/
        ├── api/              # Axios API calls
        ├── components/       # Reusable components
        ├── pages/            # LoginPage, RegisterPage, DashboardPage
        ├── store/            # Zustand global state
        ├── types/            # TypeScript interfaces
        └── utils/            # Format helper, export CSV
```

---

## 🚀 Local Development

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Gmail App Password (untuk OTP)

### 1. Clone Repository

```bash
git clone https://github.com/myfarism/finance-tracker.git
cd finance-tracker
```

### 2. Setup Backend

```bash
cd backend

# Copy env
cp .env.example .env
```

Isi file `.env`:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=finance_tracker
JWT_SECRET=your_super_secret_key_minimum_32_chars
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=youremail@gmail.com
SMTP_PASSWORD=your_google_app_password
OTP_EXPIRY_MINUTES=5
```

```bash
# Jalankan server
go run cmd/main.go
# Server berjalan di http://localhost:8080
```

### 3. Setup Frontend

```bash
cd frontend

# Copy env
cp .env.example .env
```

Isi file `.env`:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

```bash
# Install dependencies
npm install

# Jalankan dev server
npm run dev
# App berjalan di http://localhost:5173
```

---

## 🧪 Running Tests

```bash
cd backend

# Jalankan semua unit test
go test ./internal/... -v

# Dengan coverage report
go test ./internal/... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Buka visual coverage di browser
go tool cover -html=coverage.out
```

**Hasil:**
```
--- PASS: TestLogin_Success
--- PASS: TestLogin_EmailNotFound
--- PASS: TestLogin_WrongPassword
--- PASS: TestUpsertBudget_Success
--- PASS: TestGetBudgetByMonth_OverBudget
... (33 tests total)

coverage: 70.9% of statements
```

---

## 🐳 Docker (Optional)

```bash
cd backend

# Build image
docker build -t finance-tracker-api .

# Jalankan container
docker run -p 8080:8080 --env-file .env finance-tracker-api
```

---

## 📡 API Endpoints

### Auth
| Method | Endpoint | Deskripsi |
|---|---|---|
| `POST` | `/api/v1/auth/register` | Daftar akun baru + kirim OTP |
| `POST` | `/api/v1/auth/verify-otp` | Verifikasi OTP → return token |
| `POST` | `/api/v1/auth/resend-otp` | Kirim ulang OTP |
| `POST` | `/api/v1/auth/login` | Login dengan email & password |

### Transactions *(Protected)*
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/api/v1/transactions` | List transaksi (support filter & search) |
| `POST` | `/api/v1/transactions` | Tambah transaksi baru |
| `PUT` | `/api/v1/transactions/:id` | Update transaksi |
| `DELETE` | `/api/v1/transactions/:id` | Hapus transaksi |
| `GET` | `/api/v1/transactions/summary` | Ringkasan pemasukan, pengeluaran, saldo |

### Budgets *(Protected)*
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/api/v1/budgets` | List budget bulan ini |
| `POST` | `/api/v1/budgets` | Buat/update budget per kategori |
| `DELETE` | `/api/v1/budgets/:id` | Hapus budget |

---

## 🌐 Deployment

### Backend — Railway
1. Push repository ke GitHub
2. Buka [railway.app](https://railway.app) → **New Project** → **Deploy from GitHub**
3. Tambahkan **PostgreSQL** plugin
4. Set semua environment variables dari `.env`

### Frontend — Vercel
1. Buka [vercel.com](https://vercel.com) → **Import Project** dari GitHub
2. Set environment variable:
   ```
   VITE_API_URL=https://your-api.up.railway.app/api/v1
   ```
3. Deploy

---

## 📄 License

[MIT](LICENSE) © 2026 Faris
