//
//  Account.swift
//  demoBank
//

import Foundation

struct Account: Codable {
    let id: String
    let userId: String
    let balance: Double
    let createdAt: Date
    
    enum CodingKeys: String, CodingKey {
        case id
        case userId = "user_id"
        case balance
        case createdAt = "created_at"
    }
}