export interface Video {
  id: string;
  type: string;
  format: string;
  storageLocation: S3Object;
  width?: number;
  height?: number;
  duration?: number;
  bitrate?: number;
  codec?: string;
  size?: number;
  contentType?: string;
  streamInfo?: StreamInfo;
  metadata?: string;
  status?: string;
  thumbnail?: Image;
  createdAt: string;
  updatedAt: string;
}

export interface Image {
  fileName: string;
  url: string;
  storageLocation?: S3Object;
  width?: number;
  height?: number;
  size?: number;
  contentType?: string;
  metadata?: string;
}

export interface S3Object {
  bucket: string;
  key: string;
  url: string;
}

export interface StreamInfo {
  downloadUrl?: string;
  cdnPrefix?: string;
  playUrl?: string;
}

export interface PublishRule {
  publishAt?: string;
  unpublishAt?: string;
  regions?: string[];
  ageRating?: string;
}

export interface Asset {
  id: string;
  slug: string;
  title?: string;
  description?: string;
  type: string;
  genre?: string;
  genres?: string[];
  tags?: string[];
  status?: string;
  createdAt: string;
  updatedAt: string;
  metadata?: string;
  ownerId?: string;
  videos: Video[];
  publishRule?: PublishRule;
}

export interface Bucket {
  id: string;
  key: string;
  name: string;
  description?: string;
  type: string;
  status?: string;
  assetIds: string[];
  createdAt: string;
  updatedAt: string;
  assets?: Asset[];
}

export interface AssetResponse {
  assets: Asset[];
  count: number;
}

export interface BucketResponse {
  buckets: Bucket[];
  count: number;
} 