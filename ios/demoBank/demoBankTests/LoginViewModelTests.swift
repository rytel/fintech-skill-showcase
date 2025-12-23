import Testing
import Combine
import Foundation
@testable import demoBank

@MainActor
struct LoginViewModelTests {
    
    @Test func loginSuccessUpdatesState() async throws {
        let mockAuthService = MockAuthService()
        let viewModel = LoginViewModel(authService: mockAuthService)
        
        #expect(!viewModel.isLoggedIn)
        
        viewModel.login()
        
        // Wait for change since login is async and updates on MainActor
        try await Task.sleep(nanoseconds: 200_000_000)
        
        #expect(viewModel.isLoggedIn)
    }
}

class MockAuthService: AuthServiceProtocol {
    func login(credentials: LoginRequest) -> AnyPublisher<LoginResponse, Error> {
        return Just(LoginResponse(token: "test-token"))
            .setFailureType(to: Error.self)
            .eraseToAnyPublisher()
    }
}