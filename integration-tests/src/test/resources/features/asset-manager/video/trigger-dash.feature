Feature: Trigger DASH Processing

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Trigger DASH for uploaded video
    Given request
    """
    {
      "query": "mutation RequestTranscode($assetId: ID!, $videoId: ID!, $format: VideoFormat!) { requestTranscode(assetId: $assetId, videoId: $videoId, format: $format) }",
      "variables": {
        "assetId": "#(assetId)",
        "videoId": "#(videoId)",
        "format": "dash"
      }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.requestTranscode == true
    And def result = response

