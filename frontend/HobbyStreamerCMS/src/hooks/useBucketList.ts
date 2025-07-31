import { useState, useEffect } from 'react';
import { Alert } from 'react-native';
import { useAssetService } from '../services/api';
import { Bucket } from '../types/asset';
import { getAuthToken, validateTokenLocally } from '../services/api';

export function useBucketList(refreshTrigger?: number) {
  const assetService = useAssetService();
  const [buckets, setBuckets] = useState<Bucket[]>([]);
  const [selectedBucket, setSelectedBucket] = useState<Bucket | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showSuccessMessage, setShowSuccessMessage] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [updating, setUpdating] = useState(false);
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    loadBuckets();
    if (refreshTrigger && refreshTrigger > 0) {
      setShowSuccessMessage(true);
      setTimeout(() => setShowSuccessMessage(false), 3000);
    }
  }, [refreshTrigger]);

  const loadBuckets = async (isRefresh = false) => {
    try {
      if (isRefresh) {
        setRefreshing(true);
      } else {
        setLoading(true);
      }
      
      const response = await assetService.getBuckets();
      setBuckets(response.buckets);
      setError(null);
    } catch (err) {
      console.error('Error loading buckets:', err);
      setError('Failed to load buckets. Make sure the backend services are running.');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  const handleRefresh = async () => {
    await loadBuckets(true);
  };

  const handleBucketSelect = async (bucket: Bucket) => {
    setSelectedBucket(bucket);
  };

  const handleDeleteBucket = async () => {
    if (!selectedBucket) return;

    try {
      setDeleting(true);
      const token = await getAuthToken();
      let ownerId = '';
      if (token) {
        const { valid, user } = validateTokenLocally(token);
        if (valid && user && user.id) {
          ownerId = user.id;
        }
      }
      if (!ownerId) {
        Alert.alert('Error', 'Could not determine logged-in user.');
        setDeleting(false);
        return;
      }
      await assetService.deleteBucket(selectedBucket.id, ownerId);
      setSelectedBucket(null);
      await loadBuckets();
      setShowSuccessMessage(true);
      setTimeout(() => setShowSuccessMessage(false), 3000);
    } catch (err: any) {
      console.error('Error deleting bucket:', err);
      Alert.alert('Error', `Failed to delete bucket: ${err.message || 'Unknown error'}`);
    } finally {
      setDeleting(false);
    }
  };

  const handleUpdateBucket = async (field: string, value: any) => {
    if (!selectedBucket) return;

    try {
      setUpdating(true);
      
      const patches = [];
      
      if (field === 'name') {
        patches.push({ op: 'replace', path: '/name', value });
      } else if (field === 'description') {
        patches.push({ op: 'replace', path: '/description', value });
      } else if (field === 'type') {
        patches.push({ op: 'replace', path: '/type', value });
      } else if (field === 'status') {
        patches.push({ op: 'replace', path: '/status', value });
      }

      const updatedBucket = await assetService.patchBucket(selectedBucket.id, patches);
      setSelectedBucket(updatedBucket);
      
      setBuckets(prevBuckets => 
        prevBuckets.map(bucket => 
          bucket.id === selectedBucket.id ? updatedBucket : bucket
        )
      );
      
      setShowSuccessMessage(true);
      setTimeout(() => setShowSuccessMessage(false), 3000);
    } catch (err) {
      console.error('Error updating bucket:', err);
      Alert.alert('Error', 'Failed to update bucket. Please try again.');
    } finally {
      setUpdating(false);
    }
  };

  const handleCreateBucket = async (bucketData: any) => {
    try {
      setCreating(true);
      const newBucket = await assetService.createBucket(bucketData);
      setSelectedBucket(newBucket);
      await loadBuckets();
      setShowSuccessMessage(true);
      setTimeout(() => setShowSuccessMessage(false), 3000);
      return newBucket;
    } catch (err: any) {
      console.error('Error creating bucket:', err);
      Alert.alert('Error', `Failed to create bucket: ${err.message || 'Unknown error'}`);
      throw err;
    } finally {
      setCreating(false);
    }
  };

  const handleAddAssetToBucket = async (assetId: string) => {
    if (!selectedBucket) return;

    try {
      setUpdating(true);
      const token = await getAuthToken();
      let ownerId = '';
      if (token) {
        const { valid, user } = validateTokenLocally(token);
        if (valid && user && user.id) {
          ownerId = user.id;
        }
      }
      if (!ownerId) {
        Alert.alert('Error', 'Could not determine logged-in user.');
        setUpdating(false);
        return;
      }
      const success = await assetService.addAssetToBucket(selectedBucket.id, assetId, ownerId);
      if (success) {       
        const updatedBucket = await assetService.getBucket(selectedBucket.id);
        setSelectedBucket(updatedBucket);
        setShowSuccessMessage(true);
        setTimeout(() => setShowSuccessMessage(false), 3000);
      } else {
        Alert.alert('Error', 'Failed to add asset to bucket. Please try again.');
      }
    } catch (err) {
      console.error('Error adding asset to bucket:', err);
      Alert.alert('Error', 'Failed to add asset to bucket. Please try again.');
    } finally {
      setUpdating(false);
    }
  };

  const handleRemoveAssetFromBucket = async (assetId: string) => {
    if (!selectedBucket) return;

    try {
      setUpdating(true);
      const token = await getAuthToken();
      let ownerId = '';
      if (token) {
        const { valid, user } = validateTokenLocally(token);
        if (valid && user && user.id) {
          ownerId = user.id;
        }
      }
      if (!ownerId) {
        Alert.alert('Error', 'Could not determine logged-in user.');
        setUpdating(false);
        return;
      }
      
      const success = await assetService.removeAssetFromBucket(selectedBucket.id, assetId, ownerId);
      if (success) {
        const updatedBucket = await assetService.getBucket(selectedBucket.id);
        setSelectedBucket(updatedBucket);
        
        setBuckets(prevBuckets => 
          prevBuckets.map(bucket => 
            bucket.id === selectedBucket.id ? updatedBucket : bucket
          )
        );
        
        setShowSuccessMessage(true);
        setTimeout(() => setShowSuccessMessage(false), 3000);
      } else {
        Alert.alert('Error', 'Failed to remove asset from bucket. Please try again.');
      }
    } catch (err) {
      console.error('Error removing asset from bucket:', err);
      Alert.alert('Error', 'Failed to remove asset from bucket. Please try again.');
    } finally {
      setUpdating(false);
    }
  };

  return {
    buckets,
    selectedBucket,
    loading,
    refreshing,
    error,
    showSuccessMessage,
    deleting,
    updating,
    creating,
    handleRefresh,
    handleBucketSelect,
    handleDeleteBucket,
    handleUpdateBucket,
    handleCreateBucket,
    handleAddAssetToBucket,
    handleRemoveAssetFromBucket,
  };
} 