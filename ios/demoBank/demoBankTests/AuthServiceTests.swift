import Testing
import Combine
import Foundation
@testable import demoBank

@MainActor
final class AuthServiceTests {
    
    @Test func loginReturnsSuccess() async throws {
        let mockAPI = MockAPIServiceForAuth()
        let service = await AuthService(apiService: mockAPI)
        let request = LoginRequest(username: "test", password: "password")
        
        // Using AsyncPublisher to test Combine stream in a modern way
        var response: LoginResponse?
        for try await value in service.login(credentials: request).values {
            response = value
            break
        }
        
        #expect(response?.token == "mock-token")
    }
}

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
