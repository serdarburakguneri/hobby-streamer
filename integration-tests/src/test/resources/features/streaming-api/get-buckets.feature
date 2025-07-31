Feature: Get All Buckets from Streaming API

Background:
    * url streamingApiUrl
    * header Content-Type = 'application/json'

Scenario: Get all buckets
    Given path '/api/v1/buckets'
    When method GET
    Then status 200
    And match response contains { count: '#number', limit: '#number' }
    And def result = response 