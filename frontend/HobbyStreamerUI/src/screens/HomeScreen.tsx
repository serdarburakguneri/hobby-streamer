import React, { useState, useEffect } from 'react';
import { 
  View, 
  Text, 
  StyleSheet, 
  ScrollView, 
  RefreshControl, 
  Alert,
  StatusBar 
} from 'react-native';
import { BucketRow } from '../components/BucketRow';
import { Asset } from '../types/asset';
import { Bucket } from '../types/asset';
import { apiService } from '../services/api';

export const HomeScreen: React.FC = () => {
  const [buckets, setBuckets] = useState<Bucket[]>([]);
  const [allAssets, setAllAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const loadData = async () => {
    try {
      const [bucketsData, assetsData] = await Promise.all([
        apiService.getBuckets(),
        apiService.getAssets(),
      ]);
      
      // For now, just set buckets without individual assets since the endpoint is failing
      setBuckets(bucketsData);
      setAllAssets(assetsData);
    } catch (error) {
      console.error('Failed to load data:', error);
      Alert.alert('Error', 'Failed to load content. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await loadData();
    setRefreshing(false);
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleAssetPress = (asset: Asset) => {
    Alert.alert(
      asset.title || 'Asset Details',
      `Genre: ${asset.genre || 'Unknown'}\nType: ${asset.type}`,
      [
        { text: 'Cancel', style: 'cancel' },
        { text: 'Watch', onPress: () => console.log('Watch asset:', asset.slug) }
      ]
    );
  };

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <Text style={styles.loadingText}>Loading...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#000000" />
      
      <ScrollView
        style={styles.scrollView}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        {/* All Assets Section */}
        {allAssets.length > 0 && (
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>All Content</Text>
            <BucketRow
              title="Available Content"
              assets={allAssets}
              onAssetPress={handleAssetPress}
            />
          </View>
        )}

        {/* Bucket Sections */}
        {buckets.map((bucket) => (
          <View key={bucket.id} style={styles.section}>
            <Text style={styles.sectionTitle}>{bucket.name}</Text>
            <Text style={styles.sectionDescription}>{bucket.description}</Text>
            <Text style={styles.sectionInfo}>Bucket Key: {bucket.key}</Text>
          </View>
        ))}

        {buckets.length === 0 && allAssets.length === 0 && (
          <View style={styles.emptyContainer}>
            <Text style={styles.emptyText}>No content available</Text>
          </View>
        )}
      </ScrollView>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000000',
  },
  scrollView: {
    flex: 1,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#000000',
  },
  loadingText: {
    color: '#ffffff',
    fontSize: 16,
  },
  section: {
    marginVertical: 16,
    paddingHorizontal: 16,
  },
  sectionTitle: {
    color: '#ffffff',
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  sectionDescription: {
    color: '#cccccc',
    fontSize: 14,
    marginBottom: 4,
  },
  sectionInfo: {
    color: '#888888',
    fontSize: 12,
    marginBottom: 12,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: 100,
  },
  emptyText: {
    color: '#666666',
    fontSize: 16,
  },
}); 