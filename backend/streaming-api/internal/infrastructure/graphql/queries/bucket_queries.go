package queries

const GetBucketQuery = `
query GetBucket($key: String!) {
  bucketByKey(key: $key) {
    id
    key
    name
    description
    type
    status
    createdAt
    updatedAt
    assets {
      id
      slug
      title
      description
      type
      genre
      genres
      tags
      status
      metadata
      createdAt
      updatedAt
      ownerId
      videos {
        id
        type
        format
        storageLocation { bucket key url }
        width
        height
        duration
        bitrate
        codec
        size
        contentType
        streamInfo { downloadUrl cdnPrefix url }
        metadata
        status
        thumbnail { id fileName url type storageLocation { bucket key url } width height size contentType metadata createdAt updatedAt }
        createdAt
        updatedAt
        quality
        isReady
        isProcessing
        isFailed
        segmentCount
        videoCodec
        audioCodec
        avgSegmentDuration
        segments
        frameRate
        audioChannels
        audioSampleRate
        transcodingInfo { jobId progress outputUrl error completedAt }
      }
      images {
        id
        fileName
        url
        type
        storageLocation { bucket key url }
        width
        height
        size
        contentType
        metadata
        createdAt
        updatedAt
      }
      publishRule {
        publishAt
        unpublishAt
        regions
        ageRating
      }
    }
  }
}`

const GetBucketsQuery = `
query GetBuckets($limit: Int, $nextKey: String) {
  buckets(limit: $limit, nextKey: $nextKey) {
    items {
      id
      key
      name
      description
      type
      status
      createdAt
      updatedAt
      assets {
        id
        slug
        title
        description
        type
        genre
        genres
        tags
        status
        metadata
        createdAt
        updatedAt
        ownerId
        videos {
          id
          type
          format
          storageLocation { bucket key url }
          width
          height
          duration
          bitrate
          codec
          size
          contentType
          streamInfo { downloadUrl cdnPrefix url }
          metadata
          status
          thumbnail { id fileName url type storageLocation { bucket key url } width height size contentType metadata createdAt updatedAt }
          createdAt
          updatedAt
          quality
          isReady
          isProcessing
          isFailed
          segmentCount
          videoCodec
          audioCodec
          avgSegmentDuration
          segments
          frameRate
          audioChannels
          audioSampleRate
          transcodingInfo { jobId progress outputUrl error completedAt }
        }
        images {
          id
          fileName
          url
          type
          storageLocation { bucket key url }
          width
          height
          size
          contentType
          metadata
          createdAt
          updatedAt
        }
        publishRule {
          publishAt
          unpublishAt
          regions
          ageRating
        }
      }
    }
  }
}`
