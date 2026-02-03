# 🛒 Order Management API (POS System)

Order Management API adalah **Point of Sale (POS) System** yang dibangun menggunakan **Golang**, **Echo Framework**, dan **PostgreSQL**. API ini menyediakan fitur lengkap untuk mengelola penjualan, inventory, pembayaran, dan laporan bisnis.

---

## 🚀 Fitur Utama

### 👥 Autentikasi & Manajemen User
- ✅ Registrasi & Login User/Admin (JWT)
- ✅ Role-based Access Control (Admin, User/Cashier)
- ✅ Password reset & account verification
- ✅ Enkripsi password dengan bcrypt

### 📦 Manajemen Produk
- ✅ CRUD Produk (admin-only)
- ✅ Stock tracking per produk
- ✅ Upload foto produk
- ✅ Redis caching untuk performa

### 🛒 Shopping Cart & Checkout
- ✅ Tambah/edit/hapus item dari cart
- ✅ Real-time cart total calculation
- ✅ Checkout dengan konversi otomatis ke order

### 💳 Pembayaran
- ✅ **Metode Pembayaran Multiple:**
  - Cash (dengan automatic change calculation)
  - Midtrans (online payment gateway)
- ✅ Auto webhook untuk payment confirmation
- ✅ Order status auto-update saat payment berhasil

### 🧾 Receipt & Invoice
- ✅ Auto-generate receipt number (RCP20260203XXXX)
- ✅ Tax calculation (10%)
- ✅ Detail receipt items dengan harga
- ✅ Cashier & store information
- ✅ Print-ready format

### 📊 Sales Reporting (Admin)
- ✅ Sales report by date range
- ✅ Daily sales summary
- ✅ Monthly sales summary
- ✅ Payment method breakdown (cash vs Midtrans)
- ✅ Top 10 products by sales volume
- ✅ Metrics: Total sales, transactions, tax, avg transaction value, customer count

### 📧 Email & Background Jobs
- ✅ Email otomatis (welcome, verification, notifications)
- ✅ Async processing dengan Goroutine & Queue
- ✅ Photo upload processing

---

## 🧱 Teknologi & Library

- **Go 1.20+**
- [Echo v4](https://echo.labstack.com/) - Web Framework
- [GORM](https://gorm.io/) - ORM Database
- **PostgreSQL** - Database
- **Redis** - Caching & Queue
- **Midtrans** - Payment Gateway
- JWT (`github.com/golang-jwt/jwt/v5`) - Authentication
- Gomail (SMTP) - Email sending
- UUID (`github.com/google/uuid`) - ID generation
- Bcrypt - Password hashing
- Migrate - Database migration

---

## 📁 Struktur Proyek
```bash
OrderManagement-API/
├── cmd/app/ 
│   └── main.go              # Entry point aplikasi
├── internal/
│   ├── entity/              # Domain models
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── order.go
│   │   ├── cart.go
│   │   ├── receipt.go
│   │   └── sales_report.go
│   ├── repository/          # Data access layer
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── order.go
│   │   ├── cart.go
│   │   ├── receipt.go
│   │   └── sales_report.go
│   ├── service/             # Business logic
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── order.go
│   │   ├── cart.go
│   │   ├── receipt.go
│   │   └── sales_report.go
│   ├── http/
│   │   ├── handler/         # HTTP handlers
│   │   ├── binder/          # Request DTOs
│   │   ├── router/          # Route definitions
│   │   └── middleware/      # JWT & Auth middleware
│   ├── builder/             # Dependency injection
│   └── mocks/               # Test mocks
├── pkg/
│   ├── token/               # JWT helper
│   ├── email/               # Email sender (SMTP)
│   ├── encrypt/             # AES encryption
│   ├── cache/               # Redis client
│   ├── server/              # Echo config
│   ├── midtrans/            # Payment gateway
│   ├── postgres/            # DB connection
│   ├── response/            # JSON response formatter
│   └── worker/              # Goroutine workers
├── db/
│   ├── migrations/          # SQL migrations (000001-000007)
│   └── seed/                # Database seeders
├── .env                     # Environment variables
├── docker-compose.yml       # PostgreSQL & Redis
├── go.mod
└── README.md
```

---

## ⚙️ Setup & Jalankan

### 1. Clone Repo
```bash
git clone https://github.com/Kevinmajesta/OrderManagement-API.git
cd OrderManagement-API
```

### 2. Siapkan .env
```bash
# App Config
ENV="dev"
PORT="8080"

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=cuaniaga
POSTGRES_PASSWORD=cuaniaga
POSTGRES_DATABASE=cuaniaga

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET_KEY=your_secret_key_here

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Midtrans (Payment Gateway)
MIDTRANS_SERVER_KEY=your_server_key
MIDTRANS_CLIENT_KEY=your_client_key
MIDTRANS_IS_PRODUCTION=false
```

### 3. Setup Database & Cache
```bash
# Gunakan Docker
docker-compose up -d

# Atau jalankan PostgreSQL & Redis secara manual
```

### 4. Run Database Migrations
```bash
migrate -path db/migrations -database "postgres://cuaniaga:cuaniaga@localhost:5432/cuaniaga?sslmode=disable" up
```

### 5. Jalankan Aplikasi
```bash
go mod tidy
go run cmd/app/main.go
```

Server akan berjalan di `http://localhost:8080`

---

## 📚 API Endpoints

### Authentication
```
POST   /auth/register           # Register user baru
POST   /auth/login              # Login & dapatkan token
POST   /login/admin             # Login admin
POST   /admins                  # Create admin
```

### Products
```
GET    /products                # Get semua produk
GET    /products/{id}           # Get produk by ID
POST   /products                # Create produk (admin)
PUT    /products/{id}           # Update produk (admin)
DELETE /products/{id}           # Delete produk (admin)
```

### Shopping Cart
```
GET    /cart                    # Get cart user
POST   /cart/items              # Add item ke cart
PUT    /cart/items/{id}         # Update cart item
DELETE /cart/items/{id}         # Remove item dari cart
POST   /cart/checkout           # Checkout & buat order
```

### Orders
```
GET    /orders                  # Get order history user
GET    /orders/{id}             # Get detail order
```

### Receipts
```
POST   /receipts                # Generate receipt
GET    /receipts/{id}           # Get receipt by ID
GET    /receipts                # Get user's receipts
```

### Sales Reports (Admin Only)
```
POST   /reports/sales/date-range   # Report by date range
GET    /reports/sales/daily         # Daily sales report
GET    /reports/sales/monthly       # Monthly sales report
```

---

## 🔐 Database Schema

### Tables (7 migrations)
- **users** - User data & authentication
- **products** - Product inventory
- **orders** - Order transactions
- **order_items** - Detail pesanan
- **carts** - Shopping cart
- **cart_items** - Cart items
- **receipts** - Invoice/receipt
- **receipt_items** - Receipt details

---

## 📋 Workflow Cashier

1. **Login** → Dapatkan JWT token
2. **Browse Products** → Lihat katalog produk
3. **Add to Cart** → Pilih & tambah items
4. **Checkout** → Pilih metode pembayaran (cash/midtrans)
5. **Payment** → Proses pembayaran
6. **Generate Receipt** → Print/save invoice
7. **View Reports** → (Admin) Lihat sales analytics

---

## 🧪 Testing dengan Postman

Import collection Postman untuk test semua endpoint:
- [Postman Collection](https://www.postman.com/lunar-resonance-148572/workspace/kevin-work/collection/33423852-49715f15-5735-4460-9cc0-ada1fa7bb18b?action=share&creator=33423852)

Atau manual test:
```bash
# Login
curl -X POST http://localhost:8080/app/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Get Products
curl -X GET http://localhost:8080/app/api/v1/products \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 💡 Notes

- Server auto-seeds admin & demo user pada startup
- Email & photo upload berjalan asynchronously
- Payment webhook otomatis update order status
- Semua endpoint protected JWT kecuali login & register

---

## 👨‍💼 Author

**Kevin Majesta**  
E-commerce & POS System Developer



