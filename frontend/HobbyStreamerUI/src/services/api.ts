import { Asset, Bucket, AssetResponse, BucketResponse } from '../types/asset';

const API_BASE_URL = 'http://localhost:8084/api/v1';

class ApiService {
  private async fetch<T>(endpoint: string): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`);
    
    if (!response.ok) {
      throw new Error(`API request failed: ${response.status} ${response.statusText}`);
    }
    
    return response.json();
  }

  async getBuckets(limit = 10, nextKey?: string): Promise<{ buckets: Bucket[], nextKey?: string, count: number, limit: number }> {
    let url = `/buckets?limit=${limit}`;
    if (nextKey) url += `&nextKey=${encodeURIComponent(nextKey)}`;
    return this.fetch(url);
  }

  async getBucket(key: string): Promise<Bucket> {
    return this.fetch(`/buckets/${key}`);
  }

  async getAssets(): Promise<Asset[]> {
    const response: AssetResponse = await this.fetch('/assets');
    return response.assets;
  }

  async getAsset(slug: string): Promise<Asset> {
    return this.fetch(`/assets/${slug}`);
  }

  async getAssetsInBucket(bucketKey: string): Promise<Asset[]> {
    const response: AssetResponse = await this.fetch(`/buckets/${bucketKey}/assets`);
    return response.assets;
  }
}

export const apiService = new ApiService(); 