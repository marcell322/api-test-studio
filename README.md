# API Test Studio

A lightweight REST API testing tool, built as a simplified alternative to Postman for learning and small-scale development.

## Problem & Audience

Postman and similar tools are powerful but heavy for simple day-to-day API testing. API Test Studio targets backend/full-stack developers, QA engineers, and students who want to send requests, organize them into collections, and review history — without the overhead of a full-featured client.

## Current Status

- ✅ User registration & login (JWT-based auth)
- ✅ Authenticated `/me` endpoint
- ✅ Collections (full CRUD, ownership-scoped per user)
- ✅ Saved Requests (full CRUD, ownership enforced through parent collection)
- ✅ Send Request + History (executes live HTTP requests via `/api/send`, logs both successes and failures)
- ⏳ Frontend (Vue 3) — not started yet

See [SPEC.md](./SPEC.md) for the full intended feature set and [CLAUDE.md](./CLAUDE.md) for the development guidelines used when working with AI on this project.

## Tech Stack

**Backend:** Go, Gin, GORM, SQLite (via `modernc.org/sqlite`, pure-Go driver — no CGO needed)
**Planned Frontend:** Vue 3, Vite, Axios, TailwindCSS

## Architecture

Clean architecture, backend only for now:
cmd/app              → entrypoint
internal/domain          → models, repository interfaces (no framework deps)
internal/usecase          → business logic (services)
internal/adapters/http        → Gin handlers
internal/adapters/persistence → GORM implementation of repositories
internal/adapters/auth      → JWT generation/validation
internal/middleware         → Gin middleware (auth guard)
internal/config           → env-based configuration

Handlers only validate input, call a service, and return a response — no business logic or SQL lives in the handler layer. Ownership checks (a user can only touch their own resources) live in the service layer, not the handler.

## Getting Started

```bash
cp .env.example .env
go run ./cmd/app
```

Server starts on the port set in `.env` (default `:8080`).

## API Overview

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/register` | POST | – | Create a new user |
| `/api/login` | POST | – | Log in, returns JWT |
| `/api/me` | GET | ✅ | Current user info |
| `/api/collections` | GET/POST | ✅ | List / create collections |
| `/api/collections/:id` | GET/PUT/DELETE | ✅ | Get / rename / delete a collection |
| `/api/requests` | GET/POST | ✅ | List / create saved requests (`?collection_id=` to filter) |
| `/api/requests/:id` | GET/PUT/DELETE | ✅ | Get / update / delete a saved request |
| `/api/send` | POST | ✅ | Execute an HTTP request live, logs to history |
| `/api/history` | GET | ✅ | List past executed requests |
| `/api/history/:id` | GET/DELETE | ✅ | Get / delete a history entry |

## Known Limitations

- **SSRF exposure**: `/api/send` lets an authenticated user make the server issue arbitrary outbound requests. Not currently restricted to public IP ranges — a production deployment should block requests to private/internal addresses before this is exposed beyond local development.
- No automated tests yet (manual `curl`-based verification only — see `AI_COLLABORATION.md`).
- No rate limiting on `/api/send`.

## Future Improvements

- Vue 3 frontend with JSON viewer
- Import/export Postman collections
- Unit tests for services and handlers
- SSRF protection (block private IP ranges)
- Docker support
- Swagger/OpenAPI docs

## AI Collaboration

This project was built in collaboration with AI assistants. See [AI_COLLABORATION.md](./AI_COLLABORATION.md) for specifics on what was AI-generated, what was changed and why, and how it was verified.