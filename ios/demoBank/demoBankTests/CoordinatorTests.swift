import Testing
@testable import demoBank

struct CoordinatorTests {
    
    @Test func coordinatorProtocolExists() {
        // Ten test ma na celu sprawdzenie, czy protokół Coordinator jest zdefiniowany
        // i czy MockCoordinator go implementuje.
        let mockCoordinator = MockCoordinator()
        #expect(mockCoordinator != nil)
    }
}

class MockCoordinator: Coordinator {
    var childCoordinators: [Coordinator] = []
    func start() {}
}