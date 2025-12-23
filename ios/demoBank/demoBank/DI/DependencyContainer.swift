import Foundation

/// Prosty kontener Dependency Injection zarządzający czasem życia serwisów.

final class DependencyContainer {
    static let shared = DependencyContainer()
    
    private var services: [String: Any] = [:]
    
    private init() {}
    
    /// Rejestruje instancję dla danego typu.
    func register<T>(type: T.Type, component: Any) {
        let key = String(describing: type)
        services[key] = component
    }
    
    /// Zwraca zarejestrowaną instancję dla danego typu.
    func resolve<T>(type: T.Type) -> T? {
        let key = String(describing: type)
        return services[key] as? T
    }

    func setupMocks() {
        if ProcessInfo.processInfo.arguments.contains("-useMockData") {
            register(type: AuthServiceProtocol.self, component: MockAuthService())
            register(type: APIServiceProtocol.self, component: MockAPIService())
        }
    }
}
