import Foundation
import Combine

final class MockData {
    static let token = "mock-jwt-token"
    static let account = Account(
        id: "123e4567-e89b-12d3-a456-426614174000",
        userId: "test_user",
        balance: 15750.50,
        createdAt: Date()
    )
    
    static let transactions = [
        Transaction(id: "1", accountId: "acc1", type: .deposit, amount: 2000.0, createdAt: Date().addingTimeInterval(-86400), description: nil),
        Transaction(id: "2", accountId: "acc1", type: .withdrawal, amount: 150.0, createdAt: Date().addingTimeInterval(-43200), description: nil),
        Transaction(id: "3", accountId: "acc1", type: .deposit, amount: 500.0, createdAt: Date().addingTimeInterval(-3600), description: nil)
    ]
}

final class MockAuthService: AuthServiceProtocol {
    func login(credentials: LoginRequest) -> AnyPublisher<LoginResponse, Error> {
        if credentials.username == "test_user" && credentials.password == "password123" {
            return Just(LoginResponse(token: MockData.token))
                .setFailureType(to: Error.self)
                .delay(for: .seconds(0.5), scheduler: RunLoop.main)
                .eraseToAnyPublisher()
        } else {
            return Fail(error: APIError.unauthorized)
                .delay(for: .seconds(0.5), scheduler: RunLoop.main)
                .eraseToAnyPublisher()
        }
    }
}

final class MockAPIService: APIServiceProtocol {
    private func simulateNetworkDelay() async throws {
        try await Task.sleep(nanoseconds: 500_000_000) // 0.5 sekundy
    }
    
    func login(username: String, password: String) async throws -> LoginResponse {
        try await simulateNetworkDelay()
        if username == "test_user" && password == "password123" {
            return LoginResponse(token: MockData.token)
        }
        throw APIError.unauthorized
    }
    
    func fetchAccount(userId: String) async throws -> Account {
        try await simulateNetworkDelay()
        return MockData.account
    }
    
    func fetchTransactions(userId: String) async throws -> [Transaction] {
        try await simulateNetworkDelay()
        return MockData.transactions
    }
    
    func performTransaction(userId: String, type: TransactionType, amount: Double) async throws -> Account {
        try await simulateNetworkDelay()
        var newAccount = MockData.account
        if type == .withdrawal {
            newAccount.balance -= amount
        } else {
            newAccount.balance += amount
        }
        return newAccount
    }
}