# Project Purpose

**Daneizo** is a state engine for P2P lending platforms, written in Go. It helps manage the states and transitions required in peer-to-peer lending workflows efficiently and reliably.

## Key Features
- Modular state engine for P2P lending lifecycle
- Database migration and seeding support using Goose
- Environment-based configuration
- Designed for extensibility and maintainability

## Tech Stack
- **Language**: Go 1.23
- **Database**: MySQL
- **Migration Tool**: Goose
- **HTTP Router**: httprouter (julienschmidt/httprouter)
- **Query Builder**: GoQu (doug-martin/goqu)
- **Logging**: Zap (go.uber.org/zap)
- **Validation**: go-playground/validator/v10
- **ID Generation**: Snowflake (bwmarrin/snowflake)
- **Configuration**: Viper (spf13/viper)
- **Testing**: testify (stretchr/testify)
- **Mocking**: go-sqlmock (DATA-DOG/go-sqlmock), mockery