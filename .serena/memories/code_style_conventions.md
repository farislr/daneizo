# Code Style and Conventions

## Go Standards
- **Formatting**: Standard `gofmt` formatting (no custom formatting rules)
- **Naming**: Follow Go naming conventions (camelCase, PascalCase)
- **Package Names**: Short, lowercase, single word when possible
- **Import Organization**: Standard Go import grouping (std, external, internal)

## Project-Specific Conventions

### Package Structure
- **Internal packages**: Use `internal/` for private code
- **Domain modules**: Self-contained in `internal/[domain]/`
- **Shared utilities**: In `internal/pkg/pkg[name]/` with consistent naming

### Naming Patterns
- **Interfaces**: Often suffix-less (e.g., `InsertLoanStore`)
- **Implementations**: Descriptive names (e.g., `CreateProposedLoan`)
- **SQL Entities**: Suffix with `SQLEntity` (e.g., `LoanSQLEntity`)
- **Mocks**: Prefix with `Mock` (e.g., `MockCreateNewLoan`)
- **Test Suites**: Suffix with `Suite` (e.g., `loanSQLGatewaySuite`)

### Error Handling
- **Custom errors**: Use `internal/pkg/pkgerror` package
- **Error codes**: Defined in `pkgerror/code.go`
- **Error wrapping**: Follow Go 1.13+ error wrapping patterns

### Testing Patterns
- **Test suites**: Use testify/suite for complex test setups
- **Mocking**: Use mockery-generated mocks with expecter pattern
- **SQL mocking**: Use go-sqlmock for database testing
- **Test naming**: Follow `Test[Type]_[Method]` pattern

### Database Patterns
- **Query building**: Use GoQu query builder via `pkgsql` wrapper
- **Migrations**: Use Goose with numbered SQL files
- **Entities**: Separate SQL entities in `sqlentity` package

### HTTP Patterns
- **Routing**: Use httprouter for HTTP routing
- **Handlers**: Follow `pkghttp/v1` patterns for consistent HTTP handling
- **Endpoints**: Use endpoint pattern with options

### Configuration
- **Environment**: Use Viper for configuration management
- **Struct tags**: Use appropriate tags for validation and JSON marshaling