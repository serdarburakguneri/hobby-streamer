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
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const loadData = async () => {
    try {
      const bucketsData = await apiService.getBuckets();
      console.log('Loaded buckets:', bucketsData.length);
      
      const bucketsWithAssets = await Promise.all(
        bucketsData.map(async (bucket) => {
          try {
            const assets = await apiService.getAssetsInBucket(bucket.key);
            console.log(`Bucket "${bucket.name}" has ${assets.length} assets`);
            return {
              ...bucket,
              assets
            };
          } catch (error) {
            console.error(`Failed to load assets for bucket ${bucket.name}:`, error);
            return {
              ...bucket,
              assets: []
            };
          }
        })
      );
      
      setBuckets(bucketsWithAssets);
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
        {buckets.map((bucket) => (
          <BucketRow
            key={bucket.id}
            title={bucket.name}
            assets={bucket.assets || []}
            onAssetPress={handleAssetPress}
          />
        ))}

        {buckets.length === 0 && (
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