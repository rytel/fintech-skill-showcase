//
//  LoggerExtensions.swift
//  demoBank
//
//  Created for robust logging support.
//

import Foundation
import OSLog

extension Logger {
    /// Subsystem identifies the app, useful for filtering in Console.app
    private static var subsystem = Bundle.main.bundleIdentifier ?? "com.demobank.app"

    /// Logs related to API calls, Networking, and Reachability
    static let network = Logger(subsystem: subsystem, category: "Network")

    /// Logs related to Authentication, Tokens, and Keychain access
    static let auth = Logger(subsystem: subsystem, category: "Auth")

    /// Logs related to UI, ViewModels, and User Interactions
    static let ui = Logger(subsystem: subsystem, category: "UI")

    /// General application lifecycle logs
    static let app = Logger(subsystem: subsystem, category: "App")
}
