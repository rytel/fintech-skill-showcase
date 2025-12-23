import Testing
import Foundation
@testable import demoBank


struct DashboardViewModelTests {
    
    @Test func fetchAccountUpdatesState() async {
        let mockAPI = MockAPIServiceForDashboard()
        let viewModel = DashboardViewModel(apiService: mockAPI)
        
        await viewModel.fetchData()
        
        #expect(viewModel.account != nil)
        #expect(viewModel.account?.balance == 1234.56)
        #expect(!viewModel.isLoading)
    }
}


class MockAPIServiceForDashboard: APIServiceProtocol {
    func login(username: String, password: String) async throws -> LoginResponse {
        return LoginResponse(token: "token")
    }
    
    func fetchAccount(userId: String) async throws -> Account {
        return Account(id: "1", userId: "u1", balance: 1234.56, createdAt: Date())
    }
    
    func fetchTransactions(userId: String) async throws -> [Transaction] {
        return []
    }
    
    func performTransaction(userId: String, type: TransactionType, amount: Double) async throws -> Account {
        return try await fetchAccount(userId: userId)
    }
}
