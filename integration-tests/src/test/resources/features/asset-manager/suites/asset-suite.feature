@asset-manager
Feature: Asset Management Suite

Background:  
    * def timestampUtil = call read('classpath:timestamp.js')
    * def sampleUtils = call read('classpath:sample-files-utils.js')
    * def testDataHelper = call read('classpath:test-data-helper.js')

Scenario: Asset flow
    * def movieTitle = testDataHelper.getRandomMovieTitle()
    * def generatedAssetSlug = testDataHelper.generateAssetSlug(movieTitle) + '-' + timestampUtil.generateTimestamp()
    * def assetGenre = testDataHelper.getRandomGenre()
    * def assetGenres = testDataHelper.getRandomGenres(2)
    
    # Step 1: Create asset
    * def createResult = call read('classpath:features/asset-manager/asset/create-asset.feature') { assetSlug: '#(generatedAssetSlug)', assetTitle: '#(movieTitle)', assetGenre: '#(assetGenre)', assetGenres: '#(assetGenres)' }
    * print 'Create asset result =', createResult
    * def assetId = createResult.response.data.createAsset.id

    # Step 2: Upload video
    * def uploadVideoResult = call read('classpath:features/asset-manager/video/upload-video.feature') { assetId: '#(assetId)' }
    * print 'DEBUG uploadVideoResult =', uploadVideoResult
    * def videoId = uploadVideoResult.videoId
    * print 'DEBUG videoId =', videoId
    * def input = uploadVideoResult.input

    # Step 3: Trigger HLS
    * def triggerHLSResult = call read('classpath:features/asset-manager/video/trigger-hls.feature') { assetId: '#(assetId)', videoId: '#(videoId)', input: '#(input)' }
    * print 'Trigger HLS result =', triggerHLSResult
    * match triggerHLSResult.response.errors == '#notpresent'
    * match triggerHLSResult.response.data.requestTranscode == true

    # Step 3b: Trigger DASH similar to CMS flow
    * def triggerDASHResult = call read('classpath:features/asset-manager/video/trigger-dash.feature') { assetId: '#(assetId)', videoId: '#(videoId)', input: '#(input)' }
    * print 'Trigger DASH result =', triggerDASHResult
    * match triggerDASHResult.response.errors == '#notpresent'
    * match triggerDASHResult.response.data.requestTranscode == true

    # Step 4: Wait for HLS processing
    * print 'Waiting for HLS processing to complete...'
    * java.lang.Thread.sleep(5000)

    # Step 5: Check HLS status
    * def checkHLSResult = call read('classpath:features/asset-manager/video/check-hls.feature') { assetId: '#(assetId)' }
    * print 'Check HLS result =', checkHLSResult
    
    # Step 6: Add image to the asset with upload URL
    * def addImageResult = call read('classpath:features/asset-manager/image/upload-image.feature') { assetId: '#(assetId)' }
    * print 'Add image result =', addImageResult    

    # Step 7: Get asset to verify it was created
    * def getAssetResult = call read('classpath:features/asset-manager/asset/get-asset.feature') { assetId: '#(assetId)' }
    * print 'Get asset result =', getAssetResult

    # Step 8: Update asset
    * def updatedTitle = 'Updated ' + movieTitle
    * def updateResult = call read('classpath:features/asset-manager/asset/update-asset.feature') { assetId: '#(assetId)', newDescription: 'Updated description for the asset' }
    * print 'Update asset result =', updateResult
