# 💰 Crypto Wallet Service

Backend service untuk mengelola wallet cryptocurrency dengan dukungan deposit, withdraw, dan monitoring portfolio real-time menggunakan CoinGecko API.

## 🎯 Fitur

- ✅ Autentikasi user dengan JWT
- 💼 Multi-currency wallet (BTC, ETH, USDT, IDR)
- 💸 Deposit dan withdrawal
- 📊 Portfolio real-time dengan harga dari CoinGecko
- 📜 Transaction history dengan pagination
- ⚡ Redis caching untuk performa optimal
- 🐳 Docker support untuk deployment mudah

## 🧱 Tech Stack

| Komponen      | Teknologi                  |
| ------------- | -------------------------- |
| Bahasa        | Go 1.21                    |
| Database      | PostgreSQL 15              |
| ORM           | GORM                       |
| Cache         | Redis 7                    |
| Framework     | Gin Web Framework          |
| API Eksternal | CoinGecko API              |
| Auth          | JWT (golang-jwt/jwt)       |
| Container     | Docker & Docker Compose    |

## 📁 Struktur Folder

```
crypto-wallet-service/
├── cmd/
│   └── main.go                    # Entry point aplikasi
├── config/
│   └── config.go                  # Konfigurasi database, Redis, JWT
├── internal/
│   ├── models/                    # Data models
│   │   ├── user.go
│   │   ├── wallet.go
│   │   └── transaction.go
│   ├── handlers/                  # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── wallet_handler.go
│   │   └── transaction_handler.go
│   ├── services/                  # Business logic
│   │   ├── coingecko_service.go
│   │   └── wallet_service.go
│   ├── repository/                # Database operations
│   │   ├── user_repo.go
│   │   ├── wallet_repo.go
│   │   └── transaction_repo.go
│   ├── middleware/                # Middleware
│   │   └── jwt_middleware.go
│   └── routes/                    # Route configuration
│       └── routes.go
├── .env.example                   # Environment variables template
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## 🚀 Quick Start

### Prasyarat

- Go 1.21 atau lebih baru
- PostgreSQL 15
- Redis 7
- Docker & Docker Compose (opsional)

### Setup Manual

1. **Clone repository**
```bash
cd d:\projek\porto_golang
```

2. **Copy environment file**
```bash
copy .env.example .env
```

3. **Edit `.env` sesuai kebutuhan**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=crypto_wallet
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-super-secret-jwt-key
```

4. **Install dependencies**
```bash
go mod download
```

5. **Jalankan PostgreSQL dan Redis**
```bash
# Pastikan PostgreSQL dan Redis sudah running
```

6. **Run aplikasi**
```bash
go run cmd/main.go
```

Server akan berjalan di `http://localhost:8080`

### Setup dengan Docker

1. **Build dan jalankan semua services**
```bash
docker-compose up -d
```

2. **Cek logs**
```bash
docker-compose logs -f app
```

3. **Stop services**
```bash
docker-compose down
```

## 📡 API Endpoints

### Authentication

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-11-07T10:00:00Z"
  }
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

### User Profile

#### Get Current User
```http
GET /api/user/me
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2025-11-07T10:00:00Z"
}
```

### Wallet

#### Get Portfolio
```http
GET /api/wallet
Authorization: Bearer <token>
```

**Response:**
```json
{
  "assets": [
    {
      "currency": "BTC",
      "balance": 0.002,
      "price_idr": 950000000,
      "value_idr": 1900000
    },
    {
      "currency": "USDT",
      "balance": 50,
      "price_idr": 15500,
      "value_idr": 775000
    }
  ],
  "total_value_idr": 2675000
}
```

#### Deposit
```http
POST /api/wallet/deposit
Authorization: Bearer <token>
Content-Type: application/json

{
  "currency": "BTC",
  "amount": 0.001
}
```

**Response:**
```json
{
  "message": "Deposit successful",
  "currency": "BTC",
  "amount": 0.001
}
```

#### Withdraw
```http
POST /api/wallet/withdraw
Authorization: Bearer <token>
Content-Type: application/json

{
  "currency": "BTC",
  "amount": 0.0005
}
```

### Transactions

#### Get Transaction History
```http
GET /api/transactions?page=1&limit=20
Authorization: Bearer <token>
```

**Response:**
```json
{
  "transactions": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "type": "deposit",
      "currency": "BTC",
      "amount": 0.001,
      "price_at": 950000000,
      "created_at": "2025-11-07T10:00:00Z"
    }
  ],
  "pagination": {
    "total": 50,
    "page": 1,
    "limit": 20,
    "total_pages": 3
  }
}
```

## 💾 Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Wallets Table
```sql
CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    currency VARCHAR(10) NOT NULL,
    balance NUMERIC(18,8) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    amount NUMERIC(18,8) NOT NULL,
    price_at NUMERIC(18,2),
    created_at TIMESTAMP DEFAULT NOW()
);
```

## 🔧 Konfigurasi

### Environment Variables

| Variable                 | Description                    | Default                             |
| ------------------------ | ------------------------------ | ----------------------------------- |
| `SERVER_PORT`            | Port server                    | 8080                                |
| `GIN_MODE`               | Mode Gin (debug/release)       | debug                               |
| `DB_HOST`                | PostgreSQL host                | localhost                           |
| `DB_PORT`                | PostgreSQL port                | 5432                                |
| `DB_USER`                | PostgreSQL user                | postgres                            |
| `DB_PASSWORD`            | PostgreSQL password            | postgres                            |
| `DB_NAME`                | Database name                  | crypto_wallet                       |
| `REDIS_HOST`             | Redis host                     | localhost                           |
| `REDIS_PORT`             | Redis port                     | 6379                                |
| `JWT_SECRET`             | JWT signing secret             | your-super-secret-jwt-key           |
| `COINGECKO_API_URL`      | CoinGecko API base URL         | https://api.coingecko.com/api/v3    |
| `CACHE_DURATION_SECONDS` | Cache duration untuk price     | 60                                  |

## 🧪 Testing API

### Menggunakan curl

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'

# Get Portfolio (gunakan token dari login)
curl http://localhost:8080/api/wallet \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Deposit
curl -X POST http://localhost:8080/api/wallet/deposit \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"currency":"IDR","amount":1000000}'
```

## 📊 Supported Currencies

| Currency | Name            | Type   |
| -------- | --------------- | ------ |
| BTC      | Bitcoin         | Crypto |
| ETH      | Ethereum        | Crypto |
| USDT     | Tether          | Crypto |
| IDR      | Indonesian Rupiah | Fiat  |

## 🔐 Security Features

- ✅ Password hashing dengan bcrypt
- ✅ JWT token authentication
- ✅ Protected routes dengan middleware
- ✅ Input validation
- ✅ SQL injection prevention (GORM ORM)

## 🎨 Best Practices

1. **Error Handling**: Semua error di-handle dengan proper HTTP status codes
2. **Logging**: Database dan Redis connection logging
3. **Caching**: Redis cache untuk harga crypto (refresh tiap 60 detik)
4. **Transaction**: Database transaction untuk deposit/withdraw
5. **Validation**: Input validation di handler layer
6. **Architecture**: Clean architecture dengan separation of concerns

## 📝 Development

### Menambah Currency Baru

Edit `internal/services/coingecko_service.go`:
```go
coinMap := map[string]string{
    "BTC":  "bitcoin",
    "ETH":  "ethereum",
    "USDT": "tether",
    "SOL":  "solana",  // Tambahkan di sini
}
```

### Custom Cache Duration

Edit `.env`:
```env
CACHE_DURATION_SECONDS=120  # 2 menit
```

## 🐛 Troubleshooting

### Database Connection Error
```
Failed to connect to database: connection refused
```
**Solusi**: Pastikan PostgreSQL running dan kredensial di `.env` benar

### Redis Connection Error
```
Failed to connect to Redis: connection refused
```
**Solusi**: Pastikan Redis running di port yang sesuai

### CoinGecko API Error
```
Failed to fetch prices from CoinGecko
```
**Solusi**: Cek koneksi internet atau gunakan cache yang ada





---

**Happy Coding! 🚀**
