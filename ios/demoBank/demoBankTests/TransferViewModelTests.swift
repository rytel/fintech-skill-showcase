import Testing
import Foundation
@testable import demoBank


struct TransferViewModelTests {
    
    private var mockAPI: MockAPIService
    private var viewModel: TransferViewModel
    
    init() {
        mockAPI = MockAPIService()
        viewModel = TransferViewModel(apiService: mockAPI)
    }
    
    @Test func validationLogic() {
        // Initial state
        #expect(!viewModel.isRecipientValid)
        
        // Recipient validation
        viewModel.recipientName = "Jan Kowalski"
        viewModel.recipientAccount = "1234567890"
        #expect(viewModel.isRecipientValid)
        
        // Amount validation
        #expect(!viewModel.isAmountValid)
        viewModel.amount = "100.50"
        viewModel.title = "Test transfer"
        #expect(viewModel.isAmountValid)
    }
    
    @Test func stepTransitions() {
        // Recipient -> Amount
        viewModel.recipientName = "Jan"
        viewModel.recipientAccount = "1234567890"
        viewModel.nextStep()
        #expect(viewModel.currentStep == .amount)
        
        // Amount -> Confirmation
        viewModel.amount = "50"
        viewModel.title = "Title"
        viewModel.nextStep()
        #expect(viewModel.currentStep == .confirmation)
    }
    
    @Test func toastNotification() async throws {
        viewModel.recipientName = "Jan"
        viewModel.recipientAccount = "1234567890"
        viewModel.amount = "50"
        viewModel.title = "Title"
        
        // Go to confirmation
        viewModel.nextStep()
        viewModel.nextStep()
        
        viewModel.nextStep() // Triggers performTransfer
        
        // Wait for task to finish (it has a sleep)
        try await Task.sleep(nanoseconds: 2_100_000_000)
        
        #expect(viewModel.showToast)
        #expect(viewModel.toastType == .success)
        #expect(viewModel.currentStep == .success)
    }
}