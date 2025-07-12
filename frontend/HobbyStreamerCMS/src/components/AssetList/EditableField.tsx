import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, StyleSheet, ScrollView } from 'react-native';
import { AssetGenre, formatGenreName } from '../../types/asset';

interface EditableFieldProps {
  label: string;
  field: string;
  value: any;
  onUpdate: (field: string, value: any) => Promise<void>;
  placeholder?: string;
  type?: 'text' | 'genre' | 'multiGenre' | 'tags';
}

export default function EditableField({ 
  label, 
  field, 
  value, 
  onUpdate, 
  placeholder,
  type = 'text'
}: EditableFieldProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState('');
  const [selectedGenres, setSelectedGenres] = useState<string[]>([]);

  const startEditing = () => {
    setIsEditing(true);
    if (type === 'multiGenre') {
      setSelectedGenres(value || []);
      setEditValue('');
    } else if (type === 'tags') {
      setEditValue((value || []).join(', '));
    } else {
      setEditValue(value || '');
    }
  };

  const cancelEditing = () => {
    setIsEditing(false);
    setEditValue('');
    setSelectedGenres([]);
  };

  const saveEdit = async () => {
    try {
      let valueToSave: any;
      
      if (type === 'multiGenre') {
        valueToSave = selectedGenres;
      } else if (type === 'tags') {
        valueToSave = editValue.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0);
      } else {
        valueToSave = editValue;
      }
      
      await onUpdate(field, valueToSave);
      setIsEditing(false);
    } catch (error) {
      console.error('Error updating field:', error);
    }
  };

  const handleGenreToggle = (genre: string) => {
    setSelectedGenres(prev => {
      const isSelected = prev.includes(genre);
      if (isSelected) {
        return prev.filter(g => g !== genre);
      } else {
        return [...prev, genre];
      }
    });
  };

  const handleRemoveGenre = (genreToRemove: string) => {
    setSelectedGenres(prev => prev.filter(g => g !== genreToRemove));
  };

  const renderEditContent = () => {
    if (type === 'genre') {
      return (
        <ScrollView style={styles.genrePickerScroll} horizontal showsHorizontalScrollIndicator={false}>
          <View style={styles.genrePickerContainer}>
            {Object.values(AssetGenre).map((genre) => (
              <TouchableOpacity
                key={genre}
                style={[
                  styles.genrePickerOption,
                  editValue === genre && styles.genrePickerOptionSelected
                ]}
                onPress={() => setEditValue(genre)}
              >
                <Text style={[
                  styles.genrePickerOptionText,
                  editValue === genre && styles.genrePickerOptionTextSelected
                ]}>
                  {formatGenreName(genre)}
                </Text>
              </TouchableOpacity>
            ))}
          </View>
        </ScrollView>
      );
    }

    if (type === 'multiGenre') {
      return (
        <ScrollView style={styles.genrePickerScroll} showsVerticalScrollIndicator={false}>
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
              <Text style={styles.selectedGenresLabel}>Selected:</Text>
              <View style={styles.selectedGenresList}>
                {selectedGenres.map((genre) => (
                  <View key={genre} style={styles.selectedGenreItem}>
                    <Text style={styles.selectedGenreText}>{formatGenreName(genre)}</Text>
                    <TouchableOpacity
                      style={styles.removeGenreButton}
                      onPress={() => handleRemoveGenre(genre)}
                    >
                      <Text style={styles.removeGenreButtonText}>×</Text>
                    </TouchableOpacity>
                  </View>
                ))}
              </View>
            </View>
          )}
        </ScrollView>
      );
    }

    return (
      <TextInput
        style={styles.editInput}
        value={editValue}
        onChangeText={setEditValue}
        placeholder={placeholder}
        autoFocus
      />
    );
  };

  const renderDisplayValue = () => {
    if (type === 'genre') {
      return value ? formatGenreName(value) : 'Click to edit';
    }
    
    if (type === 'multiGenre') {
      return value?.length > 0 ? value.map((g: string) => formatGenreName(g)).join(', ') : 'Click to edit';
    }
    
    if (type === 'tags') {
      return value?.length > 0 ? value.join(', ') : 'Click to edit';
    }
    
    return value || 'Click to edit';
  };

  return (
    <View style={styles.container}>
      <Text style={styles.label}>{label}:</Text>
      {isEditing ? (
        <View style={styles.editContainer}>
          {renderEditContent()}
          <View style={styles.editButtons}>
            <TouchableOpacity style={styles.editButton} onPress={saveEdit}>
              <Text style={styles.editButtonText}>✓</Text>
            </TouchableOpacity>
            <TouchableOpacity style={styles.editButton} onPress={cancelEditing}>
              <Text style={styles.editButtonText}>✗</Text>
            </TouchableOpacity>
          </View>
        </View>
      ) : (
        <TouchableOpacity 
          style={styles.editableValue}
          onPress={startEditing}
          activeOpacity={0.7}
        >
          <View style={styles.editableContent}>
            <Text style={styles.detailValue}>{renderDisplayValue()}</Text>
            <Text style={styles.editIcon}>✏️</Text>
          </View>
        </TouchableOpacity>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 5,
  },
  label: {
    fontSize: 14,
    fontWeight: '500',
  },
  editContainer: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
  },
  editInput: {
    flex: 1,
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 4,
    padding: 8,
    fontSize: 14,
    backgroundColor: '#fff',
  },
  editButtons: {
    flexDirection: 'row',
    marginLeft: 8,
  },
  editButton: {
    width: 30,
    height: 30,
    borderRadius: 15,
    justifyContent: 'center',
    alignItems: 'center',
    marginLeft: 4,
  },
  editButtonText: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  editableValue: {
    flex: 1,
    padding: 8,
    borderRadius: 4,
    backgroundColor: '#f8f9fa',
    borderWidth: 1,
    borderColor: '#e9ecef',
    minHeight: 40,
    justifyContent: 'center',
  },
  editableContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  detailValue: {
    fontSize: 14,
    fontWeight: '500',
  },
  editIcon: {
    fontSize: 16,
    marginLeft: 8,
  },
  genrePickerScroll: {
    maxHeight: 120,
    marginBottom: 10,
  },
  genrePickerContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 6,
    marginBottom: 8,
  },
  genrePickerOption: {
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 12,
    backgroundColor: '#f0f0f0',
    borderWidth: 1,
    borderColor: '#ddd',
    marginBottom: 4,
  },
  genrePickerOptionSelected: {
    backgroundColor: '#007AFF',
    borderColor: '#007AFF',
  },
  genrePickerOptionText: {
    fontSize: 11,
    color: '#666',
  },
  genrePickerOptionTextSelected: {
    color: '#fff',
    fontWeight: '600',
  },
  genreGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 4,
    marginBottom: 8,
  },
  genreChip: {
    paddingHorizontal: 6,
    paddingVertical: 3,
    borderRadius: 10,
    backgroundColor: '#f0f0f0',
    borderWidth: 1,
    borderColor: '#ddd',
    marginBottom: 3,
  },
  genreChipSelected: {
    backgroundColor: '#007AFF',
    borderColor: '#007AFF',
  },
  genreChipText: {
    fontSize: 9,
    color: '#666',
  },
  genreChipTextSelected: {
    color: '#fff',
    fontWeight: '600',
  },
  selectedGenresContainer: {
    marginTop: 8,
  },
  selectedGenresLabel: {
    fontSize: 12,
    fontWeight: '600',
    marginBottom: 4,
    color: '#666',
  },
  selectedGenresList: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 4,
  },
  selectedGenreItem: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#e3f2fd',
    borderWidth: 1,
    borderColor: '#2196f3',
    borderRadius: 12,
    paddingHorizontal: 8,
    paddingVertical: 4,
  },
  selectedGenreText: {
    fontSize: 10,
    color: '#1976d2',
    marginRight: 4,
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
    fontSize: 10,
    fontWeight: 'bold',
  },
}); 