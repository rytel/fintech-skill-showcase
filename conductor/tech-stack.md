# Tech Stack

## Backend Core
- **Language:** Go 1.23
    - *Rationale:* High performance, excellent concurrency primitives, and strongly typed, making it ideal for financial systems and microservices.
- **Architecture:** Microservices
    - *Rationale:* To learn and implement independent scaling, distributed data consistency, and service-to-service communication.

## Data Storage & Integrity
- **Primary Database:** PostgreSQL 15
    - *Rationale:* Industry-standard relational database with robust support for ACID transactions and row-level locking (SELECT FOR UPDATE).
- **Migration Tool:** Custom SQL scripts (standard Go layout)
- **Database Driver:** `lib/pq`

## Security & Authentication
- **Authentication:** WebAuthn (Passkeys) & JWT (v5)
    - *Rationale:* Passkeys provide biometric-ready security for iOS, while JWTs handle stateless session management across microservices.
- **Communication:** TLS/HTTPS for all external and internal traffic.

## Infrastructure & DevOps
- **Containerization:** Docker & Docker Compose
    - *Rationale:* For consistent development environments and seamless deployment.
- **Cloud Provider:** Google Cloud Platform (GCP)
    - *Rationale:* To leverage managed services (Cloud Run/GKE, Cloud SQL) and professional infrastructure management.

## Testing & Quality
- **Unit Testing:** Go Standard Library (`testing`) & `sqlmock`
- **Integration Testing:** Dockerized environment for end-to-end flow validation.
- **API Documentation:** OpenAPI 3.0 (Swagger)
