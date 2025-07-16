import React from 'react';
import { View, Text, Image, StyleSheet, TouchableOpacity, Dimensions } from 'react-native';
import { Asset } from '../types/asset';

interface AssetCardProps {
  asset: Asset;
  onPress: (asset: Asset) => void;
}

const { width } = Dimensions.get('window');
const cardWidth = width * 0.25;
const cardHeight = cardWidth * 1.4;

export const AssetCard: React.FC<AssetCardProps> = ({ asset, onPress }) => {
  const hasVideos = Array.isArray(asset.videos) && asset.videos.length > 0;
  const thumbnailUrl = hasVideos
    ? asset.videos[0]?.thumbnail?.url || asset.videos[0]?.thumbnail?.storageLocation?.url
    : undefined;

  const imageSource = thumbnailUrl ? { uri: thumbnailUrl } : require('../../assets/video-placeholder.png');

  const title = asset.title || 'Untitled';

  return (
    <TouchableOpacity 
      style={styles.container} 
      onPress={() => onPress(asset)}
      activeOpacity={0.8}
    >
      <Text style={styles.title}>{title}</Text>
      <Image 
        source={imageSource} 
        style={styles.thumbnail}
        resizeMode="cover"
      />
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    width: cardWidth,
    height: cardHeight,
    marginHorizontal: 4,
    borderRadius: 6,
    overflow: 'hidden',
    backgroundColor: '#1a1a1a',
  },
  title: {
    color: '#ffffff',
    fontSize: 12,
    fontWeight: '500',
    padding: 8,
    textAlign: 'left',
  },
  thumbnail: {
    width: '100%',
    height: '100%',
  },
}); 