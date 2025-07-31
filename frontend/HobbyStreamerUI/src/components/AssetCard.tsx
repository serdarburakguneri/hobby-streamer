import React, { useState } from 'react';
import { View, Image, StyleSheet, TouchableOpacity, Dimensions, Text } from 'react-native';
import { Asset } from '../types/asset';
import { VideoPlayer } from './VideoPlayer';

interface AssetCardProps {
  asset: Asset;
  onPress: (asset: Asset) => void;
}

const { width } = Dimensions.get('window');
const cardWidth = width * 0.15;
const cardHeight = cardWidth * 1.4;

export const AssetCard: React.FC<AssetCardProps> = ({ asset, onPress }) => {
  const [showVideoPlayer, setShowVideoPlayer] = useState(false);
  
  const hasVideos = Array.isArray(asset.videos) && asset.videos.length > 0;
  
  // Try to get poster image first, then fall back to video thumbnail
  const posterImage = asset.images?.find(img => img.type?.toLowerCase() === 'poster');
  const thumbnailUrl = hasVideos
    ? asset.videos[0]?.thumbnail?.url || asset.videos[0]?.thumbnail?.storageLocation?.url
    : undefined;

  const imageSource = posterImage?.url 
    ? { uri: posterImage.url } 
    : thumbnailUrl 
    ? { uri: thumbnailUrl } 
    : require('../../assets/video-placeholder.png');

  const handlePress = () => {
    if (hasVideos) {
      setShowVideoPlayer(true);
    } else {
      onPress(asset);
    }
  };

  const handleCloseVideoPlayer = () => {
    setShowVideoPlayer(false);
  };

  return (
    <>
      <TouchableOpacity 
        style={styles.container} 
        onPress={handlePress}
        activeOpacity={0.7}
      >
        <Image 
          source={imageSource} 
          style={styles.thumbnail}
          resizeMode="cover"
        />
        <View style={styles.overlay} />
        {hasVideos && (
          <View style={styles.playButton}>
            <Text style={styles.playIcon}>▶</Text>
          </View>
        )}
      </TouchableOpacity>
      
      <VideoPlayer
        asset={asset}
        visible={showVideoPlayer}
        onClose={handleCloseVideoPlayer}
      />
    </>
  );
};

const styles = StyleSheet.create({
  container: {
    width: cardWidth,
    height: cardHeight,
    marginHorizontal: 12,
    borderRadius: 8,
    overflow: 'hidden',
    backgroundColor: '#2a2a2a',
    elevation: 3,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
  },
  thumbnail: {
    width: '100%',
    height: '100%',
  },
  overlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.1)',
  },
  playButton: {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: [{ translateX: -20 }, { translateY: -20 }],
    backgroundColor: 'rgba(0, 0, 0, 0.7)',
    borderRadius: 40,
    width: 60,
    height: 60,
    justifyContent: 'center',
    alignItems: 'center',
  },
  playIcon: {
    fontSize: 24,
    color: '#fff',
    marginLeft: 2,
  },
}); 