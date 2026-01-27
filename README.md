# MusicStreaming

[![Build](https://img.shields.io/github/actions/workflow/status/FGguy/MusicStreaming/CI.yaml?label=build&logo=go)](https://github.com/FGguy/MusicStreaming/actions/workflows/CI.yaml)
[![Tests](https://img.shields.io/github/actions/workflow/status/FGguy/MusicStreaming/CI.yaml?label=tests&logo=go)](https://github.com/FGguy/MusicStreaming/actions/workflows/CI.yaml)
[![Lint](https://img.shields.io/github/actions/workflow/status/FGguy/MusicStreaming/CI.yaml?label=lint&logo=go)](https://github.com/FGguy/MusicStreaming/actions/workflows/CI.yaml)
[![Security](https://img.shields.io/github/actions/workflow/status/FGguy/MusicStreaming/CI.yaml?label=security&logo=go)](https://github.com/FGguy/MusicStreaming/actions/workflows/CI.yaml)

A self-hosted music streaming server implementing the Subsonic API, built with Go.

## Table of Contents

- [About](#about)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Development](#development)
- [API Compatibility](#api-compatibility)
- [License](#license)
- [Contact](#contact)


## About

MusicStreaming is an open-source, self-hosted solution for streaming your personal music collection. It provides full control over your audio library, allowing you to stream to any device using Subsonic-compatible clients without relying on third-party cloud services.

The project is built following **hexagonal architecture** (ports and adapters) principles, ensuring clean separation of concerns, high testability, and maintainability.


## Features

- **Subsonic API Compatibility**: Use any Subsonic-compatible client (DSub, Ultrasonic, Sublime Music, etc.)
- **Multi-Format Support**: Serve local music files (MP3, FLAC, OGG, etc.) over the network
- **User Management**: Multi-user support with role-based access control (RBAC)
- **Media Browsing**: Browse artists, albums, songs, and cover art
- **Authentication**: Secure user authentication and authorization
- **Media Scanning**: Automatic scanning of music directories
- **Caching**: Redis-powered caching for improved performance
- **Type-Safe Database**: SQLC-generated code for safe database operations
- **Streaming**: HTTP/HTTPS streaming to any device
- **Lightweight**: Easy to run on home servers, Raspberry Pi, or small VPS

## Tech Stack

### Core
- **Go 1.23+**: Modern Go with generics and improved performance
- **Gin**: Fast HTTP web framework for handlers
- **PostgreSQL**: Primary data persistence with pgx driver
- **Redis**: Caching layer for improved performance

### Database & Code Generation
- **SQLC**: Type-safe SQL code generation
- **pgx/v5**: High-performance PostgreSQL driver

### Configuration & Logging
- **Viper**: Configuration management
- **slog**: Structured logging (JSON format)
- **godotenv**: Environment variable management

### Testing
- **testify**: Assertions and test utilities
- **mockery**: Mock generation for interfaces

### DevOps
- **Docker**: Containerization
- **Docker Compose**: Local development environment
- **GitHub Actions**: CI/CD pipeline
- **golangci-lint**: Code quality and linting

## Prerequisites

Before running this application, ensure you have:

- **Go 1.23+** installed ([Download](https://golang.org/dl/))
- **PostgreSQL 14+** running locally or via Docker
- **Redis 7+** running locally or via Docker
- **Make** (optional, for using Makefile commands)


## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/FGguy/MusicStreaming.git
cd MusicStreaming
```

### 2. Start Dependencies with Docker Compose

```bash
cd compose
docker-compose up -d
cd ..
```

This starts PostgreSQL and Redis containers.

### 3. Set Up Environment Variables

Create a `.env` file in the project root:

```bash
# Database
POSTGRES_CONNECTION_STRING=postgres://musicstreaming:musicstreaming@localhost:5432/musicstreaming?sslmode=disable

# Redis
REDIS_CONNECTION_STRING=localhost:6379

# Configuration file path (optional)
CONFIG_PATH=./musicstreaming.yaml
```

### 4. Initialize the Database

```bash
# Apply the database schema
psql $POSTGRES_CONNECTION_STRING -f internal/adapter/sql/tables.sql
```

### 5. Build and Run

```bash
# Build the application
go build -o musicstreaming ./cmd/app

# Run the application
./musicstreaming --loglevel info
```

The server will start on `http://localhost:8080`.

## Configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `POSTGRES_CONNECTION_STRING` | PostgreSQL connection URL | `postgres://user:pass@localhost:5432/dbname` |
| `REDIS_CONNECTION_STRING` | Redis address | `localhost:6379` |
| `CONFIG_PATH` | Path to YAML config file | `./musicstreaming.yaml` |

### Configuration File

Create `musicstreaming.yaml` for application settings:

```yaml
music-directories:
  - /path/to/music/folder1
  - /path/to/music/folder2
```

### Command-Line Flags

- `--loglevel`: Set logging level (info, debug, warn, error)
  ```bash
  ./musicstreaming --loglevel debug
  ```


## Development

### Project Setup

```bash
# Install dependencies
go mod download

# Install development tools
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/vektra/mockery/v2@latest
```

### Available Make Commands

```bash
# Build the application
make build

# Run linter
make lint

# Run unit tests
make test

# Run integration tests (requires Docker)
make integration

# Generate SQLC code from SQL queries
make sqlc

# Generate mocks for testing
make mocks

# Run all code generation
make generate
```

### Database Migrations

When modifying the database schema:

1. Update `internal/adapter/sql/tables.sql`
2. Update queries in `internal/adapter/sql/queries/`
3. Regenerate SQLC code: `sqlc generate`
4. Update repository implementations if needed

## API Compatibility

This server implements the [Subsonic API specification](http://www.subsonic.org/pages/api.jsp), allowing compatibility with a wide range of clients:

### Compatible Clients

- **Android**: DSub, Ultrasonic, substreamer
- **iOS**: play:Sub, substreamer
- **Desktop**: Sublime Music, Supersonic, Sonixd
- **Web**: Airsonic Web UI, Jamstash

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

If you have questions or suggestions, feel free to:
- Open an issue on GitHub
- Contact [@FGguy](https://github.com/FGguy) via GitHub
