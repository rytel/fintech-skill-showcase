import Foundation
import Combine
import LocalAuthentication

/// ViewModel dla widoku logowania.
@MainActor final class LoginViewModel: ObservableObject {
    @Published var username = ""
    @Published var password = ""
    @Published var isLoading = false
    @Published var errorMessage: String?
    @Published var isLoggedIn = false
    @Published var canUseBiometrics = false
    
    private let authService: AuthServiceProtocol
    private var cancellables = Set<AnyCancellable>()
    
    init(authService: AuthServiceProtocol) {
        self.authService = authService
        
        if ProcessInfo.processInfo.arguments.contains("-skipBiometrics") {
            self.canUseBiometrics = false
        } else {
            checkBiometrics()
        }
    }
    
    func checkBiometrics() {
        let context = LAContext()
        var error: NSError?
        canUseBiometrics = context.canEvaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, error: &error)
    }
    
    func loginWithBiometrics() {
        let context = LAContext()
        let reason = "Log in to your bank account"
        
        context.evaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, localizedReason: reason) { [weak self] success, error in
            DispatchQueue.main.async {
                if success {
                    self?.isLoggedIn = true
                } else if let error = error {
                    self?.errorMessage = error.localizedDescription
                }
            }
        }
    }
    
    func login() {
        isLoading = true
        errorMessage = nil
        
        let request = LoginRequest(username: username, password: password)
        
        authService.login(credentials: request)
            .receive(on: RunLoop.main)
            .sink { [weak self] completion in
                self?.isLoading = false
                if case .failure(let error) = completion {
                    self?.errorMessage = error.localizedDescription
                }
            } receiveValue: { [weak self] response in
                KeychainHelper.shared.saveToken(response.token)
                self?.isLoggedIn = true
            }
            .store(in: &cancellables)
    }
}
