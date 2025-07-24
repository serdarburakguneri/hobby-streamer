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
  images?: Image[];
  publishRule?: PublishRule;
}

export interface PublishRule {
  publishAt?: string;
  unpublishAt?: string;
  regions?: string[];
  ageRating?: string;
}

export enum VideoType {
  MAIN = 'main',
  TRAILER = 'trailer',
  BEHIND_THE_SCENES = 'behind',
  INTERVIEW = 'interview'
}

export enum VideoFormat {
  RAW = 'raw',
  HLS = 'hls',
  DASH = 'dash'
}

export interface S3Object {
  bucket: string;
  key: string;
  url: string;
}

export interface StreamInfo {
  downloadUrl?: string;
  cdnPrefix?: string;
  url?: string;
}

export interface Video {
  id: string;
  type: VideoType;
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
  metadata?: Record<string, any>;
  status?: string;
  thumbnail?: Image;
  createdAt: string;
  updatedAt: string;
}

export enum ImageType {
  THUMBNAIL = 'thumbnail',
  POSTER = 'poster',
  BANNER = 'banner',
  HERO = 'hero',
  LOGO = 'logo',
  SCREENSHOT = 'screenshot',
  BEHIND_THE_SCENES = 'behind_scenes',
  INTERVIEW = 'interview'
}

export interface Image {
  id: string;
  fileName: string;
  url: string;
  type: ImageType;
  storageLocation?: S3Object;
  width?: number;
  height?: number;
  size?: number;
  contentType?: string;
  metadata?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
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
  MOVIE = 'movie',
  TV_SHOW = 'tv_show',
  SERIES = 'series',
  SEASON = 'season',
  EPISODE = 'episode',
  DOCUMENTARY = 'documentary',
  SHORT = 'short',
  TRAILER = 'trailer',
  BONUS = 'bonus',
  BEHIND_SCENES = 'behind_scenes',
  INTERVIEW = 'interview',
  MUSIC_VIDEO = 'music_video',
  PODCAST = 'podcast',
  LIVE = 'live'
}

export interface Bucket {
  id: string;
  key: string;
  name: string;
  description?: string;
  type: BucketType;
  status?: string;
  createdAt: string;
  updatedAt: string;
  assets?: Asset[];
  ownerId?: string;
}

export enum BucketType {
  COLLECTION = 'collection',
  PLAYLIST = 'playlist',
  CATEGORY = 'category'
}

export enum BucketStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  DRAFT = 'draft'
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