# Spec: Account Service Microservice

## Overview
The Account Service is a dedicated microservice responsible for managing customer accounts, multi-currency balances, and core account operations. It replaces the account-related logic currently residing in the monolithic structure.

## Goals
- Provide a clean API for account management.
- Support multiple currencies per customer.
- Ensure strict data integrity and ACID compliance for balance updates.
- Decouple account management from other banking functions.

## Features
- **Account Management:**
    - Create new accounts for customers.
    - Retrieve account details and balances.
    - List accounts for a customer.
- **Multi-currency Support:**
    - Support for various currencies (e.g., PLN, USD, EUR).
    - Logic for validating currency codes.
- **Balance Operations:**
    - Atomic updates for deposits and withdrawals (internal API for other services).
    - Row-level locking to prevent race conditions.

## Technical Requirements
- **Language:** Go 1.23
- **Database:** PostgreSQL (dedicated schema or database for the service).
- **Authentication:** JWT validation for requests.
- **API:** RESTful with JSON (`camelCase`).
- **Communication:** Prepared for inter-service communication (REST/gRPC).

## Success Criteria
- Successful migration of account data (if applicable) or implementation of new schema.
- All account-related unit and integration tests passing.
- >80% code coverage.
- Documented API endpoints via OpenAPI.
