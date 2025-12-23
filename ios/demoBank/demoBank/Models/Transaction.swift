//
//  Transaction.swift
//  demoBank
//

import Foundation

enum TransactionType: String, Codable {
    case deposit = "DEPOSIT"
    case withdrawal = "WITHDRAWAL"
    case transferIn = "TRANSFER_IN"
    case transferOut = "TRANSFER_OUT"
}

struct Transaction: Codable {
    let id: String
    let accountId: String
    let type: TransactionType
    let amount: Double
    let createdAt: Date
    let description: String?
    
    enum CodingKeys: String, CodingKey {
        case id
        case accountId = "account_id"
        case type
        case amount
        case createdAt = "created_at"
        case description
    }
}