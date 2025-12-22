import XCTest
@testable import demoBank

final class TransactionListViewModelTests: XCTestCase {
    
    @MainActor
    func testFetchTransactionsUpdatesState() async {
        let mockAPI = MockAPIServiceForTransactions()
        let viewModel = TransactionListViewModel(apiService: mockAPI)
        
        await viewModel.fetchTransactions(userId: "de305d54-75b4-431b-adb2-eb6b9e546014")
        
        XCTAssertEqual(viewModel.transactions.count, 2)
        XCTAssertEqual(viewModel.transactions.first?.amount, 100.0)
        XCTAssertFalse(viewModel.isLoading)
    }

    @MainActor
    func testFetchTransactionsHandlesError() async {
        let mockAPI = MockAPIServiceForTransactions()
        mockAPI.shouldReturnError = true
        let viewModel = TransactionListViewModel(apiService: mockAPI)
        
        await viewModel.fetchTransactions(userId: "123")
        
        XCTAssertTrue(viewModel.transactions.isEmpty)
        XCTAssertNotNil(viewModel.errorMessage)
        XCTAssertFalse(viewModel.isLoading)
    }
}

class MockAPIServiceForTransactions: APIServiceProtocol {
    var shouldReturnError = false

    func login(username: String, password: String) async throws -> LoginResponse { return LoginResponse(token: "") }
    func fetchAccount(userId: String) async throws -> Account { 
        return Account(id: "1", userId: "1", balance: 0, createdAt: Date())
    }
    
    func fetchTransactions(userId: String) async throws -> [Transaction] {
        if shouldReturnError {
            throw NSError(domain: "test", code: 1, userInfo: [NSLocalizedDescriptionKey: "API Error"])
        }
        return [
            Transaction(id: "1", accountId: "1", type: .deposit, amount: 100.0, createdAt: Date(), description: nil),
            Transaction(id: "2", accountId: "1", type: .withdrawal, amount: -50.0, createdAt: Date().addingTimeInterval(-60), description: nil)
        ]
    }
    
    func performTransaction(userId: String, type: TransactionType, amount: Double) async throws -> Account {
        return try await fetchAccount(userId: userId)
    }
}