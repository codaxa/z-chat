# Z-Chat

A real-time chat application built with Go and WebSockets.

## Features

- Real-time messaging via WebSockets
- Client connection management with hub pattern
- Message and user validation
- HTTP REST API

## Quick Start

```bash
# Clone and run
git clone <repository-url>
cd z-chat
go mod download
go run cmd/chatserver/main.go
```

Server runs on `localhost:8080`

## API Endpoints

- `GET /health` - Health check
- `GET /ws` - WebSocket connection for real-time chat

## WebSocket Usage

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (event) => console.log('Message:', event.data);
ws.send('Hello!');
```

## Development

```bash
# Run tests
go test ./...

# Code quality
golangci-lint run
```

## Architecture

```
├── cmd/chatserver/          # Entry point
├── internal/
│   ├── hub/                # WebSocket hub & clients
│   ├── handlers/           # HTTP handlers
│   ├── domain/models/      # User & Message models
│   └── transport/http/     # Routing
```

## Dependencies

- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket support
- [go-chi/chi](https://github.com/go-chi/chi) - HTTP router
- [go-playground/validator](https://github.com/go-playground/validator) - Validation