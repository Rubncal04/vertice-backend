# Vertice Backend API

API system for managing products and orders of Vertice, built with Go and Echo framework.

## Table of Contents
- [Description](#description)
- [Requirements](#requirements)
- [Installation](#installation)
- [Running Locally](#running-locally)
- [Running Tests](#running-tests)
- [API Documentation (Swagger)](#api-documentation-swagger)
- [Docker Usage](#docker-usage)
- [Example API Usage](#example-api-usage)

## Description
Vertice Backend API is a RESTful service for managing users, products, and orders. It provides authentication, product management, and order processing endpoints, following best practices for clean architecture and API design.

## Requirements
- Go 1.18+
- (Optional) Docker
- (Optional) PostgreSQL or your preferred database

## Installation
1. Clone the repository:
   ```sh
   git clone <repo-url>
   cd vertice-backend
   ```
2. Install dependencies:
   ```sh
   go mod download
   ```
3. Copy the example environment file and configure it:
   ```sh
   cp .env.example .env
   # Edit .env with your DB credentials and JWT secret
   ```

## Running Locally
1. Run database migrations (if needed):
   ```sh
   go run cmd/main.go # The app will auto-migrate on startup
   ```
2. Start the server:
   ```sh
   go run cmd/main.go
   ```
   The API will be available at `http://localhost:8080` by default.

## Running Tests
To run all tests:
```sh
 go test ./tests/... -v
```

## API Documentation (Swagger)
Interactive API docs are available via Swagger UI:
- Start the server
- Visit: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

To regenerate the docs after editing comments:
```sh
swag init -g cmd/main.go
```

## Docker Usage
Build and run the app using Docker:

- Build all services (app & database):
  ```sh
  docker-compose up --build
  ```
- Stop services:
  ```sh
  docker-compose down
  ```
- Watch logs on real time:
  ```sh
  docker-compose logs -f
  ```
- Rebuild image of the app:
  ```sh
  docker-compose build --no-cache
  ```
- Access to the PostgreSQL database from container:
  ```sh
  docker-compose exec db psql -U $DB_USER -d $DB_NAME
  ```

## Example API Usage

### Register User
**Request**
```http
POST /api/v1/users/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```
**Success Response**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```
**Error Response**
```json
{
  "error": "email already exists"
}
```

### Login
**Request**
```http
POST /api/v1/users/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```
**Success Response**
```json
{
  "token": "<jwt-token>"
}
```
**Error Response**
```json
{
  "error": "invalid credentials"
}
```

### Create Product
**Request**
```http
POST /api/v1/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "code": "PROD001",
  "name": "Laptop",
  "description": "High performance laptop",
  "price": 1299.99,
  "stock": 10
}
```
**Success Response**
```json
{
  "id": 1,
  "code": "PROD001",
  "name": "Laptop",
  "description": "High performance laptop",
  "price": 1299.99,
  "stock": 10
}
```
**Error Response**
```json
{
  "error": "product code already exists for this user"
}
```

### Create Order
**Request**
```http
POST /api/v1/orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "items": [
    { "product_id": 1, "quantity": 2 }
  ]
}
```
**Success Response**
```json
{
  "id": 1,
  "status": "pending",
  "total_amount": 2599.98,
  "items": [
    {
      "id": 1,
      "product_id": 1,
      "product": {
        "id": 1,
        "code": "PROD001",
        "name": "Laptop",
        "price": 1299.99
      },
      "quantity": 2,
      "unit_price": 1299.99,
      "subtotal": 2599.98
    }
  ],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```
**Error Response**
```json
{
  "error": "insufficient stock"
}
```

---

Feel free to contribute or open issues for improvements!
