import Foundation
import SwiftUI
import Combine

/// Koordynator zarządzający procesem uwierzytelniania.
final class AuthCoordinator: Coordinator, ObservableObject {
    var childCoordinators: [Coordinator] = []
    
    /// Callback wywoływany po pomyślnym zalogowaniu.
    var onLoginSuccess: (() -> Void)?
    
    @Published var rootView: AnyView = AnyView(EmptyView())
    
    private let factory: ViewModelFactory
    
    init(factory: ViewModelFactory) {
        self.factory = factory
    }
    
    func start() {
        showLogin()
    }
    
    func showLogin() {
        let viewModel = factory.makeLoginViewModel()
        
        // Obserwujemy stan zalogowania w ViewModelu
        viewModel.$isLoggedIn
            .filter { $0 }
            .sink { [weak self] _ in
                self?.onLoginSuccess?()
            }
            .store(in: &cancellables)
        
        rootView = AnyView(LoginView(viewModel: viewModel))
    }
    
    private var cancellables = Set<AnyCancellable>()
}
