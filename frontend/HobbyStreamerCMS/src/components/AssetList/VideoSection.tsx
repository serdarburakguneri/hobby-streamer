import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import VideoUpload from './VideoUpload';
import { triggerTranscodeJob } from '../../services/api';
import { VideoType } from '../../types/asset';

interface VideoSectionProps {
  videos: any[] | undefined;
  onDeleteVideo: (videoType: string, videoName: string) => void;
  onUpdate: (field: string, value: any) => void;
  assetId: string;
  onUploadComplete: () => void;
  onRefreshAsset: () => void;
}

export default function VideoSection({ videos, onDeleteVideo, onUpdate, assetId, onUploadComplete, onRefreshAsset }: VideoSectionProps) {
  const [showMainUpload, setShowMainUpload] = useState(false);
  const [showTrailerUpload, setShowTrailerUpload] = useState(false);
  const [triggeringJobs, setTriggeringJobs] = useState<{[key: string]: boolean}>({});

  const handleTriggerTranscode = async (videoType: string, format: 'hls' | 'dash') => {
    const jobKey = `${videoType}-${format}`;
    setTriggeringJobs(prev => ({ ...prev, [jobKey]: true }));
    
    try {
      const video = videos?.find(v => v.type === videoType);
      if (!video?.raw?.storageLocation) {
        throw new Error('Raw video not found for this video type');
      }

      const rawVideo = video.raw.storageLocation;
      const keyParts = rawVideo.key.split('/');
      const sourceFileName = keyParts[keyParts.length - 1];
      
      if (!sourceFileName) {
        throw new Error('Invalid video file path - missing filename');
      }
      
      await triggerTranscodeJob(assetId, videoType.toLowerCase(), format, {
        bucket: rawVideo.bucket,
        key: rawVideo.key,
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

  const getVideoIcon = (type: string) => {
    switch (type.toLowerCase()) {
      case 'hls':
        return 'play-circle';
      case 'dash':
        return 'play';
      case 'raw':
        return 'videocam';
      case 'thumbnail':
        return 'image';
      default:
        return 'videocam';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'ready':
        return 'checkmark-circle';
      case 'processing':
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
      case 'processing':
        return '#2196F3';
      case 'failed':
        return '#F44336';
      default:
        return '#9E9E9E';
    }
  };

  const renderVideoFormat = (format: string, videoData: any, icon: string, videoType: string) => {
    const hasData = videoData && videoData.storageLocation;
    const status = hasData ? videoData.status : 'not_available';
    const canTrigger = format === 'hls' || format === 'dash';
    const jobKey = `${videoType}-${format}`;
    const isTriggering = triggeringJobs[jobKey] || false;
    
    return (
      <View style={styles.formatSection}>
        <View style={styles.formatHeader}>
          <View style={styles.formatInfo}>
            <Ionicons name={icon as any} size={18} color="#007AFF" />
            <Text style={styles.formatLabel}>{format.toUpperCase()}</Text>
          </View>
          {canTrigger && (
            <TouchableOpacity 
              style={[styles.triggerButton, isTriggering && styles.triggerButtonDisabled]}
              onPress={() => handleTriggerTranscode(videoType, format as 'hls' | 'dash')}
              disabled={isTriggering}
            >
              <Ionicons 
                name={isTriggering ? 'sync' : 'play'} 
                size={14} 
                color={isTriggering ? '#999' : '#fff'} 
              />
              <Text style={[styles.triggerButtonText, isTriggering && styles.triggerButtonTextDisabled]}>
                {isTriggering ? 'Converting...' : 'Convert'}
              </Text>
            </TouchableOpacity>
          )}
        </View>
        
        <View style={styles.videoStatus}>
          <Ionicons 
            name={getStatusIcon(status) as any} 
            size={14} 
            color={getStatusColor(status)} 
          />
          <Text style={styles.videoStatusText}>
            Status: {hasData ? status : 'No ' + format.toUpperCase() + ' format available'}
          </Text>
        </View>
        
        {hasData && (
          <View style={styles.videoDetail}>
            <Ionicons name="link" size={12} color="#666" />
            <Text style={styles.videoDetailText}>
              {videoData.storageLocation.url}
            </Text>
          </View>
        )}
      </View>
    );
  };

  const renderVideoContent = (videoType: VideoType, title: string, icon: string) => {
    const video = videos?.find(v => v.type === videoType);
    const hasVideo = video && (video.raw || video.hls || video.dash || video.thumbnail);

    return (
      <View style={styles.videoTypeContainer}>
        <View style={styles.videoTypeHeader}>
          <View style={styles.videoTypeInfo}>
            <Ionicons name={icon as any} size={20} color="#333" />
            <Text style={styles.videoTypeTitle}>{title}</Text>
          </View>
          <View style={styles.headerActions}>
            {hasVideo && (
              <TouchableOpacity 
                style={styles.deleteVideoButton}
                onPress={() => onDeleteVideo(videoType, videoType)}
              >
                <Ionicons name="trash" size={16} color="#ff3b30" />
              </TouchableOpacity>
            )}
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
        </View>
        
        {hasVideo ? (
          <View style={styles.formatsContainer}>
            {renderVideoFormat('raw', video.raw, 'videocam', videoType)}
            {renderVideoFormat('hls', video.hls, 'play-circle', videoType)}
            {renderVideoFormat('dash', video.dash, 'play', videoType)}
            {video.thumbnail && (
              <View style={styles.formatSection}>
                <View style={styles.formatHeader}>
                  <View style={styles.formatInfo}>
                    <Ionicons name="image" size={18} color="#007AFF" />
                    <Text style={styles.formatLabel}>THUMBNAIL</Text>
                  </View>
                </View>
                <View style={styles.videoDetail}>
                  <Ionicons name="image" size={12} color="#666" />
                  <Text style={styles.videoDetailText}>
                    {video.thumbnail.url}
                  </Text>
                </View>
              </View>
            )}
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
      
      <View style={styles.videosContainer}>
        {renderVideoContent(VideoType.MAIN, 'Main Video', 'videocam')}
        {renderVideoContent(VideoType.TRAILER, 'Trailer', 'play-circle')}
      </View>
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
    gap: 20,
  },
  videoTypeContainer: {
    backgroundColor: '#f8f9fa',
    borderRadius: 8,
    padding: 16,
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
  headerActions: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  deleteVideoButton: {
    padding: 4,
  },
  addButton: {
    backgroundColor: '#007AFF',
    width: 28,
    height: 28,
    borderRadius: 14,
    justifyContent: 'center',
    alignItems: 'center',
  },
  formatsContainer: {
    gap: 12,
  },
  formatSection: {
    backgroundColor: '#fff',
    borderRadius: 6,
    padding: 12,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  formatHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  formatInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  formatLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333',
  },
  triggerButton: {
    backgroundColor: '#28a745',
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
    gap: 4,
  },
  triggerButtonDisabled: {
    backgroundColor: '#e9ecef',
  },
  triggerButtonText: {
    fontSize: 12,
    color: '#fff',
    fontWeight: '600',
  },
  triggerButtonTextDisabled: {
    color: '#999',
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