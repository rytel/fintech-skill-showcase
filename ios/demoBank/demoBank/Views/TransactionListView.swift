import SwiftUI

struct TransactionRow: View {
    let transaction: Transaction
    
    var body: some View {
        HStack(spacing: 16) {
            ZStack {
                Circle()
                    .fill(transactionColor.opacity(0.1))
                    .frame(width: 44, height: 44)
                
                Image(systemName: transactionIcon)
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundColor(transactionColor)
            }
            
            VStack(alignment: .leading, spacing: 4) {
                Text(transactionTitle)
                    .font(.system(.subheadline, design: .rounded))
                    .fontWeight(.semibold)
                
                Text(formatDate(transaction.createdAt))
                    .font(.caption2)
                    .foregroundColor(.secondary)
            }
            
            Spacer()
            
            Text(formatAmount(transaction.amount))
                .font(.system(.subheadline, design: .rounded))
                .fontWeight(.bold)
                .foregroundColor(transaction.amount < 0 ? .primary : .green)
        }
        .padding(.vertical, 8)
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
    
    private func formatDate(_ date: Date) -> String {
        let formatter = DateFormatter()
        formatter.dateStyle = .medium
        formatter.timeStyle = .none
        return formatter.string(from: date)
    }
    
    private func formatAmount(_ amount: Double) -> String {
        let formatter = NumberFormatter()
        formatter.numberStyle = .currency
        formatter.currencyCode = "PLN"
        return formatter.string(from: NSNumber(value: amount)) ?? "\(amount) PLN"
    }
}

struct TransactionSkeletonRow: View {
    @State private var opacity: Double = 0.3
    
    var body: some View {
        HStack(spacing: 16) {
            Circle()
                .fill(Color.gray.opacity(0.2))
                .frame(width: 44, height: 44)
            
            VStack(alignment: .leading, spacing: 8) {
                RoundedRectangle(cornerRadius: 4)
                    .fill(Color.gray.opacity(0.2))
                    .frame(width: 120, height: 12)
                
                RoundedRectangle(cornerRadius: 4)
                    .fill(Color.gray.opacity(0.2))
                    .frame(width: 80, height: 8)
            }
            
            Spacer()
            
            RoundedRectangle(cornerRadius: 4)
                .fill(Color.gray.opacity(0.2))
                .frame(width: 60, height: 12)
        }
        .padding(.vertical, 8)
        .opacity(opacity)
        .onAppear {
            withAnimation(.easeInOut(duration: 1.0).repeatForever(autoreverses: true)) {
                opacity = 0.7
            }
        }
    }
}

struct TransactionListView: View {
    @ObservedObject var viewModel: TransactionListViewModel
    let userId: String
    
    var body: some View {
        VStack(spacing: 0) {
            if viewModel.isLoading && viewModel.transactions.isEmpty {
                ForEach(0..<5) { _ in
                    TransactionSkeletonRow()
                    Divider().padding(.leading, 60)
                }
            } else if let error = viewModel.errorMessage {
                VStack(spacing: 8) {
                    Image(systemName: "exclamationmark.triangle")
                        .foregroundColor(.orange)
                    Text(error)
                        .font(.caption)
                        .foregroundColor(.secondary)
                        .multilineTextAlignment(.center)
                    Button("Ponów") {
                        Task {
                            await viewModel.fetchTransactions(userId: userId)
                        }
                    }
                    .font(.caption)
                    .buttonStyle(.bordered)
                }
                .padding()
                .frame(maxWidth: .infinity)
            } else if viewModel.transactions.isEmpty {
                VStack(spacing: 12) {
                    Image(systemName: "tray")
                        .font(.largeTitle)
                        .foregroundColor(.secondary)
                    Text("Brak transakcji")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                }
                .frame(maxWidth: .infinity, minHeight: 150)
                .padding()
            } else {
                ForEach(viewModel.transactions, id: \.id) { transaction in
                    NavigationLink(destination: TransactionDetailView(transaction: transaction)) {
                        TransactionRow(transaction: transaction)
                            .contentShape(Rectangle())
                    }
                    .buttonStyle(.plain)
                    
                    if transaction.id != viewModel.transactions.last?.id {
                        Divider().padding(.leading, 60)
                    }
                }
            }
        }
        .onAppear {
            Task {
                await viewModel.fetchTransactions(userId: userId)
            }
        }
    }
}