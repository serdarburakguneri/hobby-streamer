package queries

const GetAssetsQuery = `
query GetAssets {
  assets {
    items {
      id
      slug
      title
      description
      type
      genre
      genres
      tags
      status
      createdAt
      updatedAt
      metadata
      ownerId
      videos {
        id
        label
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
