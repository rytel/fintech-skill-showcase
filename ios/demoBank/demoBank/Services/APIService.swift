//
//  APIService.swift
//  demoBank
//

import Foundation

final class APIService {
    private let session: URLSession
    private let baseURL = URL(string: "http://localhost:8080")!
    
    init(session: URLSession = .shared) {
        self.session = session
    }
    
    // Generic request performer
    private func performRequest<T: Decodable>(
        path: String,
        method: String = "GET",
        body: Data? = nil,
        requiresAuth: Bool = true
    ) async throws -> T {
        guard let url = URL(string: path, relativeTo: baseURL) else {
            throw APIError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        if requiresAuth {
            if let token = KeychainHelper.shared.getToken() {
                request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
            }
        }
        
        request.httpBody = body
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.unknown
        }
        
        guard (200...299).contains(httpResponse.statusCode) else {
            if httpResponse.statusCode == 401 {
                throw APIError.unauthorized
            }
            throw APIError.serverError(httpResponse.statusCode)
        }
        
        do {
            let decoder = JSONDecoder()
            decoder.dateDecodingStrategy = .iso8601
            return try decoder.decode(T.self, from: data)
        } catch {
            throw APIError.decodingError(error.localizedDescription)
        }
    }
    
    // Auth
    func login(username: String, password: String) async throws -> LoginResponse {
        let request = LoginRequest(username: username, password: password)
        let body = try JSONEncoder().encode(request)
        return try await performRequest(path: "/api/login", method: "POST", body: body, requiresAuth: false)
    }
    
    // Account
    func fetchAccount(userId: String) async throws -> Account {
        return try await performRequest(path: "/api/account/\(userId)")
    }
    
    // Transactions
    func fetchTransactions(userId: String) async throws -> [Transaction] {
        return try await performRequest(path: "/api/account/\(userId)/transactions")
    }
}
