import React, { useState, useEffect } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, ActivityIndicator, TextInput, Switch } from 'react-native';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { Asset } from '../../types/asset';

interface PublishSettingsProps {
  asset: Asset;
  onPublish: (publishAt: Date | null, unpublishAt: Date | null, ageRating: string) => Promise<void>;
  publishing: boolean;
}

export default function PublishSettings({ asset, onPublish, publishing }: PublishSettingsProps) {
  const [publishAt, setPublishAt] = useState<Date | null>(null);
  const [unpublishAt, setUnpublishAt] = useState<Date | null>(null);
  const [ageRating, setAgeRating] = useState('');
  const [showDatePicker, setShowDatePicker] = useState<'publishAt' | 'unpublishAt' | null>(null);
  const [tempDate, setTempDate] = useState<Date | null>(null);

  useEffect(() => {
    if (asset.publishRule) {
      const publishDate = asset.publishRule.publishAt ? new Date(asset.publishRule.publishAt) : null;
      const unpublishDate = asset.publishRule.unpublishAt ? new Date(asset.publishRule.unpublishAt) : null;
      const rating = asset.publishRule.ageRating || '';
      
      setPublishAt(publishDate);
      setUnpublishAt(unpublishDate);
      setAgeRating(rating);
    } else {
      setPublishAt(null);
      setUnpublishAt(null);
      setAgeRating('');
    }
  }, [asset]);

  const openDatePicker = (type: 'publishAt' | 'unpublishAt') => {
    const currentValue = type === 'publishAt' ? publishAt : unpublishAt;
    setTempDate(currentValue);
    setShowDatePicker(type);
  };

  const confirmDatePicker = () => {
    if (showDatePicker === 'publishAt') {
      setPublishAt(tempDate);
    } else if (showDatePicker === 'unpublishAt') {
      setUnpublishAt(tempDate);
    }
    setShowDatePicker(null);
    setTempDate(null);
  };

  const cancelDatePicker = () => {
    setShowDatePicker(null);
    setTempDate(null);
  };

  const handlePublish = async () => {
    await onPublish(publishAt, unpublishAt, ageRating);
  };

  return (
    <View style={styles.container}>
      <Text style={styles.sectionTitle}>Publish Settings</Text>
      
      <View style={styles.publishField}>
        <Text style={styles.publishLabel}>Publish Date</Text>
        <TouchableOpacity 
          style={styles.datePickerButton}
          onPress={() => openDatePicker('publishAt')}
        >
          <Text style={styles.datePickerButtonText}>
            {publishAt ? publishAt.toLocaleString() : 'Select date'}
          </Text>
        </TouchableOpacity>
      </View>

      <View style={styles.publishField}>
        <Text style={styles.publishLabel}>Unpublish Date</Text>
        <TouchableOpacity 
          style={styles.datePickerButton}
          onPress={() => openDatePicker('unpublishAt')}
        >
          <Text style={styles.datePickerButtonText}>
            {unpublishAt ? unpublishAt.toLocaleString() : 'Select date'}
          </Text>
        </TouchableOpacity>
      </View>

      <View style={styles.publishField}>
        <Text style={styles.publishLabel}>Age Rating</Text>
        <TextInput
          style={styles.publishInput}
          value={ageRating}
          onChangeText={setAgeRating}
          placeholder="e.g., PG, PG-13, R"
          placeholderTextColor="#999"
        />
      </View>

      <TouchableOpacity 
        style={[styles.updatePublishButton, publishing && styles.updatePublishButtonDisabled]}
        onPress={handlePublish}
        disabled={publishing}
      >
        {publishing ? (
          <ActivityIndicator size="small" color="#fff" />
        ) : (
          <Text style={styles.updatePublishButtonText}>Update Publish Settings</Text>
        )}
      </TouchableOpacity>

      {showDatePicker && (
        <View style={styles.modalOverlay}>
          <View style={styles.modalContainer}>
            <Text style={styles.modalTitle}>
              {showDatePicker === 'publishAt' ? 'Set Publish Date' : 'Set Unpublish Date'}
            </Text>
            <Text style={styles.modalMessage}>
              Select date and time
            </Text>
            
            <View style={styles.datePickerContainer}>
              <DatePicker
                selected={tempDate}
                onChange={(date: Date | null) => setTempDate(date)}
                showTimeSelect
                timeFormat="HH:mm"
                timeIntervals={15}
                dateFormat="MMMM d, yyyy h:mm aa"
                className="date-picker-input"
                placeholderText="Select date and time"
              />
            </View>
            
            <View style={styles.modalButtons}>
              <TouchableOpacity 
                style={[styles.modalButton, styles.cancelButton]}
                onPress={cancelDatePicker}
              >
                <Text style={styles.cancelButtonText}>Cancel</Text>
              </TouchableOpacity>
              <TouchableOpacity 
                style={[styles.modalButton, styles.publishButton]}
                onPress={confirmDatePicker}
              >
                <Text style={styles.publishButtonText}>Confirm</Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginTop: 20,
    paddingTop: 20,
    borderTopWidth: 1,
    borderTopColor: '#e0e0e0',
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  publishField: {
    marginBottom: 12,
  },
  publishLabel: {
    fontSize: 14,
    fontWeight: '600',
    marginBottom: 6,
    color: '#333',
  },
  datePickerButton: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 6,
    padding: 10,
    backgroundColor: '#fff',
  },
  datePickerButtonText: {
    fontSize: 14,
    color: '#333',
  },
  publishInput: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 6,
    padding: 10,
    fontSize: 14,
    backgroundColor: '#fff',
  },
  updatePublishButton: {
    backgroundColor: '#28a745',
    padding: 12,
    borderRadius: 6,
    alignItems: 'center',
    marginTop: 10,
  },
  updatePublishButtonDisabled: {
    backgroundColor: '#ccc',
  },
  updatePublishButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  modalOverlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 1000,
  },
  modalContainer: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 24,
    margin: 20,
    minWidth: 400,
    maxWidth: 600,
  },
  modalTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 10,
    textAlign: 'center',
  },
  modalMessage: {
    fontSize: 14,
    color: '#666',
    marginBottom: 20,
    textAlign: 'center',
    lineHeight: 20,
  },
  datePickerContainer: {
    width: '100%',
    minWidth: 360,
    marginBottom: 20,
    alignItems: 'center',
  },
  modalButtons: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: 10,
  },
  modalButton: {
    flex: 1,
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  cancelButton: {
    backgroundColor: '#f0f0f0',
  },
  publishButton: {
    backgroundColor: '#28a745',
  },
  cancelButtonText: {
    color: '#666',
    fontSize: 14,
    fontWeight: '600',
  },
  publishButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
}); 