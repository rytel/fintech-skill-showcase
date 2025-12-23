import Foundation

/// Fabryka odpowiedzialna za tworzenie ViewModels z wstrzykniętymi zależnościami.

final class ViewModelFactory {
    private let container: DependencyContainer
    
    init(container: DependencyContainer = .shared) {
        self.container = container
    }
    
    func makeLoginViewModel() -> LoginViewModel {
        let authService = container.resolve(type: AuthServiceProtocol.self) ?? AuthService()
        return LoginViewModel(authService: authService)
    }
    
    func makeDashboardViewModel() -> DashboardViewModel {
        let apiService = container.resolve(type: APIServiceProtocol.self) ?? APIService()
        return DashboardViewModel(apiService: apiService)
    }
    
    func makeTransactionListViewModel() -> TransactionListViewModel {
        let apiService = container.resolve(type: APIServiceProtocol.self) ?? APIService()
        return TransactionListViewModel(apiService: apiService)
    }
    
    func makeTransferViewModel() -> TransferViewModel {
        let apiService = container.resolve(type: APIServiceProtocol.self) ?? APIService()
        return TransferViewModel(apiService: apiService)
    }
}
