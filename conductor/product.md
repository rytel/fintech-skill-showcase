# Initial Concept
This app is being created for educational purposes to learn how to build a working bank from scratch.

# Product Vision: Edu-Bank Backend
A professional, high-integrity banking server implemented in Go, designed to serve as a secure and scalable backend for a future iOS mobile application. The project follows a **Cloud-Native approach**, utilizing containerization and cloud infrastructure to ensure production-grade reliability.

## Core Objectives
- **Operational Integrity:** Ensure 100% correctness of financial operations using ACID-compliant transactions and row-level locking.
- **Modern Security:** Implement industry-leading authentication (Passkeys/WebAuthn) and comprehensive audit logging.
- **Infinite Scalability:** Utilize a microservices architecture to allow independent scaling of core banking functions.
- **Cloud-Native Foundations:** Deployment-ready via **Docker** and optimized for **Google Cloud Platform (GCP)** to learn modern DevOps and infrastructure management.
- **Mobile-First Connectivity:** Provide a robust, well-documented RESTful API optimized for integration with iOS (Swift).

## Key Features
- **Account & Multi-currency Management:** Support for multiple account types and built-in currency exchange (FX) capabilities.
- **Transactional Engine:** Distributed handling of deposits, withdrawals, and peer-to-peer transfers with consistency guarantees.
- **Security & Audit:** 
    - Passkey-based authentication for biometric-ready login (FaceID/TouchID).
    - Immutable audit logs for every state change and access attempt.
- **Asynchronous Processing:** Background task handling for push notifications and reporting.

## Target Architecture
- **Microservices:** A distributed system of specialized services (e.g., Auth Service, Transaction Service, FX Service) containerized with Docker for seamless deployment and scaling.
