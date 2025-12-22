import Foundation
import SwiftUI
import Combine

/// Główny koordynator aplikacji zarządzający przełączaniem między flow logowania a głównym widokiem.
final class AppCoordinator: Coordinator, ObservableObject {
    var childCoordinators: [Coordinator] = []
    
    @Published var rootView: AnyView = AnyView(EmptyView())
    
    private let factory: ViewModelFactory
    
    init(factory: ViewModelFactory = ViewModelFactory()) {
        self.factory = factory
    }
    
    func start() {
        showLogin()
    }
    
    func showLogin() {
        let authCoordinator = AuthCoordinator(factory: factory)
        authCoordinator.onLoginSuccess = { [weak self] in
            self?.showMain()
            self?.childDidFinish(authCoordinator)
        }
        childCoordinators.append(authCoordinator)
        authCoordinator.start()
        rootView = AnyView(authCoordinator.rootView)
    }
    
    func showMain() {
        let mainCoordinator = MainCoordinator(factory: factory)
        childCoordinators.append(mainCoordinator)
        mainCoordinator.start()
        rootView = AnyView(MainTabView(coordinator: mainCoordinator))
    }
}
