import Foundation
import Combine

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
        isLoading = true
        errorMessage = nil
        
        do {
            let fetchedTransactions = try await apiService.fetchTransactions(userId: userId)
            self.transactions = fetchedTransactions.sorted(by: { $0.createdAt > $1.createdAt })
        } catch {
            self.errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
}
