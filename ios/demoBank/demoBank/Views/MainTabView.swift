import SwiftUI

struct MainTabView: View {
    @ObservedObject var coordinator: MainCoordinator
    
    var body: some View {
        TabView(selection: $coordinator.selectedTab) {
            DashboardView(
                viewModel: coordinator.dashboardViewModel,
                transactionListViewModel: coordinator.transactionListViewModel,
                onTransferTap: { coordinator.showTransfer = true }
            )
                .tabItem {
                    Label("Główna", systemImage: "house.fill")
                }
                .tag(MainCoordinator.Tab.dashboard)
            
            Text("Płatności")
                .tabItem {
                    Label("Płatności", systemImage: "creditcard.fill")
                }
                .tag(MainCoordinator.Tab.payments)
            
            Text("Profil")
                .tabItem {
                    Label("Profil", systemImage: "person.fill")
                }
                .tag(MainCoordinator.Tab.profile)
        }
        .sheet(isPresented: $coordinator.showTransfer, onDismiss: {
            Task {
                await coordinator.dashboardViewModel.fetchData()
                if let accountId = coordinator.dashboardViewModel.account?.id {
                    await coordinator.transactionListViewModel.fetchTransactions(userId: accountId)
                }
            }
        }) {
            let transferCoordinator = TransferCoordinator(factory: coordinator.factory)
            transferCoordinator.rootView
        }
    }
}