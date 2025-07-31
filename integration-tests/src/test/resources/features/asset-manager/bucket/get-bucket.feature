Feature: Get Bucket

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Get bucket
    Given request
    """
    {
        "query": "query GetBucket($key: String!) { bucketByKey(key: $key) { id key name type status } }",
        "variables": {
            "key": "#(bucketKey)"
        }
    }
    """
    When method POST
    Then status 200   
    And match response.errors == '#notpresent'
    And match response.data.bucketByKey != null
    And match response.data.bucketByKey.id != null
    And match response.data.bucketByKey.key != null
    And def result = response 