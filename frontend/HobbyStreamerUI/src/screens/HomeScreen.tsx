import React, { useState, useEffect } from 'react';
import { 
  View, 
  Text, 
  StyleSheet, 
  ScrollView, 
  RefreshControl, 
  Alert,
  StatusBar,
  FlatList,
  ActivityIndicator,
  Button
} from 'react-native';
import { BucketRow } from '../components/BucketRow';
import { Asset } from '../types/asset';
import { Bucket } from '../types/asset';
import { apiService } from '../services/api';

export const HomeScreen: React.FC = () => {
  const [buckets, setBuckets] = useState<Bucket[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [nextKey, setNextKey] = useState<string | undefined>(undefined);
  const [hasMore, setHasMore] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const limit = 10;

  const loadBuckets = async (append = false) => {
    if (loadingMore && append) return;
    if (!append) setLoading(true);
    if (append) setLoadingMore(true);
    try {
      const response = await apiService.getBuckets(limit, append ? nextKey : undefined);
      const newBuckets = response.buckets;
      setBuckets(append ? [...buckets, ...newBuckets] : newBuckets);
      setNextKey(response.nextKey);
      setHasMore(!!response.nextKey);
    } catch (error) {
      console.error('Failed to load data:', error);
      Alert.alert('Error', 'Failed to load content. Please try again.');
    } finally {
      if (!append) setLoading(false);
      if (append) setLoadingMore(false);
    }
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await loadBuckets(false);
    setRefreshing(false);
  };

  useEffect(() => {
    loadBuckets(false);
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

  const renderBucket = ({ item }: { item: Bucket }) => {
    if (!item.name || !item.name.trim() || !item.assets || item.assets.length === 0) {
      return null;
    }
    return (
      <BucketRow
        key={item.id}
        title={item.name}
        assets={item.assets}
        onAssetPress={handleAssetPress}
      />
    );
  };

  if (loading && !refreshing) {
    return (
      <View style={styles.loadingContainer}>
        <Text style={styles.loadingText}>Loading...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#000000" />
      <FlatList
        data={buckets}
        renderItem={renderBucket}
        keyExtractor={item => item.id}
        onEndReached={() => {
          if (hasMore && !loadingMore) loadBuckets(true);
        }}
        onEndReachedThreshold={0.5}
        ListFooterComponent={loadingMore ? <ActivityIndicator size="small" color="#fff" /> : null}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <Text style={styles.emptyText}>No content available</Text>
          </View>
        }
      />
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