Feature: Add Asset to Bucket

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Add asset to bucket
    Given request
    """
    {
        "query": "mutation AddAssetToBucket($input: AddAssetToBucketInput!) { addAssetToBucket(input: $input) }",
        "variables": {
            "input": {
                "bucketId": "#(bucketId)",
                "assetId": "#(assetId)",
                "ownerId": "#(ownerId)"
            }
        }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.addAssetToBucket == true
    And def result = response 