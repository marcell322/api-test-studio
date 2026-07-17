# AI Collaboration

## AI Tools Used

- ChatGPT (GPT-5.5) — architecture planning, initial scaffolding
- GitHub Copilot — inline code generation during implementation (see co-authored commits)

## How AI Helped

- Proposed the clean architecture folder layout (domain / usecase / adapters)
- Generated the initial Gin router and route grouping
- Drafted the JWT auth flow (token generation, validation, middleware)
- Proposed the SQLite schema for the User model
- Drafted the initial README and SPEC content

## Specific Example: Caught a Duplication Bug

While wiring up authentication, AI-generated code produced **two separate implementations** of the same JWT middleware:

- `internal/adapters/auth/jwt.go` — included its own `Middleware()` function
- `internal/middleware/auth.go` — a second, separate `AuthMiddleware()` that duplicated the same logic

Only the second one was actually registered in `router.go`; the first was dead code left over from an earlier generation pass. I caught this during review by tracing which functions were actually called from `router.go`, removed the unused `Middleware()` from `jwt.go`, and kept a single source of truth (`auth.ValidateToken` used by `middleware.AuthMiddleware`).

**Verification:** re-ran `go build ./...` to confirm nothing else referenced the removed function, then manually tested `/api/me` with and without a valid Bearer token to confirm auth still worked correctly after the cleanup.

## Specific Example: Adding Ownership Checks to Collection CRUD

The initial AI-generated draft of Collection CRUD (create/list/get/update/delete) worked correctly for a single user, but didn't check that a collection actually belonged to the requesting user before returning or modifying it. In other words, any logged-in user could read or delete another user's collection just by guessing an ID.

I asked for and added an explicit ownership check inside `CollectionService.Get`: it loads the collection, and returns `ErrForbidden` if `collection.UserID != userID`. Every other method (`Rename`, `Delete`) routes through `Get` first, so the check applies everywhere automatically instead of being duplicated per-handler.

I also introduced two shared, typed errors — `usecase.ErrNotFound` and `usecase.ErrForbidden` — so the handler layer can map them to the correct HTTP status (404 vs 403) instead of everything collapsing into a generic 400 or 500.

**Verification:** manually tested with two different user accounts — confirmed user B gets a 403 when trying to GET/PUT/DELETE a collection created by user A, and a 404 when the ID doesn't exist at all.
## Other Human Decisions

- Split validation logic (`validateRegister`, `validateLogin`) out of handlers into standalone functions per CLAUDE.md's "no business logic in handlers" rule — AI's first draft had validation inline in the handler
- Standardized error response shape (`{success, message}`) across all endpoints; AI's initial output was inconsistent between the register and login handlers
- Rewrote README to reflect actual implementation status rather than the originally AI-drafted aspirational version

## Verification Process

- Manual API testing via `curl` for register/login/me flows
- `go build ./...` after every structural change
- Code review of every AI-generated file before commit, checking against the rules in CLAUDE.md