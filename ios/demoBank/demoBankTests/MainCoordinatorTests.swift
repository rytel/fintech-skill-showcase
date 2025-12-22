import XCTest
@testable import demoBank

final class MainCoordinatorTests: XCTestCase {
    
    func testMainCoordinatorInitializesWithChildCoordinators() {
        let coordinator = MainCoordinator(factory: ViewModelFactory())
        XCTAssertTrue(coordinator.childCoordinators.isEmpty)
    }
}
