import Foundation
import Combine

@MainActor
protocol AuthServiceProtocol {
    func login(credentials: LoginRequest) -> AnyPublisher<LoginResponse, Error>
}

@MainActor
final class AuthService: AuthServiceProtocol {
    private let apiService: APIServiceProtocol
    
    init(apiService: APIServiceProtocol = APIService()) {
        self.apiService = apiService
    }
    
    func login(credentials: LoginRequest) -> AnyPublisher<LoginResponse, Error> {
        let apiService = self.apiService
        return Future { promise in
            Task {
                do {
                    let response = try await apiService.login(username: credentials.username, password: credentials.password)
                    promise(.success(response))
                } catch {
                    promise(.failure(error))
                }
            }
        }
        .eraseToAnyPublisher()
    }
}
