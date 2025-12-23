import Testing
import Combine
import Foundation
@testable import demoBank

@MainActor
struct AuthServiceTests {
    
    @Test func loginReturnsSuccess() async throws {
        let mockAPI = MockAPIServiceForAuth()
        let service = AuthService(apiService: mockAPI)
        let request = LoginRequest(username: "test", password: "password")
        
        let publisher = service.login(credentials: request)
        
        // Manual continuation to avoid AsyncPublisher sendability issues in tests
        let response: LoginResponse = try await withCheckedThrowingContinuation { continuation in
            var cancellable: AnyCancellable?
            cancellable = publisher.first().sink(
                receiveCompletion: { completion in
                    if case .failure(let error) = completion {
                        continuation.resume(throwing: error)
                    }
                    _ = cancellable
                },
                receiveValue: { value in
                    continuation.resume(returning: value)
                }
            )
        }
        
        #expect(response.token == "mock-token")
    }
}

@MainActor
class MockAPIServiceForAuth: APIServiceProtocol {
    func login(username: String, password: String) async throws -> LoginResponse {
        return LoginResponse(token: "mock-token")
    }
    func fetchAccount(userId: String) async throws -> Account { 
        return Account(id: "1", userId: "1", balance: 0, createdAt: Date())
    }
    func fetchTransactions(userId: String) async throws -> [Transaction] { return [] }
    func performTransaction(userId: String, type: TransactionType, amount: Double) async throws -> Account {
        return try await fetchAccount(userId: userId)
    }
}