# Tech Stack

## Backend Core
- **Language:** Go 1.23
    - *Rationale:* High performance, excellent concurrency primitives, and strongly typed.
- **Router:** `chi` (v5)
    - *Rationale:* Lightweight, idiomatic, and compatible with `net/http`.
- **Utilities:** `google/uuid`
    - *Rationale:* Standard for UUID generation.
- **Architecture:** Microservices
    - *Rationale:* To learn and implement independent scaling, distributed data consistency, and service-to-service communication.

## Data Storage & Integrity
- **Primary Database:** PostgreSQL 15
    - *Rationale:* Industry-standard relational database with robust support for ACID transactions and row-level locking (SELECT FOR UPDATE).
- **Migration Tool:** Custom SQL scripts (standard Go layout)
- **Database Driver:** `lib/pq`

## Mobile (iOS)
- **Language:** Swift 6
    - *Rationale:* Modern, safe, and the standard for native iOS development.
- **UI Framework:** SwiftUI
    - *Rationale:* Declarative UI framework for building modern interfaces across Apple platforms.
- **Architecture:** MVVM-C (Model-View-ViewModel-Coordinator)
    - *Rationale:* Standard architectural pattern for SwiftUI applications enhanced with Coordinators for decoupled navigation and Dependency Injection for testability.
- **Networking:** `URLSession` with `Codable`
    - *Rationale:* Native Apple frameworks for high-performance networking and JSON parsing.

## Security & Authentication
- **Authentication:** WebAuthn (Passkeys) & JWT (v5)
    - *Rationale:* Passkeys provide biometric-ready security for iOS, while JWTs handle stateless session management across microservices.
- **Communication:** TLS/HTTPS for all external and internal traffic.
- **Secure Storage:** iOS Keychain
    - *Rationale:* Hardware-backed secure storage for sensitive credentials like JWT tokens.

## Infrastructure & DevOps
- **Containerization:** Docker & Docker Compose
    - *Rationale:* For consistent development environments and seamless deployment.
- **Cloud Provider:** Google Cloud Platform (GCP)
    - *Rationale:* To leverage managed services (Cloud Run/GKE, Cloud SQL) and professional infrastructure management.

## Testing & Quality
- **Unit Testing:** Go Standard Library (`testing`), `sqlmock`, and `testify` (assert/mock)
- **iOS Testing:** Swift Testing
    - *Rationale:* The new, modern testing framework from Apple for Swift applications.
- **Integration Testing:** Dockerized environment for end-to-end flow validation.
- **API Documentation:** OpenAPI 3.0 (Swagger)