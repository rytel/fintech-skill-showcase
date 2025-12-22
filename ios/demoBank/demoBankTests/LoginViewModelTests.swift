import XCTest
import Combine
@testable import demoBank

final class LoginViewModelTests: XCTestCase {
    var cancellables: Set<AnyCancellable> = []
    
    func testLoginSuccessUpdatesState() {
        let mockAuthService = MockAuthService()
        let viewModel = LoginViewModel(authService: mockAuthService)
        
        let expectation = XCTestExpectation(description: "Login success")
        
        viewModel.$isLoggedIn
            .dropFirst()
            .sink { isLoggedIn in
                if isLoggedIn {
                    expectation.fulfill()
                }
            }
            .store(in: &cancellables)
        
        viewModel.login()
        
        wait(for: [expectation], timeout: 2.0)
        XCTAssertTrue(viewModel.isLoggedIn)
    }
}

class MockAuthService: AuthServiceProtocol {
    func login(credentials: LoginRequest) -> AnyPublisher<LoginResponse, Error> {
        return Just(LoginResponse(token: "test-token"))
            .setFailureType(to: Error.self)
            .eraseToAnyPublisher()
    }
}
