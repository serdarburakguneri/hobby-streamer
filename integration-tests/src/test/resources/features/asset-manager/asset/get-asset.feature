Feature: Get Asset

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Get asset by slug
    Given request
    """
    {
        "query": "query GetAsset($id: ID!) { asset(id: $id) { id slug title description type status ownerId createdAt updatedAt } }",
        "variables": {
            "id": "#(assetId)"
        }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.asset != null
    And match response.data.asset.id != null
    And match response.data.asset.slug != null
    And def result = response 