import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Bucket, BucketType } from '../../types/asset';

interface BucketItemProps {
  bucket: Bucket;
  isSelected: boolean;
  onSelect: (bucket: Bucket) => void;
}

export default function BucketItem({ bucket, isSelected, onSelect }: BucketItemProps) {
  const formatBucketType = (type: BucketType): string => {
    return type.replace('_', ' ').toLowerCase();
  };

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  const getBucketIcon = (type: BucketType) => {
    switch (type.toLowerCase()) {
      case 'movie_collection':
        return 'film';
      case 'series_collection':
        return 'tv';
      case 'documentary_collection':
        return 'library';
      case 'trailer_collection':
        return 'play-circle';
      case 'featured':
        return 'star';
      case 'trending':
        return 'trending-up';
      default:
        return 'folder';
    }
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

  const assetCount = bucket.assets?.length || 0;

  return (
    <TouchableOpacity
      style={[styles.container, isSelected && styles.selectedContainer]}
      onPress={() => onSelect(bucket)}
    >
      <View style={styles.header}>
        <View style={styles.iconContainer}>
          <Ionicons 
            name={getBucketIcon(bucket.type) as any} 
            size={24} 
            color={isSelected ? '#fff' : '#007AFF'} 
          />
        </View>
        <View style={styles.content}>
          <Text style={[styles.name, isSelected && styles.selectedText]}>
            {bucket.name}
          </Text>
          <View style={[styles.badge, isSelected && styles.selectedBadge]}>
            <Text style={[styles.badgeText, isSelected && styles.selectedBadgeText]}>
              {formatBucketType(bucket.type)}
            </Text>
          </View>
        </View>
        <View style={styles.statusContainer}>
          <Ionicons 
            name={getStatusIcon(bucket.status || 'active') as any} 
            size={16} 
            color={isSelected ? '#e0e0e0' : getStatusColor(bucket.status || 'active')} 
          />
        </View>
      </View>
      
      {bucket.description && (
        <Text style={[styles.description, isSelected && styles.selectedSubtext]}>
          {bucket.description}
        </Text>
      )}
      
      <View style={styles.footer}>
        <View style={styles.assetCountContainer}>
          <Ionicons 
            name="videocam" 
            size={14} 
            color={isSelected ? '#e0e0e0' : '#999'} 
          />
          <Text style={[styles.assetCount, isSelected && styles.selectedSubtext]}>
            {assetCount} asset{assetCount !== 1 ? 's' : ''}
          </Text>
        </View>
        <View style={styles.dateContainer}>
          <Ionicons 
            name="calendar" 
            size={14} 
            color={isSelected ? '#e0e0e0' : '#999'} 
          />
          <Text style={[styles.date, isSelected && styles.selectedSubtext]}>
            {formatDate(bucket.createdAt)}
          </Text>
        </View>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: 'transparent',
    borderRadius: 8,
    padding: 12,
    marginBottom: 8,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  selectedContainer: {
    backgroundColor: '#007AFF',
    borderColor: '#007AFF',
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
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },
  name: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    flex: 1,
    marginRight: 8,
  },
  selectedText: {
    color: '#fff',
  },
  badge: {
    backgroundColor: '#f0f0f0',
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 4,
  },
  selectedBadge: {
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
  },
  badgeText: {
    fontSize: 12,
    color: '#666',
    textTransform: 'capitalize',
  },
  selectedBadgeText: {
    color: '#fff',
  },
  statusContainer: {
    marginLeft: 8,
  },
  description: {
    fontSize: 14,
    color: '#666',
    marginBottom: 8,
    lineHeight: 18,
    marginLeft: 36,
  },
  selectedSubtext: {
    color: '#e0e0e0',
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginLeft: 36,
  },
  assetCountContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  assetCount: {
    fontSize: 12,
    color: '#999',
  },
  dateContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  date: {
    fontSize: 12,
    color: '#999',
  },
}); 