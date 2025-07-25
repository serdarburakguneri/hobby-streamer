import { gql, useApolloClient } from '@apollo/client';
import axios from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { Asset, AssetCreateDTO, AssetUpdateDTO, AssetPage, AssetInput, AssetType, Image, ImageType, BucketStatus } from '../types/asset';
import { API_CONFIG } from '../config/api';

const AUTH_BASE_URL = API_CONFIG.AUTH_BASE_URL;
const GRAPHQL_BASE_URL = API_CONFIG.GRAPHQL_BASE_URL;
const authApi = axios.create({
  baseURL: AUTH_BASE_URL,
  timeout: 10000,
});

let logoutCallback: (() => void) | null = null;

export const setLogoutCallback = (callback: () => void) => {
  logoutCallback = callback;
};

let isRefreshing = false;
let refreshPromise: Promise<{ accessToken: string; refreshToken: string }> | null = null;

export const refreshTokenIfNeeded = async (): Promise<string | null> => {
  try {
    const currentToken = await getAuthToken();
    if (!currentToken) {
      return null;
    }

    if (!isTokenExpiringSoon(currentToken)) {
      return currentToken;
    }

    if (isRefreshing && refreshPromise) {
      const result = await refreshPromise;
      return result.accessToken;
    }

    isRefreshing = true;
    refreshPromise = performTokenRefresh();
    
    const result = await refreshPromise;
    return result.accessToken;
  } catch (error) {
    console.error('Token refresh failed:', error);
    await clearAuthTokens();
    if (logoutCallback) {
      logoutCallback();
    }
    return null;
  } finally {
    isRefreshing = false;
    refreshPromise = null;
  }
};

const performTokenRefresh = async (): Promise<{ accessToken: string; refreshToken: string }> => {
  const refreshToken = await getRefreshToken();
  if (!refreshToken) {
    throw new Error('No refresh token available');
  }

  console.log('Refreshing token...');
  const result = await authService.refreshToken(refreshToken);
  
  await setAuthTokens(result.accessToken, result.refreshToken);
  console.log('Token refreshed successfully');
  
  return result;
};

const retryWithBackoff = async <T>(
  operation: () => Promise<T>,
  maxRetries: number = 3,
  baseDelay: number = 1000
): Promise<T> => {
  let lastError: any;
  
  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await operation();
    } catch (error: any) {
      lastError = error;
      
      if (error.graphQLErrors?.some((e: any) => e.extensions?.code === 'UNAUTHENTICATED')) {
        throw error;
      }
      
      if (attempt === maxRetries) {
        throw error;
      }
      
      const delay = baseDelay * Math.pow(2, attempt);
      await new Promise(resolve => setTimeout(resolve, delay));
    }
  }
  
  throw lastError;
};

export const getAuthToken = async (): Promise<string | null> => {
  try {
    const token = await AsyncStorage.getItem('authToken');
    console.log('getAuthToken called, found token:', !!token);
    return token;
  } catch (error) {
    console.log('getAuthToken called, AsyncStorage error:', error);
    return null;
  }
};

export const getRefreshToken = async (): Promise<string | null> => {
  try {
    const token = await AsyncStorage.getItem('refreshToken');
    console.log('getRefreshToken called, found token:', !!token);
    return token;
  } catch (error) {
    console.log('getRefreshToken called, AsyncStorage error:', error);
    return null;
  }
};

export const setAuthTokens = async (accessToken: string, refreshToken: string): Promise<void> => {
  try {
    await AsyncStorage.setItem('authToken', accessToken);
    await AsyncStorage.setItem('refreshToken', refreshToken);
    console.log('setAuthTokens called, tokens stored:', !!accessToken, !!refreshToken);
  } catch (error) {
    console.log('setAuthTokens called, AsyncStorage error:', error);
  }
};

export const setAuthToken = async (token: string): Promise<void> => {
  try {
    await AsyncStorage.setItem('authToken', token);
    console.log('setAuthToken called, token stored:', !!token);
  } catch (error) {
    console.log('setAuthToken called, AsyncStorage error:', error);
  }
};

export const clearAuthTokens = async (): Promise<void> => {
  try {
    await AsyncStorage.removeItem('authToken');
    await AsyncStorage.removeItem('refreshToken');
    console.log('clearAuthTokens called, tokens removed');
  } catch (error) {
    console.log('clearAuthTokens called, AsyncStorage error:', error);
  }
};

export const clearAuthToken = async (): Promise<void> => {
  await clearAuthTokens();
};

// Helper function to validate token locally (client-side)
export const validateTokenLocally = (token: string): { valid: boolean; user?: any; message?: string } => {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) {
      return { valid: false, message: 'Invalid token format' };
    }

    const payload = JSON.parse(atob(parts[1]));
    const now = Math.floor(Date.now() / 1000);

    // Check expiration
    if (payload.exp && payload.exp < now) {
      return { valid: false, message: 'Token expired' };
    }

    // Check not-before time
    if (payload.nbf && payload.nbf > now) {
      return { valid: false, message: 'Token not yet valid' };
    }

    // Extract user info
    const user = {
      id: payload.sub || '',
      username: payload.preferred_username || '',
      email: payload.email || '',
      roles: []
    };

    // Extract roles
    if (payload.realm_access && payload.realm_access.roles) {
      user.roles = payload.realm_access.roles;
    }

    return { valid: true, user };
  } catch (error) {
    return { valid: false, message: 'Invalid token' };
  }
};

// Helper function to check if token is about to expire (within 5 minutes)
export const isTokenExpiringSoon = (token: string): boolean => {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) {
      return true; // Consider invalid tokens as expiring soon
    }

    const payload = JSON.parse(atob(parts[1]));
    const now = Math.floor(Date.now() / 1000);
    const fiveMinutesFromNow = now + (5 * 60); // 5 minutes in seconds

    // Check if token expires within the next 5 minutes
    if (payload.exp && payload.exp < fiveMinutesFromNow) {
      return true;
    }

    return false;
  } catch (error) {
    return true; // Consider tokens with parsing errors as expiring soon
  }
};


const GET_ASSETS = gql`
  query GetAssets($limit: Int, $nextKey: String) {
    assets(limit: $limit, nextKey: $nextKey) {
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
        parent {
          id
          title
          type
        }
        publishRule {
          publishAt
          unpublishAt
          regions
          ageRating
        }
        videos {
          id
          type
          format
          storageLocation {
            bucket
            key
            url
          }
          width
          height
          duration
          bitrate
          codec
          size
          contentType
          streamInfo {
            downloadUrl
            cdnPrefix
            url
          }
          metadata
          status
          thumbnail {
            fileName
            url
            storageLocation {
              bucket
              key
              url
            }
            width
            height
            size
            contentType
            metadata
          }
          transcodingInfo {
            jobId
            progress
            outputUrl
            error
            completedAt
          }
          createdAt
          updatedAt
          quality
          isReady
          isProcessing
          isFailed
        }
        images {
          id
          fileName
          url
          type
          storageLocation {
            bucket
            key
            url
          }
          width
          height
          size
          contentType
          metadata
          streamInfo {
            downloadUrl
            cdnPrefix
            url
          }
          createdAt
          updatedAt
        }
      }
      nextKey
    }
  }
`;

const GET_BUCKETS = gql`
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
          createdAt
          updatedAt
          metadata
          ownerId
          videos {
            id
            type
            format
            storageLocation {
              bucket
              key
              url
            }
            width
            height
            duration
            bitrate
            codec
            size
            contentType
            streamInfo {
              downloadUrl
              cdnPrefix
              url
            }
            metadata
            status
            thumbnail {
              fileName
              url
              storageLocation {
                bucket
                key
                url
              }
              width
              height
              size
              contentType
              metadata
            }
            transcodingInfo {
              jobId
              progress
              outputUrl
              error
              completedAt
            }
            createdAt
            updatedAt
            quality
            isReady
            isProcessing
            isFailed
          }
          images {
            id
            fileName
            url
            type
            width
            height
            size
            contentType
            metadata
            streamInfo {
              downloadUrl
              cdnPrefix
              url
            }
            createdAt
            updatedAt
          }
        }
      }
      nextKey
    }
  }
`;

const GET_BUCKET = gql`
  query GetBucket($id: ID!) {
    bucket(id: $id) {
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
        createdAt
        updatedAt
        metadata
        ownerId
        videos {
          id
          type
          format
          storageLocation {
            bucket
            key
            url
          }
          width
          height
          duration
          bitrate
          codec
          size
          contentType
          streamInfo {
            downloadUrl
            cdnPrefix
            url
          }
          metadata
          status
          thumbnail {
            fileName
            url
            storageLocation {
              bucket
              key
              url
            }
            width
            height
            size
            contentType
            metadata
          }
          transcodingInfo {
            jobId
            progress
            outputUrl
            error
            completedAt
          }
          createdAt
          updatedAt
          quality
          isReady
          isProcessing
          isFailed
        }
        images {
          id
          fileName
          url
          type
          width
          height
          size
          contentType
          metadata
          streamInfo {
            downloadUrl
            cdnPrefix
            url
          }
          createdAt
          updatedAt
        }
      }
    }
  }
`;

const GET_ASSETS_BY_PARENT = gql`
  query GetAssetsByParent($parentId: ID!, $limit: Int, $nextKey: String) {
    assetsByParent(parentId: $parentId, limit: $limit, nextKey: $nextKey) {
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
        parent {
          id
          title
          type
        }
        publishRule {
          publishAt
          unpublishAt
          regions
          ageRating
        }
        videos {
          id
          type
          format
          storageLocation {
            bucket
            key
            url
          }
          width
          height
          duration
          bitrate
          codec
          size
          contentType
          streamInfo {
            downloadUrl
            cdnPrefix
            url
          }
          metadata
          status
          thumbnail {
            fileName
            url
            storageLocation {
              bucket
              key
              url
            }
            width
            height
            size
            contentType
            metadata
          }
          transcodingInfo {
            jobId
            progress
            outputUrl
            error
            completedAt
          }
          createdAt
          updatedAt
          quality
          isReady
          isProcessing
          isFailed
        }
      }
      nextKey
    }
  }
`;

const GET_ASSET = gql`
  query GetAsset($id: ID!) {
    asset(id: $id) {
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
      parent {
        id
        title
        type
      }
      publishRule {
        publishAt
        unpublishAt
        regions
        ageRating
      }
      images {
        id
        fileName
        url
        type
        storageLocation {
          bucket
          key
          url
        }
        width
        height
        size
        contentType
        metadata
        createdAt
        updatedAt
      }
      videos {
        id
        type
        format
        storageLocation {
          bucket
          key
          url
        }
        width
        height
        duration
        bitrate
        codec
        size
        contentType
        streamInfo {
          downloadUrl
          cdnPrefix
          url
        }
        metadata
        status
        thumbnail {
          fileName
          url
          storageLocation {
            bucket
            key
            url
          }
          width
          height
          size
          contentType
          metadata
        }
        transcodingInfo {
          jobId
          progress
          outputUrl
          error
          completedAt
        }
        createdAt
        updatedAt
        quality
        isReady
        isProcessing
        isFailed
      }
    }
  }
`;

const SEARCH_ASSETS = gql`
  query SearchAssets($query: String!, $limit: Int) {
    searchAssets(query: $query, limit: $limit) {
      items {
        id
        slug
        title
        type
      }
    }
  }
`;


const CREATE_ASSET = gql`
  mutation CreateAsset($input: CreateAssetInput!) {
    createAsset(input: $input) {
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
      parent {
        id
        title
        type
      }
      publishRule {
        publishAt
        unpublishAt
        regions
        ageRating
      }
      videos {
        id
        type
        format
        storageLocation {
          bucket
          key
          url
        }
        width
        height
        duration
        bitrate
        codec
        size
        contentType
        streamInfo {
          downloadUrl
          cdnPrefix
          url
        }
        metadata
        status
        thumbnail {
          fileName
          url
          storageLocation {
            bucket
            key
            url
          }
          width
          height
          size
          contentType
          metadata
        }
        transcodingInfo {
          jobId
          progress
          outputUrl
          error
          completedAt
        }
        createdAt
        updatedAt
        quality
        isReady
        isProcessing
        isFailed
      }
    }
  }
`;

const PATCH_ASSET = gql`
  mutation PatchAsset($id: ID!, $patches: [JSONPatch!]!) {
    patchAsset(id: $id, patches: $patches) {
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
      parent {
        id
        title
        type
      }
      publishRule {
        publishAt
        unpublishAt
        regions
        ageRating
      }
      videos {
        id
        type
        format
        storageLocation {
          bucket
          key
          url
        }
        width
        height
        duration
        bitrate
        codec
        size
        contentType
        streamInfo {
          downloadUrl
          cdnPrefix
          url
        }
        metadata
        status
        thumbnail {
          fileName
          url
          storageLocation {
            bucket
            key
            url
          }
          width
          height
          size
          contentType
          metadata
        }
        transcodingInfo {
          jobId
          progress
          outputUrl
          error
          completedAt
        }
        createdAt
        updatedAt
        quality
        isReady
        isProcessing
        isFailed
      }
    }
  }
`;

const DELETE_ASSET = gql`
  mutation DeleteAsset($id: ID!) {
    deleteAsset(id: $id)
  }
`;

const DELETE_VIDEO = gql`
  mutation DeleteVideo($assetId: ID!, $videoId: ID!) {
    deleteVideo(assetId: $assetId, videoId: $videoId) {
      id
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
          type
          format
          storageLocation {
            bucket
            key
            url
          }
          width
          height
          duration
          bitrate
          codec
          size
          contentType
          streamInfo {
            downloadUrl
            cdnPrefix
            url
          }
          metadata
          status
          thumbnail {
            fileName
            url
            storageLocation {
              bucket
              key
              url
            }
            width
            height
            size
            contentType
            metadata
          }
          transcodingInfo {
            jobId
            progress
            outputUrl
            error
            completedAt
          }
          createdAt
          updatedAt
          quality
          isReady
          isProcessing
          isFailed
        }
    }
  }
`;



const ADD_VIDEO = gql`
  mutation AddVideo($input: AddVideoInput!) {
    addVideo(input: $input) {
      id
      label
      type
      format
      storageLocation {
        bucket
        key
        url
      }
      width
      height
      duration
      bitrate
      codec
      size
      contentType
      streamInfo {
        downloadUrl
        cdnPrefix
        url
      }
      metadata
      status
      thumbnail {
        id
        fileName
        url
        type
        width
        height
        size
        contentType
        metadata
        createdAt
        updatedAt
      }
      transcodingInfo {
        jobId
        progress
        outputUrl
        error
        completedAt
      }
      createdAt
      updatedAt
      quality
      isReady
      isProcessing
      isFailed
    }
  }
`;

const CREATE_BUCKET = gql`
  mutation CreateBucket($input: CreateBucketInput!) {
    createBucket(input: $input) {
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
        createdAt
        updatedAt
        metadata
        ownerId
      }
    }
  }
`;

const UPDATE_BUCKET = gql`
  mutation UpdateBucket($input: UpdateBucketInput!) {
    updateBucket(input: $input) {
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
        createdAt
        updatedAt
        metadata
        ownerId
      }
    }
  }
`;

const DELETE_BUCKET = gql`
  mutation DeleteBucket($id: ID!, $ownerId: String!) {
    deleteBucket(input: { id: $id, ownerId: $ownerId })
  }
`;

const ADD_ASSET_TO_BUCKET = gql`
  mutation AddAssetToBucket($bucketId: ID!, $assetId: ID!, $ownerId: String!) {
    addAssetToBucket(input: { bucketId: $bucketId, assetId: $assetId, ownerId: $ownerId })
  }
`;

const REMOVE_ASSET_FROM_BUCKET = gql`
  mutation RemoveAssetFromBucket($bucketId: ID!, $assetId: ID!) {
    removeAssetFromBucket(bucketId: $bucketId, assetId: $assetId) {
      id
      key
      name
      description
      type
      status
      assetIds
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
        createdAt
        updatedAt
        metadata
        ownerId
      }
    }
  }
`;

const ADD_IMAGE = gql`
  mutation AddImage($input: AddImageInput!) {
    addImage(input: $input) {
      id
      images {
        id
        fileName
        url
        type
        width
        height
        size
        contentType
        metadata
        createdAt
        updatedAt
        streamInfo {
          downloadUrl
          cdnPrefix
          url
        }
      }
      updatedAt
    }
  }
`;

const DELETE_IMAGE = gql`
  mutation DeleteImage($assetId: ID!, $imageId: ID!) {
    deleteImage(assetId: $assetId, imageId: $imageId) {
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
      images {
        id
        fileName
        url
        type
        width
        height
        size
        contentType
        metadata
        createdAt
        updatedAt
          streamInfo {
            downloadUrl
            cdnPrefix
            url
          }
        }
    }
  }
`;

// Helper function to parse metadata JSON string to object
const parseMetadata = (metadataString?: string): Record<string, any> | undefined => {
  if (!metadataString) return undefined;
  try {
    return JSON.parse(metadataString);
  } catch (error) {
    console.warn('Failed to parse metadata JSON:', error);
    return undefined;
  }
};

// Helper function to convert Asset with string metadata to Asset with object metadata
const parseImages = (imagesData?: any): any[] => {
  if (!imagesData) return [];
  let arr = imagesData;
  if (typeof imagesData === 'string') {
    try {
      arr = JSON.parse(imagesData);
    } catch {
      return [];
    }
  }
  if (!Array.isArray(arr)) return [];
  return arr
    .filter(img => img && (img.ImageID || img.id))
    .map(img => ({
      id: img.ImageID || img.id,
      fileName: img.ImageFileName || img.fileName,
      url: img.ImageURL || img.url,
      type: img.ImageType || img.type,
      storageLocation: img.ImageStorageLoc || img.storageLocation,
      width: img.ImageWidth || img.width,
      height: img.ImageHeight || img.height,
      size: img.ImageSize || img.size,
      contentType: img.ImageContentType || img.contentType,
      metadata: img.ImageMetadata || img.metadata,
      createdAt: img.ImageCreatedAt || img.createdAt,
      updatedAt: img.ImageUpdatedAt || img.updatedAt,
      streamInfo: img.ImageStreamInfo || img.streamInfo,
    }));
};

const convertAssetMetadata = (asset: any): Asset => {
  console.log('convertAssetMetadata called with asset:', asset);
  console.log('asset.images before parsing:', asset.images);
  
  const converted = {
    ...asset,
    metadata: parseMetadata(asset.metadata),
    images: parseImages(asset.images),
  };
  
  console.log('converted asset.images after parsing:', converted.images);
  return converted;
};

// Custom hook for asset service
export const useAssetService = () => {
  const client = useApolloClient();

  return {

    getAssets: async (page = 1, limit = 20): Promise<{ assets: Asset[]; total: number }> => {
      return retryWithBackoff(async () => {
        try {
          const response = await client.query({
            query: GET_ASSETS,
            variables: { limit },
            fetchPolicy: 'no-cache',
          });
          
          console.log('Raw GraphQL response:', response.data);
          const rawAssets = response.data.assets.items || [];
          console.log('Raw assets before conversion:', rawAssets);
          const assets = rawAssets.map(convertAssetMetadata);
          return {
            assets,
            total: assets.length
          };
        } catch (error: any) {
          if (error.graphQLErrors?.some((e: any) => e.extensions?.code === 'UNAUTHENTICATED')) {
            console.log('Token expired or invalid, logging out...');
            clearAuthToken();
            if (logoutCallback) {
              logoutCallback();
            }
          }
          throw error;
        }
      });
    },


    getAsset: async (id: string): Promise<Asset> => {
      const response = await client.query({
        query: GET_ASSET,
        variables: { id },
        fetchPolicy: 'no-cache',
      });
      return convertAssetMetadata(response.data.asset);
    },


    createAsset: async (assetData: AssetCreateDTO): Promise<Asset> => {
      if (!assetData.slug) {
        throw new Error('Slug is required');
      }

      const input: AssetInput = {
        title: assetData.title,
        slug: assetData.slug,
        description: assetData.description,
        type: assetData.type,
        genre: assetData.genre,
        genres: assetData.genres,
        tags: assetData.tags,
        metadata: assetData.metadata ? JSON.stringify(assetData.metadata) : undefined,
        ownerId: assetData.ownerId,
        parentId: assetData.parentId,
      };

      console.log('AssetInput being sent to GraphQL:', input);
      console.log('Slug in input:', input.slug);
      console.log('Slug type in input:', typeof input.slug);

      const response = await client.mutate({
        mutation: CREATE_ASSET,
        variables: { input },
      });
      return convertAssetMetadata(response.data.createAsset);
    },


    updateAsset: async (id: string, assetData: AssetUpdateDTO, clearFields: string[] = []): Promise<Asset> => {
      const patches: any[] = [];
      
      if (assetData.title !== undefined) {
        if (assetData.title === null || assetData.title === '') {
          patches.push({ op: 'remove', path: '/title' });
        } else {
          patches.push({ op: 'replace', path: '/title', value: assetData.title });
        }
      }
      
      if (assetData.description !== undefined) {
        if (assetData.description === null || assetData.description === '') {
          patches.push({ op: 'remove', path: '/description' });
        } else {
          patches.push({ op: 'replace', path: '/description', value: assetData.description });
        }
      }

      if (assetData.type !== undefined) {
        if (assetData.type === null) {
          patches.push({ op: 'remove', path: '/type' });
        } else {
          patches.push({ op: 'replace', path: '/type', value: assetData.type });
        }
      }
      
      if (assetData.genre !== undefined) {
        if (assetData.genre === null || assetData.genre === '' || assetData.genre.includes('undefined')) {
          patches.push({ op: 'remove', path: '/genre' });
        } else {
          patches.push({ op: 'replace', path: '/genre', value: assetData.genre.toUpperCase() });
        }
      }
      
      if (assetData.genres !== undefined) {
        if (assetData.genres === null || assetData.genres.length === 0) {
          patches.push({ op: 'remove', path: '/genres' });
        } else {
          patches.push({ op: 'replace', path: '/genres', value: JSON.stringify(assetData.genres) });
        }
      }
      
      if (assetData.tags !== undefined) {
        if (assetData.tags === null || assetData.tags.length === 0) {
          patches.push({ op: 'remove', path: '/tags' });
        } else {
          patches.push({ op: 'replace', path: '/tags', value: JSON.stringify(assetData.tags) });
        }
      }
      
      if (assetData.metadata !== undefined) {
        if (assetData.metadata === null || Object.keys(assetData.metadata).length === 0) {
          patches.push({ op: 'remove', path: '/metadata' });
        } else {
          patches.push({ op: 'replace', path: '/metadata', value: JSON.stringify(assetData.metadata) });
        }
      }
      
      if (assetData.ownerId !== undefined) {
        if (assetData.ownerId === null || assetData.ownerId === '') {
          patches.push({ op: 'remove', path: '/ownerId' });
        } else {
          patches.push({ op: 'replace', path: '/ownerId', value: assetData.ownerId });
        }
      }
      
      if (assetData.parentId !== undefined) {
        if (assetData.parentId === null || assetData.parentId === '') {
          patches.push({ op: 'remove', path: '/parentId' });
        } else {
          patches.push({ op: 'replace', path: '/parentId', value: assetData.parentId });
        }
      }
      
      console.log('Sending patches to backend:', { id, patches, assetData });
      
      const response = await client.mutate({
        mutation: PATCH_ASSET,
        variables: { id, patches },
      });
      return convertAssetMetadata(response.data.patchAsset);
    },


    deleteAsset: async (id: string): Promise<void> => {
      await client.mutate({
        mutation: DELETE_ASSET,
        variables: { id },
      });
    },


    deleteVideo: async (assetId: string, videoId: string): Promise<Asset> => {
      const response = await client.mutate({
        mutation: DELETE_VIDEO,
        variables: { assetId, videoId },
      });
      return convertAssetMetadata(response.data.deleteVideo);
    },


    publishAsset: async (id: string, publishAt?: string | null, unpublishAt?: string | null, regions?: string[], ageRating?: string | null, clearFields: string[] = []): Promise<Asset> => {
      const patches: any[] = [];
      if (publishAt !== undefined) {
        patches.push({ op: 'replace', path: '/publishAt', value: publishAt ? new Date(publishAt).toISOString() : '' });
      }
      if (unpublishAt !== undefined) {
        patches.push({ op: 'replace', path: '/unpublishAt', value: unpublishAt ? new Date(unpublishAt).toISOString() : '' });
      }
      if (regions !== undefined) {
        patches.push({ op: 'replace', path: '/regions', value: JSON.stringify(regions) });
      }
      if (ageRating !== undefined) {
        patches.push({ op: 'replace', path: '/ageRating', value: ageRating || '' });
      }
      
      // Add remove operations for cleared fields
      clearFields.forEach(field => {
        patches.push({ op: 'remove', path: `/${field}` });
      });
      
      const response = await client.mutate({
        mutation: PATCH_ASSET,
        variables: { id, patches },
      });
      return convertAssetMetadata(response.data.patchAsset);
    },





    addVideo: async (assetId: string, videoType: string, format: string, bucket: string, key: string, url: string, contentType: string, size: number): Promise<any> => {
      const input = {
        assetId,
        label: videoType,
        format,
        bucket,
        key,
        url,
        contentType,
        size,
      };
      
      console.log('Sending addVideo input:', input);
      console.log('Input field values:', {
        assetId_type: typeof assetId, assetId_value: assetId,
        label_type: typeof videoType, label_value: videoType,
        bucket_type: typeof bucket, bucket_value: bucket,
        key_type: typeof key, key_value: key,
        url_type: typeof url, url_value: url,
        contentType_type: typeof contentType, contentType_value: contentType,
        size_type: typeof size, size_value: size,
      });
      
      const response = await client.mutate({
        mutation: ADD_VIDEO,
        variables: { input },
      });
      return response.data.addVideo;
    },


    searchAssets: async (query: string, limit = 10): Promise<{ assets: any[] }> => {
      const response = await client.query({
        query: SEARCH_ASSETS,
        variables: { query, limit },
        fetchPolicy: 'no-cache',
      });
      return { assets: response.data.searchAssets.items };
    },


    getAssetsByParent: async (parentId: string, limit = 20): Promise<{ assets: Asset[]; total: number }> => {
      try {
        const response = await client.query({
          query: GET_ASSETS_BY_PARENT,
          variables: { parentId, limit },
          fetchPolicy: 'no-cache',
        });
        
        const rawAssets = response.data.assetsByParent.items || [];
        const assets = rawAssets.map(convertAssetMetadata);
        return {
          assets,
          total: assets.length
        };
      } catch (error: any) {
        if (error.graphQLErrors?.some((e: any) => e.extensions?.code === 'UNAUTHENTICATED')) {
          console.log('Token expired or invalid, logging out...');
          clearAuthTokens();
          if (logoutCallback) {
            logoutCallback();
          }
        }
        throw error;
      }
    },


    getUploadUrl: async (fileName: string, assetId: string, videoType: string): Promise<{ url: string }> => {
      try {
        const response = await axios.post(`${API_CONFIG.API_GATEWAY_BASE_URL}/upload`, {
          fileName: fileName,
          assetId: assetId,
          videoType: videoType
        }, {
          headers: {
            'Content-Type': 'application/json',
          },
        });
        
        const presignedUrl = response.data.url;
        const localhostUrl = presignedUrl.replace('localstack:4566', 'localhost:4566');
        
        return { url: localhostUrl };
      } catch (error) {
        console.error('Error getting upload URL:', error);
        throw new Error('Failed to get upload URL');
      }
    },

    uploadFile: async (uploadUrl: string, file: File | Blob): Promise<void> => {
      await axios.put(uploadUrl, file, {
        headers: {
          'Content-Type': file instanceof File ? file.type : 'application/octet-stream',
        },
      });
    },


    getBuckets: async (limit = 20): Promise<{ buckets: any[]; total: number }> => {
      return retryWithBackoff(async () => {
        try {
          const response = await client.query({
            query: GET_BUCKETS,
            variables: { limit },
            fetchPolicy: 'no-cache',
          });
          
          const buckets = response.data.buckets.items || [];
          return {
            buckets,
            total: buckets.length
          };
        } catch (error: any) {
          if (error.graphQLErrors?.some((e: any) => e.extensions?.code === 'UNAUTHENTICATED')) {
            console.log('Token expired or invalid, logging out...');
            clearAuthTokens();
            if (logoutCallback) {
              logoutCallback();
            }
          }
          throw error;
        }
      });
    },

    getBucket: async (id: string): Promise<any> => {
      try {
        const response = await client.query({
          query: GET_BUCKET,
          variables: { id },
          fetchPolicy: 'no-cache',
        });
        return response.data.bucket;
      } catch (error: any) {
        if (error.graphQLErrors?.some((e: any) => e.extensions?.code === 'UNAUTHENTICATED')) {
          console.log('Token expired or invalid, logging out...');
          clearAuthTokens();
          if (logoutCallback) {
            logoutCallback();
          }
        }
        throw error;
      }
    },

    createBucket: async ({ key, name, description, type, ownerId, metadata, status }: { key: string, name: string, description?: string, type: string, ownerId?: string, metadata?: string, status?: string }): Promise<any> => {
      const input = {
        key: key,
        name: name,
        description: description || '',
        type: type,
        status: status || BucketStatus.DRAFT,
        ownerId: ownerId,
        metadata: metadata ? JSON.stringify(metadata) : undefined,
      };

      const response = await client.mutate({
        mutation: CREATE_BUCKET,
        variables: { input },
      });
      return response.data.createBucket;
    },

    updateBucket: async (id: string, bucketData: any): Promise<any> => {
      const input = {
        name: bucketData.name,
        description: bucketData.description || '',
        status: bucketData.status || BucketStatus.DRAFT,
      };

      const response = await client.mutate({
        mutation: UPDATE_BUCKET,
        variables: { id, input },
      });
      return response.data.updateBucket;
    },

    patchBucket: async (id: string, patches: any[]): Promise<any> => {
      const input: any = { id };
      
      patches.forEach(patch => {
        if (patch.op === 'replace') {
          const field = patch.path.substring(1);
          input[field] = patch.value;
        }
      });

      const response = await client.mutate({
        mutation: UPDATE_BUCKET,
        variables: { input },
      });
      return response.data.updateBucket;
    },

    deleteBucket: async (id: string, ownerId: string): Promise<void> => {
      await client.mutate({
        mutation: DELETE_BUCKET,
        variables: { id, ownerId },
      });
    },

    addAssetToBucket: async (bucketId: string, assetId: string, ownerId: string): Promise<boolean> => {
      const response = await client.mutate({
        mutation: ADD_ASSET_TO_BUCKET,
        variables: { bucketId, assetId, ownerId },
      });
      return response.data.addAssetToBucket;
    },

    removeAssetFromBucket: async (bucketId: string, assetId: string): Promise<any> => {
      const response = await client.mutate({
        mutation: REMOVE_ASSET_FROM_BUCKET,
        variables: { bucketId, assetId },
      });
      return response.data.removeAssetFromBucket;
    },

    getImageUploadUrl: async (fileName: string, assetId: string, imageType: ImageType): Promise<{ url: string }> => {
      try {
        const response = await axios.post(`${API_CONFIG.API_GATEWAY_BASE_URL}/image-upload`, {
          fileName: fileName,
          assetId: assetId,
          imageType: imageType
        }, {
          headers: {
            'Content-Type': 'application/json',
          },
        });
        
        const presignedUrl = response.data.url;
        const localhostUrl = presignedUrl.replace('localstack:4566', 'localhost:4566');
        
        return { url: localhostUrl };
      } catch (error) {
        console.error('Error getting image upload URL:', error);
        throw new Error('Failed to get image upload URL');
      }
    },

    addImageToAsset: async (assetId: string, imageData: Partial<Image>): Promise<Asset> => {
      console.log('addImageToAsset called with:', { assetId, imageData });
      const input = {
        assetId,
        type: imageData.type,
        fileName: imageData.fileName,
        bucket: 'content-east',
        key: `${assetId}/images/${imageData.type?.toLowerCase()}/${imageData.fileName}`,
        url: imageData.url,
        contentType: 'image/jpeg',
        size: (typeof imageData.size === 'number' && imageData.size > 0) ? imageData.size : 1,
      };
      console.log('GraphQL input:', input);
      
      const response = await client.mutate({
        mutation: ADD_IMAGE,
        variables: { input },
      });
      
      console.log('GraphQL response:', response.data);
      return response.data.addImage;
    },

    deleteImageFromAsset: async (assetId: string, imageId: string): Promise<void> => {
      await client.mutate({
        mutation: DELETE_IMAGE,
        variables: { assetId, imageId },
      });
    },
  };
};

export const assetService = {
  getAssets: async (page = 1, limit = 20): Promise<{ assets: Asset[]; total: number }> => {
    throw new Error('Use useAssetService hook instead');
  },
  getAsset: async (id: string): Promise<Asset> => {
    throw new Error('Use useAssetService hook instead');
  },
  createAsset: async (assetData: AssetCreateDTO): Promise<Asset> => {
    throw new Error('Use useAssetService hook instead');
  },
  updateAsset: async (id: string, assetData: AssetUpdateDTO): Promise<Asset> => {
    throw new Error('Use useAssetService hook instead');
  },
  deleteAsset: async (id: string): Promise<void> => {
    throw new Error('Use useAssetService hook instead');
  },
  getUploadUrl: async (fileName: string, contentType: string): Promise<any> => {
    throw new Error('Use useAssetService hook instead');
  },
  uploadFile: async (uploadUrl: string, file: any): Promise<void> => {
    throw new Error('Use useAssetService hook instead');
  },
};

export const authService = {
  login: async (username: string, password: string): Promise<{ accessToken: string; refreshToken: string }> => {
    const response = await authApi.post('/login', {
      username,
      password,
      client_id: 'asset-manager',
    });
    return { 
      accessToken: response.data.access_token,
      refreshToken: response.data.refresh_token
    };
  },


  refreshToken: async (refreshToken: string): Promise<{ accessToken: string; refreshToken: string }> => {
    const response = await authApi.post('/refresh', {
      refresh_token: refreshToken,
    });
    return { 
      accessToken: response.data.access_token,
      refreshToken: response.data.refresh_token
    };
  },


  validateToken: async (token: string): Promise<{ valid: boolean; user?: any }> => {
    try {
      console.log('validateToken called with token:', token ? `${token.substring(0, 20)}...` : 'null');
      const response = await authApi.post('/validate', { token });
      console.log('validateToken response:', response.data);
      return response.data;
    } catch (error) {
      console.log('validateToken error:', error);
      return { valid: false };
    }
  },
}; 