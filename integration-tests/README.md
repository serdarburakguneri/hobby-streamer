# Integration Tests

> Integration testing for HobbyStreamer services using Karate framework.

Tests GraphQL APIs, REST endpoints, authentication flows, service communication, health checks, file uploads, SQS message processing, data population.

## Features
GraphQL API testing, REST API testing, authentication flows, service-to-service communication, health checks, file upload testing, SQS message testing, bucket and asset lifecycle testing, data population, HTML reports, environment configuration.

## Requirements
Java 11+, Maven 3.6+, Docker Compose, AWS CLI.

## Quick Start

Start the local environment:
```bash
cd local
./build.sh
```

Run tests:
```bash
cd integration-tests
mvn test
```

## Usage

Run all tests:
```bash
mvn test
```

Run specific test suites:
```bash
mvn test -Dtest=KarateTestRunner#testAssetManager
mvn test -Dtest=KarateTestRunner#testAssetManagerBucketSuite
mvn test -Dtest=KarateTestRunner#testAssetManagerAssetSuite
mvn test -Dtest=KarateTestRunner#testDataPopulate
mvn test -Dtest=KarateTestRunner#testAuthService
mvn test -Dtest=KarateTestRunner#testStreamingApi
mvn test -Dtest=KarateTestRunner#testLambdaFunctions
```

Run with verbose output:
```bash
mvn test -Dkarate.options="--verbose"
```

## Test Structure
- **Bucket Suite**: Complete bucket lifecycle (create, update, add assets)
- **Asset Suite**: Asset management with video/image uploads
- **Data Population**: Creates buckets with multiple assets for testing
- **Auth Service**: Authentication and authorization flows
- **Streaming API**: Video streaming functionality
- **Lambda Functions**: AWS Lambda integration tests

## Reports
Test reports generated in `target/karate-reports/`. 