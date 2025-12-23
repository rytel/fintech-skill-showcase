import Foundation
import Combine
import OSLog

/// ViewModel dla listy transakcji.
final class TransactionListViewModel: ObservableObject {
    @Published var transactions: [Transaction] = []
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    private let apiService: APIServiceProtocol
    private var cancellables = Set<AnyCancellable>()
    
    init(apiService: APIServiceProtocol) {
        self.apiService = apiService
    }
    
    
    func fetchTransactions(userId: String) async {
        Logger.ui.info("Fetching transactions for user: \(userId, privacy: .public)")
        isLoading = true
        errorMessage = nil
        
        do {
            let fetchedTransactions = try await apiService.fetchTransactions(userId: userId)
            self.transactions = fetchedTransactions.sorted(by: { $0.createdAt > $1.createdAt })
            Logger.ui.info("Successfully fetched \(self.transactions.count) transactions")
        } catch {
            Logger.ui.error("Failed to fetch transactions: \(error.localizedDescription)")
            self.errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
}
