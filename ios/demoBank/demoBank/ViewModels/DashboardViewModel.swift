import Foundation
import Combine

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
    
    @MainActor
    func fetchData() async {
        isLoading = true
        errorMessage = nil
        
        do {
            // W prawdziwej aplikacji użylibyśmy realnego userId z sesji
            // Używamy UUID konta, które zostało zaserwowane w bazie danych
            let fetchedAccount = try await apiService.fetchAccount(userId: "de305d54-75b4-431b-adb2-eb6b9e546014")
            self.account = fetchedAccount
        } catch {
            self.errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
}
