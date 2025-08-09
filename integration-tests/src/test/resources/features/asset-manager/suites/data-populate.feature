@asset-manager
Feature: Data Population Suite

Background:  
    * def timestampUtil = call read('classpath:timestamp.js')
    * def sampleUtils = call read('classpath:sample-files-utils.js')
    * def testDataHelper = call read('classpath:test-data-helper.js')

Scenario: Populate database with buckets and assets using loops
    # Configuration
    * def bucketCount = 1
    * def assetsPerBucket = 5
    
    # Create buckets and assets using loops
    * eval
    """
    for (var bucketIndex = 1; bucketIndex <= bucketCount; bucketIndex++) {
      // Create bucket
      var bucketName = testDataHelper.getRandomBucketName();
      var generatedBucketKey = testDataHelper.generateBucketKey(bucketName) + '-' + timestampUtil.generateTimestamp();
      var createBucketResult = karate.call('classpath:features/asset-manager/bucket/create-bucket.feature', { bucketName: bucketName, bucketKey: generatedBucketKey });
      var bucketId = createBucketResult.response.data.createBucket.id;
      
      // Create assets for this bucket
      for (var assetIndex = 1; assetIndex <= assetsPerBucket; assetIndex++) {
        var movieTitle = testDataHelper.getRandomMovieTitle();
        var generatedAssetSlug = testDataHelper.generateAssetSlug(movieTitle) + '-' + timestampUtil.generateTimestamp();
        var assetGenre = testDataHelper.getRandomGenre();
        var assetGenres = testDataHelper.getRandomGenres(2);
        
        // Create asset
        var createAssetResult = karate.call('classpath:features/asset-manager/asset/create-asset.feature', { assetSlug: generatedAssetSlug, assetTitle: movieTitle, assetGenre: assetGenre, assetGenres: assetGenres });
        var assetId = createAssetResult.response.data.createAsset.id;
        
        // Upload video
        var addVideoResult = karate.call('classpath:features/asset-manager/video/upload-video.feature', { assetId: assetId });
        var videoId = addVideoResult.videoId;

        // Trigger HLS & DASH transcodes like CMS
        karate.call('classpath:features/asset-manager/video/trigger-hls.feature', { assetId: assetId, videoId: videoId, input: addVideoResult.input });
        karate.call('classpath:features/asset-manager/video/trigger-dash.feature', { assetId: assetId, videoId: videoId, input: addVideoResult.input });

        // Upload image (poster)
        karate.call('classpath:features/asset-manager/image/upload-image.feature', { assetId: assetId });

        // Add asset to bucket
        karate.call('classpath:features/asset-manager/bucket/add-asset-to-bucket.feature', { bucketId: bucketId, assetId: assetId, ownerId: 'test-user' });
      }
    }
    """
    
    # Summary
    * print 'Data population completed successfully!'
    * print 'Created', bucketCount, 'buckets with', assetsPerBucket, 'assets each'
    * print 'Total assets created:', bucketCount * assetsPerBucket
    * print 'Uploaded videos and images, and added each asset to its bucket'