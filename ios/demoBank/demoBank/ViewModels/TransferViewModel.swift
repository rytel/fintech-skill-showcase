import Foundation
import Combine
import OSLog

final class TransferViewModel: ObservableObject {
    enum Step {
        case recipient
        case amount
        case confirmation
        case success
    }
    
    @Published var currentStep: Step = .recipient {
        didSet {
            Logger.ui.debug("Transfer step changed to: \(String(describing: self.currentStep))")
        }
    }
    @Published var recipientName: String = ""
    @Published var recipientAccount: String = ""
    @Published var amount: String = ""
    @Published var title: String = ""
    
    @Published var isLoading: Bool = false
    @Published var errorMessage: String?
    
    // Toast Management
    @Published var showToast: Bool = false
    @Published var toastType: ToastType = .info
    @Published var toastMessage: String = ""
    
    private let apiService: APIServiceProtocol
    private var cancellables = Set<AnyCancellable>()
    
    init(apiService: APIServiceProtocol) {
        self.apiService = apiService
    }
    
    var isRecipientValid: Bool {
        !recipientName.isEmpty && recipientAccount.count >= 10
    }
    
    var isAmountValid: Bool {
        guard let amountDouble = Double(amount.replacingOccurrences(of: ",", with: ".")) else { return false }
        return amountDouble > 0 && !title.isEmpty
    }
    
    func nextStep() {
        switch currentStep {
        case .recipient:
            if isRecipientValid { currentStep = .amount }
        case .amount:
            if isAmountValid { currentStep = .confirmation }
        case .confirmation:
            performTransfer()
        default:
            break
        }
    }
    
    func previousStep() {
        switch currentStep {
        case .amount: currentStep = .recipient
        case .confirmation: currentStep = .amount
        default: break
        }
    }
    
    
    private func performTransfer() {
        Logger.ui.info("Initiating transfer of \(self.amount) to \(self.recipientName, privacy: .private)")
        isLoading = true
        errorMessage = nil
        
        guard let amountDouble = Double(amount.replacingOccurrences(of: ",", with: ".")) else {
            self.errorMessage = "Nieprawidłowa kwota"
            self.isLoading = false
            Logger.ui.warning("Transfer failed: Invalid amount entered")
            return
        }
        
        Task {
            do {
                // W tej demonstracji używamy tego samego stałego ID konta
                _ = try await apiService.performTransaction(
                    userId: "de305d54-75b4-431b-adb2-eb6b9e546014",
                    type: .withdrawal,
                    amount: amountDouble
                )
                
                Logger.ui.info("Transfer completed successfully")
                self.isLoading = false
                self.toastType = .success
                self.toastMessage = "Przelew do \(self.recipientName) został wysłany!"
                self.showToast = true
                self.currentStep = .success
                
                // Auto-hide toast after 3 seconds
                Task {
                    try? await Task.sleep(nanoseconds: 3_000_000_000)
                    self.showToast = false
                }
            } catch {
                Logger.ui.error("Transfer failed: \(error.localizedDescription)")
                self.isLoading = false
                self.errorMessage = error.localizedDescription
                self.toastType = .error
                self.toastMessage = "Błąd: \(error.localizedDescription)"
                self.showToast = true
            }
        }
    }
}
