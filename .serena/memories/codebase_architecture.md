# Codebase Architecture

## Project Structure
```
daneizo/
├── main.go                    # Entry point, calls app.Run()
├── internal/                  # Private application code
│   ├── app/                   # Application bootstrap and configuration
│   │   ├── app.go            # Main App struct with dependencies
│   │   ├── run.go            # Application startup logic
│   │   ├── config.go         # Configuration management (Viper)
│   │   ├── http_server.go    # HTTP server setup
│   │   ├── http_router.go    # Route definitions
│   │   ├── mysql_db.go       # Database connection
│   │   └── ...               # Other app components
│   ├── loan/                 # Loan domain module
│   │   ├── module.go         # Domain module interface
│   │   └── internal/         # Private domain implementation
│   │       ├── entity/       # Domain entities
│   │       ├── gateway/      # Data access layer
│   │       ├── interactor/   # Business logic
│   │       ├── usecase/      # Use case implementations
│   │       └── mocks/        # Generated mocks
│   └── pkg/                  # Shared internal packages
│       ├── pkgerror/         # Custom error handling
│       ├── pkghttp/v1/       # HTTP utilities and patterns
│       ├── pkgsql/           # SQL utilities (GoQu wrappers)
│       ├── pkguid/           # ID generation (Snowflake)
│       └── pkgmocks/         # Generated mocks for shared packages
├── migration/                 # Database migrations (Goose)
└── seed/                     # Database seed data (Goose)
```

## Architecture Patterns
- **Clean Architecture**: Domain-driven with clear separation of concerns
- **Hexagonal Architecture**: Gateway pattern for external dependencies
- **Dependency Injection**: Dependencies injected via App struct
- **Domain Modules**: Self-contained business domains (loan module)

## App Struct
Central dependency container with:
- Database connection (*sql.DB)
- Query builder (GoQu)
- Validator (*validator.Validate)
- Logger (*zap.Logger)
- HTTP router (*httprouter.Router)
- HTTP server (*http.Server)
- Configuration (*viper.Viper)
- Snowflake ID generator
- Closer functions for graceful shutdown