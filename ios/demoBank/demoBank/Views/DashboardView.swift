import SwiftUI

struct DashboardView: View {
    @ObservedObject var viewModel: DashboardViewModel
    @ObservedObject var transactionListViewModel: TransactionListViewModel
    var onTransferTap: (() -> Void)? = nil
    
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(alignment: .leading, spacing: 20) {
                    if viewModel.isLoading {
                        ProgressView("Pobieranie danych...")
                            .frame(maxWidth: .infinity, minHeight: 150)
                    } else if let account = viewModel.account {
                        AccountCard(account: account)
                            .transition(.opacity.combined(with: .move(edge: .top)))
                        
                        Text("Ostatnie transakcje")
                            .font(.headline)
                            .padding(.top, 10)
                        
                        TransactionListView(viewModel: transactionListViewModel, userId: account.id)
                    } else if let error = viewModel.errorMessage {
                        VStack(spacing: 12) {
                            Image(systemName: "exclamationmark.triangle.fill")
                                .font(.largeTitle)
                                .foregroundColor(.orange)
                            Text("Błąd: \(error)")
                                .font(.callout)
                                .multilineTextAlignment(.center)
                            Button("Spróbuj ponownie") {
                                Task {
                                    await viewModel.fetchData()
                                }
                            }
                            .buttonStyle(.bordered)
                        }
                        .frame(maxWidth: .infinity, minHeight: 150)
                    }
                    
                    Text("Szybkie akcje")
                        .font(.headline)
                        .padding(.top, 10)
                    
                    HStack(spacing: 15) {
                        Button(action: { onTransferTap?() }) {
                            QuickActionButton(title: "Przelew", icon: "arrow.up.right.circle.fill", color: .blue)
                                .accessibilityIdentifier("quick_action_transfer")
                        }
                        .buttonStyle(.plain)
                        
                        QuickActionButton(title: "BLIK", icon: "b.circle.fill", color: .pink)
                        QuickActionButton(title: "Historia", icon: "clock.fill", color: .gray)
                    }
                    
                    Spacer()
                }
                .padding()
            }
            .navigationTitle("Moje Finanse")
            .accessibilityIdentifier("dashboard_title_text")
            .accessibilityIdentifier("dashboard_nav_bar")
            .refreshable {
                await viewModel.fetchData()
            }
        }
        .onAppear {
            Task {
                await viewModel.fetchData()
            }
        }
    }
}

struct QuickActionButton: View {
    let title: String
    let icon: String
    let color: Color
    
    var body: some View {
        VStack {
            Image(systemName: icon)
                .font(.system(size: 30))
                .foregroundColor(color)
            Text(title)
                .font(.caption)
                .fontWeight(.medium)
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, 12)
        .background(color.opacity(0.1))
        .cornerRadius(12)
    }
}

#Preview {
    // Mock API for preview
    class PreviewAPI: APIServiceProtocol {
        func login(username: String, password: String) async throws -> LoginResponse { return LoginResponse(token: "") }
        func fetchAccount(userId: String) async throws -> Account { 
            return Account(id: "1", userId: "1", balance: 5432.10, createdAt: Date())
        }
        func fetchTransactions(userId: String) async throws -> [Transaction] { return [] }
        func performTransaction(userId: String, type: TransactionType, amount: Double) async throws -> Account {
            return try await fetchAccount(userId: userId)
        }
    }
    
    return DashboardView(
        viewModel: DashboardViewModel(apiService: PreviewAPI()),
        transactionListViewModel: TransactionListViewModel(apiService: PreviewAPI())
    )
}