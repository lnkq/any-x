# Any-X — Anonymous Broadcast Board

> [!NOTE]
> `README.md` & `static/index.html` were written by AI.

Any-X is a minimal anonymous message board where anyone can post short messages and they are broadcast to connected clients in real time via WebSockets. Messages are optionally persisted to a local SQLite database.

Features
- Anonymous posting — no accounts or authentication.
- Real-time broadcast over WebSockets (`/ws`).
- Optional persistence using SQLite (`./storage.db` by default).
- Lightweight frontend served from `./static`.

Quick start (local)

Prerequisites
- Go 1.20+ installed.
- On macOS, ensure CGO is enabled if building with the SQLite driver (`github.com/mattn/go-sqlite3`).

Run directly (development)

```bash
# from repository root
go run ./cmd/any-x
```

Build and run

```bash
go build -o any-x ./cmd/any-x
./any-x
```