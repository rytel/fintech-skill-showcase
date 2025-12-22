//
//  Transaction.swift
//  demoBank
//

import Foundation

enum TransactionType: String, Codable {
    case deposit = "DEPOSIT"
    case withdraw = "WITHDRAW"
}

struct Transaction: Codable {
    let id: Int
    let accountId: Int
    let type: TransactionType
    let amount: Double
    let createdAt: Date
    
    enum CodingKeys: String, CodingKey {
        case id
        case accountId = "account_id"
        case type
        case amount
        case createdAt = "created_at"
    }
}
