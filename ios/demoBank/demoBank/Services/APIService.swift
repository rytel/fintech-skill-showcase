//
//  APIService.swift
//  demoBank
//

import Foundation

final class APIService: APIServiceProtocol {
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
        
        // Debug: Print raw response
        if let jsonString = String(data: data, encoding: .utf8) {
            print("API Response for \(path): \(jsonString)")
        }
        
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
            let formatter = DateFormatter()
            formatter.calendar = Calendar(identifier: .iso8601)
            formatter.locale = Locale(identifier: "en_US_POSIX")
            formatter.timeZone = TimeZone(secondsFromGMT: 0)
            
            decoder.dateDecodingStrategy = .custom({ decoder in
                let container = try decoder.singleValueContainer()
                let dateString = try container.decode(String.self)
                
                let formats = [
                    "yyyy-MM-dd'T'HH:mm:ss.SSSSSSZ",
                    "yyyy-MM-dd'T'HH:mm:ss.SSSZ",
                    "yyyy-MM-dd'T'HH:mm:ssZ",
                    "yyyy-MM-dd'T'HH:mm:ss.SSSSSS'Z'",
                    "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'",
                    "yyyy-MM-dd'T'HH:mm:ss'Z'"
                ]
                
                for format in formats {
                    formatter.dateFormat = format
                    if let date = formatter.date(from: dateString) {
                        return date
                    }
                }
                
                throw DecodingError.dataCorruptedError(in: container, debugDescription: "Cannot decode date string \(dateString)")
            })
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
    
    func performTransaction(userId: String, type: TransactionType, amount: Double) async throws -> Account {
        let requestBody: [String: Any] = [
            "user_id": userId,
            "type": type.rawValue,
            "amount": amount
        ]
        let body = try JSONSerialization.data(withJSONObject: requestBody)
        return try await performRequest(path: "/api/transactions", method: "POST", body: body)
    }
}
