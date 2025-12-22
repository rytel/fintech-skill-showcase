import SwiftUI

struct AccountCard: View {
    let account: Account
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Text("Dostępne środki")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                Spacer()
                Image(systemName: "creditcard.fill")
                    .foregroundColor(.blue)
            }
            
            Text(formatBalance(account.balance))
                .font(.system(size: 34, weight: .bold, design: .rounded))
            
            HStack {
                Text("Konto osobiste")
                    .font(.caption)
                    .padding(.horizontal, 8)
                    .padding(.vertical, 4)
                    .background(Color.blue.opacity(0.1))
                    .cornerRadius(4)
                
                Spacer()
                
                Text("**** \(account.id.suffix(4))")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
        }
        .padding()
        .background(Color(.secondarySystemBackground))
        .cornerRadius(16)
        .shadow(color: Color.black.opacity(0.05), radius: 5, x: 0, y: 2)
    }
    
    private func formatBalance(_ amount: Double) -> String {
        let formatter = NumberFormatter()
        formatter.numberStyle = .currency
        formatter.currencyCode = "PLN" // W realnej aplikacji bralibyśmy z modelu account.currency
        return formatter.string(from: NSNumber(value: amount)) ?? "\(amount) PLN"
    }
}

#Preview {
    AccountCard(account: Account(id: "de305d54-75b4-431b-adb2-eb6b9e546014", userId: "1", balance: 12500.50, createdAt: Date()))
        .padding()
}