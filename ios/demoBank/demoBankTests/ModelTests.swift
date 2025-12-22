//
//  ModelTests.swift
//  demoBankTests
//

import Testing
import Foundation
@testable import demoBank

struct ModelTests {
    
    @Test func loginRequestEncoding() throws {
        let request = LoginRequest(username: "test_user", password: "password123")
        let encoder = JSONEncoder()
        let data = try encoder.encode(request)
        
        let json = try JSONSerialization.jsonObject(with: data) as? [String: Any]
        #expect(json?["username"] as? String == "test_user")
        #expect(json?["password"] as? String == "password123")
    }
    
    @Test func loginResponseDecoding() throws {
        let jsonString = """
        {
            "token": "eyJh... (JWT Token)"
        }
        """
        let data = jsonString.data(using: .utf8)!
        let decoder = JSONDecoder()
        let response = try decoder.decode(LoginResponse.self, from: data)
        
        #expect(response.token == "eyJh... (JWT Token)")
    }
    
    @Test func accountDecoding() throws {
        let jsonString = """
        {
            "id": 1,
            "user_id": "test_user",
            "balance": 1000,
            "created_at": "2023-01-01T12:00:00Z"
        }
        """
        let data = jsonString.data(using: .utf8)!
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        let account = try decoder.decode(Account.self, from: data)
        
        #expect(account.id == 1)
        #expect(account.userId == "test_user")
        #expect(account.balance == 1000.0)
    }
    
    @Test func transactionDecoding() throws {
        let jsonString = """
        [
            {
                "id": 1,
                "account_id": 1,
                "type": "DEPOSIT",
                "amount": 100,
                "created_at": "2023-01-01T12:00:00Z"
            }
        ]
        """
        let data = jsonString.data(using: .utf8)!
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        let transactions = try decoder.decode([Transaction].self, from: data)
        
        #expect(transactions.count == 1)
        #expect(transactions[0].id == 1)
        #expect(transactions[0].accountId == 1)
        #expect(transactions[0].type == .deposit)
        #expect(transactions[0].amount == 100.0)
    }
}