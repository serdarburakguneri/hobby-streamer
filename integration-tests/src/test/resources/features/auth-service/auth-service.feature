@auth
Feature: Auth Service API Tests

Background:
    * url authServiceUrl
    * header Content-Type = 'application/json'

Scenario: Health check
    Given path '/health'
    When method GET
    Then status 200

Scenario: Login with valid credentials
    Given path '/auth/login'
    And request
    """
    {
        "username": "admin",
        "password": "admin",
        "client_id": "asset-manager"
    }
    """
    When method POST
    Then status 200
    And match response.access_token != null
    And match response.refresh_token != null
    And match response.token_type == "Bearer"
    And match response.expires_in != null 