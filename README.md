# API Test Studio

A lightweight REST API testing tool, built as a simplified alternative to Postman for learning and small-scale development.

## Problem & Audience

Postman and similar tools are powerful but heavy for simple day-to-day API testing. API Test Studio targets backend/full-stack developers, QA engineers, and students who want to send requests, organize them into collections, and review history — without the overhead of a full-featured client.

## Current Status

This project is a work in progress. What's implemented so far:

- ✅ User registration & login (JWT-based auth)
- ✅ Authenticated `/me` endpoint
- ✅ Collections (API routes scaffolded, persistence not yet implemented)
- 🚧 Saved Requests (API routes scaffolded, persistence not yet implemented)
- 🚧 Request History (API routes scaffolded, persistence not yet implemented)
- ⏳ Frontend (Vue 3) — not started yet

See [SPEC.md](./SPEC.md) for the full intended feature set and [CLAUDE.md](./CLAUDE.md) for the development guidelines used when working with AI on this project.

## Tech Stack

**Backend:** Go, Gin, GORM, SQLite (via `modernc.org/sqlite`, pure-Go driver — no CGO needed)
**Planned Frontend:** Vue 3, Vite, Axios, TailwindCSS

## Architecture

Clean architecture, backend only for now:
cmd/app          → entrypoint
internal/domain      → models, repository interfaces (no framework deps)
internal/usecase      → business logic (services)
internal/adapters/http    → Gin handlers
internal/adapters/persistence → GORM implementation of repositories
internal/adapters/auth   → JWT generation/validation
internal/middleware     → Gin middleware (auth guard)
internal/config       → env-based configuration

Handlers only validate input, call a service, and return a response — no business logic or SQL lives in the handler layer.

## Getting Started

```bash
cp .env.example .env
go run ./cmd/app
```

Server starts on the port set in `.env` (default `:8080`).

### Try it

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"secret123"}'

curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'
```

## Future Improvements

- Finish Collections / Saved Requests / Request History persistence
- Vue 3 frontend with JSON viewer
- Import/export Postman collections
- Unit tests for services and handlers
- Docker support
- Swagger/OpenAPI docs

## AI Collaboration

This project was built in collaboration with AI assistants. See [AI_COLLABORATION.md](./AI_COLLABORATION.md) for specifics on what was AI-generated, what was changed and why, and how it was verified.