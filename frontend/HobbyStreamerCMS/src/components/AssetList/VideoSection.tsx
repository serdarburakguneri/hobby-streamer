import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';

interface VideoSectionProps {
  videos: any[];
  onDeleteVideo: (videoType: string, videoName: string) => void;
  onUpdate: (field: string, value: any) => void;
}

export default function VideoSection({ videos, onDeleteVideo, onUpdate }: VideoSectionProps) {
  const [showUploadSection, setShowUploadSection] = useState(false);

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

  return (
    <View style={styles.container}>
      <View style={styles.sectionHeader}>
        <View style={styles.titleContainer}>
          <Ionicons name="videocam" size={20} color="#333" />
          <Text style={styles.sectionTitle}>Videos</Text>
        </View>
        <TouchableOpacity 
          style={styles.addButton}
          onPress={() => setShowUploadSection(!showUploadSection)}
        >
          <Ionicons 
            name={showUploadSection ? 'remove' : 'add'} 
            size={20} 
            color="#fff" 
          />
        </TouchableOpacity>
      </View>
      
      {videos && videos.length > 0 ? (
        <View style={styles.existingVideos}>
          {videos.map((video, index) => (
            <View key={`video-${video.type}-${index}`} style={styles.videoItem}>
              <View style={styles.videoHeader}>
                <View style={styles.videoInfo}>
                  <Ionicons 
                    name={getVideoIcon(video.type) as any} 
                    size={18} 
                    color="#007AFF" 
                  />
                  <Text style={styles.videoLabel}>{video.type}</Text>
                </View>
                <TouchableOpacity 
                  style={styles.deleteVideoButton}
                  onPress={() => onDeleteVideo(video.type, video.type)}
                >
                  <Ionicons name="trash" size={16} color="#ff3b30" />
                </TouchableOpacity>
              </View>
              <View style={styles.videoStatus}>
                <Ionicons 
                  name={getStatusIcon(video.status || 'unknown') as any} 
                  size={14} 
                  color={getStatusColor(video.status || 'unknown')} 
                />
                <Text style={styles.videoStatusText}>
                  Status: {video.status || 'Unknown'}
                </Text>
              </View>
              {video.raw && (
                <View style={styles.videoDetail}>
                  <Ionicons name="link" size={12} color="#666" />
                  <Text style={styles.videoDetailText}>
                    Raw: {video.raw.storageLocation.url}
                  </Text>
                </View>
              )}
              {video.hls && (
                <View style={styles.videoDetail}>
                  <Ionicons name="link" size={12} color="#666" />
                  <Text style={styles.videoDetailText}>
                    HLS: {video.hls.storageLocation.url}
                  </Text>
                </View>
              )}
              {video.dash && (
                <View style={styles.videoDetail}>
                  <Ionicons name="link" size={12} color="#666" />
                  <Text style={styles.videoDetailText}>
                    DASH: {video.dash.storageLocation.url}
                  </Text>
                </View>
              )}
              {video.thumbnail && (
                <View style={styles.videoDetail}>
                  <Ionicons name="image" size={12} color="#666" />
                  <Text style={styles.videoDetailText}>
                    Thumbnail: {video.thumbnail.url}
                  </Text>
                </View>
              )}
            </View>
          ))}
        </View>
      ) : (
        <View style={styles.emptyState}>
          <Ionicons name="videocam-outline" size={32} color="#ccc" />
          <Text style={styles.emptyText}>No videos uploaded yet</Text>
        </View>
      )}
      
      {showUploadSection && (
        <View style={styles.uploadSection}>
          <Text style={styles.uploadTitle}>Upload New Video</Text>
          <Text style={styles.uploadDescription}>
            Video upload functionality will be implemented here
          </Text>
        </View>
      )}
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
  addButton: {
    backgroundColor: '#007AFF',
    width: 32,
    height: 32,
    borderRadius: 16,
    justifyContent: 'center',
    alignItems: 'center',
  },
  existingVideos: {
    gap: 12,
  },
  videoItem: {
    backgroundColor: '#f8f9fa',
    borderRadius: 6,
    padding: 12,
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
  videoLabel: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
  },
  deleteVideoButton: {
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
  uploadSection: {
    marginTop: 16,
    padding: 16,
    backgroundColor: '#f8f9fa',
    borderRadius: 6,
  },
  uploadTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    marginBottom: 8,
  },
  uploadDescription: {
    fontSize: 14,
    color: '#666',
  },
}); 