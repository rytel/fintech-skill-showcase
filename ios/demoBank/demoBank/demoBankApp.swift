//
//  demoBankApp.swift
//  demoBank
//
//  Created by Rafal Rytel on 22/12/2025.
//

import SwiftUI

@main
struct demoBankApp: App {
    @StateObject private var appCoordinator = AppCoordinator()
    
    init() {
        DependencyContainer.shared.setupMocks()
    }
    
    var body: some Scene {
        WindowGroup {
            appCoordinator.rootView
                .onAppear {
                    appCoordinator.start()
                }
        }
    }
}