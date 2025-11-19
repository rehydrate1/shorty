# Shorty

[![Go Report Card](https://goreportcard.com/badge/github.com/rehydrate1/shorty)](https://goreportcard.com/report/github.com/rehydrate1/shorty)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Shorty is a high-performance URL shortening service written in Go.

It is designed to demonstrate **Clean Architecture** principles, **Dependency Injection**, and robust **Unit Testing** with mocks. The project uses modern Go idioms and tooling suitable for production-grade applications.

## Tech Stack

-   **Language:** Go 1.22+
-   **Web Framework:** [Gin](https://github.com/gin-gonic/gin) (Fast HTTP web framework)
-   **Storage:** PostgreSQL
-   **Database Driver:** [pgx](https://github.com/jackc/pgx) (High-performance driver)
-   **Configuration:** [cleanenv](https://github.com/ilyakaznacheev/cleanenv) (12-Factor App compliant)
-   **Logging:** `log/slog` (Structured JSON logging)
-   **Testing:** `testing` (Standard lib) + [testify](https://github.com/stretchr/testify) (Assertions & Mocks)

## Architecture

The project follows the **Standard Go Project Layout**:

```text
.
├── cmd/shorty/      # Entry point (main.go), dependency injection & wiring
├── internal/
│   ├── config/      # Configuration loading logic
│   ├── http-server/ # HTTP Handlers (Controllers)
│   ├── lib/         # Shared libraries (e.g., random generator)
│   └── storage/     # Storage layer (Repository pattern implementation)
├── .env.example     # Template for environment variables
└── .env             # Local environment variables (not committed)
```

## Getting Started

### Prerequisites

-   **Go** 1.22 or higher
-   **PostgreSQL** instance

### 1. Database Setup

Create a database (e.g., `shorty_db`) and execute the following SQL to create the table:

```sql
CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    short_key VARCHAR(10) NOT NULL UNIQUE,
    original_url TEXT NOT NULL
);
CREATE INDEX idx_short_key ON links(short_key);
```

### 2. Configuration

Copy the example configuration file:

```sh
cp .env.example .env
```

Open `.env` and populate `DATABASE_DSN` with your PostgreSQL credentials.

### 3. Run

```sh
go mod tidy
go run cmd/shorty/main.go
```

You should see a log message:
`{"time":"...","level":"INFO","msg":"Starting server","url":"http://localhost:8080"}`

## Running Tests

The project includes unit tests for handlers using **mocks** to isolate the logic from the database.

```sh
go test -v ./internal/...
```

## API Usage

### 1. Shorten a URL

**Request:**
```bash
curl -X POST http://localhost:8080/shorten \
     -H "Content-Type: application/json" \
     -d '{"url": "https://www.ozon.ru"}'
```

**Response:**
```json
{
  "short_url": "http://localhost:8080/Ab3dE1"
}
```

### 2. Redirect

Open the short URL in your browser or use curl:

```bash
curl -v http://localhost:8080/Ab3dE1
```

**Response:**
`HTTP/1.1 302 Found` with `Location: https://www.ozon.ru`

## License

This project is licensed under the MIT License.
