# Plan: Account Service Microservice

## Phase 1: Foundation & Data Model [checkpoint: 6807f88]
- [x] Task: Design the PostgreSQL schema for accounts and balances. ac81f6f
- [x] Task: Write tests for the account repository (ACID transactions, locking). 731bdbc
- [x] Task: Implement the account repository in Go. 7c17e13
- [x] Task: Conductor - User Manual Verification 'Phase 1: Foundation & Data Model' (Protocol in workflow.md)

## Phase 2: Core Logic & Service Layer
- [ ] Task: Write tests for account creation and retrieval logic.
- [ ] Task: Implement service layer for account management.
- [ ] Task: Write tests for multi-currency balance validation.
- [ ] Task: Implement balance update logic with safety checks.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Core Logic & Service Layer' (Protocol in workflow.md)

## Phase 3: API & Connectivity
- [ ] Task: Define OpenAPI specification for the Account Service.
- [ ] Task: Write tests for REST handlers (JSON mapping, error responses).
- [ ] Task: Implement REST API handlers using Go.
- [ ] Task: Integrate JWT authentication for protected endpoints.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: API & Connectivity' (Protocol in workflow.md)

## Phase 4: Integration & Cleanup
- [ ] Task: Write integration tests for the full account lifecycle.
- [ ] Task: Implement health check endpoints and logging.
- [ ] Task: Refactor existing monolithic references to use the new Account Service (mocking/proxying).
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Integration & Cleanup' (Protocol in workflow.md)
