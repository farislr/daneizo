# Daneizo

Daneizo is state engine for p2p lending.

## Setup

- Install goose -> go install github.com/pressly/goose/v3/cmd/goose@latest
- Copy .env.example to .env
- Fill in the .env file

## Migrations

- goose -dir migration/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" up

## Seed

- goose -no-versioning -dir seed/ mysql "user:pass@tcp(localhost:3306)/daneizo?parseTime=true" up
