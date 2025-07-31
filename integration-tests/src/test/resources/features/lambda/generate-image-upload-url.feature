Feature: Generate Image Upload URL Lambda

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url apiGatewayUrl
    * header Content-Type = 'application/json'
    * header Authorization = 'Bearer ' + authToken

Scenario: Generate image upload URL
    Given path '/image-upload'
    And request
    """
    {
        "assetId": "#(assetId)",
        "fileName": "#(fileName)",
        "imageType": "#(imageType)"
    }
    """
    When method POST
    Then status 200
    And match response.url != null
    And def result = response 