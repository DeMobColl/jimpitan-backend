# Jimpitan Backend - Go Implementation

Implementasi backend Golang menggantikan Google Apps Script dengan database MySQL.

## 📋 Struktur Proyek

```
backend-go/
├── cmd/
│   └── server/
│       └── main.go              # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── db.go                # Database connection
│   ├── handlers/
│   │   ├── auth_handler.go      # Auth endpoints
│   │   ├── user_handler.go      # User CRUD endpoints
│   │   ├── customer_handler.go  # Customer management endpoints
│   │   └── transaction_handler.go # Transaction endpoints
│   ├── middleware/
│   │   └── auth.go              # JWT authentication middleware
│   ├── models/
│   │   └── models.go            # Data models/structs
│   ├── services/
│   │   ├── auth_service.go      # Authentication business logic
│   │   ├── user_service.go      # User management logic
│   │   ├── customer_service.go  # Customer management logic
│   │   └── transaction_service.go # Transaction processing logic
│   └── utils/
│       └── crypto.go            # Hashing, token generation, ID generation
├── migrations/
│   ├── 001_initial_schema.sql   # Database schema setup
│   └── 002_add_indexes.sql      # Performance indexes
├── .env.example                 # Environment variables template
├── go.mod                       # Go module definition
├── Makefile                     # Build and deployment commands
└── README.md                    # This file
```

## 🚀 Installation & Setup

### Prerequisites
- Go 1.21 or higher
- MySQL 8.0 or higher
- (Optional) `air` for hot reload during development

### 1. Clone & Setup

```bash
cd backend-go
cp .env.example .env
# Edit .env with your MySQL credentials
```

### 2. Download Dependencies

```bash
make deps
# or manually:
go mod download
```

### 3. Database Setup

Create MySQL database:
```bash
mysql -u root -p
CREATE DATABASE jimpitan CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'jimpitan'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON jimpitan.* TO 'jimpitan'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

Run migrations:
```bash
make migrate
# or manually:
mysql -h localhost -u jimpitan -p jimpitan < migrations/001_initial_schema.sql
mysql -h localhost -u jimpitan -p jimpitan < migrations/002_add_indexes.sql
```

### 4. Start Backend

**Production mode:**
```bash
make build
make run
```

**Development mode (with hot reload):**
```bash
make dev
```

Server will start on `http://localhost:8080`

## 📡 API Endpoints

### Authentication (Public)

```
POST /api/login
GET  /api/verifyToken?token=xxx
POST /api/logout (Protected)
```

### Users (Protected)

```
GET    /api/users                           # List all users (Admin)
POST   /api/users                           # Create user (Admin)
PUT    /api/users?id=USR-001                # Update user (Admin)
DELETE /api/users?id=USR-001                # Delete user (Admin)
GET    /api/users/activity?user_id=USR-001 # Get user transactions
POST   /api/users/password                  # Change own password
POST   /api/users/bulk-delete               # Bulk delete users (Admin)
```

### Customers (Protected)

```
GET    /api/customers                       # List all customers
POST   /api/customers                       # Create customer
PUT    /api/customers?id=CUST-001           # Update customer
DELETE /api/customers?id=CUST-001           # Delete customer
GET    /api/customers/qr?qr_hash=abc123    # Get customer by QR
GET    /api/customers/history?customer_id=CUST-001 # Customer history
POST   /api/customers/bulk-delete           # Bulk delete customers
```

### Transactions (Protected)

```
GET  /api/transactions             # List all transactions
POST /api/transactions             # Submit new transaction
DELETE /api/transactions?id=0001   # Delete transaction
```

## 🔐 Authentication

Menggunakan JWT (JSON Web Tokens) dengan implementasi:

1. **Login** - User kirim username + password → Backend generate JWT token
2. **Token Storage** - Frontend simpan token di localStorage
3. **Protected Routes** - Kirim token via `Authorization: Bearer <token>` header
4. **Token Expiry** - Default 7 hari (configurable via `JWT_EXPIRY_HOURS`)

Token digenerate menggunakan RS256 signing method.

## 🗄️ Database Schema

### Users Table
```
id: VARCHAR(20) - USR-001, USR-002, ...
name: VARCHAR(255)
role: ENUM('admin', 'petugas')
username: VARCHAR(100) UNIQUE
password_hash: VARCHAR(255) - SHA-256
token: VARCHAR(255) UNIQUE
token_expiry: DATETIME
last_login: DATETIME
created_at: DATETIME
updated_at: DATETIME
deleted_at: DATETIME (soft delete)
```

### Customers Table
```
id: VARCHAR(20) - CUST-001, CUST-002, ...
blok: VARCHAR(50)
nama: VARCHAR(255)
qr_hash: VARCHAR(10) UNIQUE
total_setoran: DECIMAL(12,2)
last_transaction: DATETIME
created_at: DATETIME
updated_at: DATETIME
deleted_at: DATETIME (soft delete)
```

### Transactions Table
```
id: VARCHAR(20) - 0001, 0002, ...
timestamp: DATETIME
customer_id: VARCHAR(20) - FK to customers
blok: VARCHAR(50) - Denormalized
nama: VARCHAR(255) - Denormalized
nominal: DECIMAL(12,2)
user_id: VARCHAR(20) - FK to users
petugas: VARCHAR(255) - Denormalized
created_at: DATETIME
deleted_at: DATETIME (soft delete)
```

## 📝 Environment Variables

```
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=jimpitan
DB_PASSWORD=your_password
DB_NAME=jimpitan

# Server
PORT=8080
ENV=development

# JWT
JWT_SECRET=your-very-secure-secret-key
JWT_EXPIRY_HOURS=168

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

## 🔄 Migration from Google Apps Script

### Key Changes

1. **Database**: Google Sheets → MySQL
2. **Authentication**: Session tokens → JWT tokens
3. **API Format**: JSONP GET requests → RESTful JSON endpoints
4. **Password Hashing**: JavaScript SHA-256 → Go crypto/sha256
5. **Token Generation**: 32-char random string → JWT with claims

### Data Migration

Spreadsheet data can be migrated using:
1. Export sheets to CSV
2. Write importer script (see `/scripts/migrate-from-sheets.go`)
3. Run: `go run scripts/migrate-from-sheets.go`

## 🧪 Testing

### Unit Tests
```bash
make test
```

### Manual Testing with cURL

**Login:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

**Get Customers:**
```bash
curl -X GET "http://localhost:8080/api/customers" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Create Transaction:**
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "customer_id":"CUST-001",
    "blok":"A1",
    "nama":"John Doe",
    "nominal":50000,
    "user_id":"USR-001",
    "petugas":"Petugas A"
  }'
```

## 📚 Dependencies

- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/gorilla/mux` - HTTP routing
- `github.com/rs/cors` - CORS handling
- `golang.org/x/crypto` - Password hashing

## 🚨 Common Issues

### "Failed to connect to database"
- Check MySQL is running
- Verify credentials in `.env`
- Check database exists: `mysql -u root -p -e "SHOW DATABASES;"`

### "Port 8080 already in use"
- Change `PORT` in `.env`
- Or kill existing process: `lsof -i :8080`

### "Token not found in request"
- Ensure Authorization header format: `Authorization: Bearer <token>`
- For GET requests, can also use `?token=<token>` query param

## 📖 Additional Resources

- API Documentation: See `../JimpReact/docs/GOLANG_API.md`
- Frontend Integration: See `../JimpReact/docs/FRONTEND_INTEGRATION.md`
- Database Schema Diagram: See `../JimpReact/docs/DATABASE_SCHEMA.md`
