# iOS App API Integration Guide

This server has been expanded to support iOS application development and testing.

## Base URL
`http://localhost:8080` (Local)

## Test Environment Setup

To ensure a predictable state for UI tests, use the Reset endpoint before running your tests.

### Reset / Seed Database
**POST** `/api/test/reset`
- **Description**: Truncates all tables and creates a default user `test_user` with balance `1000.0`.
- **Response**: `200 OK` "Środowisko testowe zresetowane..."

## Authentication (JWT)

### Login (Mock)
**POST** `/api/login`
- **Request Body**:
  ```json
  {
    "username": "test_user",
    "password": "password123"
  }
  ```
- **Response**:
  ```json
  {
    "token": "eyJh... (JWT Token)"
  }
  ```
- **Note**: This is a mock login. Only `test_user` / `password123` works. The token is valid for 24 hours.

## Account & Transactions

### Get Account Balance
**GET** `/api/account/{user_id}`
- **Response**:
  ```json
  {
    "id": 1,
    "user_id": "test_user",
    "balance": 1000,
    "created_at": "..."
  }
  ```

### Get Transaction History
**GET** `/api/account/{user_id}/transactions`
- **Response**:
  ```json
  [
    {
      "id": 1,
      "account_id": 1,
      "type": "DEPOSIT",
      "amount": 100,
      "created_at": "..."
    },
    ...
  ]
  ```

### Create Transaction
**POST** `/api/transactions`
- **Request Body**:
  ```json
  {
    "user_id": "test_user",
    "type": "DEPOSIT", // or "WITHDRAW"
    "amount": 50.0
  }
  ```
- **Response**: Updated Account Object.

## Notes for iOS Developer
- Use the `Reset` endpoint to start fresh.
- Store the JWT token from `/api/login` in the Keychain.
- Send the token in the header: `Authorization: Bearer <token>` (Server has `AuthMiddleware` ready but endpoints are currently open for easier initial integration. Use `AuthMiddleware` to wrap handlers in `main.go` when ready to enforce security).
