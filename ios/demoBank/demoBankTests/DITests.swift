import Testing
import Foundation
@testable import demoBank


@MainActor
struct DITests {
    
    @Test func dependencyContainerRegistersAndResolvesService() {
        let container = DependencyContainer.shared
        let mockService = MockAPIService()
        
        container.register(type: APIServiceProtocol.self, component: mockService)
        
        let resolvedService = container.resolve(type: APIServiceProtocol.self)
        #expect(resolvedService is MockAPIService)
    }
}


class MockAPIService: APIServiceProtocol {
    func login(username: String, password: String) async throws -> LoginResponse {
        return LoginResponse(token: "")
    }
    
    func fetchAccount(userId: String) async throws -> Account {
        return Account(id: "0", userId: "", balance: 0, createdAt: Date())
    }
    
    func fetchTransactions(userId: String) async throws -> [Transaction] {
        return []
    }
    
    func performTransaction(userId: String, type: TransactionType, amount: Double) async throws -> Account {
        return Account(id: "0", userId: "", balance: 0, createdAt: Date())
    }
}
