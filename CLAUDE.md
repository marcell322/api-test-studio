# CLAUDE.md

This document provides development guidelines for AI assistants contributing to this project.

---

# Project Overview

Project Name:

API Test Studio

Architecture:

Frontend (Vue 3)
↓

REST API

↓

Go (Gin)

↓

SQLite

---

# Development Principles

Always prioritize

- Readability
- Maintainability
- Simplicity

Avoid unnecessary complexity.

---

# Backend Rules

Use

- Go 1.24+
- Gin Framework
- GORM ORM
- SQLite

Project structure

internal/

    handlers/

    services/

    repositories/

    middleware/

    models/

Business logic should never be placed inside handlers.

Handlers should only

- Validate request
- Call service
- Return response

---

# API Design

RESTful conventions

Example

GET /collections

POST /collections

GET /collections/:id

PUT /collections/:id

DELETE /collections/:id

Return JSON only.

---

# Response Format

Success

{
    "success": true,
    "data": {}
}

Error

{
    "success": false,
    "message": "Invalid request"
}

---

# Database Rules

Never perform SQL directly inside handlers.

Always access the database through repositories.

Models should remain simple.

---

# Validation

Always validate

- Required fields
- Email format
- URL format
- Empty request body
- Invalid HTTP method

Return meaningful error messages.

---

# Error Handling

Never panic.

Return proper HTTP status codes.

Examples

200 OK

201 Created

400 Bad Request

401 Unauthorized

404 Not Found

500 Internal Server Error

---

# Frontend Rules

Use Vue 3 Composition API.

Separate

Pages

Components

Services

Reusable UI belongs in components.

API calls belong in services.

---

# Coding Style

Keep functions small.

Prefer composition over inheritance.

Avoid duplicated code.

Use descriptive variable names.

Comment only when necessary.

---

# Security

Hash passwords using bcrypt.

Never store plaintext passwords.

Protect authenticated endpoints using JWT.

Validate user input.

---

# Testing

Before committing

- Verify login
- Verify CRUD operations
- Verify API request sending
- Verify database persistence

---

# Git Commit Guidelines

Write meaningful commit messages.

Good examples

feat: implement JWT authentication

feat: add request history

fix: validate empty URL

refactor: simplify request service

docs: update README

Avoid

update

fix

changes

done

---

# AI Collaboration Guidelines

AI-generated code should

- Be reviewed
- Be tested manually
- Be refactored when necessary

Never merge AI-generated code without verification.

Document significant AI contributions in AI_COLLABORATION.md.