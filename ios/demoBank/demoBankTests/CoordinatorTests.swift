import Testing
@testable import demoBank


@MainActor
final class CoordinatorTests {
    
    @Test func coordinatorProtocolExists() {
        // Ten test ma na celu sprawdzenie, czy protokół Coordinator jest zdefiniowany
        // i czy MockCoordinator go implementuje.
        let mockCoordinator = MockCoordinator()
        #expect(true)
    }
}


@MainActor
class MockCoordinator: Coordinator {
    var childCoordinators: [Coordinator] = []
    func start() {}
}