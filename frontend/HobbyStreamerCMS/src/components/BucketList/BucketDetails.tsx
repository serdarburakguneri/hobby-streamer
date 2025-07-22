import React, { useState } from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Bucket, BucketType, Asset } from '../../types/asset';
import AddAssetToBucket from './AddAssetToBucket';
import EditableField from '../AssetList/EditableField';

interface BucketDetailsProps {
  bucket: Bucket | null;
  onUpdate: (field: string, value: any) => void;
  onRemoveAsset: (assetId: string) => void;
  onAddAsset: (assetId: string) => void;
  updating: boolean;
}

export default function BucketDetails({ 
  bucket, 
  onUpdate, 
  onRemoveAsset, 
  onAddAsset,
  updating 
}: BucketDetailsProps) {
  const [showAssetList, setShowAssetList] = useState(true);
  const [showAddAssetModal, setShowAddAssetModal] = useState(false);
  const [showRemoveConfirmation, setShowRemoveConfirmation] = useState(false);
  const [assetToRemove, setAssetToRemove] = useState<string | null>(null);

  if (!bucket) {
    return (
      <View style={styles.emptyContainer}>
        <Ionicons name="document-outline" size={32} color="#ccc" />
        <Text style={styles.emptyText}>Select a bucket to view details</Text>
      </View>
    );
  }

  const formatBucketType = (type: BucketType): string => {
    return type.replace('_', ' ').toLowerCase();
  };

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  const handleRemoveAsset = (assetId: string) => {
    console.log('handleRemoveAsset called with assetId:', assetId);
    setAssetToRemove(assetId);
    setShowRemoveConfirmation(true);
  };

  const confirmRemoveAsset = () => {
    if (assetToRemove) {
      console.log('Remove confirmed, calling onRemoveAsset with:', assetToRemove);
      onRemoveAsset(assetToRemove);
      setShowRemoveConfirmation(false);
      setAssetToRemove(null);
    }
  };

  const cancelRemoveAsset = () => {
    setShowRemoveConfirmation(false);
    setAssetToRemove(null);
  };

  const getStatusIcon = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'active':
        return 'checkmark-circle';
      case 'inactive':
        return 'pause-circle';
      case 'draft':
        return 'create';
      default:
        return 'help-circle';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'active':
        return '#4CAF50';
      case 'inactive':
        return '#FF9800';
      case 'draft':
        return '#9E9E9E';
      default:
        return '#9E9E9E';
    }
  };

  return (
    <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Ionicons name="information-circle" size={20} color="#333" />
          <Text style={styles.sectionTitle}>Basic Information</Text>
        </View>
        
        <EditableField
          label="Name"
          field="name"
          value={bucket.name}
          onUpdate={onUpdate}
          placeholder="Enter bucket name"
        />
        
        <EditableField
          label="Description"
          field="description"
          value={bucket.description || ''}
          onUpdate={onUpdate}
          placeholder="Enter bucket description"
        />
        
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Type:</Text>
          <View style={styles.typeContainer}>
            <Text style={styles.typeText}>{formatBucketType(bucket.type)}</Text>
            <View style={styles.typeBadge}>
              <Text style={styles.typeBadgeText}>{bucket.type}</Text>
            </View>
          </View>
        </View>
        
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Status:</Text>
          <View style={styles.statusContainer}>
            <Ionicons 
              name={getStatusIcon(bucket.status || 'active') as any} 
              size={16} 
              color={getStatusColor(bucket.status || 'active')} 
            />
            <Text style={styles.statusText}>{bucket.status || 'active'}</Text>
          </View>
        </View>
        
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Created:</Text>
          <Text style={styles.detailValue}>{formatDate(bucket.createdAt)}</Text>
        </View>
        
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Last Updated:</Text>
          <Text style={styles.detailValue}>{formatDate(bucket.updatedAt)}</Text>
        </View>
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Ionicons name="videocam" size={20} color="#333" />
          <Text style={styles.sectionTitle}>Assets ({bucket.assets?.length || 0})</Text>
          <View style={styles.sectionActions}>
            <TouchableOpacity
              style={styles.addButton}
              onPress={() => setShowAddAssetModal(true)}
              disabled={updating}
            >
              <Ionicons name="add" size={16} color="#fff" />
              <Text style={styles.addButtonText}>Add Asset</Text>
            </TouchableOpacity>
            <TouchableOpacity
              style={styles.toggleButton}
              onPress={() => setShowAssetList(!showAssetList)}
            >
              <Ionicons 
                name={showAssetList ? 'chevron-up' : 'chevron-down'} 
                size={16} 
                color="#fff" 
              />
              <Text style={styles.toggleButtonText}>
                {showAssetList ? 'Hide' : 'Show'}
              </Text>
            </TouchableOpacity>
          </View>
        </View>

        {showAssetList && (
          <View style={styles.assetList}>
            {bucket.assets && bucket.assets.length > 0 ? (
              bucket.assets.map((asset: Asset) => (
                <View key={asset.id} style={styles.assetItem}>
                  <View style={styles.assetInfo}>
                    <View style={styles.assetHeader}>
                      <Ionicons name="videocam" size={16} color="#007AFF" />
                      <Text style={styles.assetTitle}>{asset.title || `Asset ${asset.id}`}</Text>
                    </View>
                    <Text style={styles.assetType}>{asset.type}</Text>
                    {asset.genre && (
                      <Text style={styles.assetGenre}>{asset.genre}</Text>
                    )}
                  </View>
                  <TouchableOpacity
                    style={[styles.removeButton, updating && styles.disabledButton]}
                    onPress={() => {
                      console.log('Remove button pressed for asset:', asset.id, 'updating:', updating);
                      if (!updating) {
                        handleRemoveAsset(asset.id);
                      }
                    }}
                    disabled={updating}
                  >
                    {updating ? (
                      <ActivityIndicator size="small" color="#ff3b30" />
                    ) : (
                      <>
                        <Ionicons name="trash" size={14} color="#ff3b30" />
                        <Text style={styles.removeButtonText}>Remove</Text>
                      </>
                    )}
                  </TouchableOpacity>
                </View>
              ))
            ) : (
              <View style={styles.emptyAssetsContainer}>
                <Ionicons name="videocam-outline" size={32} color="#ccc" />
                <Text style={styles.emptyAssetsText}>No assets in this bucket</Text>
              </View>
            )}
          </View>
        )}
      </View>

      {updating && (
        <View style={styles.loadingOverlay}>
          <ActivityIndicator size="large" color="#007AFF" />
          <Text style={styles.loadingText}>Updating bucket...</Text>
        </View>
      )}

      <AddAssetToBucket
        bucketId={bucket.id}
        assets={bucket.assets || []}
        onAddAsset={onAddAsset}
        onClose={() => setShowAddAssetModal(false)}
        visible={showAddAssetModal}
      />

      {showRemoveConfirmation && (
        <View style={styles.modalOverlay}>
          <View style={styles.modalContainer}>
            <Text style={styles.modalTitle}>Remove Asset</Text>
            <Text style={styles.modalMessage}>
              Are you sure you want to remove this asset from the bucket?
            </Text>
            <View style={styles.modalButtons}>
              <TouchableOpacity 
                style={[styles.modalButton, styles.cancelButton]}
                onPress={cancelRemoveAsset}
              >
                <Text style={styles.cancelButtonText}>Cancel</Text>
              </TouchableOpacity>
              <TouchableOpacity 
                style={[styles.modalButton, styles.confirmDeleteButton]}
                onPress={confirmRemoveAsset}
                disabled={updating}
              >
                {updating ? (
                  <ActivityIndicator size="small" color="#fff" />
                ) : (
                  <Text style={styles.confirmDeleteButtonText}>Remove</Text>
                )}
              </TouchableOpacity>
            </View>
          </View>
        </View>
      )}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  emptyText: {
    fontSize: 16,
    color: '#666',
    textAlign: 'center',
    marginTop: 8,
  },
  section: {
    backgroundColor: '#fff',
    borderRadius: 8,
    padding: 16,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  sectionActions: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  addButton: {
    backgroundColor: '#4CAF50',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 6,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  addButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
  },
  field: {
    marginBottom: 16,
    alignItems: 'flex-end',
  },
  detailRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 5,
  },
  detailLabel: {
    fontSize: 14,
    fontWeight: '500',
  },
  detailValue: {
    fontSize: 14,
    fontWeight: '500',
  },
  labelContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
    marginBottom: 4,
  },
  label: {
    fontSize: 14,
    fontWeight: '600',
    color: '#666',
  },
  value: {
    fontSize: 16,
    color: '#333',
  },
  typeContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  typeText: {
    fontSize: 16,
    color: '#333',
    textTransform: 'capitalize',
  },
  typeBadge: {
    backgroundColor: '#f0f0f0',
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
  },
  typeBadgeText: {
    fontSize: 12,
    color: '#666',
    fontWeight: '600',
  },
  statusContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  statusText: {
    fontSize: 16,
    color: '#333',
    textTransform: 'capitalize',
  },
  toggleButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 6,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  toggleButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  assetList: {
    marginTop: 8,
  },
  assetItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 12,
    backgroundColor: '#f8f9fa',
    borderRadius: 6,
    marginBottom: 8,
  },
  assetInfo: {
    flex: 1,
  },
  assetHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
    marginBottom: 4,
  },
  assetTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
  },
  assetType: {
    fontSize: 14,
    color: '#666',
    marginBottom: 2,
  },
  assetGenre: {
    fontSize: 12,
    color: '#999',
  },
  removeButton: {
    backgroundColor: '#ff3b30',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 6,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  removeButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  disabledButton: {
    opacity: 0.5,
  },
  emptyAssetsContainer: {
    alignItems: 'center',
    padding: 20,
  },
  modalOverlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 1000,
  },
  modalContainer: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 24,
    margin: 20,
    minWidth: 400,
    maxWidth: 600,
  },
  modalTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 10,
    textAlign: 'center',
  },
  modalMessage: {
    fontSize: 14,
    color: '#666',
    marginBottom: 20,
    textAlign: 'center',
    lineHeight: 20,
  },
  modalButtons: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: 10,
  },
  modalButton: {
    flex: 1,
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  cancelButton: {
    backgroundColor: '#f0f0f0',
  },
  confirmDeleteButton: {
    backgroundColor: '#ff3b30',
  },
  cancelButtonText: {
    color: '#666',
    fontSize: 14,
    fontWeight: '600',
  },
  confirmDeleteButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  emptyAssetsText: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
    fontStyle: 'italic',
    marginTop: 8,
  },
  loadingOverlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(255, 255, 255, 0.8)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    marginTop: 8,
    fontSize: 16,
    color: '#666',
  },
}); 