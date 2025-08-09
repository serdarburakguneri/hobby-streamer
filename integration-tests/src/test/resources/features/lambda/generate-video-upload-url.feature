Feature: Generate Video Upload URL Lambda

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url apiGatewayUrl
    * header Content-Type = 'application/json'
    * header Authorization = 'Bearer ' + authToken

Scenario: Generate video upload URL
    Given path '/upload'
    And request
    """
    {
        "assetId": "#(assetId)",
        "fileName": "#(fileName)",
        "videoType": "#(videoType)"
    }
    """
    When method POST
    Then status 200
    And match response.url != null
    And def result = response

