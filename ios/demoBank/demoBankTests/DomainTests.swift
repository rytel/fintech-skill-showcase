//
//  DomainTests.swift
//  demoBankTests
//

import Testing
import Foundation
@testable import demoBank

@Suite(.serialized)
@MainActor
struct DomainTests {
    
    // MARK: - Keychain Tests
    
    @Test func keychainOperations() {
        let helper = KeychainHelper.shared
        let testToken = "test-token-\(UUID().uuidString)"
        
        helper.deleteToken()
        #expect(helper.getToken() == nil)
        
        helper.saveToken(testToken)
        #expect(helper.getToken() == testToken)
        
        helper.deleteToken()
        #expect(helper.getToken() == nil)
    }
    
    // MARK: - API Service Tests
    
    private let apiService: APIService
    
    init() {
        let configuration = URLSessionConfiguration.ephemeral
        configuration.protocolClasses = [MockURLProtocol.self]
        let session = URLSession(configuration: configuration)
        self.apiService = APIService(session: session)
    }
    
    @Test func loginSuccess() async throws {
        let mockToken = "fake-jwt-token"
        let responseData = "{\"token\": \"\(mockToken)\"}".data(using: .utf8)!
        
        MockURLProtocol.requestHandler = { request in
            let response = HTTPURLResponse(url: request.url!, statusCode: 200, httpVersion: nil, headerFields: nil)!
            return (response, responseData)
        }
        
        let response = try await apiService.login(username: "test", password: "123")
        #expect(response.token == mockToken)
    }
    
    @Test func fetchAccountSuccess() async throws {
        let responseData =
        """
        {
            "id": "de305d54-75b4-431b-adb2-eb6b9e546014",
            "user_id": "test_user",
            "balance": 1000.0,
            "created_at": "2023-01-01T12:00:00Z"
        }
        """
        .data(using: .utf8)!
        
        MockURLProtocol.requestHandler = { request in
            let authHeader = request.value(forHTTPHeaderField: "Authorization")
            #expect(authHeader == "Bearer some-token")
            let response = HTTPURLResponse(url: request.url!, statusCode: 200, httpVersion: nil, headerFields: nil)!
            return (response, responseData)
        }
        
        // Ważne: Zapisujemy token bezpośrednio przed testem w tym samym wątku (dzięki .serialized)
        KeychainHelper.shared.saveToken("some-token")
        
        let account = try await apiService.fetchAccount(userId: "test_user")
        #expect(account.balance == 1000.0)
        
        KeychainHelper.shared.deleteToken()
    }
    
    @Test func apiUnauthorizedError() async throws {
        MockURLProtocol.requestHandler = { request in
            let response = HTTPURLResponse(url: request.url!, statusCode: 401, httpVersion: nil, headerFields: nil)!
            return (response, Data())
        }
        
        await #expect(throws: APIError.unauthorized) {
            try await apiService.login(username: "test", password: "wrong")
        }
    }
}
