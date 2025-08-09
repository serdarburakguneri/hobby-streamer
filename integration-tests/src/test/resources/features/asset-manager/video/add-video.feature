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
    * def result = response    

    Given request
    """
    {
      "query": "query GetAsset($id: ID!) { asset(id: $id) { videos { id storageLocation { key } } } }",
      "variables": { "id": "#(assetId)" }
    }
    """
    And header Authorization = 'Bearer ' + authToken
    And header Content-Type = 'application/json'
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    * def vids = response.data.asset.videos
    * def videoId =
    """
    function(){
      for (var i = 0; i < vids.length; i++) {
        if (vids[i].storageLocation && vids[i].storageLocation.key == karate.get('videoKey')) {
          return vids[i].id;
        }
      }
      return null;
    }()
    """