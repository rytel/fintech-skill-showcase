import Foundation

struct TransferRequest: Codable, Sendable {
    let recipientName: String
    let recipientAccount: String
    let amount: Double
    let title: String
    
    enum CodingKeys: String, CodingKey {
        case recipientName = "recipient_name"
        case recipientAccount = "recipient_account"
        case amount
        case title
    }
}
