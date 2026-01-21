# Any-X

Anonymous real-time broadcast board. Clients connect over WebSocket and receive new posts instantly. A small static UI is served from `./static`. Posts can be stored in SQLite and replayed on connect.

## Features

- Anonymous posting (no accounts/auth).
- Real-time updates over WebSocket at `/ws`.
- Optional persistence + history replay (SQLite).
- Static frontend served from `/`.

## Quick start

### Requirements

- Go 1.20+
- CGO enabled (needed for `github.com/mattn/go-sqlite3`)

### Run

```bash
go run ./cmd/any-x
```

Open:

- http://localhost:8080

### Build

```bash
go build -o any-x ./cmd/any-x
./any-x
```

## How it works

- HTTP serves `./static` at `/`.
- WebSocket endpoint: `ws://localhost:8080/ws`
	- Send a plain text message.
	- Server broadcasts a JSON `Post` to other clients.
	- With storage enabled, a new client receives up to 50 latest posts.

`Post` fields:

- `content` (string)
- `created_at` (timestamp)

## Configuration

Config is loaded from YAML (default: `config/local.yaml`).

- `-config path/to/config.yaml` or `CONFIG_PATH=path/to/config.yaml`

Settings:

| Key | Env | Default | Description |
| --- | --- | --- | --- |
| `env` | `ENV` | `local` | Environment name (currently used for logging setup). |
| `storage_path` | `STORAGE_PATH` | `./storage.db` | Path to SQLite DB. Set to empty to disable persistence. |
| `http_server.address` | `HTTP_ADDRESS` | `localhost:8080` | Bind address. |
| `http_server.timeout` | `HTTP_TIMEOUT` | `4s` | Read/write timeout. |
| `http_server.idle_timeout` | `HTTP_IDLE_TIMEOUT` | `60s` | Idle timeout. |

Example:

```bash
HTTP_ADDRESS=0.0.0.0:8080 STORAGE_PATH=./storage.db go run ./cmd/any-x
```

## Project layout

- `cmd/any-x`: entrypoint
- `internal/app`: HTTP server + routing
- `internal/web-app/handlers/websocket`: WebSocket handler
- `internal/storage/sqlite`: SQLite persistence implementation
- `static/`: frontend
