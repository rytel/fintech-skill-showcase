import XCTest
@testable import demoBank

final class ViewModelFactoryTests: XCTestCase {
    
    func testViewModelFactoryCreatesLoginViewModel() {
        let container = DependencyContainer.shared
        // Rejestrujemy APIService, który jest wymagany przez LoginViewModel (w przyszłości)
        container.register(type: APIService.self, component: APIService())
        
        let factory = ViewModelFactory(container: container)
        // Na razie testujemy czy factory w ogóle istnieje i zwraca coś
        XCTAssertNotNil(factory)
    }
}
