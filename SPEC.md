# API Test Studio - Project Specification

## Project Overview

API Test Studio is a lightweight web application for testing REST APIs.

The application enables developers to organize API requests into collections, send HTTP requests, inspect responses, and review request history through a clean web interface.

This project is designed as a simplified alternative to Postman for learning and small-scale development.

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

- Register an account
- Login securely
- Create API request collections
- Save frequently used requests
- Send HTTP requests
- View formatted JSON responses
- Store request history

---

# Functional Requirements

## Authentication

### Register

- User can register using:
  - Username
  - Email
  - Password

### Login

- Authenticate using JWT

### Logout

- Remove authentication token

---

## Collections

Users can

- Create collection
- Rename collection
- Delete collection

Example:

Development APIs

Production APIs

School Project

---

## Saved Requests

Each request contains

- Request Name
- HTTP Method
- URL
- Headers
- Request Body
- Collection

Supported Methods

- GET
- POST
- PUT
- PATCH
- DELETE

---

## API Testing

Users can

- Send request
- View status code
- View response headers
- View response body
- View response time

---

## Request History

Store

- Method
- URL
- Status Code
- Response Time
- Timestamp

---

## JSON Viewer

Pretty-print JSON response

Support expand/collapse view

---

# Non-functional Requirements

- Responsive UI
- RESTful API
- JWT Authentication
- SQLite database
- Clear project structure
- Easy deployment

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
- Axios
- TailwindCSS

Version Control

- Git

---

# Future Improvements

- Import Postman Collection
- Export Collection
- Environment Variables
- Dark Mode
- Docker Support
- Unit Tests
- API Documentation (Swagger)

---

# Project Timeline

Phase 1

- Backend setup
- Authentication
- Database

Phase 2

- Collections
- Saved Requests
- Request History

Phase 3

- Frontend
- Testing
- Documentation

---

# Success Criteria

The project is considered complete when users can

- Login successfully
- Create collections
- Save API requests
- Send HTTP requests
- View formatted responses
- Review request history