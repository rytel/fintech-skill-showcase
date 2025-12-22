import Foundation
import Combine

protocol AuthServiceProtocol {
    func login(credentials: LoginRequest) -> AnyPublisher<LoginResponse, Error>
}

final class AuthService: AuthServiceProtocol {
    private let apiService: APIServiceProtocol
    
    init(apiService: APIServiceProtocol = APIService()) {
        self.apiService = apiService
    }
    
    func login(credentials: LoginRequest) -> AnyPublisher<LoginResponse, Error> {
        return Future { [weak self] promise in
            Task {
                do {
                    guard let response = try await self?.apiService.login(username: credentials.username, password: credentials.password) else {
                        promise(.failure(APIError.unknown))
                        return
                    }
                    promise(.success(response))
                } catch {
                    promise(.failure(error))
                }
            }
        }
        .eraseToAnyPublisher()
    }
}
