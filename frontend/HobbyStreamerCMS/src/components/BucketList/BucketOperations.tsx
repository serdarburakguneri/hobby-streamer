import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
  Modal,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Bucket, BucketType } from '../../types/asset';

interface BucketOperationsProps {
  bucket: Bucket | null;
  onDelete: () => void;
  onUpdate: (field: string, value: any) => void;
  deleting: boolean;
  updating: boolean;
}

export default function BucketOperations({ 
  bucket, 
  onDelete, 
  onUpdate, 
  deleting, 
  updating 
}: BucketOperationsProps) {
  if (!bucket) {
    return (
      <View style={styles.emptyContainer}>
        <Ionicons name="settings-outline" size={32} color="#ccc" />
        <Text style={styles.emptyText}>Select a bucket to view operations</Text>
      </View>
    );
  }

  const handleDelete = () => {
    onDelete();
  };

  const handleStatusToggle = () => {
    let newStatus: string;
    switch (bucket.status?.toLowerCase()) {
      case 'draft':
        newStatus = 'active';
        break;
      case 'active':
        newStatus = 'inactive';
        break;
      case 'inactive':
        newStatus = 'active';
        break;
      default:
        newStatus = 'active';
    }
    onUpdate('status', newStatus);
  };

  const formatBucketType = (type: BucketType): string => {
    return type.replace('_', ' ').toLowerCase();
  };

  const getStatusIcon = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'active':
        return 'checkmark-circle';
      case 'inactive':
        return 'pause-circle';
      case 'draft':
        return 'create';
      default:
        return 'help-circle';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'active':
        return '#4CAF50';
      case 'inactive':
        return '#FF9800';
      case 'draft':
        return '#9E9E9E';
      default:
        return '#9E9E9E';
    }
  };

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
          onPress={handleDelete}
          disabled={deleting}
        >
          {deleting ? (
            <ActivityIndicator size="small" color="#fff" />
          ) : (
            <>
              <Ionicons name="trash" size={16} color="#fff" />
              <Text style={styles.deleteButtonText}>Delete Bucket</Text>
            </>
          )}
        </TouchableOpacity>
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Ionicons name="information-circle" size={20} color="#333" />
          <Text style={styles.sectionTitle}>Status</Text>
        </View>
        <View style={styles.statusContainer}>
          <Ionicons 
            name={getStatusIcon(bucket.status || 'draft') as any} 
            size={16} 
            color={getStatusColor(bucket.status || 'draft')} 
          />
          <Text style={styles.statusText}>{bucket.status || 'draft'}</Text>
        </View>
        <TouchableOpacity
          style={[
            styles.statusButton,
            bucket.status?.toLowerCase() === 'active' ? styles.deactivateButton : styles.activateButton
          ]}
          onPress={handleStatusToggle}
          disabled={updating}
        >
          {updating ? (
            <ActivityIndicator size="small" color="#fff" />
          ) : (
            <>
              <Ionicons 
                name={bucket.status?.toLowerCase() === 'active' ? 'pause-circle' : 'checkmark-circle'} 
                size={16} 
                color="#fff" 
              />
              <Text style={styles.statusButtonText}>
                {bucket.status?.toLowerCase() === 'draft' ? 'Activate' :
                 bucket.status?.toLowerCase() === 'active' ? 'Deactivate' :
                 bucket.status?.toLowerCase() === 'inactive' ? 'Activate' : 'Activate'}
              </Text>
            </>
          )}
        </TouchableOpacity>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  emptyText: {
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
  statusContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 12,
  },
  statusText: {
    fontSize: 16,
    color: '#333',
    textTransform: 'capitalize',
  },
  statusButton: {
    paddingVertical: 10,
    paddingHorizontal: 16,
    borderRadius: 6,
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'center',
    gap: 8,
  },
  activateButton: {
    backgroundColor: '#4CAF50',
  },
  deactivateButton: {
    backgroundColor: '#FF9800',
  },
  statusButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  infoItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  infoLabelContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  infoLabel: {
    fontSize: 14,
    color: '#666',
  },
  infoValue: {
    fontSize: 14,
    color: '#333',
    fontWeight: '500',
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
}); 