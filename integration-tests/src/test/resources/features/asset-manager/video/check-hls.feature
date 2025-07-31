Feature: Check HLS Status

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Check HLS status for asset video
    Given request
    """
    {
        "query": "query GetAssetWithVideos($id: ID!) { asset(id: $id) { id title videos { id label format status } } }",
        "variables": {
            "id": "#(assetId)"
        }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.asset != null
    And def result = response