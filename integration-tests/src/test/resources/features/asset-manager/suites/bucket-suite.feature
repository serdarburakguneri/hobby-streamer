@asset-manager
Feature: Bucket Management Suite

Background:  
    * def timestampUtil = call read('classpath:timestamp.js')
    * def sampleUtils = call read('classpath:sample-files-utils.js')
    * def testDataHelper = call read('classpath:test-data-helper.js')

Scenario: Bucket flow
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
    
    # Step 8: Add asset to bucket
    * def addAssetResult = call read('classpath:features/asset-manager/bucket/add-asset-to-bucket.feature') { bucketId: '#(bucketId)', assetId: '#(assetId)', ownerId: 'test-user' }
    * print 'Add asset to bucket result =', addAssetResult
    
    # Step 9: Get bucket to verify all data
    * def getBucketResult = call read('classpath:features/asset-manager/bucket/get-bucket.feature') { bucketKey: '#(generatedBucketKey)' }
    * print 'Get bucket result =', getBucketResult
    
    # Step 10: Remove asset from bucket
    * def removeAssetResult = call read('classpath:features/asset-manager/bucket/remove-asset-from-bucket.feature') { bucketId: '#(bucketId)', assetId: '#(assetId)' }
    * print 'Remove asset from bucket result =', removeAssetResult
    
    # Step 11: Delete bucket
    * def deleteBucketResult = call read('classpath:features/asset-manager/bucket/delete-bucket.feature') { bucketId: '#(bucketId)' }
    * print 'Delete bucket result =', deleteBucketResult





   