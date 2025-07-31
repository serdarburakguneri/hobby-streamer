import React, { useState, useMemo } from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
  RefreshControl,
} from 'react-native';
import Layout from '../components/Layout';
import AssetItem from '../components/AssetList/AssetItem';
import AssetDetails from '../components/AssetList/AssetDetails';
import AssetOperations from '../components/AssetList/AssetOperations';
import { useAssetList } from '../hooks/useAssetList';
import { Asset } from '../types/asset';

interface AssetListScreenProps {
  onCreateAsset: () => void;
  refreshTrigger?: number;
}

export default function AssetListScreen({ onCreateAsset, refreshTrigger }: AssetListScreenProps) {
  const [showDeleteConfirmation, setShowDeleteConfirmation] = useState(false);

  const {
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
    handleRefresh,
    handleRefreshSelectedAsset,
    handleAssetSelect,
    handleDeleteAsset,
    handleUpdateAsset,
    handlePublishAsset,
  } = useAssetList(refreshTrigger);





  const performDelete = async () => {
    await handleDeleteAsset();
    setShowDeleteConfirmation(false);
  };

  const renderAssetItem = ({ item }: { item: Asset }) => (
    <AssetItem
      asset={item}
      isSelected={selectedAsset?.id === item.id}
      onSelect={handleAssetSelect}
    />
  );



  if (loading) {
    return (
      <View style={styles.centerContainer}>
        <ActivityIndicator size="large" color="#007AFF" />
        <Text style={styles.loadingText}>Loading assets...</Text>
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.centerContainer}>
        <Text style={styles.errorText}>{error}</Text>
        <TouchableOpacity style={styles.retryButton} onPress={() => handleRefresh()}>
          <Text style={styles.retryButtonText}>Retry</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <Layout>
      {showSuccessMessage && (
        <View style={styles.successMessage}>
          <Text style={styles.successText}>✓ Operation completed successfully!</Text>
        </View>
      )}
      
      {showDeleteConfirmation && (
        <View style={styles.modalOverlay}>
          <View style={styles.modalContainer}>
            <Text style={styles.modalTitle}>Delete Asset</Text>
            <Text style={styles.modalMessage}>
              Are you sure you want to delete "{selectedAsset?.title || `Asset ${selectedAsset?.id}`}"? This action cannot be undone.
            </Text>
            <View style={styles.modalButtons}>
              <TouchableOpacity 
                style={[styles.modalButton, styles.cancelButton]}
                onPress={() => setShowDeleteConfirmation(false)}
              >
                <Text style={styles.cancelButtonText}>Cancel</Text>
              </TouchableOpacity>
              <TouchableOpacity 
                style={[styles.modalButton, styles.confirmDeleteButton]}
                onPress={performDelete}
                disabled={deleting}
              >
                {deleting ? (
                  <ActivityIndicator size="small" color="#fff" />
                ) : (
                  <Text style={styles.confirmDeleteButtonText}>Delete</Text>
                )}
              </TouchableOpacity>
            </View>
          </View>
        </View>
      )}





      <View style={styles.headerSection}>
        <Text style={styles.pageTitle}>Assets</Text>
        <TouchableOpacity style={styles.createButton} onPress={onCreateAsset}>
          <Text style={styles.createButtonText}>Create Asset</Text>
        </TouchableOpacity>
      </View>

      <View style={styles.mainContent}>
        <View style={styles.leftPanel}>
          <FlatList
            data={assets}
            keyExtractor={(item) => item.id}
            renderItem={renderAssetItem}
            style={styles.leftPanelList}
            contentContainerStyle={styles.leftPanelListContent}
            refreshControl={
              <RefreshControl
                refreshing={refreshing}
                onRefresh={handleRefresh}
                colors={['#007AFF']}
                tintColor="#007AFF"
              />
            }
            ListEmptyComponent={
              <View style={styles.emptyContainer}>
                <Text style={styles.emptyText}>No assets found</Text>
                <Text style={styles.emptySubtext}>
                  Create your first asset to get started
                </Text>
              </View>
            }
          />
        </View>

        <View style={styles.middlePanel}>
          <AssetDetails
            asset={selectedAsset}
            onUpdate={handleUpdateAsset}
            onSelectChild={handleAssetSelect}
            children={children}
            childrenLoading={childrenLoading}
            onRefresh={handleRefreshSelectedAsset}
          />
        </View>

        <View style={styles.rightPanel}>
          <AssetOperations
            asset={selectedAsset}
            onDelete={() => setShowDeleteConfirmation(true)}
            onPublish={handlePublishAsset}
            deleting={deleting}
            publishing={publishing}
          />
        </View>
      </View>
    </Layout>
  );
}

const styles = StyleSheet.create({
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#f5f5f5',
    padding: 20,
  },
  createButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 6,
  },
  createButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  headerSection: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#ffffff',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  pageTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#333',
  },
  mainContent: {
    flexDirection: 'row',
    flex: 1,
    backgroundColor: '#f0f2f5',
  },
  leftPanel: {
    flex: 1,
    padding: 10,
    backgroundColor: '#f0f2f5',
  },
  middlePanel: {
    flex: 2,
    padding: 10,
    backgroundColor: '#ffffff',
  },
  rightPanel: {
    flex: 1,
    padding: 10,
    backgroundColor: '#fafbfc',
  },
  panelTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  loadingText: {
    marginTop: 10,
    fontSize: 16,
    color: '#666',
  },
  errorText: {
    fontSize: 16,
    color: '#d32f2f',
    textAlign: 'center',
    marginBottom: 20,
  },
  retryButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 20,
    paddingVertical: 10,
    borderRadius: 8,
  },
  retryButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  emptyContainer: {
    alignItems: 'center',
    padding: 40,
  },
  emptyText: {
    fontSize: 18,
    fontWeight: '600',
    color: '#666',
    marginBottom: 10,
  },
  emptySubtext: {
    fontSize: 14,
    color: '#999',
    textAlign: 'center',
  },
  leftPanelList: {
    backgroundColor: '#f0f2f5',
  },
  leftPanelListContent: {
    backgroundColor: '#f0f2f5',
  },
  successMessage: {
    backgroundColor: '#4CAF50',
    padding: 12,
    alignItems: 'center',
  },
  successText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
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
}); 