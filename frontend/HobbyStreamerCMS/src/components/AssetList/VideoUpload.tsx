import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Alert, ActivityIndicator } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import * as DocumentPicker from 'expo-document-picker';
import axios from 'axios';
import { useAssetService } from '../../services/api';
import { VideoType } from '../../types/asset';
import { API_CONFIG } from '../../config/api';

interface VideoUploadProps {
  assetId: string;
  onUploadComplete: () => void;
  onCancel: () => void;
  onRefreshAsset: () => void;
}

export default function VideoUpload({ assetId, onUploadComplete, onCancel, onRefreshAsset }: VideoUploadProps) {
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [selectedVideoType, setSelectedVideoType] = useState<VideoType | null>(null);
  const { getUploadUrl, uploadFile, addVideo } = useAssetService();

  const handleFilePick = async () => {
    try {
      const result = await DocumentPicker.getDocumentAsync({
        type: ['video/*'],
        copyToCacheDirectory: true,
      });

      if (result.canceled) {
        return;
      }

      const file = result.assets[0];
      if (!file) {
        Alert.alert('Error', 'No file selected');
        return;
      }

      await uploadVideo(file);
    } catch (error) {
      console.error('Error picking file:', error);
      Alert.alert('Error', 'Failed to pick file');
    }
  };

  const uploadVideo = async (file: DocumentPicker.DocumentPickerAsset) => {
    if (!file.uri || !file.name) {
      Alert.alert('Error', 'Invalid file');
      return;
    }

    if (!selectedVideoType) {
      Alert.alert('Error', 'Please select a video type');
      return;
    }

    setUploading(true);
    setUploadProgress(0);

    try {
      const fileName = `${assetId}/${Date.now()}_${file.name}`;
      
      const { url: uploadUrl } = await getUploadUrl(fileName);
      
      const response = await fetch(file.uri);
      const blob = await response.blob();
      
      await uploadFile(uploadUrl, blob);
      
      const bucket = 'raw-storage';
      const key = fileName;
      const url = `${API_CONFIG.LOCALSTACK_BASE_URL}/${bucket}/${key}`;
      
      await addVideo(assetId, selectedVideoType, bucket, key, url, file.mimeType || 'video/mp4', file.size || 0);
      
      setUploadProgress(100);
      Alert.alert('Success', 'Video uploaded successfully');
      onRefreshAsset();
      onUploadComplete();
    } catch (error) {
      console.error('Upload error:', error);
      Alert.alert('Error', 'Failed to upload video. Please try again.');
    } finally {
      setUploading(false);
      setUploadProgress(0);
    }
  };

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Upload Video</Text>
        <TouchableOpacity onPress={onCancel} style={styles.cancelButton}>
          <Ionicons name="close" size={20} color="#666" />
        </TouchableOpacity>
      </View>

      {uploading ? (
        <View style={styles.uploadingContainer}>
          <ActivityIndicator size="large" color="#007AFF" />
          <Text style={styles.uploadingText}>Uploading video...</Text>
          <View style={styles.progressBar}>
            <View style={[styles.progressFill, { width: `${uploadProgress}%` }]} />
          </View>
          <Text style={styles.progressText}>{uploadProgress}%</Text>
        </View>
      ) : (
        <View style={styles.uploadContainer}>
          <Text style={styles.sectionTitle}>Select Video Type</Text>
          <View style={styles.videoTypeContainer}>
            <TouchableOpacity
              style={[
                styles.videoTypeButton,
                selectedVideoType === VideoType.MAIN && styles.videoTypeButtonSelected
              ]}
              onPress={() => setSelectedVideoType(VideoType.MAIN)}
            >
              <Ionicons 
                name="videocam" 
                size={24} 
                color={selectedVideoType === VideoType.MAIN ? '#fff' : '#007AFF'} 
              />
              <Text style={[
                styles.videoTypeText,
                selectedVideoType === VideoType.MAIN && styles.videoTypeTextSelected
              ]}>
                Main Video
              </Text>
            </TouchableOpacity>
            
            <TouchableOpacity
              style={[
                styles.videoTypeButton,
                selectedVideoType === VideoType.TRAILER && styles.videoTypeButtonSelected
              ]}
              onPress={() => setSelectedVideoType(VideoType.TRAILER)}
            >
              <Ionicons 
                name="play-circle" 
                size={24} 
                color={selectedVideoType === VideoType.TRAILER ? '#fff' : '#007AFF'} 
              />
              <Text style={[
                styles.videoTypeText,
                selectedVideoType === VideoType.TRAILER && styles.videoTypeTextSelected
              ]}>
                Trailer
              </Text>
            </TouchableOpacity>
          </View>

          {selectedVideoType && (
            <TouchableOpacity style={styles.uploadButton} onPress={handleFilePick}>
              <Ionicons name="cloud-upload" size={32} color="#007AFF" />
              <Text style={styles.uploadButtonText}>Select Video File</Text>
              <Text style={styles.uploadDescription}>
                Choose a video file to upload (MP4, MOV, AVI, etc.)
              </Text>
            </TouchableOpacity>
          )}
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
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
  },
  cancelButton: {
    padding: 4,
  },
  uploadContainer: {
    padding: 20,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    marginBottom: 16,
    textAlign: 'center',
  },
  videoTypeContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: 24,
  },
  videoTypeButton: {
    alignItems: 'center',
    padding: 16,
    borderWidth: 2,
    borderColor: '#007AFF',
    borderRadius: 8,
    backgroundColor: '#f8f9fa',
    minWidth: 120,
  },
  videoTypeButtonSelected: {
    backgroundColor: '#007AFF',
  },
  videoTypeText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#007AFF',
    marginTop: 8,
  },
  videoTypeTextSelected: {
    color: '#fff',
  },
  uploadButton: {
    alignItems: 'center',
    padding: 20,
    borderWidth: 2,
    borderColor: '#007AFF',
    borderStyle: 'dashed',
    borderRadius: 8,
    backgroundColor: '#f8f9fa',
  },
  uploadButtonText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#007AFF',
    marginTop: 8,
    marginBottom: 4,
  },
  uploadDescription: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
  },
  uploadingContainer: {
    alignItems: 'center',
    padding: 20,
  },
  uploadingText: {
    fontSize: 16,
    color: '#333',
    marginTop: 12,
    marginBottom: 16,
  },
  progressBar: {
    width: '100%',
    height: 8,
    backgroundColor: '#e0e0e0',
    borderRadius: 4,
    overflow: 'hidden',
    marginBottom: 8,
  },
  progressFill: {
    height: '100%',
    backgroundColor: '#007AFF',
  },
  progressText: {
    fontSize: 14,
    color: '#666',
  },
}); 