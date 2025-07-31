import React, { useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, ActivityIndicator, TextInput } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Asset } from '../../types/asset';
import PublishSettings from './PublishSettings';

interface AssetOperationsProps {
  asset: Asset | null;
  onDelete: () => void;
  onPublish: (publishAt: Date | null, unpublishAt: Date | null, ageRating: string) => Promise<void>;
  deleting: boolean;
  publishing: boolean;
}

export default function AssetOperations({ 
  asset, 
  onDelete, 
  onPublish,
  deleting,
  publishing
}: AssetOperationsProps) {
  if (!asset) {
    return (
      <View style={styles.placeholder}>
        <Ionicons name="settings-outline" size={32} color="#ccc" />
        <Text style={styles.placeholderText}>Select an asset to perform operations</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.titleContainer}>
        <Ionicons name="settings" size={24} color="#333" />
        <Text style={styles.title}>Operations</Text>
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Ionicons name="warning" size={20} color="#ff3b30" />
          <Text style={styles.sectionTitle}>Danger Zone</Text>
        </View>
        <TouchableOpacity 
          style={[styles.deleteButton, deleting && styles.deleteButtonDisabled]}
          onPress={onDelete}
          disabled={deleting}
          activeOpacity={0.7}
        >
          {deleting ? (
            <ActivityIndicator size="small" color="#fff" />
          ) : (
            <>
              <Ionicons name="trash" size={16} color="#fff" />
              <Text style={styles.deleteButtonText}>Delete Asset</Text>
            </>
          )}
        </TouchableOpacity>
      </View>



      <PublishSettings 
        asset={asset}
        onPublish={onPublish}
        publishing={publishing}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  placeholder: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  placeholderText: {
    fontSize: 16,
    color: '#666',
    textAlign: 'center',
    marginTop: 8,
  },
  titleContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 16,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#333',
  },
  section: {
    backgroundColor: '#fff',
    borderRadius: 8,
    padding: 16,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  sectionHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 12,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
  },
  deleteButton: {
    backgroundColor: '#ff3b30',
    paddingVertical: 12,
    paddingHorizontal: 16,
    borderRadius: 6,
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'center',
    gap: 8,
  },
  deleteButtonDisabled: {
    backgroundColor: '#ccc',
  },
  deleteButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  timestampsSection: {
    backgroundColor: '#fff',
    borderRadius: 8,
    padding: 16,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  detailRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  detailLabelContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  detailLabel: {
    fontSize: 14,
    color: '#666',
  },
  detailValue: {
    fontSize: 14,
    color: '#333',
    fontWeight: '500',
  },
}); 