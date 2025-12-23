//
//  KeychainHelper.swift
//  demoBank
//

import Foundation
import Security


final class KeychainHelper {
    static let shared = KeychainHelper()
    private init() {}
    
    private var hasResetForUIInThisRun = false
    
    func save(_ data: Data, service: String, account: String) {
        let query = [
            kSecValueData: data,
            kSecAttrService: service,
            kSecAttrAccount: account,
            kSecClass: kSecClassGenericPassword,
            kSecAttrAccessible: kSecAttrAccessibleAfterFirstUnlock
        ] as [CFString: Any]

        // Add item to keychain
        let status = SecItemAdd(query as CFDictionary, nil)

        if status == errSecDuplicateItem {
            // Item already exists, so update it.
            let query = [
                kSecAttrService: service,
                kSecAttrAccount: account,
                kSecClass: kSecClassGenericPassword
            ] as [CFString: Any]

            let attributesToUpdate = [kSecValueData: data] as [CFString: Any]

            SecItemUpdate(query as CFDictionary, attributesToUpdate as CFDictionary)
        }
    }
    
    func read(service: String, account: String) -> Data? {
        let query = [
            kSecAttrService: service,
            kSecAttrAccount: account,
            kSecClass: kSecClassGenericPassword,
            kSecReturnData: kCFBooleanTrue as Any,
            kSecMatchLimit: kSecMatchLimitOne
        ] as [CFString: Any]

        var result: AnyObject?
        SecItemCopyMatching(query as CFDictionary, &result)

        return result as? Data
    }
    
    func delete(service: String, account: String) {
        let query = [
            kSecAttrService: service,
            kSecAttrAccount: account,
            kSecClass: kSecClassGenericPassword
        ] as [CFString: Any]

        SecItemDelete(query as CFDictionary)
    }
}

// Helper for Token specifically
extension KeychainHelper {
    private static let tokenService = "com.demobank.auth"
    private static let tokenAccount = "jwt-token"
    
    func saveToken(_ token: String) {
        if let data = token.data(using: .utf8) {
            save(data, service: KeychainHelper.tokenService, account: KeychainHelper.tokenAccount)
        }
    }
    
    func getToken() -> String? {
        // Reset only ONCE per app launch if argument is present
        if ProcessInfo.processInfo.arguments.contains("-resetKeychain") && !hasResetForUIInThisRun {
            deleteToken()
            hasResetForUIInThisRun = true
            return nil
        }
        
        if let data = read(service: KeychainHelper.tokenService, account: KeychainHelper.tokenAccount) {
            return String(data: data, encoding: .utf8)
        }
        return nil
    }
    
    func deleteToken() {
        delete(service: KeychainHelper.tokenService, account: KeychainHelper.tokenAccount)
    }
}
