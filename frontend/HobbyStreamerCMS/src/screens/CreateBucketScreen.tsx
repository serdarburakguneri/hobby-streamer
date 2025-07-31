import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TextInput,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  Alert,
} from 'react-native';
import Layout from '../components/Layout';
import { BucketType, BucketStatus } from '../types/asset';
import { useAssetService, getAuthToken, validateTokenLocally } from '../services/api';

interface CreateBucketScreenProps {
  onBack: () => void;
  onBucketCreated: (bucket: any) => void;
}

export default function CreateBucketScreen({ onBack, onBucketCreated }: CreateBucketScreenProps) {
  const assetService = useAssetService();
  const [key, setKey] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [type, setType] = useState<BucketType>(BucketType.COLLECTION);
  const [creating, setCreating] = useState(false);

  const bucketTypes = [
    { value: BucketType.COLLECTION, label: 'Collection' },
    { value: BucketType.PLAYLIST, label: 'Playlist' },
    { value: BucketType.CATEGORY, label: 'Category' },
  ];

  const handleCreate = async () => {
    console.log('Create button clicked');
    if (!name.trim()) {
      Alert.alert('Error', 'Bucket name is required');
      return;
    }

    if (!key.trim()) {
      Alert.alert('Error', 'Bucket key is required');
      return;
    }

    try {
      setCreating(true);
      let ownerId = '';
      const token = await getAuthToken();
      if (token) {
        const { valid, user } = validateTokenLocally(token);
        if (valid && user && user.id) {
          ownerId = user.id;
        }
      }
      if (!ownerId) {
        Alert.alert('Error', 'Could not determine logged-in user.');
        setCreating(false);
        return;
      }
      const newBucket = await assetService.createBucket({
        key: key.trim(),
        name: name.trim(),
        description: description.trim(),
        type,
        status: BucketStatus.DRAFT,
        ownerId,
      });
      onBucketCreated(newBucket);
      onBack();
    } catch (error: any) {
      console.error('Error creating bucket:', error);
      Alert.alert('Error', `Failed to create bucket: ${error.message || 'Unknown error'}`);
    } finally {
      setCreating(false);
    }
  };

  const handleCancel = () => {
    console.log('Cancel button clicked');
    if (name.trim() || description.trim() || key.trim()) {
      Alert.alert(
        'Discard Changes',
        'Are you sure you want to discard your changes?',
        [
          { text: 'Cancel', style: 'cancel' },
          { text: 'Discard', style: 'destructive', onPress: onBack },
        ]
      );
    } else {
      onBack();
    }
  };

  return (
    <Layout
      headerTitle="Create New Bucket"
      headerLeft={
        <TouchableOpacity style={styles.backButton} onPress={onBack}>
          <Text style={styles.backButtonText}>← Back</Text>
        </TouchableOpacity>
      }
      headerRight={
        <TouchableOpacity
          style={[styles.createButton, (!name.trim() || !key.trim() || creating) && styles.createButtonDisabled]}
          onPress={handleCreate}
          disabled={!name.trim() || !key.trim() || creating}
        >
          {creating ? (
            <ActivityIndicator size="small" color="#fff" />
          ) : (
            <Text style={styles.createButtonText}>Create</Text>
          )}
        </TouchableOpacity>
      }
    >
      <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
        <View style={styles.form}>
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Basic Information</Text>
            
            <View style={styles.field}>
              <Text style={styles.label}>Key *</Text>
              <TextInput
                style={styles.input}
                value={key}
                onChangeText={setKey}
                placeholder="Enter bucket key"
                maxLength={100}
              />
            </View>

            <View style={styles.field}>
              <Text style={styles.label}>Name *</Text>
              <TextInput
                style={styles.input}
                value={name}
                onChangeText={setName}
                placeholder="Enter bucket name"
                maxLength={100}
              />
            </View>

            <View style={styles.field}>
              <Text style={styles.label}>Description</Text>
              <TextInput
                style={[styles.input, styles.textArea]}
                value={description}
                onChangeText={setDescription}
                placeholder="Enter bucket description"
                multiline
                numberOfLines={4}
                maxLength={500}
              />
            </View>
          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Configuration</Text>
            
            <View style={styles.field}>
              <Text style={styles.label}>Type</Text>
              <View style={styles.optionsContainer}>
                {bucketTypes.map((option) => (
                  <TouchableOpacity
                    key={option.value}
                    style={[
                      styles.option,
                      type === option.value && styles.selectedOption,
                    ]}
                    onPress={() => setType(option.value)}
                  >
                    <Text
                      style={[
                        styles.optionText,
                        type === option.value && styles.selectedOptionText,
                      ]}
                    >
                      {option.label}
                    </Text>
                  </TouchableOpacity>
                ))}
              </View>
            </View>


          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Preview</Text>
            <View style={styles.preview}>
              <Text style={styles.previewTitle}>{name || 'Bucket Name'}</Text>
              {description && (
                <Text style={styles.previewDescription}>{description}</Text>
              )}
              <View style={styles.previewMeta}>
                <Text style={styles.previewMetaText}>Type: {bucketTypes.find(t => t.value === type)?.label}</Text>
                <Text style={styles.previewMetaText}>Status: Draft</Text>
              </View>
            </View>
          </View>
        </View>
      </ScrollView>
    </Layout>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  form: {
    padding: 20,
  },
  section: {
    backgroundColor: '#fff',
    borderRadius: 8,
    padding: 16,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 16,
  },
  field: {
    marginBottom: 20,
  },
  label: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    marginBottom: 8,
  },
  input: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 6,
    padding: 12,
    fontSize: 16,
    backgroundColor: '#fff',
  },
  textArea: {
    height: 100,
    textAlignVertical: 'top',
  },
  optionsContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  option: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 6,
    borderWidth: 1,
    borderColor: '#ddd',
    backgroundColor: '#fff',
  },
  selectedOption: {
    backgroundColor: '#007AFF',
    borderColor: '#007AFF',
  },
  optionText: {
    fontSize: 14,
    color: '#333',
    fontWeight: '500',
  },
  selectedOptionText: {
    color: '#fff',
  },
  preview: {
    backgroundColor: '#f8f9fa',
    borderRadius: 6,
    padding: 16,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  previewTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 8,
  },
  previewDescription: {
    fontSize: 14,
    color: '#666',
    marginBottom: 12,
    lineHeight: 20,
  },
  previewMeta: {
    flexDirection: 'row',
    gap: 16,
  },
  previewMetaText: {
    fontSize: 12,
    color: '#999',
  },
  backButton: {
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
    borderRadius: 6,
  },
  backButtonText: {
    color: '#007AFF',
    fontSize: 16,
    fontWeight: '600',
  },
  createButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 20,
    paddingVertical: 12,
    borderRadius: 6,
  },
  createButtonDisabled: {
    backgroundColor: '#ccc',
  },
  createButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
}); 