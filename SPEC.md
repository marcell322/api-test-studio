# API Test Studio - Project Specification

## Project Overview

API Test Studio is a lightweight web application for testing REST APIs.

The application enables developers to organize API requests into collections, send HTTP requests, inspect responses, and review request history through a clean web interface.

This project is designed as a simplified alternative to Postman for learning and small-scale development.

**Status: functionally complete for the scope below.** See [README.md](./README.md) for current status by feature and known limitations.

---

# Problem Statement

Developers frequently need to test REST APIs during development.

Existing tools such as Postman provide extensive functionality but can be overwhelming for simple use cases.

This project aims to provide a lightweight, easy-to-use API testing tool focused on the most commonly used features.

---

# Target Users

- Backend Developers
- Full Stack Developers
- QA Engineers
- Students learning REST APIs

---

# Goals

The application should allow users to:

- Register an account ✅
- Login securely ✅
- Create API request collections ✅
- Save frequently used requests ✅
- Send HTTP requests ✅
- View formatted JSON responses ✅
- Store request history ✅

---

# Functional Requirements

## Authentication

### Register

- User can register using: ✅
  - Username
  - Email
  - Password

### Login

- Authenticate using JWT ✅

### Logout

- Remove authentication token ✅ (clears token client-side; token remains valid server-side until natural expiry — no server-side revocation/blocklist)

---

## Collections

Users can ✅

- Create collection
- Rename collection
- Delete collection

Each collection is scoped to the user that created it; users cannot see or modify each other's collections.

Example:

Development APIs

Production APIs

School Project

---

## Saved Requests

Each request contains ✅

- Request Name
- HTTP Method
- URL
- Headers
- Request Body
- Collection

Supported Methods ✅

- GET
- POST
- PUT
- PATCH
- DELETE

Ownership is enforced through the parent collection: a saved request can only be created in, viewed, or modified within a collection the current user owns.

---

## API Testing

Users can ✅

- Send request
- View status code
- View response headers
- View response body
- View response time

Every executed request — success or failure (network error, timeout, DNS failure) — is logged to history.

---

## Request History

Store ✅

- Method
- URL
- Status Code
- Response Time
- Timestamp

Failed requests are stored with `status_code: 0` and an error message, rather than being dropped.

---

## JSON Viewer

Pretty-print JSON response ✅

Support expand/collapse view ❌ **not implemented** — response body is pretty-printed but not an interactive collapsible tree. Listed under Future Improvements below.

---

# Non-functional Requirements

- Responsive UI — ⚠️ desktop-first; not yet tested/tuned for mobile breakpoints
- RESTful API ✅
- JWT Authentication ✅
- SQLite database ✅
- Clear project structure ✅
- Easy deployment — ⚠️ runs locally via `go run` / `npm run dev`; no Docker or deployment config yet

---

# Technology Stack

Backend

- Go
- Gin
- GORM
- SQLite

Frontend

- Vue 3
- Vite
- Pinia
- Vue Router
- Axios
- TailwindCSS

Version Control

- Git

---

# Known Limitations

See [README.md](./README.md#known-limitations) for the full list — notably: `/api/send` has no SSRF protection (can be pointed at internal/private addresses), and CORS is currently allowlisted only for the local Vite dev server.

---

# Future Improvements

- Import Postman Collection
- Export Collection
- Environment Variables
- JSON Viewer expand/collapse tree
- Docker Support
- Unit Tests (backend services/handlers, frontend components)
- API Documentation (Swagger)
- SSRF protection on `/api/send`
- Mobile-responsive layout

---

# Project Timeline

Phase 1 ✅

- Backend setup
- Authentication
- Database

Phase 2 ✅

- Collections
- Saved Requests
- Request History

Phase 3 ✅

- Frontend
- Manual end-to-end testing
- Documentation

---

# Success Criteria

The project is considered complete when users can

- Login successfully ✅
- Create collections ✅
- Save API requests ✅
- Send HTTP requests ✅
- View formatted responses ✅
- Review request history ✅

All success criteria met as of this revision.