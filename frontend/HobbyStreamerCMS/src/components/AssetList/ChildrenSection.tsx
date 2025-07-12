import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native';
import { Asset, AssetType } from '../../types/asset';

interface ChildrenSectionProps {
  asset: Asset;
  children: Asset[];
  childrenLoading: boolean;
  onSelectChild: (child: Asset) => void;
}

export default function ChildrenSection({ 
  asset, 
  children, 
  childrenLoading, 
  onSelectChild 
}: ChildrenSectionProps) {
  if (asset.type !== AssetType.SERIES && asset.type !== AssetType.SEASON) {
    return null;
  }

  return (
    <View style={styles.container}>
      <Text style={styles.sectionTitle}>
        {asset.type === AssetType.SERIES ? 'Seasons' : 'Episodes'}
      </Text>
      
      {childrenLoading ? (
        <View style={styles.loading}>
          <ActivityIndicator size="small" color="#007AFF" />
          <Text style={styles.loadingText}>
            Loading {asset.type === AssetType.SERIES ? 'seasons' : 'episodes'}...
          </Text>
        </View>
      ) : children.length > 0 ? (
        <View style={styles.childrenList}>
          {children.map((child) => (
            <TouchableOpacity
              key={child.id}
              style={styles.childItem}
              onPress={() => onSelectChild(child)}
            >
              <Text style={styles.childTitle}>{child.title || `Asset ${child.id}`}</Text>
              <Text style={styles.childDetails}>
                {child.type} • {child.genre || 'No genre'}
              </Text>
              <Text style={styles.childStatus}>Status: {child.status || 'Unknown'}</Text>
              <Text style={styles.childDate}>
                Created: {new Date(child.createdAt).toLocaleDateString()}
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      ) : (
        <Text style={styles.noChildrenText}>
          No {asset.type === AssetType.SERIES ? 'seasons' : 'episodes'} found
        </Text>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginBottom: 20,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  loading: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  loadingText: {
    fontSize: 14,
    color: '#666',
    marginLeft: 10,
  },
  childrenList: {
    marginTop: 10,
  },
  childItem: {
    backgroundColor: '#f8f9fa',
    borderWidth: 1,
    borderColor: '#e9ecef',
    borderRadius: 8,
    padding: 12,
    marginBottom: 8,
  },
  childTitle: {
    fontSize: 14,
    fontWeight: '600',
    marginBottom: 4,
    color: '#333',
  },
  childDetails: {
    fontSize: 12,
    color: '#007AFF',
    marginBottom: 2,
    fontWeight: '500',
  },
  childStatus: {
    fontSize: 12,
    color: '#666',
    marginBottom: 2,
  },
  childDate: {
    fontSize: 11,
    color: '#999',
  },
  noChildrenText: {
    fontSize: 14,
    color: '#999',
    fontStyle: 'italic',
    textAlign: 'center',
    padding: 20,
  },
}); 