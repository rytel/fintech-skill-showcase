# Product Guidelines

## 1. Professional & Educational Tone
- **Code Style:** All code must be production-grade, adhering to Go's best practices (idiomatic Go).
- **Explanatory Comments:** While the code remains professional, complex architectural decisions, financial logic, and design patterns should be accompanied by "Why" comments to support the educational goal of the project.

## 2. Security & Data Privacy (Bank-Grade)
- **Defense in Depth:** Implement a multi-layered security approach:
    - **Internal:** Detailed logging with full context (RequestID, Trace) for debugging.
    - **External:** Standardized, opaque error codes for the API to prevent information leakage.
    - **Business Logic:** Strict "Fail-Fast" validation on all inputs before processing.
- **Data Masking:** Sensitive Personal Identifiable Information (PII) and financial secrets must be masked in logs.
- **Modern Auth:** Native support for Passkeys (WebAuthn) to facilitate biometric authentication on iOS.

## 3. Development Methodology
- **Test-Driven Development (TDD):** Every feature must start with a failing test. This ensures high reliability and helps document the expected behavior of financial operations.
- **Domain-Driven Design (DDD):** Use precise banking terminology in the codebase (e.g., `Customer` instead of `User`, `Ledger` instead of `Transactions`).

## 4. API Standards (Mobile-First)
- **JSON Convention:** Use `camelCase` for all API request and response fields to ensure seamless integration with Swift (iOS).
- **Documentation:** Maintain a comprehensive OpenAPI (Swagger) specification for all endpoints.
- **Idempotency:** Implement idempotency keys for all state-changing operations (transfers, deposits) to prevent duplicate transactions.

## 5. Reliability & Scalability
- **ACID Compliance:** No compromises on data integrity. Use database transactions and row-level locking for all monetary movements.
- **Cloud-Native:** The system must be containerized (Docker) and designed to run efficiently on Google Cloud Platform (GCP).
