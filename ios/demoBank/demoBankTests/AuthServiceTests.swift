import XCTest
import Combine
@testable import demoBank

final class AuthServiceTests: XCTestCase {
    var cancellables: Set<AnyCancellable> = []
    
    func testLoginReturnsSuccess() {
        let mockAPI = MockAPIServiceForAuth()
        let service = AuthService(apiService: mockAPI)
        let expectation = XCTestExpectation(description: "Login response")
        
        let request = LoginRequest(username: "test", password: "password")
        
        service.login(credentials: request)
            .sink { completion in
                if case .failure = completion {
                    XCTFail("Login should not fail")
                }
            } receiveValue: { response in
                XCTAssertEqual(response.token, "mock-token")
                expectation.fulfill()
            }
            .store(in: &cancellables)
        
        wait(for: [expectation], timeout: 2.0)
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