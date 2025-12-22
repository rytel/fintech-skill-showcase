import SwiftUI

struct TransactionDetailView: View {
    let transaction: Transaction
    @Environment(\.dismiss) private var dismiss
    
    var body: some View {
        ScrollView {
            VStack(spacing: 24) {
                // Header with icon and amount
                VStack(spacing: 16) {
                    ZStack {
                        Circle()
                            .fill(transactionColor.opacity(0.1))
                            .frame(width: 80, height: 80)
                        
                        Image(systemName: transactionIcon)
                            .font(.system(size: 32, weight: .bold))
                            .foregroundColor(transactionColor)
                    }
                    
                    VStack(spacing: 4) {
                        Text(transactionTitle)
                            .font(.headline)
                            .foregroundColor(.secondary)
                        
                        Text(formatAmount(transaction.amount))
                            .font(.system(size: 36, weight: .bold, design: .rounded))
                    }
                }
                .padding(.top, 20)
                
                // Detail List
                VStack(spacing: 0) {
                    DetailRow(title: "Data i godzina", value: formatFullDate(transaction.createdAt))
                    Divider().padding(.leading)
                    DetailRow(title: "Typ operacji", value: transactionTitle)
                    Divider().padding(.leading)
                    DetailRow(title: "Numer konta", value: transaction.accountId)
                    
                    if let description = transaction.description {
                        Divider().padding(.leading)
                        DetailRow(title: "Opis", value: description)
                    }
                    
                    Divider().padding(.leading)
                    DetailRow(title: "Status", value: "Zakończona", valueColor: .green)
                    Divider().padding(.leading)
                    DetailRow(title: "ID Transakcji", value: transaction.id)
                }
                .background(Color(.secondarySystemBackground))
                .cornerRadius(12)
                
                Button(action: {
                    // Potencjalna akcja, np. ponów przelew
                }) {
                    Text("Ponów przelew")
                        .font(.headline)
                        .foregroundColor(.white)
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.blue)
                        .cornerRadius(12)
                }
                .padding(.top, 10)
                
                Spacer()
            }
            .padding()
        }
        .navigationTitle("Szczegóły transakcji")
        .navigationBarTitleDisplayMode(.inline)
    }
    
    private var transactionIcon: String {
        switch transaction.type {
        case .deposit: return "arrow.down.left"
        case .withdrawal: return "arrow.up.right"
        case .transferIn: return "arrow.right.circle"
        case .transferOut: return "arrow.left.circle"
        }
    }
    
    private var transactionColor: Color {
        switch transaction.type {
        case .deposit, .transferIn: return .green
        case .withdrawal, .transferOut: return .red
        }
    }
    
    private var transactionTitle: String {
        switch transaction.type {
        case .deposit: return "Wpłata"
        case .withdrawal: return "Wypłata"
        case .transferIn: return "Przelew przychodzący"
        case .transferOut: return "Przelew wychodzący"
        }
    }
    
    private func formatAmount(_ amount: Double) -> String {
        let formatter = NumberFormatter()
        formatter.numberStyle = .currency
        formatter.currencyCode = "PLN"
        return formatter.string(from: NSNumber(value: amount)) ?? "\(amount) PLN"
    }
    
    private func formatFullDate(_ date: Date) -> String {
        let formatter = DateFormatter()
        formatter.dateStyle = .long
        formatter.timeStyle = .medium
        return formatter.string(from: date)
    }
}

struct DetailRow: View {
    let title: String
    let value: String
    var valueColor: Color = .primary
    
    var body: some View {
        HStack {
            Text(title)
                .foregroundColor(.secondary)
            Spacer()
            Text(value)
                .fontWeight(.medium)
                .foregroundColor(valueColor)
                .multilineTextAlignment(.trailing)
        }
        .padding()
    }
}