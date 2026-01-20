# Dynamic Accounting SaaS Backend

A production-grade multi-tenant accounting SaaS backend for health professionals with dynamic form engine and calculation capabilities.

## Features

- ✅ Multi-tenant architecture with clinic isolation
- ✅ Dynamic form builder (no hard-coded forms)
- ✅ Formula-based calculation engine
- ✅ JWT authentication
- ✅ PostgreSQL database with versioned forms
- ✅ RESTful API

## Prerequisites

- Go 1.25.3 or higher
- PostgreSQL 12 or higher
- Make (optional, for convenience commands)

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Database Setup

#### Option A: Using Docker Compose (Recommended for Development)

Start PostgreSQL using Docker Compose:

```bash
docker-compose up -d
```

This will:
- Start a PostgreSQL 15 container
- Create the `acareca` database automatically
- Run migrations from `migration/` on first startup
- Expose PostgreSQL on port 5432

To stop the database:

```bash
docker-compose down
```

To stop and remove volumes (⚠️ deletes all data):

```bash
docker-compose down -v
```

#### Option B: Manual PostgreSQL Setup

Create a PostgreSQL database:

```sql
CREATE DATABASE acareca;
```

Run the migrations from the `migration/` directory:

```bash
psql -U your_user -d acareca -f migration/20260120061252_user.sql
psql -U your_user -d acareca -f migration/20260120061308_auth_provider.sql
psql -U your_user -d acareca -f migration/20260120061318_session.sql
```

Or using the psql command line:

```bash
psql -h localhost -U postgres -d acareca < migration/20260120061252_user.sql
psql -h localhost -U postgres -d acareca < migration/20260120061308_auth_provider.sql
psql -h localhost -U postgres -d acareca < migration/20260120061318_session.sql
```

### 3. Environment Configuration

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

Edit `.env`:

   ```env
   PORT=8080
   # For Docker Compose setup:
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/acareca?sslmode=disable
   # For manual PostgreSQL setup, use your own credentials:
   # DATABASE_URL=postgres://user:password@localhost/acareca?sslmode=disable
   JWT_SECRET=your-secret-key-change-in-production
   JWT_EXPIRY=24
   ENV=development
   ```

**Important**: Change `JWT_SECRET` to a secure random string in production. You can generate one using:

```bash
openssl rand -base64 32
```

### 4. Run the Server

```bash
go run main.go
```

Or build and run:

```bash
go build -o bin/server main.go
./bin/server
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check

```bash
GET /health
```

### Form Management

**Create Form**
```bash
POST /api/clinics/:clinicId/forms
Authorization: Bearer <token>
Content-Type: application/json

{
  "form_key": "bas_quarterly",
  "name": "BAS Quarterly",
  "type": "government",
  "sections": [
    {
      "section_key": "gst",
      "title": "Goods and Services Tax",
      "order_index": 1,
      "fields": [
        {
          "field_key": "G1",
          "label": "Total Sales",
          "field_type": "currency",
          "required": true,
          "order_index": 1,
          "metadata": {
            "includes_gst": true
          }
        },
        {
          "field_key": "1A",
          "label": "GST on Sales",
          "field_type": "calculated",
          "readonly": true,
          "order_index": 2,
          "formula": {
            "expression": "G1 * 0.10",
            "dependencies": ["G1"]
          }
        }
      ]
    }
  ]
}
```

**Get Form**
```bash
GET /api/forms/:formId
Authorization: Bearer <token>
```

**Calculate Form**
```bash
POST /api/forms/:formId/calculate
Authorization: Bearer <token>
Content-Type: application/json

{
  "values": {
    "G1": 120000.00
  }
}
```

**Create Entry**
```bash
POST /api/forms/:formId/entries
Authorization: Bearer <token>
Content-Type: application/json

{
  "form_id": "uuid-here",
  "period": "2025-Q1",
  "values": {
    "G1": "120000"
  }
}
```

## Project Structure

```
backend/
├── main.go                  # Application entry point
├── cmd/
│   └── server.go           # Server initialization
├── config/                  # Configuration management
│   ├── config.go           # Configuration loading
│   ├── database.go         # Database connection & migrations
│   └── redis.go            # Redis client
├── internal/
│   ├── domain/             # Domain handlers (auth, user)
│   ├── middleware/         # HTTP middleware (auth)
│   ├── model/              # Data models
│   └── service/            # Business logic
├── migration/              # Database migrations
├── util/                   # Utility functions
└── go.mod                  # Go dependencies
```

## Key Concepts

### Dynamic Forms

Forms are **not hard-coded**. They are created dynamically via the API:

1. **Form**: Top-level container (e.g., "BAS Quarterly")
2. **Section**: Grouping of related fields (e.g., "GST Section")
3. **Field**: Individual input/display field (e.g., "Total Sales")
4. **Formula**: Calculation expression for calculated fields

### Formula Engine

The calculation engine uses the `govaluate` library for safe expression evaluation:

- Supports: `+`, `-`, `*`, `/`, parentheses
- Field references: Use field keys (e.g., `G1`, `1A`)
- Example: `G1 * 0.10` calculates 10% GST

### Multi-Tenancy

- Each clinic has isolated data
- Users can belong to multiple clinics with different roles
- Tenant isolation enforced via middleware

## Development

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

## Next Steps

Phase 1 includes:
- ✅ Core infrastructure
- ✅ Dynamic form engine
- ✅ Calculation engine
- ✅ Basic API endpoints

Future phases will add:
- Authentication endpoints (register/login)
- Clinic management
- Entry submission workflow
- Audit logging
- Reporting APIs

## License

MIT

