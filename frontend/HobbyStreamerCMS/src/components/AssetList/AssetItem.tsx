import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Asset } from '../../types/asset';

interface AssetItemProps {
  asset: Asset;
  isSelected: boolean;
  onSelect: (asset: Asset) => void;
}

export default function AssetItem({ asset, isSelected, onSelect }: AssetItemProps) {
  const getAssetIcon = (type: string) => {
    switch (type.toLowerCase()) {
      case 'movie':
        return 'film';
      case 'series':
        return 'tv';
      case 'documentary':
        return 'document-text';
      case 'trailer':
        return 'play-circle';
      default:
        return 'videocam';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'published':
        return 'checkmark-circle';
      case 'draft':
        return 'create';
      case 'processing':
        return 'sync';
      default:
        return 'help-circle';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'published':
        return '#4CAF50';
      case 'draft':
        return '#FF9800';
      case 'processing':
        return '#2196F3';
      default:
        return '#9E9E9E';
    }
  };

  return (
    <TouchableOpacity 
      style={[
        styles.assetItem,
        isSelected && styles.selectedAssetItem
      ]}
      onPress={() => onSelect(asset)}
    >
      <View style={styles.header}>
        <View style={styles.iconContainer}>
          <Ionicons 
            name={getAssetIcon(asset.type) as any} 
            size={24} 
            color={isSelected ? '#fff' : '#007AFF'} 
          />
        </View>
        <View style={styles.content}>
          <Text style={[styles.assetTitle, isSelected && styles.selectedText]}>
            {asset.title || `Asset ${asset.id}`}
          </Text>
          <Text style={[styles.assetDetails, isSelected && styles.selectedSubtext]}>
            {asset.type} • {asset.genre || 'No genre'}
          </Text>
        </View>
        <View style={styles.statusContainer}>
          <Ionicons 
            name={getStatusIcon(asset.status || 'unknown') as any} 
            size={16} 
            color={isSelected ? '#e0e0e0' : getStatusColor(asset.status || 'unknown')} 
          />
          <Text style={[styles.assetStatus, isSelected && styles.selectedSubtext]}>
            {asset.status || 'Unknown'}
          </Text>
        </View>
      </View>
      <Text style={[styles.assetDate, isSelected && styles.selectedSubtext]}>
        Created: {new Date(asset.createdAt).toLocaleDateString()}
      </Text>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  assetItem: {
    backgroundColor: 'transparent',
    padding: 15,
    marginBottom: 10,
    borderRadius: 8,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  selectedAssetItem: {
    backgroundColor: '#007AFF',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 8,
  },
  iconContainer: {
    marginRight: 12,
  },
  content: {
    flex: 1,
  },
  assetTitle: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 4,
    color: '#333',
  },
  selectedText: {
    color: '#fff',
  },
  assetDetails: {
    fontSize: 14,
    color: '#007AFF',
    fontWeight: '500',
  },
  selectedSubtext: {
    color: '#e0e0e0',
  },
  statusContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  assetStatus: {
    fontSize: 12,
    color: '#666',
    marginLeft: 4,
  },
  assetDate: {
    fontSize: 12,
    color: '#999',
    marginLeft: 36,
  },
}); 