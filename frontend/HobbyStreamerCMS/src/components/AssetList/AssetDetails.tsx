import React from 'react';
import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { Asset } from '../../types/asset';
import EditableField from './EditableField';
import VideoSection from './VideoSection';
import ImageUpload from './ImageUpload';
import ChildrenSection from './ChildrenSection';

interface AssetDetailsProps {
  asset: Asset | null;
  onUpdate: (field: string, value: any) => Promise<void>;
  onSelectChild: (child: Asset) => void;
  children: Asset[];
  childrenLoading: boolean;
  onRefresh?: () => void;
}

export default function AssetDetails({ 
  asset, 
  onUpdate, 
  onSelectChild,
  children,
  childrenLoading,
  onRefresh
}: AssetDetailsProps) {
  if (!asset) {
    return (
      <View style={styles.placeholder}>
        <Text style={styles.placeholderText}>Select an asset to view details</Text>
      </View>
    );
  }

  return (
    <ScrollView style={styles.container}>
      <EditableField
        label="Title"
        field="title"
        value={asset.title}
        onUpdate={onUpdate}
        placeholder="Enter title"
      />
      
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Basic Information</Text>
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>ID:</Text>
          <Text style={styles.detailValue}>{asset.id}</Text>
        </View>
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Slug:</Text>
          <Text style={styles.detailValue}>{asset.slug}</Text>
        </View>
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Type:</Text>
          <Text style={styles.detailValue}>{asset.type}</Text>
        </View>
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Status:</Text>
          <Text style={[styles.detailValue, styles.statusText]}>
            {asset.status || 'Unknown'}
          </Text>
        </View>
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Created:</Text>
          <Text style={styles.detailValue}>{new Date(asset.createdAt).toLocaleString()}</Text>
        </View>
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Updated:</Text>
          <Text style={styles.detailValue}>{new Date(asset.updatedAt).toLocaleString()}</Text>
        </View>
        
        <EditableField
          label="Description"
          field="description"
          value={asset.description}
          onUpdate={onUpdate}
          placeholder="Enter description"
        />
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Genres</Text>
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Primary Genre:</Text>
          <Text style={styles.detailValue}>{asset.genre || 'Not set'}</Text>
        </View>
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Additional Genres:</Text>
          <Text style={styles.detailValue}>{asset.genres?.length > 0 ? asset.genres.join(', ') : 'None'}</Text>
        </View>
        <Text style={styles.noteText}>Note: Genre editing is not yet supported with the new API</Text>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Tags</Text>
        <View style={styles.detailRow}>
          <Text style={styles.detailLabel}>Tags:</Text>
          <Text style={styles.detailValue}>{asset.tags?.length > 0 ? asset.tags.join(', ') : 'None'}</Text>
        </View>
        <Text style={styles.noteText}>Note: Tag editing is not yet supported with the new API</Text>
      </View>

      <VideoSection 
        asset={asset}
        onVideoAdded={() => {
          if (onRefresh) {
            onRefresh();
          }
        }}
      />

      <ImageUpload 
        asset={asset}
        onImageAdded={() => {
          if (onRefresh) {
            onRefresh();
          }
        }}
      />

      <ChildrenSection
        asset={asset}
        children={children}
        childrenLoading={childrenLoading}
        onSelectChild={onSelectChild}
      />
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  placeholder: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  placeholderText: {
    fontSize: 16,
    color: '#999',
  },
  section: {
    marginBottom: 20,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 10,
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
  statusText: {
    fontWeight: 'bold',
  },
  noteText: {
    fontSize: 12,
    color: '#666',
    fontStyle: 'italic',
    marginTop: 4,
  },
}); 