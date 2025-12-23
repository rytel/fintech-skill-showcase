# Fintech Skill Showcase

A complete Fullstack banking system demonstrating advanced engineering skills in building a secure financial backend and a modern iOS mobile application.

The project consists of two main modules:
- **Backend (`/backend`)**: A Go-based server ensuring ACID transactions and providing a REST API.
- **iOS (`/ios`)**: A SwiftUI mobile application implementing the client interface.

---

## 🏗 Architecture

### 1. Backend (Go)
A professional API server written in Go (1.23), implementing the bank's core business logic.
- **Key Features**:
    - **ACID Transactions**: Row-level locking (`SELECT ... FOR UPDATE`) to prevent race conditions during financial operations.
    - **Architecture**: Clean Architecture / Hexagonal with full Dependency Injection.
    - **Database**: PostgreSQL 15.
    - **Security**: JWT Tokens (demonstration version), SQL migrations on startup.

### 2. iOS (SwiftUI)
A modern client application written in Swift 5+.
- **Key Features**:
    - **UI**: SwiftUI.
    - **Patterns**: MVVM + Coordinator (navigation) + Dependency Injection Container.
    - **Networking**: Modern `async/await` with robust error handling.
    - **Testing**: Unit Tests (XCTest) and UI Tests.

---

## ⚙️ Setup Instructions (Local Development)

To run the system locally, you will need:
- Docker & Docker Compose
- Xcode 15+ (for iOS)
- Go 1.23+ (optional, for local backend development without Docker)

### Step 1: Run the Backend
The mobile application requires a running API. The backend must be started first.

```bash
cd backend
# Run in "detached" mode (background) to free up the terminal
docker-compose up -d --build
```
The server will be available at `http://localhost:8080` after a few moments.

### Step 2: Run the iOS Application
1. Open the project in Xcode:
   ```bash
   open ios/demoBank/demoBank.xcodeproj
   ```
2. Select a simulator (e.g., iPhone 16 Pro).
3. Run the app using **Cmd+R**.

The app connects to `localhost:8080`. You can log in using the test credentials:
- **Username:** `test_user`
- **Password:** `password123`

---

## 🧪 Testing

### Backend
Backend tests are self-contained (Docker Compose handles the database for integration tests).

```bash
cd backend
make test
```

### 📱 iOS (Important!)
⚠️ **NOTE:** UI tests (`demoBankUITests`) perform real network requests and log into the system (Integration Tests).

**For iOS UI tests to pass, the backend MUST be running in the background.** If the backend is not active, UI tests will fail due to connection errors.

**Correct testing procedure:**

1. **Ensure the backend is running:**
   ```bash
   cd backend && docker-compose up -d
   ```

2. **Run iOS tests:**
   - **In Xcode:** Press **Cmd+U**.
   - **In Terminal:**
     ```bash
     xcodebuild test -project ios/demoBank/demoBank.xcodeproj \
     -scheme demoBank \
     -destination 'platform=iOS Simulator,name=iPhone 16 Pro'
     ```

3. **(Optional) Clean up after testing:**
   ```bash
   cd backend && docker-compose down
   ```