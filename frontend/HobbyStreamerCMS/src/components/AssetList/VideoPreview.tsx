import React, { useState, useRef, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Modal, Dimensions, Platform } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Video, ResizeMode } from 'expo-av';
import Hls from 'hls.js';
import * as dashjs from 'dashjs';
import { API_CONFIG } from '../../config/api';

interface VideoPreviewProps {
  video: any;
  visible: boolean;
  onClose: () => void;
}

const { width: screenWidth, height: screenHeight } = Dimensions.get('window');

const convertS3UrlToHttp = (url: string): string => {
  if (!url) return url;
  
  if (url.startsWith('s3://')) {
    const s3Path = url.replace('s3://', '');
    const [bucket, ...keyParts] = s3Path.split('/');
    const key = keyParts.join('/');
    return `${API_CONFIG.CDN_BASE_URL}/${key}`;
  }
  
  return url;
};

export default function VideoPreview({ video, visible, onClose }: VideoPreviewProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const videoRef = useRef<HTMLVideoElement>(null);
  const hlsRef = useRef<Hls | null>(null);
  const dashRef = useRef<dashjs.MediaPlayerClass | null>(null);

  const getVideoUrl = () => {
    if (video?.streamInfo?.url) return video.streamInfo.url;
    if (video?.streamInfo?.cdnPrefix && video?.storageLocation?.key) {
      const prefix = video.streamInfo.cdnPrefix.replace(/\/$/, '');
      const key = video.storageLocation.key.replace(/^\//, '');
      return `${prefix}/${key}`;
    }
    if (video?.storageLocation?.url) return convertS3UrlToHttp(video.storageLocation.url);
    if (video?.storageLocation?.key) return `${API_CONFIG.CDN_BASE_URL}/${video.storageLocation.key}`;
    return null;
  };

  const videoUrl = getVideoUrl();

  useEffect(() => {
    if (Platform.OS === 'web' && videoUrl && videoRef.current) {
      const videoElement = videoRef.current;
      
      if (video.format === 'hls') {
        if (Hls.isSupported()) {
          if (hlsRef.current) {
            hlsRef.current.destroy();
          }
          
          const hls = new Hls();
          hlsRef.current = hls;
          
          hls.loadSource(videoUrl);
          hls.attachMedia(videoElement);
          
          hls.on(Hls.Events.MANIFEST_PARSED, () => {
            console.log('HLS manifest parsed successfully');
            setIsLoading(false);
            setIsPlaying(true);
          });
          
          hls.on(Hls.Events.ERROR, (event, data) => {
            console.error('HLS error:', data);
            setError(`HLS error: ${data.type} - ${data.details}`);
          });
        } else if (videoElement.canPlayType('application/vnd.apple.mpegurl')) {
          videoElement.src = videoUrl;
          videoElement.addEventListener('loadedmetadata', () => {
            setIsLoading(false);
            setIsPlaying(true);
          });
        } else {
          setError('HLS is not supported in this browser');
        }
      } else if (video.format === 'dash') {
        if (dashRef.current) {
          dashRef.current.destroy();
        }
        
        const dash = dashjs.MediaPlayer().create();
        dashRef.current = dash;
        
        dash.initialize(videoElement, videoUrl, false);
        
        dash.on(dashjs.MediaPlayer.events.MANIFEST_LOADED, () => {
          console.log('DASH manifest loaded successfully');
          setIsLoading(false);
          setIsPlaying(true);
        });
        
        dash.on(dashjs.MediaPlayer.events.ERROR, (error: any) => {
          console.error('DASH error:', error);
          setError(`DASH error: ${error.error}`);
        });
      } else {
        videoElement.src = videoUrl;
        videoElement.addEventListener('loadedmetadata', () => {
          console.log('Video loaded successfully');
          setIsLoading(false);
          setIsPlaying(true);
        });
        videoElement.addEventListener('error', (e: Event) => {
          console.error('Video error:', e);
          setError('Failed to load video');
        });
      }
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
  }, [videoUrl, video.format]);

  const handleLoadStart = () => {
    console.log('Video load started:', videoUrl);
    setIsLoading(true);
    setError(null);
  };

  const handleLoad = () => {
    console.log('Video loaded successfully:', videoUrl);
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
    } else if (error?.code) {
      switch (error.code) {
        case 'NETWORK_ERROR':
          errorMessage = 'Network error: Unable to connect to video server';
          break;
        case 'MEDIA_ERROR':
          errorMessage = 'Media error: Video format not supported or corrupted';
          break;
        case 'DECODE_ERROR':
          errorMessage = 'Decode error: Video cannot be decoded';
          break;
        default:
          errorMessage = `Video error (${error.code}): ${error.message || 'Unknown error'}`;
      }
    }
    
    setError(errorMessage);
    console.error('Video error details:', {
      error,
      videoUrl,
      videoObject: video,
      errorCode: error?.code,
      errorMessage: error?.message
    });
  };

  const handleEnd = () => {
    setIsPlaying(false);
  };

  const togglePlayPause = () => {
    setIsPlaying(!isPlaying);
  };

  const getVideoTypeLabel = (type: string) => {
    switch (type) {
      case 'main':
        return 'Main Video';
      case 'trailer':
        return 'Trailer';
      case 'behind':
        return 'Behind the Scenes';
      case 'interview':
        return 'Interview';
      default:
        return type;
    }
  };

  if (video.format === 'raw') {
    return (
      <Modal visible={visible} transparent animationType="fade">
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>Video Preview</Text>
              <TouchableOpacity onPress={onClose} style={styles.closeButton}>
                <Ionicons name="close" size={24} color="#666" />
              </TouchableOpacity>
            </View>
            <View style={styles.errorContainer}>
              <Ionicons name="warning" size={48} color="#ff6b6b" />
              <Text style={styles.errorText}>Raw video cannot be previewed</Text>
              <Text style={styles.errorSubtext}>Please transcode to DASH or HLS format first</Text>
            </View>
          </View>
        </View>
      </Modal>
    );
  }

  if (!videoUrl) {
    return (
      <Modal visible={visible} transparent animationType="fade">
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>Video Preview</Text>
              <TouchableOpacity onPress={onClose} style={styles.closeButton}>
                <Ionicons name="close" size={24} color="#666" />
              </TouchableOpacity>
            </View>
            <View style={styles.errorContainer}>
              <Ionicons name="warning" size={48} color="#ff6b6b" />
              <Text style={styles.errorText}>No video URL available</Text>
              <Text style={styles.errorSubtext}>This video cannot be previewed</Text>
            </View>
          </View>
        </View>
      </Modal>
    );
  }

  // Native DASH preview is not widely supported in expo-av; guide user to use HLS
  if (Platform.OS !== 'web' && video.format === 'dash') {
    return (
      <Modal visible={visible} transparent animationType="fade">
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>Video Preview</Text>
              <TouchableOpacity onPress={onClose} style={styles.closeButton}>
                <Ionicons name="close" size={24} color="#666" />
              </TouchableOpacity>
            </View>
            <View style={styles.errorContainer}>
              <Ionicons name="warning" size={48} color="#ff6b6b" />
              <Text style={styles.errorText}>DASH preview is not supported on this device</Text>
              <Text style={styles.errorSubtext}>Please use an HLS rendition or open on web</Text>
            </View>
          </View>
        </View>
      </Modal>
    );
  }

  return (
    <Modal visible={visible} transparent animationType="fade">
      <View style={styles.modalOverlay}>
        <View style={styles.modalContent}>
          <View style={styles.modalHeader}>
            <View style={styles.titleContainer}>
              <Text style={styles.modalTitle}>
                {getVideoTypeLabel(video.type)} - {video.format.toUpperCase()}
              </Text>
              <Text style={styles.videoId}>ID: {video.id}</Text>
            </View>
            <TouchableOpacity onPress={onClose} style={styles.closeButton}>
              <Ionicons name="close" size={24} color="#666" />
            </TouchableOpacity>
          </View>
          
          <View style={styles.videoContainer}>
            {isLoading && (
              <View style={styles.loadingContainer}>
                <Ionicons name="sync" size={32} color="#007AFF" />
                <Text style={styles.loadingText}>Loading video...</Text>
              </View>
            )}
            
            {error && (
              <View style={styles.errorContainer}>
                <Ionicons name="warning" size={48} color="#ff6b6b" />
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
                  console.log('Playback status update:', status);
                  if (status.isLoaded) {
                    setIsLoading(false);
                    if (status.didJustFinish) {
                      setIsPlaying(false);
                    }
                  } else if (status.error) {
                    console.error('Playback error:', status.error);
                    setError(`Playback error: ${status.error}`);
                  }
                }}
              />
            )}
            
            {!isLoading && !error && (
              <TouchableOpacity style={styles.playPauseButton} onPress={togglePlayPause}>
                <Ionicons 
                  name={isPlaying ? 'pause' : 'play'} 
                  size={32} 
                  color="#fff" 
                />
              </TouchableOpacity>
            )}
          </View>
          
          <View style={styles.videoInfo}>
            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Status:</Text>
              <Text style={[styles.infoValue, { color: video.status === 'ready' ? '#4CAF50' : '#666' }]}>
                {video.status || 'unknown'}
              </Text>
            </View>
            {video.size && (
              <View style={styles.infoRow}>
                <Text style={styles.infoLabel}>Size:</Text>
                <Text style={styles.infoValue}>{Math.round(video.size / 1024 / 1024)} MB</Text>
              </View>
            )}
            {video.duration && (
              <View style={styles.infoRow}>
                <Text style={styles.infoLabel}>Duration:</Text>
                <Text style={styles.infoValue}>{Math.round(video.duration)}s</Text>
              </View>
            )}
            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>URL:</Text>
              <Text style={styles.infoValue} numberOfLines={2}>
                {videoUrl}
              </Text>
            </View>
          </View>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.8)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  modalContent: {
    backgroundColor: '#fff',
    borderRadius: 12,
    width: screenWidth * 0.9,
    maxHeight: screenHeight * 0.8,
    overflow: 'hidden',
  },
  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  titleContainer: {
    flex: 1,
  },
  modalTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
  },
  videoId: {
    fontSize: 12,
    color: '#666',
    marginTop: 2,
  },
  closeButton: {
    padding: 4,
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
    marginTop: 8,
    fontSize: 14,
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
    marginTop: 8,
    textAlign: 'center',
  },
  errorSubtext: {
    color: '#ccc',
    fontSize: 14,
    marginTop: 4,
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
    transform: [{ translateX: -20 }, { translateY: -20 }],
    backgroundColor: 'rgba(0, 0, 0, 0.6)',
    borderRadius: 40,
    width: 80,
    height: 80,
    justifyContent: 'center',
    alignItems: 'center',
  },
  videoInfo: {
    padding: 16,
  },
  infoRow: {
    flexDirection: 'row',
    marginBottom: 8,
  },
  infoLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333',
    width: 80,
  },
  infoValue: {
    fontSize: 14,
    color: '#666',
    flex: 1,
  },
}); 