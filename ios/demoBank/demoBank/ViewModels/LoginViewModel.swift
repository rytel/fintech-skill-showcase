import Foundation
import Combine
import LocalAuthentication
import OSLog

/// ViewModel dla widoku logowania.
final class LoginViewModel: ObservableObject {
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
            Logger.auth.info("Biometrics skipped via launch argument")
            self.canUseBiometrics = false
        } else {
            checkBiometrics()
        }
    }
    
    func checkBiometrics() {
        let context = LAContext()
        var error: NSError?
        canUseBiometrics = context.canEvaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, error: &error)
        Logger.auth.debug("Biometrics availability checked: \(self.canUseBiometrics)")
    }
    
    func loginWithBiometrics() {
        let context = LAContext()
        let reason = "Log in to your bank account"
        
        Logger.auth.info("Initiating biometric login")
        context.evaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, localizedReason: reason) { [weak self] success, error in
            DispatchQueue.main.async {
                if success {
                    Logger.auth.info("Biometric login successful")
                    self?.isLoggedIn = true
                } else if let error = error {
                    Logger.auth.error("Biometric login failed: \(error.localizedDescription)")
                    self?.errorMessage = error.localizedDescription
                }
            }
        }
    }
    
    func login() {
        Logger.auth.info("Attempting standard login for user: \(self.username, privacy: .private)")
        isLoading = true
        errorMessage = nil
        
        let request = LoginRequest(username: username, password: password)
        
        authService.login(credentials: request)
            .receive(on: RunLoop.main)
            .sink { [weak self] completion in
                self?.isLoading = false
                if case .failure(let error) = completion {
                    Logger.auth.error("Login failed: \(error.localizedDescription)")
                    self?.errorMessage = error.localizedDescription
                }
            } receiveValue: { [weak self] response in
                Logger.auth.info("Login successful, saving token")
                KeychainHelper.shared.saveToken(response.token)
                self?.isLoggedIn = true
            }
            .store(in: &cancellables)
    }
}
