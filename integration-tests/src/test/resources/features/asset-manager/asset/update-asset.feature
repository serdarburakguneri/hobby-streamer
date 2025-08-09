Feature: Update Asset

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Update asset
    Given request
    """
    {
        "query": "mutation UpdateAssetDescription($id: ID!, $description: String!) { updateAssetDescription(id: $id, description: $description) { id slug title description type status ownerId updatedAt } }",
        "variables": {
            "id": "#(assetId)",
            "description": "#(newDescription)"
        }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.updateAssetDescription != null   
    And match response.data.updateAssetDescription.description == '#(newDescription)'
    And def result = response 