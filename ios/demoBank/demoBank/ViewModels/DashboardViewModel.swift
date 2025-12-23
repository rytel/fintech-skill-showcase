import Foundation
import Combine
import OSLog

/// ViewModel dla widoku podsumowania kont (Dashboard).
final class DashboardViewModel: ObservableObject {
    @Published var account: Account?
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    private let apiService: APIServiceProtocol
    private var cancellables = Set<AnyCancellable>()
    
    init(apiService: APIServiceProtocol) {
        self.apiService = apiService
    }
    
    
    func fetchData() async {
        Logger.ui.info("Fetching dashboard data...")
        isLoading = true
        errorMessage = nil
        
        do {
            // W prawdziwej aplikacji użylibyśmy realnego userId z sesji
            // Używamy UUID konta, które zostało zaserwowane w bazie danych
            let accountID = "de305d54-75b4-431b-adb2-eb6b9e546014"
            let fetchedAccount = try await apiService.fetchAccount(userId: accountID)
            self.account = fetchedAccount
            Logger.ui.info("Successfully fetched account: \(accountID, privacy: .public)")
        } catch {
            Logger.ui.error("Failed to fetch dashboard data: \(error.localizedDescription)")
            self.errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
}
