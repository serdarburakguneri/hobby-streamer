import { useState, useEffect } from 'react';
import { Alert } from 'react-native';
import { useAssetService } from '../services/api';
import { Asset, AssetType } from '../types/asset';

export function useAssetList(refreshTrigger?: number) {
  const assetService = useAssetService();
  const [assets, setAssets] = useState<Asset[]>([]);
  const [selectedAsset, setSelectedAsset] = useState<Asset | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showSuccessMessage, setShowSuccessMessage] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [publishing, setPublishing] = useState(false);
  const [updating, setUpdating] = useState(false);
  const [children, setChildren] = useState<Asset[]>([]);
  const [childrenLoading, setChildrenLoading] = useState(false);

  useEffect(() => {
    const fetchAndSetFirstAsset = async () => {
      await loadAssets();
      if (assets.length > 0) {
        const fullAsset = await assetService.getAsset(assets[0].id);
        setSelectedAsset(fullAsset);
      }
    };
    fetchAndSetFirstAsset();
    if (refreshTrigger && refreshTrigger > 0) {
      setShowSuccessMessage(true);
      setTimeout(() => setShowSuccessMessage(false), 3000);
    }
  }, [refreshTrigger]);

  const loadAssets = async (isRefresh = false) => {
    try {
      if (isRefresh) {
        setRefreshing(true);
      } else {
        setLoading(true);
      }
      
      const response = await assetService.getAssets();
      setAssets(response.assets);
      setError(null);
    } catch (err) {
      console.error('Error loading assets:', err);
      setError('Failed to load assets. Make sure the backend services are running.');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  const handleRefresh = async () => {
    await loadAssets(true);
  };

  const handleRefreshSelectedAsset = async () => {
    if (selectedAsset) {
      const updated = await assetService.getAsset(selectedAsset.id);
      setSelectedAsset(updated);
    }
  };

  const handleAssetSelect = async (asset: Asset) => {
    const fullAsset = await assetService.getAsset(asset.id);
    setSelectedAsset(fullAsset);
    setChildren([]);
    if (fullAsset.type === AssetType.SERIES || fullAsset.type === AssetType.SEASON) {
      await loadChildren(fullAsset.id);
    }
  };

  const loadChildren = async (parentId: string) => {
    try {
      setChildrenLoading(true);
      const response = await assetService.getAssetsByParent(parentId);
      setChildren(response.assets);
    } catch (err) {
      console.error('Error loading children:', err);
      Alert.alert('Error', 'Failed to load children. Please try again.');
    } finally {
      setChildrenLoading(false);
    }
  };

  const handleDeleteAsset = async () => {
    if (!selectedAsset) return;
    
    try {
      setDeleting(true);
      await assetService.deleteAsset(selectedAsset.id);
      setSelectedAsset(null);
      await loadAssets();
      setShowSuccessMessage(true);
      setTimeout(() => setShowSuccessMessage(false), 3000);
    } catch (err: any) {
      console.error('Error deleting asset:', err);
      Alert.alert('Error', `Failed to delete asset: ${err.message || 'Unknown error'}`);
    } finally {
      setDeleting(false);
    }
  };



  const handleUpdateAsset = async (field: string, value: any) => {
    if (!selectedAsset) return;

    try {
      setUpdating(true);
      const updateData: any = {};

      if (field === 'title') {
        updateData.title = value;
      } else if (field === 'description') {
        updateData.description = value;
      } else if (field === 'primaryGenre') {
        updateData.genre = value;
      } else if (field === 'additionalGenres') {
        updateData.genres = value;
      } else if (field === 'tags') {
        updateData.tags = value;
      }

      const updatedAsset = await assetService.updateAsset(selectedAsset.id, updateData);
      setSelectedAsset(updatedAsset);
      await loadAssets();
      setShowSuccessMessage(true);
      setTimeout(() => setShowSuccessMessage(false), 3000);
    } catch (err) {
      console.error('Error updating asset:', err);
      Alert.alert('Error', 'Failed to update asset. Please try again.');
    } finally {
      setUpdating(false);
    }
  };

  const handlePublishAsset = async (publishAt: Date | null, unpublishAt: Date | null, ageRating: string) => {
    if (!selectedAsset) return;

    try {
      setPublishing(true);
      
      const clearFields: string[] = [];
      let publishAtToSend = publishAt ? publishAt.toISOString() : undefined;
      let unpublishAtToSend = unpublishAt ? unpublishAt.toISOString() : undefined;
      let ageRatingToSend = ageRating || undefined;
      
      if (selectedAsset.publishRule?.publishAt && !publishAt) {
        publishAtToSend = undefined;
        clearFields.push('publishAt');
      }
      
      if (selectedAsset.publishRule?.unpublishAt && !unpublishAt) {
        unpublishAtToSend = undefined;
        clearFields.push('unpublishAt');
      }
      
      if (selectedAsset.publishRule?.ageRating && !ageRating) {
        ageRatingToSend = undefined;
        clearFields.push('ageRating');
      }
      
      const updatedAsset = await assetService.publishAsset(
        selectedAsset.id,
        publishAtToSend,
        unpublishAtToSend,
        undefined,
        ageRatingToSend,
        clearFields
      );
      
      setSelectedAsset(updatedAsset);
      await loadAssets();
      setShowSuccessMessage(true);
      setTimeout(() => setShowSuccessMessage(false), 3000);
    } catch (err) {
      console.error('Error publishing asset:', err);
      Alert.alert('Error', 'Failed to publish asset. Please try again.');
    } finally {
      setPublishing(false);
    }
  };

  return {
    assets,
    selectedAsset,
    loading,
    refreshing,
    error,
    showSuccessMessage,
    deleting,
    publishing,
    updating,
    children,
    childrenLoading,
    loadAssets,
    handleRefresh,
    handleRefreshSelectedAsset,
    handleAssetSelect,
    handleDeleteAsset,
    handleUpdateAsset,
    handlePublishAsset,
  };
} 