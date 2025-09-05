# Essential Commands

## Development Commands

### Build and Run
```bash
go build -o daneizo .                    # Build binary
./daneizo                                # Run the application
go run main.go                           # Build and run directly
```

### Testing
```bash
go test ./...                            # Run all tests
go test -v ./...                         # Run tests with verbose output
go test -race ./...                      # Run tests with race detection
go test ./internal/pkg/pkguid            # Run tests for specific package
go test -run TestSpecificFunction        # Run specific test
```

### Code Quality
```bash
go fmt ./...                             # Format all Go files
go vet ./...                             # Static analysis for bugs
go mod tidy                              # Clean up dependencies
go mod verify                            # Verify dependencies
```

### Mock Generation
```bash
mockery --all                            # Generate all mocks (uses .mockery.yaml config)
```

## Database Commands

### Migrations
```bash
# Run migrations
goose -dir migration/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" up

# Rollback migration
goose -dir migration/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" down

# Check migration status
goose -dir migration/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" status
```

### Seeding
```bash
# Seed database with initial data
goose -no-versioning -dir seed/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" up
```

## Environment Setup
```bash
cp .env.example .env                     # Copy environment template
# Edit .env file with your database credentials and configuration
```