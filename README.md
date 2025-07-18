# Z-Chat

A real-time chat application built with Go and WebSockets.

## Features

- **Real-time messaging** via WebSockets
- **Hub pattern** for client management
- **Message validation** and broadcasting
- **REST API** with health monitoring

## Quick Start

```bash
git clone https://github.com/codaxa/z-chat
cd z-chat
go mod download
go run cmd/chatserver/main.go
```

Server runs on `localhost:8080`

## Development

```bash
# Run tests
go test ./...

# Hot reload (requires Air)
go install github.com/cosmtrek/air@v1.49.0
air

# Code quality
golangci-lint run
```

## Architecture

```
cmd/chatserver/     # Application entry point
internal/
├── config/         # Configuration
├── hub/            # WebSocket hub & client management
├── handlers/       # HTTP & WebSocket handlers
├── domain/models/  # User & Message models with validation
└── transport/http/ # Routing & middleware
```

## Dependencies

- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket implementation
- [go-chi/chi](https://github.com/go-chi/chi) - HTTP router & middleware
- [go-playground/validator](https://github.com/go-playground/validator) - Struct validation