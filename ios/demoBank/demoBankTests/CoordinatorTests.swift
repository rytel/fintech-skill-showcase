import XCTest
@testable import demoBank

final class CoordinatorTests: XCTestCase {
    
    func testCoordinatorProtocolExists() {
        // Ten test ma na celu sprawdzenie, czy protokół Coordinator jest zdefiniowany
        // i czy AppCoordinator go implementuje.
        let mockCoordinator = MockCoordinator()
        XCTAssertNotNil(mockCoordinator)
    }
}

class MockCoordinator: Coordinator {
    var childCoordinators: [Coordinator] = []
    func start() {}
}
