# Daneizo

Daneizo is a state engine for P2P lending platforms, written in Go. It helps manage the states and transitions required in peer-to-peer lending workflows efficiently and reliably.

## Features

- Modular state engine for P2P lending lifecycle
- Database migration and seeding support using Goose
- Environment-based configuration
- Designed for extensibility and maintainability

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.18 or higher recommended)
- [Goose](https://github.com/pressly/goose) for database migrations
- MySQL database

### Setup

1. **Install Goose**

    ```sh
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

2. **Clone the repository**

    ```sh
    git clone https://github.com/farislr/daneizo.git
    cd daneizo
    ```

3. **Copy and configure environment variables**

    ```sh
    cp .env.example .env
    ```

    Edit the `.env` file and fill in the required values (e.g., database credentials).

### Database Migration

Run database migrations using Goose:

```sh
goose -dir migration/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" up
```

### Seed Database

Seed the database with initial data:

```sh
goose -no-versioning -dir seed/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" up
```

## Usage

_This section can be expanded with code examples, API documentation, or usage instructions as the project evolves._

## Contributing

Contributions are welcome! Please open issues or pull requests for bug fixes, features, or documentation improvements.

## License

[MIT License](LICENSE)

## Contact

For questions or support, open an issue or contact [@farislr](https://github.com/farislr).
