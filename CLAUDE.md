# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Daneizo is a state engine for P2P lending platforms written in Go 1.23. It manages states and transitions in peer-to-peer lending workflows using clean architecture principles with domain-driven design.

## Essential Commands

### Development Workflow
```bash
# Build and run
go build -o daneizo .
go run main.go

# Testing
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -race ./...              # Race detection
go test ./internal/pkg/pkguid    # Specific package

# Code quality
go fmt ./...                     # Format code
go vet ./...                     # Static analysis
go mod tidy                      # Clean dependencies

# Mock generation (uses .mockery.yaml)
mockery --all
```

### Database Operations
```bash
# Migrations
goose -dir migration/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" up
goose -dir migration/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" down
goose -dir migration/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" status

# Seeding
goose -no-versioning -dir seed/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" up
```

## Architecture

### High-Level Structure
```
main.go → internal/app/app.go (App struct) → domain modules
```

The `App` struct in `internal/app/app.go` is the central dependency container containing:
- Database connection (*sql.DB)
- GoQu query builder wrapper
- Validator, logger, HTTP router/server
- Viper configuration, Snowflake ID generator

### Domain Organization
- **Clean Architecture**: Domains in `internal/[domain]/internal/`
- **Loan Domain**: Main business domain at `internal/loan/`
  - `entity/`: Domain entities and SQL entities
  - `gateway/`: Data access layer
  - `interactor/`: Business logic
  - `usecase/`: Use case implementations
  - `mocks/`: Generated mocks

### Shared Packages
- `internal/pkg/pkgerror/`: Custom error handling with codes
- `internal/pkg/pkghttp/v1/`: HTTP utilities and endpoint patterns
- `internal/pkg/pkgsql/`: GoQu query builder wrappers
- `internal/pkg/pkguid/`: Snowflake ID generation

## Key Technologies

- **HTTP**: httprouter for routing
- **Database**: MySQL with GoQu query builder
- **Migrations**: Goose for schema management
- **Testing**: testify with suite pattern, go-sqlmock for DB testing
- **Mocking**: mockery with expecter pattern
- **Logging**: Zap structured logging
- **Config**: Viper for environment-based configuration

## Testing Patterns

- Use testify suites for complex test setups (e.g., `loanSQLGatewaySuite`)
- Generate mocks with mockery using `.mockery.yaml` configuration
- SQL testing uses go-sqlmock for database mocking
- Test naming: `Test[Type]_[Method]` pattern

## Development Guidelines

### Code Style
- Standard `gofmt` formatting
- Packages: `internal/pkg/pkg[name]/` naming pattern
- SQL entities: Suffix with `SQLEntity`
- Mocks: Prefix with `Mock`
- Follow Go naming conventions consistently

### Pre-commit Checklist
1. Run `go fmt ./...` and `go vet ./...`
2. Run `go test ./...` and `go test -race ./...`
3. Run `mockery --all` if interfaces changed
4. Run `go mod tidy`
5. Test database migrations if schema changed
6. Verify `go build .` succeeds

### Environment Setup
```bash
cp .env.example .env
# Configure database credentials in .env
```

The application expects MySQL database configuration and runs HTTP server on port 8081 by default.