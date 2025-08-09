Feature: Upload Image File

Background:
    * def sampleUtils = call read('classpath:sample-files-utils.js')

Scenario: Upload image file
    * def imageInfo = sampleUtils.getSampleImageInfo()
    * def imageFileName = imageInfo.fileName
    * def imageContentType = imageInfo.contentType
    * def imageSize = imageInfo.size
    * def imageBucket = 'content-east'
    * def imageKey = assetId + '/images/poster/' + imageFileName
    * def uploadUrlResult = call read('classpath:features/lambda/generate-image-upload-url.feature') { assetId: '#(assetId)', fileName: '#(imageFileName)', imageType: 'poster' }
    * def rawUrl = uploadUrlResult.result.url
    * def urlParts = rawUrl.split('?')
    * def basePath = urlParts[0]
    * def queryParams = urlParts[1]
    * def modifiedBasePath = basePath.replace('localstack:4566', 'localhost:4566')
    * def finalUploadUrl = modifiedBasePath + '?' + queryParams
    * url finalUploadUrl
    * header Content-Type = '#(imageContentType)'
    * request read('classpath:samples/images/' + imageFileName)
    When method PUT
    Then status 200
    * def addImageResult = call read('classpath:features/asset-manager/image/add-image.feature') { assetId: '#(assetId)', imageType: 'poster', fileName: '#(imageFileName)', imageBucket: '#(imageBucket)', imageKey: '#(imageKey)', imageUrl: '#(modifiedBasePath)', contentType: '#(imageContentType)', imageSize: '#(imageSize)' }
    * print 'Add image result =', addImageResult