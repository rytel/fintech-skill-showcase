//
//  APIError.swift
//  demoBank
//

import Foundation

enum APIError: Error, LocalizedError, Equatable {
    case invalidURL
    case networkError(String)
    case serverError(Int)
    case decodingError(String)
    case unauthorized
    case unknown
    
    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid URL."
        case .networkError(let message):
            return "Network error: \(message)"
        case .serverError(let code):
            return "Server returned error code: \(code)"
        case .decodingError(let message):
            return "Failed to decode response: \(message)"
        case .unauthorized:
            return "Unauthorized access. Please login again."
        case .unknown:
            return "An unknown error occurred."
        }
    }
}