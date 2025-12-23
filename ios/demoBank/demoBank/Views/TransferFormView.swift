import SwiftUI

struct TransferFormView: View {
    @ObservedObject var viewModel: TransferViewModel
    @Environment(\.dismiss) private var dismiss
    
    var body: some View {
        NavigationView {
            VStack {
                // Step Indicator
                HStack(spacing: 8) {
                    StepCircle(step: 1, currentStep: stepNumber, title: "Odbiorca")
                    Capsule().fill(stepNumber > 1 ? Color.blue : Color.gray.opacity(0.3)).frame(height: 2)
                    StepCircle(step: 2, currentStep: stepNumber, title: "Kwota")
                    Capsule().fill(stepNumber > 2 ? Color.blue : Color.gray.opacity(0.3)).frame(height: 2)
                    StepCircle(step: 3, currentStep: stepNumber, title: "Podsumowanie")
                }
                .padding(.horizontal)
                .padding(.top, 20)
                
                content
                    .padding(.top, 30)
                
                Spacer()
                
                if viewModel.currentStep != .success {
                    HStack {
                        if viewModel.currentStep != .recipient {
                            Button("Wstecz") {
                                withAnimation {
                                    viewModel.previousStep()
                                }
                            }
                            .buttonStyle(.bordered)
                            .accessibilityIdentifier("transfer_back_button")
                        }
                        
                        Button(action: {
                            withAnimation {
                                viewModel.nextStep()
                            }
                        }) {
                            if viewModel.isLoading {
                                ProgressView()
                                    .progressViewStyle(CircularProgressViewStyle(tint: .white))
                            } else {
                                Text(viewModel.currentStep == .confirmation ? "Potwierdzam i wysyłam" : "Dalej")
                            }
                        }
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(isNextDisabled ? Color.gray : Color.blue)
                        .foregroundColor(.white)
                        .cornerRadius(12)
                        .disabled(isNextDisabled || viewModel.isLoading)
                        .accessibilityIdentifier("transfer_next_button")
                    }
                    .padding()
                }
            }
            .navigationTitle("Nowy przelew")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    if viewModel.currentStep == .success {
                        Button("Zamknij") { dismiss() }
                            .accessibilityIdentifier("transfer_close_button")
                    }
                }
            }
            .overlay(alignment: .top) {
                if viewModel.showToast {
                    ToastView(type: viewModel.toastType, message: viewModel.toastMessage)
                        .transition(.move(edge: .top).combined(with: .opacity))
                        .padding(.top, 10)
                        .zIndex(1)
                        .accessibilityIdentifier("toast_view")
                }
            }
            .animation(.spring(), value: viewModel.showToast)
        }
    }
    
    @ViewBuilder
    private var content: some View {
        switch viewModel.currentStep {
        case .recipient:
            recipientStep
        case .amount:
            amountStep
        case .confirmation:
            confirmationStep
        case .success:
            successStep
        }
    }
    
    private var recipientStep: some View {
        VStack(alignment: .leading, spacing: 20) {
            VStack(alignment: .leading, spacing: 8) {
                Text("Nazwa odbiorcy")
                    .font(.caption)
                    .foregroundColor(.secondary)
                TextField("Imię i nazwisko lub nazwa firmy", text: $viewModel.recipientName)
                    .padding()
                    .background(Color(.secondarySystemBackground))
                    .cornerRadius(10)
                    .accessibilityIdentifier("recipient_name_field")
            }
            
            VStack(alignment: .leading, spacing: 8) {
                Text("Numer konta")
                    .font(.caption)
                    .foregroundColor(.secondary)
                TextField("PL00 0000 0000 0000 0000 0000 0000", text: $viewModel.recipientAccount)
                    .keyboardType(.numberPad)
                    .padding()
                    .background(Color(.secondarySystemBackground))
                    .cornerRadius(10)
                    .accessibilityIdentifier("recipient_account_field")
            }
        }
        .padding(.horizontal)
    }
    
    private var amountStep: some View {
        VStack(alignment: .leading, spacing: 20) {
            VStack(alignment: .leading, spacing: 8) {
                Text("Kwota")
                    .font(.caption)
                    .foregroundColor(.secondary)
                HStack {
                    TextField("0,00", text: $viewModel.amount)
                        .font(.system(size: 24, weight: .bold, design: .rounded))
                        .keyboardType(.decimalPad)
                        .accessibilityIdentifier("transfer_amount_field")
                    Text("PLN")
                        .fontWeight(.bold)
                }
                .padding()
                .background(Color(.secondarySystemBackground))
                .cornerRadius(10)
            }
            
            VStack(alignment: .leading, spacing: 8) {
                Text("Tytuł przelewu")
                    .font(.caption)
                    .foregroundColor(.secondary)
                TextField("np. Opłata za fakturę", text: $viewModel.title)
                    .padding()
                    .background(Color(.secondarySystemBackground))
                    .cornerRadius(10)
                    .accessibilityIdentifier("transfer_title_field")
            }
        }
        .padding(.horizontal)
    }
    
    private var confirmationStep: some View {
        VStack(spacing: 0) {
            DetailRow(title: "Odbiorca", value: viewModel.recipientName)
            Divider()
            DetailRow(title: "Numer konta", value: viewModel.recipientAccount)
            Divider()
            DetailRow(title: "Kwota", value: "\(viewModel.amount) PLN", valueColor: .blue)
            Divider()
            DetailRow(title: "Tytuł", value: viewModel.title)
        }
        .background(Color(.secondarySystemBackground))
        .cornerRadius(12)
        .padding(.horizontal)
        .accessibilityIdentifier("transfer_confirmation_view")
    }
    
    private var successStep: some View {
        VStack(spacing: 20) {
            Image(systemName: "checkmark.circle.fill")
                .font(.system(size: 80))
                .foregroundColor(.green)
                .accessibilityIdentifier("transfer_success_icon")
            
            Text("Przelew wysłany!")
                .font(.title2)
                .fontWeight(.bold)
            
            Text("Kwota \(viewModel.amount) PLN została przekazana do realizacji dla \(viewModel.recipientName).")
                .multilineTextAlignment(.center)
                .foregroundColor(.secondary)
                .padding(.horizontal)
            
            Button("Powrót do pulpitu") {
                dismiss()
            }
            .padding()
            .buttonStyle(.borderedProminent)
            .accessibilityIdentifier("back_to_dashboard_button")
        }
    }
    
    private var stepNumber: Int {
        switch viewModel.currentStep {
        case .recipient: return 1
        case .amount: return 2
        case .confirmation, .success: return 3
        }
    }
    
    private var isNextDisabled: Bool {
        switch viewModel.currentStep {
        case .recipient: return !viewModel.isRecipientValid
        case .amount: return !viewModel.isAmountValid
        default: return false
        }
    }
}

struct StepCircle: View {
    let step: Int
    let currentStep: Int
    let title: String
    
    var body: some View {
        VStack(spacing: 4) {
            ZStack {
                Circle()
                    .fill(step <= currentStep ? Color.blue : Color.gray.opacity(0.3))
                    .frame(width: 30, height: 30)
                
                if step < currentStep {
                    Image(systemName: "checkmark")
                        .font(.caption.bold())
                        .foregroundColor(.white)
                } else {
                    Text("\(step)")
                        .font(.caption.bold())
                        .foregroundColor(.white)
                }
            }
            Text(title)
                .font(.system(size: 10))
                .foregroundColor(step <= currentStep ? .primary : .secondary)
        }
    }
}
