import React, { useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Alert, ScrollView } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Asset, VideoType, VideoFormat } from '../../types/asset';
import { useAssetService } from '../../services/api';
import VideoUpload from './VideoUpload';
import VideoPreview from './VideoPreview';

interface VideoSectionProps {
  asset: Asset;
  onVideoAdded: () => void;
}

const VIDEO_TYPES: { key: VideoType; label: string; icon: string }[] = [
  { key: VideoType.MAIN, label: 'Main', icon: 'videocam' },
  { key: VideoType.TRAILER, label: 'Trailer', icon: 'play-circle' },
];

const VIDEO_FORMATS: { key: VideoFormat; label: string; icon: string }[] = [
  { key: VideoFormat.RAW, label: 'Raw', icon: 'film' },
  { key: VideoFormat.HLS, label: 'HLS', icon: 'play' },
  { key: VideoFormat.DASH, label: 'DASH', icon: 'play' },
];

export default function VideoSection({ asset, onVideoAdded }: VideoSectionProps) {
  const [selectedVideoType, setSelectedVideoType] = useState<VideoType>(VideoType.MAIN);
  
  const getDefaultFormat = () => {
    const mainVideos = asset.videos?.filter(video => video.type === VideoType.MAIN) || [];
    const dashVideo = mainVideos.find(video => video.format === VideoFormat.DASH && video.status === 'ready');
    const hlsVideo = mainVideos.find(video => video.format === VideoFormat.HLS && video.status === 'ready');
    
    if (dashVideo) return VideoFormat.DASH;
    if (hlsVideo) return VideoFormat.HLS;
    return VideoFormat.RAW;
  };
  
  const [selectedFormat, setSelectedFormat] = useState<VideoFormat>(getDefaultFormat());
  const [showUpload, setShowUpload] = useState(false);
  const [previewVideo, setPreviewVideo] = useState<any>(null);
  const [showPreview, setShowPreview] = useState(false);
  const assetService = useAssetService();

  const getVideosByTypeAndFormat = (videoType: VideoType, format: VideoFormat) => {
    return asset.videos?.filter(video => 
      video.type === videoType && video.format === format
    ) || [];
  };

  const deleteVideo = async (videoId: string) => {
    try {
      await assetService.deleteVideo(asset.id, videoId);
      Alert.alert('Success', 'Video deleted successfully');
      onVideoAdded();
    } catch (error) {
      console.error('Error deleting video:', error);
      Alert.alert('Error', 'Failed to delete video');
    }
  };

  const triggerTranscoding = async (videoType: VideoType, format: VideoFormat) => {
    try {
      const rawVideo = asset.videos?.find(video => 
        video.type === videoType && video.format === VideoFormat.RAW
      );

      if (!rawVideo) {
        Alert.alert('Error', `No raw ${videoType.toLowerCase()} video found. Please upload a raw video first.`);
        return;
      }

      if (format === VideoFormat.HLS) {        
        const input = `s3://${rawVideo.storageLocation.bucket}/${rawVideo.storageLocation.key}`;
        await assetService.triggerHLSTranscode(asset.id, rawVideo.id, input);
        Alert.alert('Success', 'HLS transcoding triggered successfully');
      } else if (format === VideoFormat.DASH) {
        const input = `s3://${rawVideo.storageLocation.bucket}/${rawVideo.storageLocation.key}`;
        await assetService.triggerDASHTranscode(asset.id, rawVideo.id, input);
        Alert.alert('Success', 'DASH transcoding triggered successfully');
      }
            
      onVideoAdded();
    } catch (error) {
      console.error('Error triggering transcoding:', error);
      Alert.alert('Error', 'Failed to trigger transcoding');
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'ready':
        return 'checkmark-circle';
      case 'transcoding':
        return 'sync';
      case 'analyzing':
        return 'search';
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
      case 'transcoding':
        return '#2196F3';
      case 'analyzing':
        return '#FF9800';
      case 'failed':
        return '#F44336';
      default:
        return '#9E9E9E';
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDuration = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60);
    
    if (hours > 0) {
      return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
    }
    return `${minutes}:${secs.toString().padStart(2, '0')}`;
  };

  const handlePreviewVideo = (video: any) => {
    setPreviewVideo(video);
    setShowPreview(true);
  };

  const handleClosePreview = () => {
    setShowPreview(false);
    setPreviewVideo(null);
  };

  return (
    <View style={styles.container}>
      <Text style={styles.sectionTitle}>Videos</Text>
      
      <View style={styles.tabContainer}>
        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
          {VIDEO_TYPES.map((videoType) => (
            <TouchableOpacity
              key={videoType.key}
              style={[
                styles.tab,
                selectedVideoType === videoType.key && styles.activeTab
              ]}
              onPress={() => setSelectedVideoType(videoType.key)}
            >
              <Ionicons 
                name={videoType.icon as any} 
                size={16} 
                color={selectedVideoType === videoType.key ? '#fff' : '#007AFF'} 
              />
              <Text style={[
                styles.tabText,
                selectedVideoType === videoType.key && styles.activeTabText
              ]}>
                {videoType.label}
              </Text>
            </TouchableOpacity>
          ))}
        </ScrollView>
      </View>

      <View style={styles.content}>
        <View style={styles.formatTabContainer}>
          <ScrollView horizontal showsHorizontalScrollIndicator={false}>
            {VIDEO_FORMATS.map((format) => (
              <TouchableOpacity
                key={format.key}
                style={[
                  styles.formatTab,
                  selectedFormat === format.key && styles.activeFormatTab
                ]}
                onPress={() => setSelectedFormat(format.key)}
              >
                <Ionicons 
                  name={format.icon as any} 
                  size={14} 
                  color={selectedFormat === format.key ? '#fff' : '#666'} 
                />
                <Text style={[
                  styles.formatTabText,
                  selectedFormat === format.key && styles.activeFormatTabText
                ]}>
                  {format.label}
                </Text>
              </TouchableOpacity>
            ))}
          </ScrollView>
        </View>

        <View style={styles.uploadSection}>
          {selectedFormat === VideoFormat.RAW ? (
            <>
              <TouchableOpacity
                style={styles.uploadButton}
                onPress={() => setShowUpload(true)}
              >
                <Ionicons name="cloud-upload" size={24} color="#007AFF" />
                <Text style={styles.uploadButtonText}>
                  Upload {selectedVideoType.toLowerCase()} video
                </Text>
              </TouchableOpacity>
              
              {showUpload && (
                <VideoUpload
                  assetId={asset.id}
                  videoType={selectedVideoType}
                  onUploadComplete={() => {
                    setShowUpload(false);
                    onVideoAdded();
                  }}
                  onCancel={() => setShowUpload(false)}
                  onRefreshAsset={onVideoAdded}
                />
              )}
            </>
          ) : (
            <TouchableOpacity
              style={styles.transcodeButton}
              onPress={() => triggerTranscoding(selectedVideoType, selectedFormat)}
            >
              <Ionicons name="play-circle" size={24} color="#007AFF" />
              <Text style={styles.transcodeButtonText}>
                Trigger {selectedFormat.toUpperCase()} Transcoding
              </Text>
            </TouchableOpacity>
          )}
        </View>

        <View style={styles.videosList}>
          {getVideosByTypeAndFormat(selectedVideoType, selectedFormat).map((video) => (
            <View key={video.id} style={styles.videoItem}>
                             <View style={styles.videoInfo}>
                 <View style={styles.videoHeader}>
                   <Text style={styles.videoName}>
                     Filename: {video.storageLocation?.key?.split('/').pop() || `Video ${video.id}`}
                   </Text>
                 </View>
                <View style={styles.videoDetails}>
                  <Text style={styles.videoDetail}>
                    Status: {video.status || 'Unknown'}
                  </Text>
                  <Text style={styles.videoDetail}>
                    Duration: {video.duration ? formatDuration(video.duration) : '0'}
                  </Text>
                  <Text style={styles.videoDetail}>
                    Size: {video.size ? formatFileSize(video.size) : '0'}
                  </Text>
                  <Text style={styles.videoDetail}>
                    Resolution: {video.width && video.height ? `${video.width}x${video.height}` : '0'}
                  </Text>
                  <Text style={styles.videoDetail}>
                    Bitrate: {video.bitrate ? `${Math.round(video.bitrate / 1000)} kbps` : '0'}
                  </Text>
                  <Text style={styles.videoDetail}>
                    Codec: {video.codec && video.codec.trim() !== '' ? video.codec : '0'}
                  </Text>
                  <Text style={styles.videoDate}>
                    Created: {new Date(video.createdAt).toLocaleDateString()}
                  </Text>
                </View>
              </View>
                             <View style={styles.videoActions}>
                 {(video.format === 'hls' || video.format === 'dash') && video.status === 'ready' && (
                   <TouchableOpacity
                     style={styles.actionButton}
                     onPress={() => handlePreviewVideo(video)}
                   >
                     <Ionicons name="eye" size={14} color="#007AFF" />
                     <Text style={styles.actionButtonText}>View</Text>
                   </TouchableOpacity>
                 )}
                 {video.format === 'raw' && (
                   <View style={styles.actionButton}>
                     <Ionicons name="information-circle" size={14} color="#666" />
                     <Text style={[styles.actionButtonText, { color: '#666' }]}>Raw</Text>
                   </View>
                 )}
                 <TouchableOpacity
                   style={[styles.actionButton, styles.deleteActionButton]}
                   onPress={() => deleteVideo(video.id)}
                 >
                   <Ionicons name="trash" size={14} color="#FF3B30" />
                   <Text style={[styles.actionButtonText, styles.deleteActionText]}>Delete</Text>
                 </TouchableOpacity>
               </View>
            </View>
          ))}
          
          {getVideosByTypeAndFormat(selectedVideoType, selectedFormat).length === 0 && (
            <Text style={styles.noVideos}>
              No {selectedVideoType.toLowerCase()} {selectedFormat.toLowerCase()} videos
            </Text>
          )}
        </View>
      </View>

      {previewVideo && (
        <VideoPreview
          video={previewVideo}
          visible={showPreview}
          onClose={handleClosePreview}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginTop: 20,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 15,
    color: '#333',
  },
  tabContainer: {
    marginBottom: 15,
  },
  tab: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 12,
    paddingVertical: 8,
    marginRight: 8,
    borderRadius: 20,
    backgroundColor: '#f0f0f0',
    gap: 6,
  },
  activeTab: {
    backgroundColor: '#007AFF',
  },
  tabText: {
    fontSize: 14,
    fontWeight: '500',
    color: '#007AFF',
  },
  activeTabText: {
    color: '#fff',
  },
  content: {
    backgroundColor: '#f8f9fa',
    borderRadius: 8,
    padding: 15,
  },
  formatTabContainer: {
    marginBottom: 15,
  },
  formatTab: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 10,
    paddingVertical: 6,
    marginRight: 6,
    borderRadius: 16,
    backgroundColor: '#e9ecef',
    gap: 4,
  },
  activeFormatTab: {
    backgroundColor: '#0056b3',
  },
  formatTabText: {
    fontSize: 12,
    fontWeight: '500',
    color: '#666',
  },
  activeFormatTabText: {
    color: '#fff',
  },
  uploadSection: {
    marginBottom: 15,
  },
  uploadButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 12,
    paddingHorizontal: 20,
    backgroundColor: '#fff',
    borderRadius: 8,
    borderWidth: 2,
    borderColor: '#007AFF',
    borderStyle: 'dashed',
    gap: 8,
  },
  uploadButtonText: {
    fontSize: 16,
    fontWeight: '500',
    color: '#007AFF',
  },
  transcodeButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 12,
    paddingHorizontal: 20,
    backgroundColor: '#fff',
    borderRadius: 8,
    borderWidth: 2,
    borderColor: '#007AFF',
    borderStyle: 'dashed',
    gap: 8,
  },
  transcodeButtonText: {
    fontSize: 16,
    fontWeight: '500',
    color: '#007AFF',
  },
  videosList: {
    gap: 10,
  },
  videoItem: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    justifyContent: 'space-between',
    padding: 12,
    backgroundColor: '#fff',
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  videoInfo: {
    flex: 1,
  },
  videoHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  videoName: {
    fontSize: 14,
    fontWeight: '500',
    color: '#333',
    flex: 1,
  },
  statusContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  videoStatus: {
    fontSize: 12,
    color: '#666',
    marginLeft: 4,
  },
  videoDetails: {
    gap: 2,
  },
  videoDetail: {
    fontSize: 12,
    color: '#666',
  },
  videoDate: {
    fontSize: 12,
    color: '#999',
    marginTop: 4,
  },
  videoActions: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  actionButton: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 10,
    paddingVertical: 6,
    borderRadius: 6,
    backgroundColor: '#f8f9fa',
    borderWidth: 1,
    borderColor: '#e9ecef',
    gap: 4,
  },
  actionButtonText: {
    fontSize: 12,
    fontWeight: '500',
    color: '#007AFF',
  },
  deleteActionButton: {
    backgroundColor: '#fff5f5',
    borderColor: '#fed7d7',
  },
  deleteActionText: {
    color: '#FF3B30',
  },
  noVideos: {
    textAlign: 'center',
    color: '#999',
    fontStyle: 'italic',
    padding: 20,
  },
}); 