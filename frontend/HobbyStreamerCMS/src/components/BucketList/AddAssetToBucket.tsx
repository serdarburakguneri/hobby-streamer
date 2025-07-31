import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  TextInput,
  ScrollView,
  ActivityIndicator,
  Modal,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useAssetService } from '../../services/api';
import { Asset } from '../../types/asset';

interface AddAssetToBucketProps {
  bucketId: string;
  assets: Asset[];
  onAddAsset: (assetId: string) => void;
  onClose: () => void;
  visible: boolean;
}

export default function AddAssetToBucket({
  bucketId,
  assets,
  onAddAsset,
  onClose,
  visible,
}: AddAssetToBucketProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<Asset[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [dropdownVisible, setDropdownVisible] = useState(false);
  const assetService = useAssetService();

  useEffect(() => {
    if (searchQuery.trim().length > 0) {
      setSearchLoading(true);
      assetService
        .searchAssets(searchQuery)
        .then((results: any) => {
          const assetIds = assets.map((a) => a.id);
          const filteredResults = results.assets.filter((asset: Asset) =>
            !assetIds.includes(asset.id)
          );
          setSearchResults(filteredResults);
          setDropdownVisible(true);
        })
        .catch(() => {
          setSearchResults([]);
        })
        .finally(() => {
          setSearchLoading(false);
        });
    } else {
      setSearchResults([]);
      setDropdownVisible(false);
    }
  }, [searchQuery, assets]);

  const handleAddAsset = (asset: Asset) => {
    onAddAsset(asset.id);
    setSearchQuery('');
    setSearchResults([]);
    setDropdownVisible(false);
    onClose();
  };

  const handleClose = () => {
    setSearchQuery('');
    setSearchResults([]);
    setDropdownVisible(false);
    onClose();
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={handleClose}
    >
      <View style={styles.modalOverlay}>
        <View style={styles.modalContainer}>
          <View style={styles.modalHeader}>
            <Ionicons name="add-circle" size={24} color="#333" />
            <Text style={styles.modalTitle}>Add Asset to Bucket</Text>
            <TouchableOpacity onPress={handleClose} style={styles.closeButton}>
              <Ionicons name="close" size={24} color="#666" />
            </TouchableOpacity>
          </View>

          <Text style={styles.modalMessage}>
            Search for assets to add to this bucket
          </Text>

          <View style={styles.searchContainer}>
            <TextInput
              style={styles.searchInput}
              placeholder="Type asset title..."
              value={searchQuery}
              onChangeText={(text) => {
                setSearchQuery(text);
                if (searchResults.length > 0) setDropdownVisible(true);
              }}
              autoFocus
            />
            {searchLoading && (
              <ActivityIndicator size="small" color="#007AFF" style={styles.searchLoading} />
            )}
          </View>

          {dropdownVisible && searchResults.length > 0 && (
            <View style={styles.resultsContainer}>
              <ScrollView style={styles.resultsList} showsVerticalScrollIndicator={false}>
                {searchResults.map((asset) => (
                  <TouchableOpacity
                    key={asset.id}
                    style={styles.resultItem}
                    onPress={() => handleAddAsset(asset)}
                  >
                    <View style={styles.assetInfo}>
                      <Text style={styles.assetTitle}>{asset.title || `Asset ${asset.id}`}</Text>
                      <Text style={styles.assetType}>{asset.type}</Text>
                      {asset.genre && (
                        <Text style={styles.assetGenre}>{asset.genre}</Text>
                      )}
                    </View>
                    <Ionicons name="add" size={20} color="#007AFF" />
                  </TouchableOpacity>
                ))}
              </ScrollView>
            </View>
          )}

          {dropdownVisible && searchResults.length === 0 && !searchLoading && searchQuery.trim().length > 0 && (
            <View style={styles.noResultsContainer}>
              <Ionicons name="search" size={32} color="#ccc" />
              <Text style={styles.noResultsText}>No assets found</Text>
              <Text style={styles.noResultsSubtext}>
                Try a different search term
              </Text>
            </View>
          )}

          <View style={styles.modalButtons}>
            <TouchableOpacity
              style={styles.cancelButton}
              onPress={handleClose}
            >
              <Text style={styles.cancelButtonText}>Cancel</Text>
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  modalContainer: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 24,
    margin: 20,
    minWidth: 400,
    maxWidth: 600,
    maxHeight: '80%',
  },
  modalHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
  },
  modalTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
    flex: 1,
    marginLeft: 8,
  },
  closeButton: {
    padding: 4,
  },
  modalMessage: {
    fontSize: 14,
    color: '#666',
    marginBottom: 20,
  },
  searchContainer: {
    position: 'relative',
    marginBottom: 16,
  },
  searchInput: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
  },
  searchLoading: {
    position: 'absolute',
    right: 12,
    top: 12,
  },
  resultsContainer: {
    maxHeight: 300,
    marginBottom: 16,
  },
  resultsList: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
  },
  resultItem: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  assetInfo: {
    flex: 1,
  },
  assetTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    marginBottom: 4,
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
  noResultsContainer: {
    alignItems: 'center',
    padding: 20,
  },
  noResultsText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#666',
    marginTop: 8,
    marginBottom: 4,
  },
  noResultsSubtext: {
    fontSize: 14,
    color: '#999',
  },
  modalButtons: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
  },
  cancelButton: {
    backgroundColor: '#f0f0f0',
    paddingHorizontal: 20,
    paddingVertical: 10,
    borderRadius: 8,
  },
  cancelButtonText: {
    color: '#666',
    fontSize: 16,
    fontWeight: '600',
  },
}); 