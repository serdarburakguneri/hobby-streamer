Feature: Create Bucket

Background:
    * call read('classpath:features/asset-manager/auth-helper.feature')
    * def authToken = access_token
    * url assetManagerUrl
    * header Authorization = 'Bearer ' + authToken

Scenario: Create bucket simple    
  
    Given request
    """
    {
      "query": "mutation CreateBucket($input: CreateBucketInput!) { createBucket(input: $input) { id key name type status } }",
      "variables": {
        "input": {
          "key": '#(bucketKey)',
          "name": '#(bucketName)',
          "description": "A curated collection of amazing content",
          "type": "collection",
          "ownerId": "test-user",
          "status": "draft"
        }
      }
    }
    """
    When method POST
    Then status 200   
    And match response.errors == '#notpresent'
    And match response.data.createBucket != null
    And match response.data.createBucket.key == '#(bucketKey)'
    And match response.data.createBucket.name == '#(bucketName)'
    And match response.data.createBucket.type == 'collection'
    And match response.data.createBucket.status == 'draft'
    And def result = response
   