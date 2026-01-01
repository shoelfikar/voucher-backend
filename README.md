# Voucher Management System - Backend API

Backend API for Voucher Management System built with Go, Gin, and PostgreSQL.

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt v5)
- **Validation**: go-playground/validator v10
- **CORS**: gin-contrib/cors
- **Configuration**: Viper

## Project Structure

```
backend/
├── cmd/api/              # Application entry point
│   └── main.go
├── internal/
│   ├── config/           # Configuration loader
│   ├── delivery/http/    # HTTP handlers, middleware, router
│   ├── domain/           # Domain entities, interfaces
│   ├── repository/       # Repository implementations
│   └── service/          # Business logic
├── pkg/                  # Reusable packages
│   ├── database/         # Database connection
│   ├── jwt/              # JWT utilities
│   └── utils/            # Common utilities
├── migrations/           # Database migration files
├── .env.example          # Example environment variables
├── Makefile             # Build automation
└── go.mod               # Go dependencies
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Installation

1. **Clone the repository**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   make install
   # or
   go mod download
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```

   Edit `.env` with your configuration:
   ```env
   PORT=8080
   GIN_MODE=debug

   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=voucher_db
   DB_SSLMODE=disable

   JWT_SECRET=your-super-secret-key-change-this
   JWT_EXPIRATION=24h

   ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
   ```

4. **Create PostgreSQL database**
   ```bash
   createdb voucher_db
   # or using psql
   psql -U postgres -c "CREATE DATABASE voucher_db;"
   ```

5. **Run database migrations**
   ```bash
   make migrate-up
   ```

   The application also supports auto-migration on startup using GORM.

## Running the Application

### Using Make (Recommended)

```bash
# Run the application
make run

# Build the binary
make build

# Run the binary
./bin/voucher-api
```

### Using Go directly

```bash
# Run
go run cmd/api/main.go

# Build
go build -o bin/voucher-api cmd/api/main.go
```

The server will start on `http://localhost:8080` (or the port specified in `.env`).

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Authentication (Public)
- `POST /api/v1/login` - User login (dummy validation)

### Vouchers (Protected - requires JWT)
- `GET /api/v1/vouchers` - Get all vouchers (with pagination, search, sort)
- `GET /api/v1/vouchers/:id` - Get voucher by ID
- `POST /api/v1/vouchers` - Create new voucher
- `PUT /api/v1/vouchers/:id` - Update voucher
- `DELETE /api/v1/vouchers/:id` - Delete voucher (soft delete)

### CSV Operations (Protected - requires JWT)
- `POST /api/v1/vouchers/upload-csv` - Import vouchers from CSV file
- `GET /api/v1/vouchers/export` - Export vouchers to CSV file

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

To get a token, login with any email and password (minimum 6 characters):

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'
```

## CSV Format

For importing vouchers, use this CSV format:

```csv
voucher_code,discount_percent,expiry_date
DISC10,10.00,2025-12-31
SAVE20,20.00,2025-06-30
PROMO50,50.00,2025-03-15
```

**Validation Rules:**
- `voucher_code`: Required, max 50 characters, must be unique
- `discount_percent`: Required, must be between 1-100
- `expiry_date`: Required, format YYYY-MM-DD, must be today or in the future

## Development

### Available Make Commands

```bash
make help          # Show all available commands
make run           # Run the application
make build         # Build the binary
make test          # Run tests
make test-coverage # Run tests with coverage
make clean         # Clean build artifacts
make install       # Install dependencies
make format        # Format code
make lint          # Run linter
make migrate-up    # Run migrations
make migrate-down  # Rollback migrations
```

### Database Migrations

Migrations are stored in the `migrations/` directory.

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Force migration to specific version
make migrate-force v=1
```

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for specific package
go test -v ./internal/service/...
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| GIN_MODE | Gin mode (debug/release) | debug |
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | PostgreSQL user | postgres |
| DB_PASSWORD | PostgreSQL password | postgres |
| DB_NAME | Database name | voucher_db |
| DB_SSLMODE | SSL mode | disable |
| JWT_SECRET | JWT secret key | (required) |
| JWT_EXPIRATION | JWT expiration time | 24h |
| ALLOWED_ORIGINS | CORS allowed origins | http://localhost:5173 |

## Production Deployment

1. Set `GIN_MODE=release` in production
2. Use a strong `JWT_SECRET`
3. Enable SSL for database (`DB_SSLMODE=require`)
4. Configure proper CORS origins
5. Use environment-specific `.env` files
6. Run behind a reverse proxy (nginx, traefik)

## License

MIT
