import React, { useState, useRef, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Modal, Dimensions, Platform, ScrollView } from 'react-native';
import { Video, ResizeMode } from 'expo-av';
import Hls from 'hls.js';
import * as dashjs from 'dashjs';
import { Asset } from '../types/asset';

interface VideoPlayerProps {
  asset: Asset;
  visible: boolean;
  onClose: () => void;
}

const { width: screenWidth, height: screenHeight } = Dimensions.get('window');

export const VideoPlayer: React.FC<VideoPlayerProps> = ({ asset, visible, onClose }) => {
  const [isPlaying, setIsPlaying] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedVideo, setSelectedVideo] = useState<any>(null);
  const videoRef = useRef<HTMLVideoElement>(null);
  const hlsRef = useRef<Hls | null>(null);
  const dashRef = useRef<dashjs.MediaPlayerClass | null>(null);

  const getBestVideo = () => {
    if (!asset.videos || asset.videos.length === 0) {
      return null;
    }

    const videos = asset.videos;
    
    const hlsVideo = videos.find(v => v.format === 'hls' && v.status === 'ready');
    if (hlsVideo) return hlsVideo;
    
    const dashVideo = videos.find(v => v.format === 'dash' && v.status === 'ready');
    if (dashVideo) return dashVideo;
    
    const rawVideo = videos.find(v => v.format === 'raw' && v.status === 'ready');
    if (rawVideo) return rawVideo;
    
    return videos[0];
  };

  const getVideoUrl = (video: any) => {
    if (video.streamInfo?.playUrl) {
      return video.streamInfo.playUrl;
    }
    if (video.storageLocation?.url) {
      return video.storageLocation.url;
    }
    return null;
  };

  const getAgeRatingLabel = (ageRating: string) => {
    switch (ageRating?.toLowerCase()) {
      case 'g':
        return 'G - General Audience';
      case 'pg':
        return 'PG - Parental Guidance';
      case 'pg-13':
        return 'PG-13 - Parental Guidance (13+)';
      case 'r':
        return 'R - Restricted (17+)';
      case 'nc-17':
        return 'NC-17 - Adults Only';
      default:
        return ageRating || 'Not Rated';
    }
  };

  const getGenreLabel = (genre: string) => {
    if (!genre) return 'Unknown Genre';
    return genre.charAt(0).toUpperCase() + genre.slice(1);
  };

  useEffect(() => {
    if (visible) {
      const bestVideo = getBestVideo();
      setSelectedVideo(bestVideo);
      setIsLoading(true);
      setError(null);
    }
  }, [visible, asset]);

  useEffect(() => {
    if (!selectedVideo || Platform.OS !== 'web' || !videoRef.current) {
      return;
    }

    const videoUrl = getVideoUrl(selectedVideo);
    if (!videoUrl) {
      setError('No video URL available');
      setIsLoading(false);
      return;
    }

    const videoElement = videoRef.current;
    
    if (selectedVideo.format === 'hls') {
      if (Hls.isSupported()) {
        if (hlsRef.current) {
          hlsRef.current.destroy();
        }
        
        const hls = new Hls();
        hlsRef.current = hls;
        
        hls.loadSource(videoUrl);
        hls.attachMedia(videoElement);
        
        hls.on(Hls.Events.MANIFEST_PARSED, () => {
          setIsLoading(false);
          setIsPlaying(true);
        });
        
        hls.on(Hls.Events.ERROR, (event, data) => {
          setError(`HLS error: ${data.type} - ${data.details}`);
          setIsLoading(false);
        });
      } else if (videoElement.canPlayType('application/vnd.apple.mpegurl')) {
        videoElement.src = videoUrl;
        videoElement.addEventListener('loadedmetadata', () => {
          setIsLoading(false);
          setIsPlaying(true);
        });
      } else {
        setError('HLS is not supported in this browser');
        setIsLoading(false);
      }
    } else if (selectedVideo.format === 'dash') {
      if (dashRef.current) {
        dashRef.current.destroy();
      }
      
      const dash = dashjs.MediaPlayer().create();
      dashRef.current = dash;
      
      dash.initialize(videoElement, videoUrl, false);
      
      dash.on(dashjs.MediaPlayer.events.MANIFEST_LOADED, () => {
        setIsLoading(false);
        setIsPlaying(true);
      });
      
      dash.on(dashjs.MediaPlayer.events.ERROR, (error: any) => {
        setError(`DASH error: ${error.error}`);
        setIsLoading(false);
      });
    } else {
      videoElement.src = videoUrl;
      videoElement.addEventListener('loadedmetadata', () => {
        setIsLoading(false);
        setIsPlaying(true);
      });
      videoElement.addEventListener('error', () => {
        setError('Failed to load video');
        setIsLoading(false);
      });
    }
    
    return () => {
      if (hlsRef.current) {
        hlsRef.current.destroy();
        hlsRef.current = null;
      }
      if (dashRef.current) {
        dashRef.current.destroy();
        dashRef.current = null;
      }
    };
  }, [selectedVideo]);

  const handleLoadStart = () => {
    setIsLoading(true);
    setError(null);
  };

  const handleLoad = () => {
    setIsLoading(false);
    setIsPlaying(true);
  };

  const handleError = (error: any) => {
    setIsLoading(false);
    let errorMessage = 'Failed to load video';
    
    if (error?.message) {
      errorMessage = error.message;
    } else if (error?.error) {
      errorMessage = error.error;
    } else if (typeof error === 'string') {
      errorMessage = error;
    }
    
    setError(errorMessage);
  };

  const togglePlayPause = () => {
    setIsPlaying(!isPlaying);
  };

  const getFormatLabel = (format: string) => {
    switch (format.toLowerCase()) {
      case 'hls':
        return 'HLS';
      case 'dash':
        return 'DASH';
      case 'raw':
        return 'RAW';
      default:
        return format.toUpperCase();
    }
  };

  const getStatusColor = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'ready':
        return '#4CAF50';
      case 'pending':
      case 'analyzing':
      case 'transcoding':
        return '#2196F3';
      case 'failed':
        return '#F44336';
      default:
        return '#9E9E9E';
    }
  };

  if (!selectedVideo) {
    return (
      <Modal visible={visible} transparent animationType="fade">
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>{asset.title || 'Video Player'}</Text>
              <TouchableOpacity onPress={onClose} style={styles.closeButton}>
                <Text style={styles.closeButtonText}>✕</Text>
              </TouchableOpacity>
            </View>
            <View style={styles.errorContainer}>
              <Text style={styles.errorText}>No playable video available</Text>
              <Text style={styles.errorSubtext}>This asset has no videos ready for playback</Text>
            </View>
          </View>
        </View>
      </Modal>
    );
  }

  const videoUrl = getVideoUrl(selectedVideo);

  return (
    <Modal visible={visible} transparent animationType="fade">
      <View style={styles.modalOverlay}>
        <View style={styles.modalContent}>
          <View style={styles.modalHeader}>
            <View style={styles.titleContainer}>
              <Text style={styles.modalTitle}>{asset.title || 'Video Player'}</Text>
              <View style={styles.videoInfo}>
                <Text style={styles.formatLabel}>{getFormatLabel(selectedVideo.format)}</Text>
                <View style={[styles.statusDot, { backgroundColor: getStatusColor(selectedVideo.status) }]} />
                <Text style={styles.statusText}>{selectedVideo.status || 'unknown'}</Text>
              </View>
            </View>
            <TouchableOpacity onPress={onClose} style={styles.closeButton}>
              <Text style={styles.closeButtonText}>✕</Text>
            </TouchableOpacity>
          </View>
          
          <View style={styles.videoContainer}>
            {isLoading && (
              <View style={styles.loadingContainer}>
                <Text style={styles.loadingText}>Loading video...</Text>
              </View>
            )}
            
            {error && (
              <View style={styles.errorContainer}>
                <Text style={styles.errorText}>{error}</Text>
                <TouchableOpacity style={styles.retryButton} onPress={() => setIsLoading(true)}>
                  <Text style={styles.retryButtonText}>Retry</Text>
                </TouchableOpacity>
              </View>
            )}
            
            {Platform.OS === 'web' ? (
              <video
                ref={videoRef}
                style={styles.video}
                controls
                onEnded={() => setIsPlaying(false)}
              />
            ) : (
              <Video
                source={{ uri: videoUrl }}
                style={styles.video}
                resizeMode={ResizeMode.CONTAIN}
                shouldPlay={isPlaying}
                useNativeControls
                isLooping={false}
                onLoadStart={handleLoadStart}
                onLoad={handleLoad}
                onError={handleError}
                onPlaybackStatusUpdate={status => {
                  if (status.isLoaded) {
                    setIsLoading(false);
                    if (status.didJustFinish) {
                      setIsPlaying(false);
                    }
                  } else if (status.error) {
                    setError(`Playback error: ${status.error}`);
                  }
                }}
              />
            )}
            
            {!isLoading && !error && (
              <TouchableOpacity style={styles.playPauseButton} onPress={togglePlayPause}>
                <Text style={styles.playPauseIcon}>
                  {isPlaying ? '⏸' : '▶'}
                </Text>
              </TouchableOpacity>
            )}
          </View>
          
          <ScrollView style={styles.assetDetails}>
            <View style={styles.detailsSection}>
              <Text style={styles.sectionTitle}>Asset Information</Text>
              
              {asset.description && (
                <View style={styles.detailRow}>
                  <Text style={styles.detailLabel}>Description:</Text>
                  <Text style={styles.detailValue}>{asset.description}</Text>
                </View>
              )}
              
              {asset.genre && (
                <View style={styles.detailRow}>
                  <Text style={styles.detailLabel}>Genre:</Text>
                  <Text style={styles.detailValue}>{getGenreLabel(asset.genre)}</Text>
                </View>
              )}
              
              {asset.genres && asset.genres.length > 0 && (
                <View style={styles.detailRow}>
                  <Text style={styles.detailLabel}>Genres:</Text>
                  <Text style={styles.detailValue}>{asset.genres.join(', ')}</Text>
                </View>
              )}
              
              {asset.publishRule?.ageRating && (
                <View style={styles.detailRow}>
                  <Text style={styles.detailLabel}>Age Rating:</Text>
                  <Text style={styles.detailValue}>{getAgeRatingLabel(asset.publishRule.ageRating)}</Text>
                </View>
              )}
              
              {asset.publishRule?.isPublic !== undefined && (
                <View style={styles.detailRow}>
                  <Text style={styles.detailLabel}>Public:</Text>
                  <Text style={[styles.detailValue, { color: asset.publishRule.isPublic ? '#4CAF50' : '#F44336' }]}>
                    {asset.publishRule.isPublic ? 'Yes' : 'No'}
                  </Text>
                </View>
              )}
            </View>
            
            <View style={styles.detailsSection}>
              <Text style={styles.sectionTitle}>Video Details</Text>
              
              {selectedVideo.duration && (
                <View style={styles.detailRow}>
                  <Text style={styles.detailLabel}>Duration:</Text>
                  <Text style={styles.detailValue}>{Math.round(selectedVideo.duration)}s</Text>
                </View>
              )}
              
              {selectedVideo.size && (
                <View style={styles.detailRow}>
                  <Text style={styles.detailLabel}>Size:</Text>
                  <Text style={styles.detailValue}>{Math.round(selectedVideo.size / 1024 / 1024)} MB</Text>
                </View>
              )}
              
              {selectedVideo.width && selectedVideo.height && (
                <View style={styles.detailRow}>
                  <Text style={styles.detailLabel}>Resolution:</Text>
                  <Text style={styles.detailValue}>{selectedVideo.width}x{selectedVideo.height}</Text>
                </View>
              )}
              
              {videoUrl && (
                <View style={styles.detailRow}>
                  <Text style={styles.detailLabel}>URL:</Text>
                  <Text style={styles.detailValue} numberOfLines={2}>
                    {videoUrl}
                  </Text>
                </View>
              )}
            </View>
          </ScrollView>
        </View>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.8)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  modalContent: {
    backgroundColor: '#1a1a1a',
    borderRadius: 12,
    width: screenWidth * 0.95,
    maxHeight: screenHeight * 0.9,
    overflow: 'hidden',
  },
  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#333',
  },
  titleContainer: {
    flex: 1,
  },
  modalTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#ffffff',
  },
  videoInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: 4,
  },
  formatLabel: {
    fontSize: 12,
    color: '#007AFF',
    fontWeight: '600',
    marginRight: 8,
  },
  statusDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginRight: 4,
  },
  statusText: {
    fontSize: 12,
    color: '#999',
  },
  closeButton: {
    padding: 8,
  },
  closeButtonText: {
    fontSize: 20,
    color: '#999',
    fontWeight: 'bold',
  },
  videoContainer: {
    position: 'relative',
    backgroundColor: '#000',
    aspectRatio: 16 / 9,
  },
  video: {
    flex: 1,
  },
  loadingContainer: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#000',
  },
  loadingText: {
    color: '#fff',
    fontSize: 16,
  },
  errorContainer: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#000',
  },
  errorText: {
    color: '#fff',
    fontSize: 16,
    textAlign: 'center',
    marginBottom: 8,
  },
  errorSubtext: {
    color: '#ccc',
    fontSize: 14,
    textAlign: 'center',
  },
  retryButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 6,
    marginTop: 12,
  },
  retryButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  playPauseButton: {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: [{ translateX: -25 }, { translateY: -25 }],
    backgroundColor: 'rgba(0, 0, 0, 0.6)',
    borderRadius: 50,
    width: 80,
    height: 80,
    justifyContent: 'center',
    alignItems: 'center',
  },
  playPauseIcon: {
    fontSize: 32,
    color: '#fff',
  },
  assetDetails: {
    maxHeight: 300,
  },
  detailsSection: {
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#333',
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#ffffff',
    marginBottom: 12,
  },
  detailRow: {
    flexDirection: 'row',
    marginBottom: 8,
    alignItems: 'flex-start',
  },
  detailLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: '#999',
    width: 100,
    flexShrink: 0,
  },
  detailValue: {
    fontSize: 14,
    color: '#ffffff',
    flex: 1,
    lineHeight: 20,
  },
}); 