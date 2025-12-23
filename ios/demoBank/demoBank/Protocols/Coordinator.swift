import Foundation
import SwiftUI

/// Podstawowy protokół dla wszystkich koordynatorów w aplikacji.
@MainActor
protocol Coordinator: AnyObject {
    /// Lista koordynatorów podrzędnych, zapobiegająca ich dealokacji.
    var childCoordinators: [Coordinator] { get set }
    
    /// Rozpoczyna przepływ koordynatora.
    func start()
}

extension Coordinator {
    /// Pomocnicza metoda do usuwania koordynatora podrzędnego.
    func childDidFinish(_ child: Coordinator) {
        for (index, coordinator) in childCoordinators.enumerated() {
            if coordinator === child {
                childCoordinators.remove(at: index)
                break
            }
        }
    }
}
