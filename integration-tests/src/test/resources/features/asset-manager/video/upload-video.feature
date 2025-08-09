Feature: Upload Video File

Background:
    * def sampleUtils = call read('classpath:sample-files-utils.js')

Scenario: Upload video file
    * def videoInfo = sampleUtils.getSampleVideoInfo()
    * def videoFileName = videoInfo.fileName
    * def videoContentType = videoInfo.contentType
    * def videoSize = videoInfo.size
    * def videoBucket = 'content-east'
    * def videoKey = assetId + '/source/' + videoFileName
    * def uploadUrlResult = call read('classpath:features/lambda/generate-video-upload-url.feature') { assetId: '#(assetId)', fileName: '#(videoFileName)', videoType: 'main' }
    * def rawUrl = uploadUrlResult.result.url
    * def uploadUrl = rawUrl.replace('localstack:4566', 'localhost:4566')
    # Build object URL without query for upsert equality and pass that to addVideo
    * def parts = rawUrl.split('\?')
    * def objectUrl = parts[0].replace('localstack:4566', 'localhost:4566')
    * url uploadUrl
    * header Content-Type = '#(videoContentType)'
    * request read('classpath:samples/videos/' + videoFileName)
    When method PUT
    Then status 200
    * def addVideoResult = call read('classpath:features/asset-manager/video/add-video.feature') { assetId: '#(assetId)', videoLabel: 'main', videoFormat: 'raw', videoBucket: '#(videoBucket)', videoKey: '#(videoKey)', videoUrl: '#(objectUrl)', contentType: '#(videoContentType)', videoSize: '#(videoSize)' }
    * print 'Add video result =', addVideoResult
    * def videoId = addVideoResult.videoId
    * def result = addVideoResult
    * def input = objectUrl