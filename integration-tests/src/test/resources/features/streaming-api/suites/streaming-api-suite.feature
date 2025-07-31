@streaming-api
Feature: Streaming API Suite

Background:  
    * def timestampUtil = call read('classpath:timestamp.js')
    * def sampleUtils = call read('classpath:sample-files-utils.js')
    * def testDataHelper = call read('classpath:test-data-helper.js')

Scenario: Streaming API flow
    * def bucketName = testDataHelper.getRandomBucketName()
    * def generatedBucketKey = testDataHelper.generateBucketKey(bucketName) + '-' + timestampUtil.generateTimestamp()
    * def movieTitle = testDataHelper.getRandomMovieTitle()
    * def generatedAssetSlug = testDataHelper.generateAssetSlug(movieTitle) + '-' + timestampUtil.generateTimestamp()
    * def assetGenre = testDataHelper.getRandomGenre()
    * def assetGenres = testDataHelper.getRandomGenres(2)
    
    # Step 1: Create bucket
    * def createBucketResult = call read('classpath:features/asset-manager/bucket/create-bucket.feature') { bucketName: '#(bucketName)', bucketKey: '#(generatedBucketKey)' }
    * print 'Create bucket result =', createBucketResult
    * def bucketId = createBucketResult.response.data.createBucket.id

    # Step 2: Create asset
    * def createAssetResult = call read('classpath:features/asset-manager/asset/create-asset.feature') { assetSlug: '#(generatedAssetSlug)', assetTitle: '#(movieTitle)', assetGenre: '#(assetGenre)', assetGenres: '#(assetGenres)' }
    * print 'Create asset result =', createAssetResult
    * def assetId = createAssetResult.response.data.createAsset.id
    
    # Step 3: Wait for data to propagate to streaming API
    * def waitTime = 5000
    * java.lang.Thread.sleep(waitTime)
    
    # Step 4: Query streaming API to get all buckets
    * def getAllBucketsResult = call read('classpath:features/streaming-api/get-buckets.feature')
    * print 'Streaming API all buckets result =', getAllBucketsResult    

