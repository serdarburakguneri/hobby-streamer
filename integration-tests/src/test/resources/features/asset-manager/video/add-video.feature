Feature: Add Video to Asset

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Add video to asset
    Given request
    """
    {
        "query": "mutation AddVideo($input: AddVideoInput!) { addVideo(input: $input) { id label type format storageLocation { bucket key url } width height duration bitrate codec size contentType status quality isReady isProcessing isFailed createdAt updatedAt } }",
        "variables": {
            "input": {
                "assetId": "#(assetId)",
                "label": "#(videoLabel)",
                "format": "#(videoFormat)",
                "bucket": "#(videoBucket)",
                "key": "#(videoKey)",
                "url": "#(videoUrl)",
                "contentType": "#(contentType)",
                "size": "#(videoSize)"
            }
        }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.addVideo != null
    And match response.data.addVideo.label == '#(videoLabel)'
    And match response.data.addVideo.format == '#(videoFormat)'
    And match response.data.addVideo.status == 'pending'
    And def result = response 