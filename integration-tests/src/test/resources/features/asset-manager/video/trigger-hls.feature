Feature: Trigger HLS Processing

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url apiGatewayUrl
    * header Content-Type = 'application/json'

Scenario: Trigger HLS for uploaded video
    Given path '/hls-job-requested'
    And request
    """
    {
        "assetId": "#(assetId)",
        "videoId": "#(videoId)",
        "input": "#(input)"
    }
    """
    When method POST
    Then status 200
    And match response.message == "HLS job requested successfully"
    And def result = response