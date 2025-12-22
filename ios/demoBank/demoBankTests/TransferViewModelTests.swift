import XCTest
@testable import demoBank

final class TransferViewModelTests: XCTestCase {
    
    private var mockAPI: MockAPIService!
    private var viewModel: TransferViewModel!
    
    override func setUp() {
        super.setUp()
        mockAPI = MockAPIService()
        viewModel = TransferViewModel(apiService: mockAPI)
    }
    
    override func tearDown() {
        viewModel = nil
        mockAPI = nil
        super.tearDown()
    }
    
    func testValidationLogic() {
        // Initial state
        XCTAssertFalse(viewModel.isRecipientValid)
        
        // Recipient validation
        viewModel.recipientName = "Jan Kowalski"
        viewModel.recipientAccount = "1234567890"
        XCTAssertTrue(viewModel.isRecipientValid)
        
        // Amount validation
        XCTAssertFalse(viewModel.isAmountValid)
        viewModel.amount = "100.50"
        viewModel.title = "Test transfer"
        XCTAssertTrue(viewModel.isAmountValid)
    }
    
    func testStepTransitions() {
        // Recipient -> Amount
        viewModel.recipientName = "Jan"
        viewModel.recipientAccount = "1234567890"
        viewModel.nextStep()
        XCTAssertEqual(viewModel.currentStep, .amount)
        
        // Amount -> Confirmation
        viewModel.amount = "50"
        viewModel.title = "Title"
        viewModel.nextStep()
        XCTAssertEqual(viewModel.currentStep, .confirmation)
    }
    
    @MainActor
    func testToastNotification() async {
        viewModel.recipientName = "Jan"
        viewModel.recipientAccount = "1234567890"
        viewModel.amount = "50"
        viewModel.title = "Title"
        
        // Go to confirmation
        viewModel.nextStep()
        viewModel.nextStep()
        
        viewModel.nextStep() // Triggers performTransfer
        
        // Wait for task to finish (it has a sleep)
        try? await Task.sleep(nanoseconds: 2_000_000_000)
        
        XCTAssertTrue(viewModel.showToast)
        XCTAssertEqual(viewModel.toastType, .success)
        XCTAssertEqual(viewModel.currentStep, .success)
    }
}
