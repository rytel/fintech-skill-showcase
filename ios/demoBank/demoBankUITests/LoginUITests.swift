import XCTest

final class LoginUITests: BaseUITest {
    
    
    func testSuccessfulLogin() throws {
        // 1. Sprawdź czy jesteśmy na ekranie logowania
        let welcomeText = app.staticTexts["Welcome to demoBank"]
        XCTAssertTrue(welcomeText.waitForExistence(timeout: 5))
        
        // 2. Wykonaj logowanie
        login(user: "test_user", pass: "password123")
        
        // 3. Sprawdź czy jesteśmy na Dashboardzie (napis "Moje Finanse")
        let dashboardTitle = app.navigationBars["Moje Finanse"]
        XCTAssertTrue(dashboardTitle.waitForExistence(timeout: 10), "Dashboard should appear after login")
    }
    
    
    func testInvalidLoginFeedback() throws {
        // 1. Wykonaj logowanie z błędnym hasłem
        login(user: "test_user", pass: "wrong_password", expectError: true)
        
        // 2. Sprawdź czy pojawił się komunikat o błędzie
        let errorMessage = app.staticTexts["error_message"]
        XCTAssertTrue(errorMessage.waitForExistence(timeout: 10), "Error message should be visible")
    }
}
