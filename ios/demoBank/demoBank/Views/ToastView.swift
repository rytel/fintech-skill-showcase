import SwiftUI

enum ToastType {
    case success
    case error
    case info
    
    var color: Color {
        switch self {
        case .success: return .green
        case .error: return .red
        case .info: return .blue
        }
    }
    
    var icon: String {
        switch self {
        case .success: return "checkmark.circle.fill"
        case .error: return "exclamationmark.circle.fill"
        case .info: return "info.circle.fill"
        }
    }
}

struct ToastView: View {
    let type: ToastType
    let message: String
    
    var body: some View {
        HStack(spacing: 12) {
            Image(systemName: type.icon)
                .foregroundColor(.white)
            
            Text(message)
                .font(.subheadline)
                .fontWeight(.medium)
                .foregroundColor(.white)
            
            Spacer()
        }
        .padding()
        .background(type.color)
        .cornerRadius(12)
        .shadow(radius: 4)
        .padding(.horizontal)
    }
}

#Preview {
    VStack {
        ToastView(type: .success, message: "Przelew wysłany pomyślnie!")
        ToastView(type: .error, message: "Wystąpił błąd podczas przesyłania.")
        ToastView(type: .info, message: "Trwa przetwarzanie...")
    }
}
