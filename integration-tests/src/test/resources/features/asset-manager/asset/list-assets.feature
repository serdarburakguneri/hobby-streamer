Feature: List Assets

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: List assets
    Given request
    """
    {
        "query": "query ListAssets { assets { items { id slug title description type status ownerId createdAt updatedAt } nextKey hasMore } }",
        "variables": {}
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.assets != null
    And match response.data.assets.items != null
    And def result = response 