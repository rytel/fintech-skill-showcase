import Foundation
import SwiftUI
import Combine

final class TransferCoordinator: Coordinator, ObservableObject {
    var childCoordinators: [Coordinator] = []
    
    @Published var rootView: AnyView = AnyView(EmptyView())
    
    private let factory: ViewModelFactory
    var onFinish: (() -> Void)?
    
    init(factory: ViewModelFactory) {
        self.factory = factory
        start()
    }
    
    func start() {
        showTransferForm()
    }
    
    func showTransferForm() {
        let viewModel = factory.makeTransferViewModel()
        rootView = AnyView(TransferFormView(viewModel: viewModel))
    }
}
