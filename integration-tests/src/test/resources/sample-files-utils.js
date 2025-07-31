function fn() {
    var sampleVideos = [
        {
            fileName: 'sample_video_1.mp4',
            contentType: 'video/mp4',
            size: 122880,
            bucket: 'content-east',
            key: 'samples/sample_video_1.mp4'
        }        
    ];
    
    var sampleImages = [
        {
            fileName: 'sample_image_1.png',
            contentType: 'image/png',
            size: 2726297,
            bucket: 'content-east',
            key: 'samples/sample_image_1.png'
        },
        {
            fileName: 'sample_image_2.png',
            contentType: 'image/png',
            size: 2726297,
            bucket: 'content-east',
            key: 'samples/sample_image_2.png'
        },
        {
            fileName: 'sample_image_3.png',
            contentType: 'image/png',
            size: 2726297,
            bucket: 'content-east',
            key: 'samples/sample_image_3.png'
        },
        {
            fileName: 'sample_image_4.png',
            contentType: 'image/png',
            size: 2726297,
            bucket: 'content-east',
            key: 'samples/sample_image_4.png'
        },
        {
            fileName: 'sample_image_5.png',
            contentType: 'image/png',
            size: 2726297,
            bucket: 'content-east',
            key: 'samples/sample_image_5.png'
        },
        {
            fileName: 'sample_image_6.png',
            contentType: 'image/png',
            size: 2726297,
            bucket: 'content-east',
            key: 'samples/sample_image_6.png'
        },
        {
            fileName: 'sample_image_7.png',
            contentType: 'image/png',
            size: 2726297,
            bucket: 'content-east',
            key: 'samples/sample_image_7.png'
        },
        {
            fileName: 'sample_image_8.png',
            contentType: 'image/png',
            size: 2726297,
            bucket: 'content-east',
            key: 'samples/sample_image_8.png'
        },
        {
            fileName: 'sample_image_9.png',
            contentType: 'image/png',
            size: 2726297,
            bucket: 'content-east',
            key: 'samples/sample_image_9.png'
        },
        {
            fileName: 'sample_image_10.png',
            contentType: 'image/png',
            size: 2726297,
            bucket: 'content-east',
            key: 'samples/sample_image_10.png'
        }
    ];
    
    var getRandomSampleVideo = function() {
        return sampleVideos[Math.floor(Math.random() * sampleVideos.length)];
    };
    
    var getRandomSampleImage = function() {
        return sampleImages[Math.floor(Math.random() * sampleImages.length)];
    };
    
    var getSampleVideoInfo = function() {
        return getRandomSampleVideo();
    };
    
    var getSampleImageInfo = function() {
        return getRandomSampleImage();
    };
    
    return { 
        getSampleVideoInfo: getSampleVideoInfo,
        getSampleImageInfo: getSampleImageInfo,
        getRandomSampleVideo: getRandomSampleVideo,
        getRandomSampleImage: getRandomSampleImage
    };
} 