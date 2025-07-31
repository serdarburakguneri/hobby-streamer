Feature: Upload Video File

Background:
    * def sampleUtils = call read('classpath:sample-files-utils.js')

Scenario: Upload video file
    * def videoInfo = sampleUtils.getSampleVideoInfo()
    * def videoFileName = videoInfo.fileName
    * def videoContentType = videoInfo.contentType
    * def videoSize = videoInfo.size
    * def videoUploadResult = call read('classpath:features/lambda/generate-video-upload-url.feature') { assetId: '#(assetId)', fileName: '#(videoFileName)', videoType: 'main' }
    * print 'Video upload URL result =', videoUploadResult
    * print 'Actual upload URL =', videoUploadResult.result.url
    * def videoBucket = 'content-east'
    * def videoKey = assetId + '/source/' + videoFileName
     # too lazy for etc host defs
    * def uploadUrl = videoUploadResult.result.url.replace('localstack:4566', 'localhost:4566')
    * print 'Modified upload URL =', uploadUrl
    * url uploadUrl
    * header Content-Type = '#(videoContentType)'
    * request read('classpath:samples/videos/' + videoFileName)
    When method PUT
    Then status 200
    * def addVideoResult = call read('classpath:features/asset-manager/video/add-video.feature') { assetId: '#(assetId)', videoLabel: 'main', videoFormat: 'raw', videoBucket: '#(videoBucket)', videoKey: '#(videoKey)', videoUrl: '#(videoUploadResult.result.url)', contentType: '#(videoContentType)', videoSize: '#(videoSize)' }
    * print 'Add video result =', addVideoResult
    * def result = addVideoResult 