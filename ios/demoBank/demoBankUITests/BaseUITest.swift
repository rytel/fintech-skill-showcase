import XCTest


@MainActor
class BaseUITest: XCTestCase {
    var app: XCUIApplication!
    
    override func setUpWithError() throws {
        continueAfterFailure = false
        app = XCUIApplication()
        app.launchArguments.append("-skipBiometrics")
        app.launchArguments.append("-resetKeychain")
        app.launchArguments.append("-useMockData")
        app.launch()
    }
    
    override func tearDownWithError() throws {
        app = nil
    }
    
    func login(user: String = "test_user", pass: String = "password123", expectError: Bool = false) {
        XCTContext.runActivity(named: "Login Journey") { _ in
            let userField = app.textFields["username_field"]
            let passField = app.secureTextFields["password_field"]
            let loginButton = app.buttons["login_button"]
            
            XCTAssertTrue(userField.waitForExistence(timeout: 10), "Login screen should appear")
            
            userField.tap()
            userField.typeText(user)
            
            passField.tap()
            passField.typeText(pass)
            
            XCTAssertTrue(loginButton.isEnabled, "Login button should be enabled")
            loginButton.tap()
            
            if !expectError {
                // Sprawdź czy nie ma błędu
                let errorText = app.staticTexts["error_message"]
                if errorText.exists {
                    XCTFail("Login failed with error: \(errorText.label)")
                }
            }
        }
    }
}
