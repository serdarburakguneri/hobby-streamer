import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TextInput,
  TouchableOpacity,
  ScrollView,
  Alert,
  ActivityIndicator,
} from 'react-native';
import { useAssetService } from '../services/api';
import { AssetCreateDTO, AssetType, AssetGenre, formatGenreName } from '../types/asset';
import Layout from '../components/Layout';

interface CreateAssetScreenProps {
  onBack: () => void;
  onAssetCreated: () => void;
}

export default function CreateAssetScreen({ onBack, onAssetCreated }: CreateAssetScreenProps) {
  const assetService = useAssetService();
  const [formData, setFormData] = useState<AssetCreateDTO>({
    title: '',
    slug: '',
    description: '',
    type: AssetType.MOVIE,
    genre: '',
    genres: [],
    tags: [],
    metadata: {},
    ownerId: 'admin',
    parentId: undefined,
  });
  const [loading, setLoading] = useState(false);
  const [selectedGenres, setSelectedGenres] = useState<string[]>([]);
  const [seasonQuery, setSeasonQuery] = useState('');
  const [seasonResults, setSeasonResults] = useState<any[]>([]);
  const [seasonSearchLoading, setSeasonSearchLoading] = useState(false);
  const [seasonDropdownVisible, setSeasonDropdownVisible] = useState(false);

  const [seriesQuery, setSeriesQuery] = useState('');
  const [seriesResults, setSeriesResults] = useState<any[]>([]);
  const [seriesSearchLoading, setSeriesSearchLoading] = useState(false);
  const [seriesDropdownVisible, setSeriesDropdownVisible] = useState(false);
  useEffect(() => {
    let active = true;
    if (
      formData.type === AssetType.EPISODE &&
      seasonQuery.trim().length > 0
    ) {
      setSeasonSearchLoading(true);
      assetService
        .searchAssets(seasonQuery)
        .then((results: any) => {
          if (active) {
            setSeasonResults(
              results.assets.filter((a: any) => a.type === AssetType.SEASON)
            );
            setSeasonDropdownVisible(true);
          }
        })
        .catch(() => {
          if (active) setSeasonResults([]);
        })
        .finally(() => {
          if (active) setSeasonSearchLoading(false);
        });
    } else {
      setSeasonResults([]);
      setSeasonDropdownVisible(false);
    }
    return () => {
      active = false;
    };
  }, [seasonQuery, formData.type]);


  useEffect(() => {
    let active = true;
    if (
      formData.type === AssetType.SEASON &&
      seriesQuery.trim().length > 0
    ) {
      setSeriesSearchLoading(true);
      assetService
        .searchAssets(seriesQuery)
        .then((results: any) => {
          if (active) {
            setSeriesResults(
              results.assets.filter((a: any) => a.type === AssetType.SERIES)
            );
            setSeriesDropdownVisible(true);
          }
        })
        .catch(() => {
          if (active) setSeriesResults([]);
        })
        .finally(() => {
          if (active) setSeriesSearchLoading(false);
        });
    } else {
      setSeriesResults([]);
      setSeriesDropdownVisible(false);
    }
    return () => {
      active = false;
    };
  }, [seriesQuery, formData.type]);

  const handleInputChange = (field: keyof AssetCreateDTO, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleTagsChange = (value: string) => {
    const tags = value.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0);
    setFormData(prev => ({
      ...prev,
      tags,
    }));
  };

  const handleGenreToggle = (genre: string) => {
    setSelectedGenres(prev => {
      const isSelected = prev.includes(genre);
      if (isSelected) {
        const newGenres = prev.filter(g => g !== genre);
        setFormData(prevData => ({ ...prevData, genres: newGenres }));
        return newGenres;
      } else {
        const newGenres = [...prev, genre];
        setFormData(prevData => ({ ...prevData, genres: newGenres }));
        return newGenres;
      }
    });
  };

  const handleSeasonSelect = (season: any) => {
    setFormData(prev => ({ ...prev, parentId: season.id }));
    setSeasonQuery(season.title);
    setSeasonDropdownVisible(false);
  };

  const handleSeriesSelect = (series: any) => {
    setFormData(prev => ({ ...prev, parentId: series.id }));
    setSeriesQuery(series.title);
    setSeriesDropdownVisible(false);
  };

  const handleSubmit = async () => {
    console.log('Create button clicked');
    if (!formData.title?.trim()) {
      Alert.alert('Error', 'Title is required');
      return;
    }
    if (!formData.slug?.trim()) {
      Alert.alert('Error', 'Slug is required');
      return;
    }

    try {
      setLoading(true);
      console.log('Creating asset with data:', formData);
      console.log('Slug value being sent:', formData.slug);
      console.log('Slug type:', typeof formData.slug);
      console.log('Parent ID being sent:', formData.parentId);
      const result = await assetService.createAsset(formData);
      console.log('Asset created successfully:', result);
      console.log('Created asset parent:', result.parent);
      
      onAssetCreated();
      onBack();
      
    } catch (error) {
      console.error('Error creating asset:', error);
      Alert.alert('Error', 'Failed to create asset. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout
      headerTitle="Create New Asset"
      headerLeft={
        <TouchableOpacity style={styles.backButton} onPress={onBack}>
          <Text style={styles.backButtonText}>← Back</Text>
        </TouchableOpacity>
      }
      headerRight={
        <TouchableOpacity
          style={[styles.createButton, loading && styles.createButtonDisabled]}
          onPress={handleSubmit}
          disabled={loading}
        >
          {loading ? (
            <ActivityIndicator size="small" color="#fff" />
          ) : (
            <Text style={styles.createButtonText}>Create</Text>
          )}
        </TouchableOpacity>
      }
    >
      <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
        <View style={styles.form}>
          <View style={styles.inputGroup}>
            <Text style={styles.label}>Title *</Text>
            <TextInput
              style={styles.input}
              value={formData.title}
              onChangeText={(value) => handleInputChange('title', value)}
              placeholder="Enter asset title"
              placeholderTextColor="#999"
            />
          </View>

          <View style={styles.inputGroup}>
            <Text style={styles.label}>Slug *</Text>
            <TextInput
              style={styles.input}
              value={formData.slug}
              onChangeText={(value) => handleInputChange('slug', value)}
              placeholder="Enter unique slug (e.g. my-movie-title)"
              placeholderTextColor="#999"
              autoCapitalize="none"
              autoCorrect={false}
            />
          </View>

          <View style={styles.inputGroup}>
            <Text style={styles.label}>Description</Text>
            <TextInput
              style={[styles.input, styles.textArea]}
              value={formData.description}
              onChangeText={(value) => handleInputChange('description', value)}
              placeholder="Enter asset description"
              placeholderTextColor="#999"
              multiline
              numberOfLines={4}
              textAlignVertical="top"
            />
          </View>



          <View style={styles.inputGroup}>
            <Text style={styles.label}>Type</Text>
            <View style={styles.pickerContainer}>
              {Object.values(AssetType).map((type) => (
                <TouchableOpacity
                  key={type}
                  style={[
                    styles.pickerOption,
                    formData.type === type && styles.pickerOptionSelected
                  ]}
                  onPress={() => {
                    handleInputChange('type', type);
                    // Reset parentId and search queries if not episode or season
                    if (type !== AssetType.EPISODE && type !== AssetType.SEASON) {
                      setFormData(prev => ({ ...prev, parentId: undefined }));
                      setSeasonQuery('');
                      setSeasonResults([]);
                      setSeriesQuery('');
                      setSeriesResults([]);
                    }
                  }}
                >
                  <Text style={[
                    styles.pickerOptionText,
                    formData.type === type && styles.pickerOptionTextSelected
                  ]}>
                    {type}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          {/* Bind to TV Show field for EPISODE */}
          {formData.type === AssetType.EPISODE && (
            <View style={styles.inputGroup}>
              <Text style={styles.label}>Bind to TV Show</Text>
              <TextInput
                style={styles.input}
                value={seasonQuery}
                onChangeText={text => {
                  setSeasonQuery(text);
                  setFormData(prev => ({ ...prev, parentId: undefined }));
                }}
                placeholder="Type season title..."
                placeholderTextColor="#999"
                onFocus={() => {
                  if (seasonResults.length > 0) setSeasonDropdownVisible(true);
                }}
              />
              {seasonSearchLoading && <ActivityIndicator size="small" color="#007AFF" />}
              {seasonDropdownVisible && seasonResults.length > 0 && (
                <View style={{ backgroundColor: '#fff', borderWidth: 1, borderColor: '#ddd', borderRadius: 8, marginTop: 4, maxHeight: 120 }}>
                  <ScrollView>
                    {seasonResults.map((season) => (
                      <TouchableOpacity
                        key={season.id}
                        style={{ padding: 10, borderBottomWidth: 1, borderBottomColor: '#eee' }}
                        onPress={() => handleSeasonSelect(season)}
                      >
                        <Text style={{ fontSize: 16 }}>{season.title}</Text>
                        <Text style={{ fontSize: 12, color: '#888' }}>ID: {season.id}</Text>
                      </TouchableOpacity>
                    ))}
                  </ScrollView>
                </View>
              )}
              {formData.parentId && (
                <Text style={{ color: '#007AFF', marginTop: 4 }}>Selected season ID: {formData.parentId}</Text>
              )}
            </View>
          )}

          {/* Bind to Series field for SEASON */}
          {formData.type === AssetType.SEASON && (
            <View style={styles.inputGroup}>
              <Text style={styles.label}>Bind to Series</Text>
              <TextInput
                style={styles.input}
                value={seriesQuery}
                onChangeText={text => {
                  setSeriesQuery(text);
                  setFormData(prev => ({ ...prev, parentId: undefined }));
                }}
                placeholder="Type series title..."
                placeholderTextColor="#999"
                onFocus={() => {
                  if (seriesResults.length > 0) setSeriesDropdownVisible(true);
                }}
              />
              {seriesSearchLoading && <ActivityIndicator size="small" color="#007AFF" />}
              {seriesDropdownVisible && seriesResults.length > 0 && (
                <View style={{ backgroundColor: '#fff', borderWidth: 1, borderColor: '#ddd', borderRadius: 8, marginTop: 4, maxHeight: 120 }}>
                  <ScrollView>
                    {seriesResults.map((series) => (
                      <TouchableOpacity
                        key={series.id}
                        style={{ padding: 10, borderBottomWidth: 1, borderBottomColor: '#eee' }}
                        onPress={() => handleSeriesSelect(series)}
                      >
                        <Text style={{ fontSize: 16 }}>{series.title}</Text>
                        <Text style={{ fontSize: 12, color: '#888' }}>ID: {series.id}</Text>
                      </TouchableOpacity>
                    ))}
                  </ScrollView>
                </View>
              )}
              {formData.parentId && (
                <Text style={{ color: '#007AFF', marginTop: 4 }}>Selected series ID: {formData.parentId}</Text>
              )}
            </View>
          )}

          <View style={styles.inputGroup}>
            <Text style={styles.label}>Primary Genre</Text>
            <View style={styles.pickerContainer}>
              {Object.values(AssetGenre).map((genre) => (
                <TouchableOpacity
                  key={genre}
                  style={[
                    styles.pickerOption,
                    formData.genre === genre && styles.pickerOptionSelected
                  ]}
                  onPress={() => handleInputChange('genre', genre)}
                >
                  <Text style={[
                    styles.pickerOptionText,
                    formData.genre === genre && styles.pickerOptionTextSelected
                  ]}>
                    {formatGenreName(genre)}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          <View style={styles.inputGroup}>
            <Text style={styles.label}>Additional Genres</Text>
            <Text style={styles.subLabel}>Select from predefined genres</Text>
            
            <View style={styles.genreGrid}>
              {Object.values(AssetGenre).map((genre) => (
                <TouchableOpacity
                  key={genre}
                  style={[
                    styles.genreChip,
                    selectedGenres.includes(genre) && styles.genreChipSelected
                  ]}
                  onPress={() => handleGenreToggle(genre)}
                >
                  <Text style={[
                    styles.genreChipText,
                    selectedGenres.includes(genre) && styles.genreChipTextSelected
                  ]}>
                    {formatGenreName(genre)}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>

            {selectedGenres.length > 0 && (
              <View style={styles.selectedGenresContainer}>
                <Text style={styles.selectedGenresLabel}>Selected Genres:</Text>
                <View style={styles.selectedGenresList}>
                  {selectedGenres.map((genre) => (
                    <View key={genre} style={styles.selectedGenreItem}>
                      <Text style={styles.selectedGenreText}>{formatGenreName(genre)}</Text>
                      <TouchableOpacity
                        style={styles.removeGenreButton}
                        onPress={() => handleGenreToggle(genre)}
                      >
                        <Text style={styles.removeGenreButtonText}>×</Text>
                      </TouchableOpacity>
                    </View>
                  ))}
                </View>
              </View>
            )}
          </View>

          <View style={styles.inputGroup}>
            <Text style={styles.label}>Tags (comma-separated)</Text>
            <TextInput
              style={styles.input}
              value={formData.tags?.join(', ')}
              onChangeText={handleTagsChange}
              placeholder="tag1, tag2, tag3"
              placeholderTextColor="#999"
            />
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
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  backButton: {
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
    borderRadius: 6,
  },
  backButtonText: {
    fontSize: 16,
    color: '#007AFF',
    fontWeight: '600',
  },
  headerTitleContainer: {
    flex: 1,
    alignItems: 'center',
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: 'bold',
  },
  headerSubtitle: {
    fontSize: 12,
    color: '#666',
    marginTop: 2,
  },
  content: {
    flex: 1,
  },
  form: {
    padding: 20,
  },
  inputGroup: {
    marginBottom: 20,
  },
  label: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 8,
    color: '#333',
  },
  input: {
    backgroundColor: '#fff',
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
    color: '#333',
  },
  textArea: {
    height: 100,
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
    fontSize: 14,
    fontWeight: '600',
    color: '#fff',
  },
  pickerContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  pickerOption: {
    backgroundColor: '#f5f5f5',
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 6,
    paddingHorizontal: 12,
    paddingVertical: 8,
  },
  pickerOptionSelected: {
    backgroundColor: '#007AFF',
    borderColor: '#007AFF',
  },
  pickerOptionText: {
    fontSize: 14,
    color: '#666',
    fontWeight: '500',
  },
  pickerOptionTextSelected: {
    color: '#fff',
    fontWeight: '600',
  },
  subLabel: {
    fontSize: 14,
    color: '#666',
    marginBottom: 12,
  },
  genreGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
    marginBottom: 16,
  },
  genreChip: {
    backgroundColor: '#f5f5f5',
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 16,
    paddingHorizontal: 12,
    paddingVertical: 6,
  },
  genreChipSelected: {
    backgroundColor: '#007AFF',
    borderColor: '#007AFF',
  },
  genreChipText: {
    fontSize: 12,
    color: '#666',
    fontWeight: '500',
  },
  genreChipTextSelected: {
    color: '#fff',
    fontWeight: '600',
  },
  selectedGenresContainer: {
    marginTop: 8,
  },
  selectedGenresLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333',
    marginBottom: 8,
  },
  selectedGenresList: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  selectedGenreItem: {
    backgroundColor: '#e3f2fd',
    borderWidth: 1,
    borderColor: '#2196f3',
    borderRadius: 16,
    paddingHorizontal: 12,
    paddingVertical: 6,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  selectedGenreText: {
    fontSize: 12,
    color: '#1976d2',
    fontWeight: '500',
  },
  removeGenreButton: {
    width: 16,
    height: 16,
    borderRadius: 8,
    backgroundColor: '#f44336',
    justifyContent: 'center',
    alignItems: 'center',
  },
  removeGenreButtonText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: 'bold',
  },
}); 