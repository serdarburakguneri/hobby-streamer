import React from 'react';
import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { AssetCard } from './AssetCard';
import { Asset } from '../types/asset';

interface BucketRowProps {
  title: string;
  assets: Asset[];
  onAssetPress: (asset: Asset) => void;
}

export const BucketRow: React.FC<BucketRowProps> = ({ title, assets, onAssetPress }) => {
  const displayTitle = title && title.trim() ? title : 'Untitled Bucket';
  if (!displayTitle) return null;
  console.log(`BucketRow "${displayTitle}" rendering with ${assets.length} assets`);

  return (
    <View style={styles.container}>
      <Text style={styles.title}>{displayTitle}</Text>
      {assets.length > 0 ? (
        <ScrollView 
          horizontal 
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={styles.scrollContent}
          decelerationRate="fast"
          snapToInterval={120}
        >
          {assets.map((asset) => (
            <AssetCard 
              key={asset.id} 
              asset={asset} 
              onPress={onAssetPress}
            />
          ))}
        </ScrollView>
      ) : (
        <View style={styles.emptyState}>
          <Text style={styles.emptyText}>No videos in this bucket</Text>
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginVertical: 24,
  },
  title: {
    color: '#ffffff',
    fontSize: 24,
    fontWeight: '700',
    marginBottom: 16,
    marginLeft: 16,
    letterSpacing: 0.5,
    textShadowColor: 'rgba(0, 0, 0, 0.5)',
    textShadowOffset: { width: 0, height: 1 },
    textShadowRadius: 2,
  },
  scrollContent: {
    paddingHorizontal: 16,
  },
  emptyState: {
    paddingHorizontal: 16,
    paddingVertical: 20,
  },
  emptyText: {
    color: '#666666',
    fontSize: 14,
    fontStyle: 'italic',
  },
}); 