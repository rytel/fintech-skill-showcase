import Foundation

@MainActor
protocol APIServiceProtocol {
    func login(username: String, password: String) async throws -> LoginResponse
    func fetchAccount(userId: String) async throws -> Account
    func fetchTransactions(userId: String) async throws -> [Transaction]
    func performTransaction(userId: String, type: TransactionType, amount: Double) async throws -> Account
}
