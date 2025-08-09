Feature: Delete Bucket

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Delete bucket
    Given request
    """
    {
        "query": "mutation DeleteBucket($id: ID!) { deleteBucket(id: $id) }",
        "variables": {
            "id": "#(bucketId)"
        }
    }
    """
    When method POST
    Then status 200
    And match response.errors == '#notpresent'
    And match response.data.deleteBucket == true
    And def result = response 