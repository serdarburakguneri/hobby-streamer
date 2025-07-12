import React, { useState, useRef } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native';
import { VideoType } from '../../types/asset';

interface VideoUploadProps {
  onUpdate: (field: string, value: any) => Promise<void>;
}

export default function VideoUpload({ onUpdate }: VideoUploadProps) {
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [videoType, setVideoType] = useState<VideoType>(VideoType.MAIN);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedFile(file);
    }
  };

  const handleUploadVideo = async () => {
    if (!selectedFile) {
      alert('Please select a file to upload');
      return;
    }

    try {
      setUploading(true);
      setUploadProgress(0);

          const interval = setInterval(() => {
        setUploadProgress(prev => {
          if (prev >= 90) {
            clearInterval(interval);
            return 90;
          }
          return prev + 10;
        });
      }, 200);

      await new Promise(resolve => setTimeout(resolve, 2000));
      
      clearInterval(interval);
      setUploadProgress(100);

      setSelectedFile(null);
      setVideoType(VideoType.MAIN);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }

      alert('Video uploaded successfully!');
    } catch (error) {
      console.error('Upload error:', error);
      alert('Failed to upload video. Please try again.');
    } finally {
      setUploading(false);
      setUploadProgress(0);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.uploadTitle}>Upload New Video</Text>
    
      <View style={styles.uploadField}>
        <Text style={styles.uploadLabel}>Video Type:</Text>
        <View style={styles.videoTypePicker}>
          {Object.values(VideoType).map((type) => (
            <TouchableOpacity
              key={type}
              style={[
                styles.videoTypeOption,
                videoType === type && styles.videoTypeOptionSelected
              ]}
              onPress={() => setVideoType(type)}
            >
              <Text style={[
                styles.videoTypeOptionText,
                videoType === type && styles.videoTypeOptionTextSelected
              ]}>
                {type.replace('_', ' ')}
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      </View>

      <View style={styles.uploadField}>
        <Text style={styles.uploadLabel}>Select File:</Text>
        <input
          ref={fileInputRef}
          type="file"
          accept="video/*"
          onChange={handleFileSelect}
          style={styles.fileInput}
        />
        {selectedFile && (
          <Text style={styles.selectedFileName}>{selectedFile.name}</Text>
        )}
      </View>

      {uploading && (
        <View style={styles.uploadProgress}>
          <Text style={styles.uploadProgressText}>Uploading... {uploadProgress}%</Text>
          <View style={styles.progressBar}>
            <View style={[styles.progressFill, { width: `${uploadProgress}%` }]} />
          </View>
        </View>
      )}

      <TouchableOpacity
        style={[styles.uploadButton, (!selectedFile || uploading) && styles.uploadButtonDisabled]}
        onPress={handleUploadVideo}
        disabled={!selectedFile || uploading}
      >
        {uploading ? (
          <ActivityIndicator size="small" color="#fff" />
        ) : (
          <Text style={styles.uploadButtonText}>Upload Video</Text>
        )}
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginTop: 15,
    padding: 15,
    backgroundColor: '#f8f9fa',
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#e9ecef',
  },
  uploadTitle: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 15,
    color: '#333',
  },
  uploadField: {
    marginBottom: 12,
  },
  uploadLabel: {
    fontSize: 14,
    fontWeight: '600',
    marginBottom: 6,
    color: '#333',
  },
  fileInput: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 6,
    padding: 8,
    fontSize: 14,
    backgroundColor: '#fff',
    width: '100%',
  },
  selectedFileName: {
    fontSize: 12,
    color: '#666',
    marginTop: 4,
    fontStyle: 'italic',
  },
  videoTypePicker: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
    marginTop: 8,
  },
  videoTypeOption: {
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderRadius: 6,
    backgroundColor: '#f0f0f0',
    borderWidth: 1,
    borderColor: '#ddd',
  },
  videoTypeOptionSelected: {
    backgroundColor: '#007AFF',
    borderColor: '#007AFF',
  },
  videoTypeOptionText: {
    fontSize: 12,
    color: '#666',
    fontWeight: '500',
  },
  videoTypeOptionTextSelected: {
    color: '#fff',
    fontWeight: '600',
  },
  uploadProgress: {
    marginVertical: 10,
  },
  uploadProgressText: {
    fontSize: 12,
    color: '#666',
    marginBottom: 5,
  },
  progressBar: {
    height: 4,
    backgroundColor: '#e9ecef',
    borderRadius: 2,
    overflow: 'hidden',
  },
  progressFill: {
    height: '100%',
    backgroundColor: '#007AFF',
    borderRadius: 2,
  },
  uploadButton: {
    backgroundColor: '#007AFF',
    padding: 12,
    borderRadius: 6,
    alignItems: 'center',
    marginTop: 10,
  },
  uploadButtonDisabled: {
    backgroundColor: '#ccc',
  },
  uploadButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
}); 