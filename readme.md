# 🛒 Order Management API

Order Management API adalah RESTful API sederhana yang dibangun menggunakan **Golang**, **Echo Framework**, dan **PostgreSQL**. API ini menyediakan fitur untuk mengelola produk, pemesanan oleh customer, autentikasi user/admin, serta pengiriman email dan upload gambar produk menggunakan Goroutine & Queue.

---

## 🚀 Fitur Utama

- ✅ **Autentikasi User & Admin (JWT)**
- 🛍️ **CRUD Produk (admin-only)**
- 📦 **Pemesanan Produk oleh Customer**
- 🔄 **Update Status Pesanan**
- 📜 **Riwayat Pesanan Customer**
- 📧 **Email Otomatis (Welcome & Verifikasi)**
- 🖼️ **Upload Foto Produk via Queue**
- 📬 **Worker Goroutine Email & Foto**
- 🧪 Middleware, Validator, Error Handler

---

## 🧱 Teknologi & Library

- **Go 1.20+**
- [Echo](https://echo.labstack.com/)
- [GORM](https://gorm.io/)
- PostgreSQL
- Redis (untuk cache)
- JWT (`github.com/golang-jwt/jwt/v5`)
- Gomail (SMTP)
- UUID (`github.com/google/uuid`)
- Bcrypt for password hashing

---

## 📁 Struktur Proyek
OrderManagement-API/
├── cmd/app/ # Entry point
├── internal/
│ ├── handler/ # HTTP handler (controller)
│ ├── service/ # Business logic
│ ├── repository/ # Database access
│ ├── entity/ # Entity/model definitions
│ ├── binder/ # Request binding
│ ├── response/ # JSON response
│ ├── middleware/ # JWT middleware
│ └── seeder/ # Seeder user/admin/product
├── pkg/
│ ├── token/ # JWT Token helper
│ ├── email/ # Email sender (SMTP)
│ ├── encrypt/ # AES encryption
│ ├── cache/ # Redis integration
│ ├── server/ # Echo server config
│ └── worker/ # Goroutine email/photo queue
├── assets/photos/ # Uploaded product images
├── .env
├── go.mod
└── README.md

## ⚙️ Setup & Jalankan

### 1. Clone Repo
git clone https://github.com/namakamu/OrderManagement-API.git
cd OrderManagement-API
### 2. Siapkan .env
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=yourpassword
POSTGRES_DB=ordermanagement

# SMTP (email)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# JWT
JWT_SECRET_KEY=your_jwt_secret_key

# Redis
REDIS_ADDR=localhost:6379
### 3. Jalankan Database & Redis
Gunakan Docker PostgreSQL dan Redis. Dengan cara docker-compose up
### 4. Jalankan Aplikasi
```bash
go run cmd/app/main.go 
```
### 5. Seeder
Seeder akan otomatis dijalankan saat server aktif

## 💡 Notes

= Gunakan tool seperti Postman atau Insomnia untuk test API.
= Email & upload berjalan secara asynchronous menggunakan goroutine
- Link Postman : [Link Postman](https://www.postman.com/lunar-resonance-148572/workspace/kevin-work/collection/33423852-49715f15-5735-4460-9cc0-ada1fa7bb18b?action=share&source=copy-link&creator=33423852)



