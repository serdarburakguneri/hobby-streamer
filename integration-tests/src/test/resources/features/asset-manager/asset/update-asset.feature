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
        "query": "mutation PatchAsset($id: ID!, $patches: [JSONPatch!]!) { patchAsset(id: $id, patches: $patches) { id slug title description type status ownerId updatedAt } }",
        "variables": {
            "id": "#(assetId)",
            "patches": [               
                {"op": "replace", "path": "/description", "value": "#(newDescription)"}
            ]
        }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.patchAsset != null   
    And match response.data.patchAsset.description == '#(newDescription)'
    And def result = response 