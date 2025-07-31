Feature: List Assets in Bucket

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: List assets in bucket
    Given request
    """
    {
        "query": "query GetBucketAssets($bucketKey: String!) { bucketByKey(key: $bucketKey) { id key name type status assets { id slug title type status } } }",
        "variables": {
            "bucketKey": "#(bucketKey)"
        }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.bucketByKey != null
    And match response.data.bucketByKey.assets != null
    And def result = response 