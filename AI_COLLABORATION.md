# AI Collaboration

## Summary

Three complete, tested vertical slices were built collaboratively with AI across this project:

1. **Auth** (register/login/me) — initial AI-generated scaffold, cleaned up a duplicate JWT middleware implementation during review (see above)
2. **Collections** — full CRUD with per-user ownership enforcement
3. **Saved Requests** — full CRUD with ownership enforced through the parent collection (cross-resource check, not just direct ownership)
4. **Send Request + History** — live HTTP execution with failure-path logging

Each feature followed the same review discipline: read the AI-generated code against CLAUDE.md's
rules before committing, manually tested the happy path and the failure/security paths (404, 403,
network failures) via curl, and fixed issues found during that testing (e.g. verifying test tokens
actually belonged to different users before trusting a 403/leak result, rather than assuming the
code was right or wrong from the first symptom).

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

## Verified: Collections CRUD

Tested manually via curl for create, list, get-by-id, update, and delete, including:
- 404 when requesting a non-existent collection ID
- 403 when requesting another user's collection (ownership check in `CollectionService.Get`)
- Confirmed via `/api/me` that the two test tokens belonged to distinct users before testing the 403 case, to rule out a test setup mistake before assuming a code bug

## Verified: Saved Requests CRUD

Tested manually via curl for create, get, update, delete, and list (both global and scoped to a collection).
Specifically verified the collection→user ownership chain: registered a second user (bob), confirmed
via /api/me that the token belonged to a distinct user_id before testing, then confirmed that user 2
cannot GET or POST a saved request against a collection owned by user 1 (403, no row created).

## Verified: Send Request + History

Tested successful request execution (httpbin.org, 200) and failure handling (invalid DNS host, 502
returned to caller). Confirmed both successful and failed attempts are persisted to history —
failed attempts log status_code: 0 with the underlying error message, which required explicitly
logging inside the error branch of HistoryService.Send rather than only after a successful response.

## Known limitation: SSRF exposure

/api/send allows an authenticated user to make the server issue arbitrary outbound HTTP requests.
Not currently restricted to public IP ranges — a production deployment should block requests to
private/internal addresses (127.0.0.1, 10.0.0.0/8, 169.254.169.254, etc.) before this is exposed
beyond local development.

## Other Human Decisions

- Split validation logic (`validateRegister`, `validateLogin`) out of handlers into standalone functions per CLAUDE.md's "no business logic in handlers" rule — AI's first draft had validation inline in the handler
- Standardized error response shape (`{success, message}`) across all endpoints; AI's initial output was inconsistent between the register and login handlers
- Rewrote README to reflect actual implementation status rather than the originally AI-drafted aspirational version

## Verification Process

- Manual API testing via `curl` for register/login/me flows
- `go build ./...` after every structural change
- Code review of every AI-generated file before commit, checking against the rules in CLAUDE.md

## Frontend: A Different Collaboration Pattern

Unlike the backend, which was built commit-by-commit with review at each step, the Vue 3 frontend
was AI-generated in a single large pass (Login/Register views, Pinia auth store, sidebar with
collection tree, request builder, response viewer, history panel — roughly 15 files at once).
Worth being upfront about: this means less granular review than the backend got. What I did verify
before committing:

- The AI-generated project actually builds (`npm install && npm run build`) with zero errors before
  I accepted any of it — no untested scaffold got merged.
- Manually exercised every screen against the real backend: register → login → create collection →
  add saved request → send a live request → check history → replay from history → log out.
- Found the CORS gap myself by actually running the app in a browser rather than assuming the API
  and frontend would talk to each other — the backend had no CORS middleware, so the first real
  browser test would have failed with every request blocked. Added `middleware/cors.go` to fix it.

## Verified: End-to-End (Frontend + Backend)

Ran both servers together (`go run ./cmd/app` + `npm run dev`) and manually tested the full user
flow in a browser rather than just curl: registration, login persistence across refresh (JWT in
localStorage), creating a collection, saving a request into it, sending a live request and viewing
the formatted response, and viewing/replaying from history. Confirmed a saved request's headers
round-trip correctly between the backend's JSON-string storage and the frontend's key/value editor
UI, since that's a boundary where a format mismatch would silently break with no error.