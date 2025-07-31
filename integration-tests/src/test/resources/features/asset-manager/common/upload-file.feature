Feature: Upload File Helper

Background:
    * def sampleUtils = call read('classpath:sample-files-utils.js')

Scenario: Upload video file
    * def videoInfo = sampleUtils.getSampleVideoInfo()
    * def videoFile = read('classpath:samples/videos/' + videoInfo.fileName)
    * url uploadUrl
    * header Content-Type = '#(videoInfo.contentType)'
    When method PUT
    Then status 200

Scenario: Upload image file
    * def imageInfo = sampleUtils.getSampleImageInfo()
    * def imageFile = read('classpath:samples/images/' + imageInfo.fileName)
    * url uploadUrl
    * header Content-Type = '#(imageInfo.contentType)'
    When method PUT
    Then status 200 