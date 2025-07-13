export interface Asset {
  id: string;
  slug: string;
  title?: string;
  description?: string;
  type: AssetType;
  genre?: string;
  genres?: string[];
  tags?: string[];
  status?: string;
  createdAt: string;
  updatedAt: string;
  metadata?: Record<string, any>;
  ownerId?: string;
  parent?: Asset;
  children?: Asset[];
  buckets?: Bucket[];
  videos?: Video[];
  publishRule?: PublishRule;
}

export interface PublishRule {
  isPublic: boolean;
  publishAt?: string;
  unpublishAt?: string;
  regions?: string[];
  ageRating?: string;
}

export enum VideoType {
  MAIN = 'MAIN',
  TRAILER = 'TRAILER',
  BEHIND_THE_SCENES = 'BEHIND_THE_SCENES',
  INTERVIEW = 'INTERVIEW'
}

export interface S3Object {
  bucket: string;
  key: string;
  url: string;
}

export interface StreamInfo {
  downloadUrl?: string;
  cdnPrefix?: string;
}

export interface VideoVariant {
  storageLocation: S3Object;
  width?: number;
  height?: number;
  duration?: number;
  bitrate?: number;
  codec?: string;
  size?: number;
  contentType?: string;
  streamInfo?: StreamInfo;
  metadata?: Record<string, any>;
  status?: string;
}

export interface Video {
  type: VideoType;
  raw?: VideoVariant;
  hls?: VideoVariant;
  dash?: VideoVariant;
  thumbnail?: Image;
}

export interface Image {
  fileName: string;
  url: string;
  storageLocation?: S3Object;
  width?: number;
  height?: number;
  size?: number;
  contentType?: string;
  metadata?: Record<string, any>;
}


export enum AssetStatus {
  DRAFT = 'draft',
  PUBLISHED = 'published'
}


export enum VideoStatus {
  PENDING = 'pending',
  ANALYZING = 'analyzing',
  TRANSCODING = 'transcoding',
  READY = 'ready',
  FAILED = 'failed'
}


export enum AssetGenre {
  ACTION = 'action',
  DRAMA = 'drama',
  COMEDY = 'comedy',
  HORROR = 'horror',
  SCI_FI = 'sci_fi',
  ROMANCE = 'romance',
  THRILLER = 'thriller',
  FANTASY = 'fantasy',
  DOCUMENTARY = 'documentary',
  MUSIC = 'music',
  NEWS = 'news',
  SPORTS = 'sports',
  KIDS = 'kids',
  EDUCATIONAL = 'educational'
}


export const formatGenreName = (genre: string): string => {
  return genre.replace('_', ' ').toUpperCase();
};



export enum AssetType {
  MOVIE = 'MOVIE',
  SERIES = 'SERIES',
  SEASON = 'SEASON',
  EPISODE = 'EPISODE',
  DOCUMENTARY = 'DOCUMENTARY',
  MUSIC = 'MUSIC',
  PODCAST = 'PODCAST',
  TRAILER = 'TRAILER',
  BEHIND_THE_SCENES = 'BEHIND_THE_SCENES',
  INTERVIEW = 'INTERVIEW'
}

export interface Bucket {
  id: string;
  key: string;
  name: string;
  description?: string;
  type: BucketType;
  status?: string;
  assetIds?: string[];
  createdAt: string;
  updatedAt: string;
  assets?: Asset[];
}

export enum BucketType {
  COLLECTION = 'COLLECTION',
  PLAYLIST = 'PLAYLIST',
  CATEGORY = 'CATEGORY'
}

export enum BucketStatus {
  ACTIVE = 'ACTIVE',
  INACTIVE = 'INACTIVE',
  DRAFT = 'DRAFT'
}

export interface AssetPage {
  items: Asset[];
  nextKey?: string;
}

export interface BucketPage {
  items: Bucket[];
  nextKey?: string;
}

export interface AssetInput {
  slug: string;
  title?: string;
  description?: string;
  type: AssetType;
  genre?: string;
  genres?: string[];
  tags?: string[];
  metadata?: string;
  ownerId?: string;
  parentId?: string;
}

export interface BucketInput {
  key: string;
  name: string;
  description?: string;
  type: BucketType;
  status?: string;
  assetIds?: string[];
}

export interface AssetCreateDTO {
  title?: string;
  slug?: string;
  description?: string;
  type: AssetType;
  genre?: string;
  genres?: string[];
  tags?: string[];
  metadata?: Record<string, any>;
  ownerId?: string;
  parentId?: string;
}

export interface AssetUpdateDTO {
  title?: string;
  description?: string;
  type?: AssetType;
  genre?: string;
  genres?: string[];
  tags?: string[];
  metadata?: Record<string, any>;
  ownerId?: string;
  parentId?: string;
  status?: string;
}

export interface UploadResponse {
  uploadUrl: string;
  assetId: number;
  fileName: string;
}

export interface AssetUpdateInput {
  title?: string | null;
  description?: string | null;
  type?: AssetType | null;
  genre?: string | null;
  genres?: string[] | null;
  tags?: string[] | null;
  metadata?: string | null;
  ownerId?: string | null;
  parentId?: string | null;
  clearFields?: string[];
} 