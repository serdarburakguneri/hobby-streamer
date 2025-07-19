import React, { useState } from 'react';
import { View, Text, Image, TouchableOpacity, StyleSheet, Modal, Dimensions, ActivityIndicator } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Image as ImageType } from '../../types/asset';

interface ImagePreviewProps {
  image: ImageType;
  onDelete: (imageId: string) => void;
}

const { width: screenWidth, height: screenHeight } = Dimensions.get('window');

export default function ImagePreview({ image, onDelete }: ImagePreviewProps) {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [imageLoading, setImageLoading] = useState(false);
  const [imageError, setImageError] = useState(false);

  const openModal = () => {
    setImageLoading(true);
    setImageError(false);
    setIsModalVisible(true);
  };
  const closeModal = () => setIsModalVisible(false);

  return (
    <View style={styles.container}>
      <TouchableOpacity style={styles.imageContainer} onPress={openModal}>
        <Image source={{ uri: image.url }} style={styles.image} resizeMode="cover" />
        <View style={styles.overlay}>
          <Ionicons name="eye" size={20} color="#fff" />
        </View>
      </TouchableOpacity>
      
      <View style={styles.imageInfo}>
        <Text style={styles.imageName} numberOfLines={1}>{image.fileName}</Text>
        <Text style={styles.imageDate}>
          {new Date(image.createdAt).toLocaleDateString()}
        </Text>
      </View>
      
      <TouchableOpacity
        style={styles.deleteButton}
        onPress={() => onDelete(image.id)}
      >
        <Ionicons name="trash" size={16} color="#FF3B30" />
      </TouchableOpacity>

      <Modal
        visible={isModalVisible}
        transparent={true}
        animationType="fade"
        onRequestClose={closeModal}
      >
        <View style={styles.modalOverlay}>
          <TouchableOpacity style={styles.modalBackground} onPress={closeModal} />
          <View style={styles.modalContent}>
            <TouchableOpacity style={styles.closeButton} onPress={closeModal}>
              <Ionicons name="close" size={24} color="#fff" />
            </TouchableOpacity>
            <View style={styles.modalImageContainer}>
              {imageLoading && (
                <View style={styles.loadingContainer}>
                  <ActivityIndicator size="large" color="#007AFF" />
                  <Text style={styles.loadingText}>Loading image...</Text>
                </View>
              )}
              {imageError && (
                <View style={styles.errorContainer}>
                  <Ionicons name="image-outline" size={48} color="#999" />
                  <Text style={styles.errorText}>Failed to load image</Text>
                  <Text style={styles.errorUrl}>{image.url}</Text>
                </View>
              )}
              <Image 
                source={{ uri: image.url }} 
                style={styles.modalImage} 
                resizeMode="contain"
                onLoadStart={() => setImageLoading(true)}
                onLoad={() => setImageLoading(false)}
                onError={() => {
                  setImageLoading(false);
                  setImageError(true);
                }}
              />
            </View>
            <View style={styles.modalInfo}>
              <Text style={styles.modalFileName}>{image.fileName}</Text>
              <Text style={styles.modalType}>{image.type}</Text>
              {image.width && image.height && (
                <Text style={styles.modalDimensions}>
                  {image.width} × {image.height}
                </Text>
              )}
              {image.size && (
                <Text style={styles.modalSize}>
                  {(image.size / 1024 / 1024).toFixed(2)} MB
                </Text>
              )}
            </View>
          </View>
        </View>
      </Modal>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#f8f9fa',
    borderRadius: 8,
    padding: 12,
    marginBottom: 8,
  },
  imageContainer: {
    position: 'relative',
    marginRight: 12,
  },
  image: {
    width: 60,
    height: 60,
    borderRadius: 6,
  },
  overlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.3)',
    borderRadius: 6,
    justifyContent: 'center',
    alignItems: 'center',
  },
  imageInfo: {
    flex: 1,
  },
  imageName: {
    fontSize: 14,
    fontWeight: '500',
    color: '#333',
    marginBottom: 2,
  },
  imageDate: {
    fontSize: 12,
    color: '#666',
  },
  deleteButton: {
    padding: 8,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.9)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  modalBackground: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
  },
  modalContent: {
    width: screenWidth * 0.9,
    height: screenHeight * 0.8,
    backgroundColor: '#fff',
    borderRadius: 12,
    overflow: 'hidden',
  },
  closeButton: {
    position: 'absolute',
    top: 16,
    right: 16,
    zIndex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    borderRadius: 20,
    padding: 8,
  },
  modalImageContainer: {
    width: '100%',
    height: '80%',
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#f0f0f0',
  },
  modalImage: {
    width: '100%',
    height: '100%',
  },
  loadingContainer: {
    position: 'absolute',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 1,
  },
  loadingText: {
    marginTop: 8,
    fontSize: 14,
    color: '#666',
  },
  errorContainer: {
    position: 'absolute',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 1,
  },
  errorText: {
    marginTop: 8,
    fontSize: 14,
    color: '#999',
    marginBottom: 4,
  },
  errorUrl: {
    fontSize: 12,
    color: '#999',
    textAlign: 'center',
    paddingHorizontal: 16,
  },
  modalInfo: {
    padding: 16,
    backgroundColor: '#f8f9fa',
  },
  modalFileName: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    marginBottom: 4,
  },
  modalType: {
    fontSize: 14,
    color: '#007AFF',
    marginBottom: 2,
  },
  modalDimensions: {
    fontSize: 12,
    color: '#666',
    marginBottom: 2,
  },
  modalSize: {
    fontSize: 12,
    color: '#666',
  },
}); 