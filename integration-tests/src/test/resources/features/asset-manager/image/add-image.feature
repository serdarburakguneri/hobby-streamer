Feature: Add Image to Asset

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Add image to asset
    Given request
    """
    {
        "query": "mutation AddImage($input: AddImageInput!) { addImage(input: $input) { id slug title description type status ownerId createdAt updatedAt images { id fileName url type storageLocation { bucket key url } width height size contentType createdAt updatedAt } } }",
        "variables": {
            "input": {
                "assetId": "#(assetId)",
                "type": "#(imageType)",
                "fileName": "#(fileName)",
                "bucket": "#(imageBucket)",
                "key": "#(imageKey)",
                "url": "#(imageUrl)",
                "contentType": "#(contentType)",
                "size": "#(imageSize)"
            }
        }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.addImage != null
    And match response.data.addImage.images != null
    And def result = response 