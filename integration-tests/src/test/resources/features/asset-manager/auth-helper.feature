Feature: Authentication Helper

Background:
    * url 'http://localhost:8080'
    * header Content-Type = 'application/json'

Scenario: Login and get token
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
    And def access_token = response.access_token
    And def refresh_token = response.refresh_token
    And def token_type = response.token_type
    And def expires_in = response.expires_in 