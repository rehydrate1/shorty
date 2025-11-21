
# Shorty

[![Go Report Card](https://goreportcard.com/badge/github.com/rehydrate1/shorty?v=1)](https://goreportcard.com/report/github.com/rehydrate1/shorty)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Shorty is a high-performance URL shortening service written in Go.

It is designed to demonstrate **Clean Architecture** principles, **Dependency Injection**, and robust **Unit Testing** with mocks. The project uses modern Go idioms and tooling suitable for production-grade applications.

## Tech Stack

-   **Language:** Go 1.22+
-   **Web Framework:** [Gin](https://github.com/gin-gonic/gin) (Fast HTTP web framework)
-   **Storage:** PostgreSQL
-   **Database Driver:** [pgx](https://github.com/jackc/pgx) (High-performance driver)
-   **Migrations:** [Goose](https://github.com/pressly/goose)
-   **Containerization:** Docker & Docker Compose
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
├── migrations/      # SQL migrations (Goose)
├── .env.example     # Template for environment variables
├── docker-compose.yml
└── Dockerfile
```

## Quick Start (Docker)

The easiest way to run the project is using Docker Compose. It will set up the app, database, and apply migrations automatically.

### 1. Configuration
Copy the example configuration file:
```sh
cp .env.example .env
```

### 2. Run
```sh
docker-compose up --build
```

The server will start at `http://localhost:8080`.

---

## Local Development (Manual)

If you prefer to run Go locally without Docker:

### 1. Prerequisites
-   Go 1.22+
-   PostgreSQL
-   [Goose](https://github.com/pressly/goose) (for migrations)

### 2. Database Setup
Create a database (e.g., `shorty_db`) and apply migrations:

```sh
# Install goose if needed: go install github.com/pressly/goose/v3/cmd/goose@latest
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://user:password@localhost:5432/shorty_db?sslmode=disable"
goose -dir ./migrations up
```

### 3. Run App
Ensure your `.env` file points to your local Postgres instance, then:

```sh
go mod tidy
go run cmd/shorty/main.go
```

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
