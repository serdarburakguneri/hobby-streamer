Feature: Update Bucket

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Update bucket status
    Given request
    """
    {
        "query": "mutation UpdateBucket($input: UpdateBucketInput!) { updateBucket(input: $input) { id key name type status } }",
        "variables": {
            "input": {
                "id": "#(bucketId)",
                "status": "#(newStatus)"
            }
        }
    }
    """
    When method POST
    Then status 200   
    And match response.errors == '#notpresent'
    And match response.data.updateBucket != null
    And match response.data.updateBucket.id != null
    And match response.data.updateBucket.key == '#(bucketKey)'
    And match response.data.updateBucket.name == '#(bucketName)'
    And match response.data.updateBucket.type == '#(bucketType)'
    And match response.data.updateBucket.status == '#(newStatus)'
    And def result = response 