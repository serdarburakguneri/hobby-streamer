Feature: Create Asset

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Create asset      
    Given request
    """
    {
        "query": "mutation CreateAsset($input: CreateAssetInput!) { createAsset(input: $input) { id slug title type status genre genres } }",
        "variables": {
            "input": {
                "slug": "#(assetSlug)",
                "title": "#(assetTitle)",
                "description": "An amazing cinematic experience that will captivate audiences",
                "type": "movie",
                "genre": "#(assetGenre)",
                "genres": "#(assetGenres)",
                "ownerId": "test-user"
            }
        }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.createAsset != null
    And match response.data.createAsset.slug == '#(assetSlug)'
    And match response.data.createAsset.title == '#(assetTitle)'
    And match response.data.createAsset.type == 'movie'
    And def result = response 