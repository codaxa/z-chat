# Z-Chat

A real-time chat application built with Go, WebSockets, and a PostgreSQL backend.

## Features

- **Real-time messaging** via WebSockets
- **User Authentication** with JWT
- **REST API** for registration, login, and message history
- **PostgreSQL Integration** for data persistence
- **Database Migrations** managed by Atlas and GORM
- **Hub Pattern** for concurrent client management
- **Configuration Management** using environment variables
- **Containerized** development environment with Docker

## Prerequisites

Before you begin, ensure you have the following installed:
- [Go](https://golang.org/doc/install) (version 1.23 or newer)
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [Task](https://taskfile.dev/installation/) for running project tasks
- [Atlas](https://atlasgo.io/cli/getting-started/setting-up) for database migrations

## Quick Start

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/codaxa/z-chat
    cd z-chat
    ```

2.  **Configure your environment:**
    Create a `.env` file in the root directory. You can use the following template, replacing the placeholder values with your local configuration.
    ```env
    # PostgreSQL Settings
    DB_USER=admin
    DB_PASSWORD=your_secure_password
    DB_NAME=z-chat
    DB_HOST=localhost
    DB_PORT=5432

    # For Docker Compose
    POSTGRES_USER=${DB_USER}
    POSTGRES_PASSWORD=${DB_PASSWORD}
    POSTGRES_DB=${DB_NAME}

    # JWT Settings
    JWT_SECRET=your_strong_jwt_secret
    ```

3.  **Start the database:**
    ```bash
    docker-compose up -d
    ```

4.  **Set up the database schema:**
    This command will create the database, run migrations, and get everything ready.
    ```bash
    task setup
    ```

5.  **Run the application:**
    ```bash
    task run
    ```
    The server will start on `localhost:8080`.

## Development

This project uses [Taskfile](https://taskfile.dev/) as a command runner.

-   **Run tests:**
    ```bash
    task test
    ```

-   **Run with hot-reloading (requires Air):**
    ```bash
    go install github.com/cosmtrek/air@latest
    task dev-air
    ```

-   **Run code quality checks:**
    ```bash
    task lint
    ```

-   **Create a new database migration:**
    After making changes to GORM models in `internal/domain/models/`, run:
    ```bash
    task migrate-new
    ```

-   **See all available tasks:**
    ```bash
    task --list
    ```

## Architecture

```
cmd/chatserver/     # Application entry point
internal/
├── config/         # Configuration management
├── context/        # Context keys for request-scoped values
├── domain/         # Core domain models and repository interfaces
├── handlers/       # HTTP and WebSocket handlers
├── hub/            # WebSocket hub & client management
├── middleware/     # HTTP middleware (e.g., authentication)
├── services/       # Business logic (e.g., authentication service)
├── storage/        # Database implementation (PostgreSQL)
└── transport/http/ # HTTP routing and server setup
migrations/         # Database migration files
```

## Dependencies

-   [go-chi/chi](https://github.com/go-chi/chi) - HTTP router & middleware
-   [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket implementation
-   [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit
-   [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - JWT implementation
-   [go-playground/validator](https://github.com/go-playground/validator) - Struct validation
-   [GORM](https://gorm.io/) - ORM for schema generation
-   [Atlas](https://atlasgo.io/) - Database schema migration tool