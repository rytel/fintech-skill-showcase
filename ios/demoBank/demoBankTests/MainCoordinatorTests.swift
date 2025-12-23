import Testing
@testable import demoBank


@MainActor
final class MainCoordinatorTests {
    
    @Test func mainCoordinatorInitializesWithChildCoordinators() {
        let coordinator = MainCoordinator(factory: ViewModelFactory())
        #expect(coordinator.childCoordinators.isEmpty)
    }
}