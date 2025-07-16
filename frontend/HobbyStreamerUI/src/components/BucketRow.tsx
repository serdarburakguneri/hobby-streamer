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
  if (assets.length === 0) {
    return null;
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>{title}</Text>
      <ScrollView 
        horizontal 
        showsHorizontalScrollIndicator={false}
        contentContainerStyle={styles.scrollContent}
      >
        {assets.map((asset) => (
          <AssetCard 
            key={asset.id} 
            asset={asset} 
            onPress={onAssetPress}
          />
        ))}
      </ScrollView>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginVertical: 16,
  },
  title: {
    color: '#ffffff',
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 12,
    marginLeft: 16,
  },
  scrollContent: {
    paddingHorizontal: 8,
  },
}); 