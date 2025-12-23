import Foundation
import SwiftUI
import Combine

/// Koordynator zarządzający głównym widokiem z zakładkami (Dashboard, Payments, Profile).
@MainActor final class MainCoordinator: Coordinator, ObservableObject {
    var childCoordinators: [Coordinator] = []
    let factory: ViewModelFactory
    
    @Published var selectedTab: Tab = .dashboard
    @Published var showTransfer: Bool = false
    
    // Persistent ViewModels for the main flow
    let dashboardViewModel: DashboardViewModel
    let transactionListViewModel: TransactionListViewModel
    
    enum Tab {
        case dashboard
        case payments
        case profile
    }
    
    init(factory: ViewModelFactory = ViewModelFactory()) {
        self.factory = factory
        self.dashboardViewModel = factory.makeDashboardViewModel()
        self.transactionListViewModel = factory.makeTransactionListViewModel()
    }
    
    func start() {
        // Inicjalizacja pod-koordynatorów dla każdej zakładki może odbyć się tutaj
    }
}
