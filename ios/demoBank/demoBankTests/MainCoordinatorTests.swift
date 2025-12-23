import Testing
@testable import demoBank

@MainActor
struct MainCoordinatorTests {
    
    @Test func mainCoordinatorInitializesWithChildCoordinators() {
        let coordinator = MainCoordinator(factory: ViewModelFactory())
        #expect(coordinator.childCoordinators.isEmpty)
    }
}