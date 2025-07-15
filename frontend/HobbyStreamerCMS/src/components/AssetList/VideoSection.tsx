import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Alert, ScrollView } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import VideoUpload from './VideoUpload';
import { triggerTranscodeJob } from '../../services/api';
import { VideoType, VideoFormat } from '../../types/asset';

interface VideoSectionProps {
  videos: any[] | undefined;
  onDeleteVideo: (videoId: string) => void;
  onUpdate: (field: string, value: any) => void;
  assetId: string;
  onUploadComplete: () => void;
  onRefreshAsset: () => void;
}

export default function VideoSection({ videos, onDeleteVideo, onUpdate, assetId, onUploadComplete, onRefreshAsset }: VideoSectionProps) {
  const [showMainUpload, setShowMainUpload] = useState(false);
  const [showTrailerUpload, setShowTrailerUpload] = useState(false);
  const [triggeringJobs, setTriggeringJobs] = useState<{[key: string]: boolean}>({});

  const handleTriggerTranscode = async (videoId: string, format: 'hls' | 'dash') => {
    const jobKey = `${videoId}-${format}`;
    setTriggeringJobs(prev => ({ ...prev, [jobKey]: true }));
    
    try {
      const video = videos?.find(v => v.id === videoId);
      if (!video?.storageLocation) {
        throw new Error('Video storage location not found');
      }

      const storageLocation = video.storageLocation;
      const keyParts = storageLocation.key.split('/');
      const sourceFileName = keyParts[keyParts.length - 1];
      
      if (!sourceFileName) {
        throw new Error('Invalid video file path - missing filename');
      }
      
      await triggerTranscodeJob(assetId, videoId, format, {
        bucket: storageLocation.bucket,
        key: storageLocation.key,
        sourceFileName
      });
      Alert.alert('Success', `${format.toUpperCase()} transcoding job triggered successfully`);
      onRefreshAsset();
    } catch (error: any) {
      console.error('Failed to trigger transcoding job:', error);
      Alert.alert('Error', `Failed to trigger ${format.toUpperCase()} transcoding job: ${error.message || 'Unknown error'}`);
    } finally {
      setTriggeringJobs(prev => ({ ...prev, [jobKey]: false }));
    }
  };

  const getVideoIcon = (format: string) => {
    switch (format.toLowerCase()) {
      case 'hls':
        return 'play-circle';
      case 'dash':
        return 'play';
      case 'raw':
        return 'videocam';
      default:
        return 'videocam';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'ready':
        return 'checkmark-circle';
      case 'pending':
      case 'analyzing':
      case 'transcoding':
        return 'sync';
      case 'failed':
        return 'close-circle';
      default:
        return 'help-circle';
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

  const getVideoTypeLabel = (type: string) => {
    switch (type) {
      case 'MAIN':
        return 'Main Video';
      case 'TRAILER':
        return 'Trailer';
      case 'BEHIND_THE_SCENES':
        return 'Behind the Scenes';
      case 'INTERVIEW':
        return 'Interview';
      default:
        return type;
    }
  };

  const renderVideoItem = (video: any) => {
    const jobKey = `${video.id}-${video.format}`;
    const isTriggering = triggeringJobs[jobKey] || false;
    const canTrigger = video.format === 'raw' && (video.type === 'MAIN' || video.type === 'TRAILER');
    
    return (
      <View key={video.id} style={styles.videoItem}>
        <View style={styles.videoHeader}>
          <View style={styles.videoInfo}>
            <Ionicons name={getVideoIcon(video.format) as any} size={18} color="#007AFF" />
            <Text style={styles.videoFormat}>{video.format.toUpperCase()}</Text>
            <Text style={styles.videoType}>({getVideoTypeLabel(video.type)})</Text>
          </View>
          <View style={styles.videoActions}>
            {canTrigger && (
              <View style={styles.transcodeButtons}>
                <TouchableOpacity 
                  style={[styles.transcodeButton, isTriggering && styles.transcodeButtonDisabled]}
                  onPress={() => handleTriggerTranscode(video.id, 'hls')}
                  disabled={isTriggering}
                >
                  <Ionicons 
                    name={isTriggering ? 'sync' : 'play'} 
                    size={12} 
                    color={isTriggering ? '#999' : '#fff'} 
                  />
                  <Text style={[styles.transcodeButtonText, isTriggering && styles.transcodeButtonTextDisabled]}>
                    HLS
                  </Text>
                </TouchableOpacity>
                <TouchableOpacity 
                  style={[styles.transcodeButton, isTriggering && styles.transcodeButtonDisabled]}
                  onPress={() => handleTriggerTranscode(video.id, 'dash')}
                  disabled={isTriggering}
                >
                  <Ionicons 
                    name={isTriggering ? 'sync' : 'play'} 
                    size={12} 
                    color={isTriggering ? '#999' : '#fff'} 
                  />
                  <Text style={[styles.transcodeButtonText, isTriggering && styles.transcodeButtonTextDisabled]}>
                    DASH
                  </Text>
                </TouchableOpacity>
              </View>
            )}
            <TouchableOpacity 
              style={styles.deleteButton}
              onPress={() => onDeleteVideo(video.id)}
            >
              <Ionicons name="trash" size={16} color="#ff3b30" />
            </TouchableOpacity>
          </View>
        </View>
        
        <View style={styles.videoStatus}>
          <Ionicons 
            name={getStatusIcon(video.status) as any} 
            size={14} 
            color={getStatusColor(video.status)} 
          />
          <Text style={styles.videoStatusText}>
            Status: {video.status || 'unknown'}
          </Text>
        </View>
        
        <View style={styles.videoDetail}>
          <Ionicons name="link" size={12} color="#666" />
          <Text style={styles.videoDetailText}>
            {video.storageLocation?.url || 'No URL available'}
          </Text>
        </View>

        {video.size && (
          <View style={styles.videoDetail}>
            <Ionicons name="information-circle" size={12} color="#666" />
            <Text style={styles.videoDetailText}>
              Size: {Math.round(video.size / 1024 / 1024)} MB
            </Text>
          </View>
        )}

        {video.duration && (
          <View style={styles.videoDetail}>
            <Ionicons name="time" size={12} color="#666" />
            <Text style={styles.videoDetailText}>
              Duration: {Math.round(video.duration)}s
            </Text>
          </View>
        )}
      </View>
    );
  };

  const renderVideoTypeSection = (videoType: VideoType, title: string, icon: string) => {
    const typeVideos = videos?.filter(v => v.type === videoType) || [];
    const hasVideos = typeVideos.length > 0;

    return (
      <View style={styles.videoTypeContainer}>
        <View style={styles.videoTypeHeader}>
          <View style={styles.videoTypeInfo}>
            <Ionicons name={icon as any} size={20} color="#333" />
            <Text style={styles.videoTypeTitle}>{title}</Text>
            <Text style={styles.videoCount}>({typeVideos.length} videos)</Text>
          </View>
          <TouchableOpacity 
            style={styles.addButton}
            onPress={() => {
              if (videoType === VideoType.MAIN) {
                setShowMainUpload(!showMainUpload);
                setShowTrailerUpload(false);
              } else {
                setShowTrailerUpload(!showTrailerUpload);
                setShowMainUpload(false);
              }
            }}
          >
            <Ionicons 
              name={((videoType === VideoType.MAIN && showMainUpload) || (videoType === VideoType.TRAILER && showTrailerUpload)) ? 'remove' : 'add'} 
              size={16} 
              color="#fff" 
            />
          </TouchableOpacity>
        </View>
        
        {hasVideos ? (
          <View style={styles.videosList}>
            {typeVideos.map(renderVideoItem)}
          </View>
        ) : (
          <View style={styles.emptyState}>
            <Ionicons name="videocam-outline" size={24} color="#ccc" />
            <Text style={styles.emptyText}>No {title.toLowerCase()} uploaded yet</Text>
          </View>
        )}
        
        {((videoType === VideoType.MAIN && showMainUpload) || (videoType === VideoType.TRAILER && showTrailerUpload)) && (
          <VideoUpload
            assetId={assetId}
            videoType={videoType}
            onUploadComplete={() => {
              onUploadComplete();
              if (videoType === VideoType.MAIN) {
                setShowMainUpload(false);
              } else {
                setShowTrailerUpload(false);
              }
            }}
            onCancel={() => {
              if (videoType === VideoType.MAIN) {
                setShowMainUpload(false);
              } else {
                setShowTrailerUpload(false);
              }
            }}
            onRefreshAsset={onRefreshAsset}
          />
        )}
      </View>
    );
  };

  return (
    <View style={styles.container}>
      <View style={styles.sectionHeader}>
        <View style={styles.titleContainer}>
          <Ionicons name="videocam" size={20} color="#333" />
          <Text style={styles.sectionTitle}>Videos</Text>
        </View>
      </View>
      
      <ScrollView style={styles.videosContainer}>
        {renderVideoTypeSection(VideoType.MAIN, 'Main Video', 'videocam')}
        {renderVideoTypeSection(VideoType.TRAILER, 'Trailer', 'play-circle')}
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: '#fff',
    borderRadius: 8,
    padding: 16,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  titleContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
  },
  videosContainer: {
    maxHeight: 600,
  },
  videoTypeContainer: {
    backgroundColor: '#f8f9fa',
    borderRadius: 8,
    padding: 16,
    marginBottom: 16,
  },
  videoTypeHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  videoTypeInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  videoTypeTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
  },
  videoCount: {
    fontSize: 14,
    color: '#666',
  },
  addButton: {
    backgroundColor: '#007AFF',
    width: 28,
    height: 28,
    borderRadius: 14,
    justifyContent: 'center',
    alignItems: 'center',
  },
  videosList: {
    gap: 12,
  },
  videoItem: {
    backgroundColor: '#fff',
    borderRadius: 6,
    padding: 12,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  videoHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  videoInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  videoFormat: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333',
  },
  videoType: {
    fontSize: 12,
    color: '#666',
  },
  videoActions: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  transcodeButtons: {
    flexDirection: 'row',
    gap: 4,
  },
  transcodeButton: {
    backgroundColor: '#28a745',
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 6,
    paddingVertical: 3,
    borderRadius: 3,
    gap: 2,
  },
  transcodeButtonDisabled: {
    backgroundColor: '#e9ecef',
  },
  transcodeButtonText: {
    fontSize: 10,
    color: '#fff',
    fontWeight: '600',
  },
  transcodeButtonTextDisabled: {
    color: '#999',
  },
  deleteButton: {
    padding: 4,
  },
  videoStatus: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    marginBottom: 6,
  },
  videoStatusText: {
    fontSize: 14,
    color: '#666',
  },
  videoDetail: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    marginBottom: 2,
  },
  videoDetailText: {
    fontSize: 12,
    color: '#666',
    flex: 1,
  },
  emptyState: {
    alignItems: 'center',
    padding: 20,
  },
  emptyText: {
    fontSize: 14,
    color: '#666',
    marginTop: 8,
    textAlign: 'center',
  },
}); 