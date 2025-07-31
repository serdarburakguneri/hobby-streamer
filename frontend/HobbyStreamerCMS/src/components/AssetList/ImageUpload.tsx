import React, { useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Alert, ScrollView } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import * as ImagePicker from 'expo-image-picker';
import { Asset, ImageType } from '../../types/asset';
import { useAssetService } from '../../services/api';
import ImagePreview from './ImagePreview';

interface ImageUploadProps {
  asset: Asset;
  onImageAdded: () => void;
}

const IMAGE_TYPES: { key: ImageType; label: string; icon: string }[] = [
  { key: ImageType.POSTER, label: 'Poster', icon: 'image' },
  { key: ImageType.BANNER, label: 'Banner', icon: 'image' },
  { key: ImageType.HERO, label: 'Hero', icon: 'image' },
  { key: ImageType.LOGO, label: 'Logo', icon: 'image' },
  { key: ImageType.SCREENSHOT, label: 'Screenshot', icon: 'image' },
];

export default function ImageUpload({ asset, onImageAdded }: ImageUploadProps) {
  const [selectedTab, setSelectedTab] = useState<ImageType>(ImageType.POSTER);
  const [uploading, setUploading] = useState(false);
  const assetService = useAssetService();



  const pickImage = async (imageType: ImageType) => {
    try {
      const result = await ImagePicker.launchImageLibraryAsync({
        mediaTypes: ImagePicker.MediaTypeOptions.Images,
        allowsEditing: true,
        aspect: imageType === ImageType.BANNER ? [16, 9] : imageType === ImageType.POSTER ? [2, 3] : imageType === ImageType.HERO ? [16, 9] : [1, 1],
        quality: 0.8,
      });

      if (!result.canceled && result.assets[0]) {
        await uploadImage(result.assets[0], imageType);
      }
    } catch (error) {
      console.error('Error picking image:', error);
      Alert.alert('Error', 'Failed to pick image');
    }
  };

  const uploadImage = async (imageAsset: ImagePicker.ImagePickerAsset, imageType: ImageType) => {
    if (!imageAsset.uri) return;

    setUploading(true);
    try {
      const fileName = `${imageType.toLowerCase()}_${Date.now()}.jpg`;
      
      const response = await assetService.getImageUploadUrl(fileName, asset.id, imageType);
      
      const imageResponse = await fetch(imageAsset.uri);
      const imageBlob = await imageResponse.blob();

      await assetService.uploadFile(response.url, imageBlob);

      const imageData = {
        url: `http://localhost:8083/cdn/${asset.id}/images/${imageType.toLowerCase()}/${fileName}`,
        type: imageType,
        fileName,
        size: imageBlob.size,
      };

      await assetService.addImageToAsset(asset.id, imageData);
      Alert.alert('Success', `${imageType} image uploaded successfully`);
      onImageAdded();
    } catch (error) {
      console.error('Error uploading image:', error);
      Alert.alert('Error', 'Failed to upload image');
    } finally {
      setUploading(false);
    }
  };

  const deleteImage = async (imageId: string) => {
    try {
      await assetService.deleteImageFromAsset(asset.id, imageId);
      Alert.alert('Success', 'Image deleted successfully');
      onImageAdded();
    } catch (error) {
      console.error('Error deleting image:', error);
      Alert.alert('Error', 'Failed to delete image');
    }
  };

  const getImagesByType = (type: ImageType) => {
    return asset.images?.filter(img => img.type === type) || [];
  };

  return (
    <View style={styles.container}>
      <Text style={styles.sectionTitle}>Images</Text>
      
      <View style={styles.tabContainer}>
        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
          {IMAGE_TYPES.map((imageType) => (
            <TouchableOpacity
              key={imageType.key}
              style={[
                styles.tab,
                selectedTab === imageType.key && styles.activeTab
              ]}
              onPress={() => setSelectedTab(imageType.key)}
            >
              <Ionicons 
                name={imageType.icon as any} 
                size={16} 
                color={selectedTab === imageType.key ? '#fff' : '#007AFF'} 
              />
              <Text style={[
                styles.tabText,
                selectedTab === imageType.key && styles.activeTabText
              ]}>
                {imageType.label}
              </Text>
            </TouchableOpacity>
          ))}
        </ScrollView>
      </View>

      <View style={styles.content}>
        <View style={styles.uploadSection}>
          <TouchableOpacity
            style={styles.uploadButton}
            onPress={() => pickImage(selectedTab)}
            disabled={uploading}
          >
            <Ionicons name="cloud-upload" size={24} color="#007AFF" />
            <Text style={styles.uploadButtonText}>
              {uploading ? 'Uploading...' : `Upload ${selectedTab.toLowerCase()}`}
            </Text>
          </TouchableOpacity>
        </View>

        <View style={styles.imagesList}>
          {getImagesByType(selectedTab).map((image) => (
            <ImagePreview
              key={image.id}
              image={image}
              onDelete={deleteImage}
            />
          ))}
          
          {getImagesByType(selectedTab).length === 0 && (
            <Text style={styles.noImages}>No {selectedTab.toLowerCase()} images</Text>
          )}
        </View>
      </View>
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
  imagesList: {
    gap: 10,
  },
  noImages: {
    textAlign: 'center',
    color: '#999',
    fontStyle: 'italic',
    padding: 20,
  },
}); 