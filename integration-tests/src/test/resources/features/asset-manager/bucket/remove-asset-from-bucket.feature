Feature: Remove Asset from Bucket

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Remove asset from bucket
    Given request
    """
    {
        "query": "mutation RemoveAssetFromBucket($input: RemoveAssetFromBucketInput!) { removeAssetFromBucket(input: $input) }",
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
    And match response.data.removeAssetFromBucket == true
    And def result = response 