import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TextInput,
  TouchableOpacity,
  Alert,
  ActivityIndicator,
  Linking,
  Platform,
} from 'react-native';
import { Button, Card, IconButton, Chip } from 'react-native-paper';
import * as Location from 'expo-location';
import * as Speech from 'expo-speech';
import { Ionicons } from '@expo/vector-icons';
import { useAppContext } from '../context/AppContext';
import { extractMemoryInfo, queryMemories, localMemorySearch } from '../utils/llmService';

/**
 * LostToFound Screen - AI-powered memory assistant
 * Allows users to record and query item locations using voice or text
 */
const LostToFoundScreen = () => {
  const {
    memories,
    addMemory,
    deleteMemory,
    conversationHistory,
    addToConversationHistory,
    clearConversationHistory,
    isPlaying,
    pauseSound,
    resumeSound,
    duckVolume,
    restoreVolume,
    settings,
  } = useAppContext();

  const [input, setInput] = useState('');
  const [isRecording, setIsRecording] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [response, setResponse] = useState('');
  const [currentLocation, setCurrentLocation] = useState(null);
  const [showMemories, setShowMemories] = useState(false);

  useEffect(() => {
    requestPermissions();
  }, []);

  /**
   * Request necessary permissions
   */
  const requestPermissions = async () => {
    try {
      // Request location permission
      const { status: locationStatus } = await Location.requestForegroundPermissionsAsync();
      if (locationStatus !== 'granted') {
        Alert.alert('Permission Denied', 'Location permission is needed to remember where you place items.');
      }

      // Get current location
      const location = await Location.getCurrentPositionAsync({});
      setCurrentLocation(location);
    } catch (error) {
      console.error('Permission error:', error);
    }
  };

  /**
   * Handle voice input (simplified - in real app, use speech recognition library)
   */
  const handleVoiceInput = async () => {
    Alert.alert(
      'Voice Input',
      'Voice recognition would be implemented here using expo-speech or a third-party library like @react-native-voice/voice. For now, please use text input.',
      [{ text: 'OK' }]
    );
  };

  /**
   * Handle user input (store or query)
   */
  const handleSubmit = async () => {
    if (!input.trim()) {
      Alert.alert('Error', 'Please enter a message');
      return;
    }

    setIsProcessing(true);
    setResponse('');

    // Handle audio priority
    if (isPlaying) {
      if (settings.priorityBehavior === 'pause') {
        await pauseSound();
      } else if (settings.priorityBehavior === 'duck') {
        await duckVolume();
      }
    }

    try {
      // Add user message to conversation history
      await addToConversationHistory('user', input);

      // Extract information using LLM
      const extractResult = await extractMemoryInfo(input, conversationHistory);

      if (extractResult.success && extractResult.data) {
        const data = extractResult.data;

        if (data.action === 'query') {
          // User is querying for an item
          await handleQuery(input);
        } else if (data.item && data.location) {
          // User is storing an item
          await handleStore(data);
        } else {
          // Unclear intent, try querying
          await handleQuery(input);
        }
      } else {
        // LLM failed, use fallback
        const fallbackResult = localMemorySearch(input, memories);
        setResponse(fallbackResult.content);
        await addToConversationHistory('assistant', fallbackResult.content);
        speakResponse(fallbackResult.content);
      }
    } catch (error) {
      console.error('Error processing input:', error);
      Alert.alert('Error', 'Failed to process your request. Please try again.');
    } finally {
      setIsProcessing(false);
      setInput('');

      // Restore audio
      if (isPlaying) {
        if (settings.priorityBehavior === 'pause') {
          await resumeSound();
        } else if (settings.priorityBehavior === 'duck') {
          await restoreVolume();
        }
      }
    }
  };

  /**
   * Handle storing a new item
   */
  const handleStore = async (data) => {
    try {
      // Get current location if available
      let coords = null;
      if (currentLocation) {
        coords = {
          latitude: currentLocation.coords.latitude,
          longitude: currentLocation.coords.longitude,
        };
      }

      // Create memory entry
      const memory = {
        item: data.item,
        location: data.location,
        coordinates: coords,
      };

      await addMemory(memory);

      const responseText = `Got it! I've remembered that your ${data.item} is at ${data.location}.`;
      setResponse(responseText);
      await addToConversationHistory('assistant', responseText);
      speakResponse(responseText);
    } catch (error) {
      console.error('Error storing memory:', error);
      Alert.alert('Error', 'Failed to store memory');
    }
  };

  /**
   * Handle querying for an item
   */
  const handleQuery = async (query) => {
    try {
      // Try LLM query first
      const llmResult = await queryMemories(query, memories, conversationHistory);

      if (llmResult.success) {
        setResponse(llmResult.content);
        await addToConversationHistory('assistant', llmResult.content);
        speakResponse(llmResult.content);
      } else {
        // Fallback to local search
        const fallbackResult = localMemorySearch(query, memories);
        setResponse(fallbackResult.content);
        await addToConversationHistory('assistant', fallbackResult.content);
        speakResponse(fallbackResult.content);
      }
    } catch (error) {
      console.error('Error querying memories:', error);
      const fallbackResult = localMemorySearch(query, memories);
      setResponse(fallbackResult.content);
      speakResponse(fallbackResult.content);
    }
  };

  /**
   * Speak response using text-to-speech
   */
  const speakResponse = (text) => {
    Speech.speak(text, {
      language: 'en',
      pitch: 1.0,
      rate: 0.9,
    });
  };

  /**
   * Open location in Google Maps
   */
  const openInMaps = (coordinates) => {
    if (!coordinates) {
      Alert.alert('No Location', 'No coordinates available for this item');
      return;
    }

    const { latitude, longitude } = coordinates;
    const url = Platform.select({
      ios: `maps:0,0?q=${latitude},${longitude}`,
      android: `geo:0,0?q=${latitude},${longitude}`,
    });

    Linking.openURL(url).catch(() => {
      Alert.alert('Error', 'Unable to open maps');
    });
  };

  /**
   * Delete a memory
   */
  const handleDeleteMemory = (id) => {
    Alert.alert(
      'Delete Memory',
      'Are you sure you want to delete this memory?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Delete',
          style: 'destructive',
          onPress: () => deleteMemory(id),
        },
      ]
    );
  };

  /**
   * Clear conversation history
   */
  const handleClearHistory = () => {
    Alert.alert(
      'Clear History',
      'Are you sure you want to clear the conversation history?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Clear',
          style: 'destructive',
          onPress: () => {
            clearConversationHistory();
            setResponse('');
          },
        },
      ]
    );
  };

  return (
    <View style={styles.container}>
      <ScrollView style={styles.content}>
        {/* Instructions */}
        <Card style={styles.card}>
          <Card.Content>
            <Text style={styles.instructionTitle}>How to use:</Text>
            <Text style={styles.instruction}>
              • Say or type: "I'm leaving my keys on the kitchen table"
            </Text>
            <Text style={styles.instruction}>
              • Ask: "Where did I keep my keys?"
            </Text>
          </Card.Content>
        </Card>

        {/* Response Display */}
        {response ? (
          <Card style={styles.responseCard}>
            <Card.Content>
              <Text style={styles.responseTitle}>Response:</Text>
              <Text style={styles.responseText}>{response}</Text>
            </Card.Content>
          </Card>
        ) : null}

        {/* Memories List Toggle */}
        <TouchableOpacity
          style={styles.toggleButton}
          onPress={() => setShowMemories(!showMemories)}
        >
          <Text style={styles.toggleButtonText}>
            {showMemories ? 'Hide' : 'Show'} Stored Memories ({memories.length})
          </Text>
          <Ionicons
            name={showMemories ? 'chevron-up' : 'chevron-down'}
            size={20}
            color="#6200ee"
          />
        </TouchableOpacity>

        {/* Memories List */}
        {showMemories && (
          <View style={styles.memoriesContainer}>
            {memories.length === 0 ? (
              <Card style={styles.card}>
                <Card.Content>
                  <Text style={styles.emptyText}>No memories stored yet</Text>
                </Card.Content>
              </Card>
            ) : (
              memories.map((memory) => (
                <Card key={memory.id} style={styles.memoryCard}>
                  <Card.Content>
                    <View style={styles.memoryHeader}>
                      <Text style={styles.memoryItem}>{memory.item}</Text>
                      <IconButton
                        icon="delete"
                        size={20}
                        onPress={() => handleDeleteMemory(memory.id)}
                      />
                    </View>
                    <Text style={styles.memoryLocation}>📍 {memory.location}</Text>
                    <Text style={styles.memoryTime}>
                      {new Date(memory.timestamp).toLocaleString()}
                    </Text>
                    {memory.coordinates && (
                      <Button
                        mode="outlined"
                        onPress={() => openInMaps(memory.coordinates)}
                        style={styles.navigateButton}
                        icon="map-marker"
                      >
                        Navigate
                      </Button>
                    )}
                  </Card.Content>
                </Card>
              ))
            )}
          </View>
        )}

        {/* Clear History Button */}
        {conversationHistory.length > 0 && (
          <Button
            mode="text"
            onPress={handleClearHistory}
            style={styles.clearButton}
          >
            Clear Conversation History
          </Button>
        )}
      </ScrollView>

      {/* Input Section */}
      <View style={styles.inputContainer}>
        <TextInput
          style={styles.input}
          placeholder="Type or use voice..."
          value={input}
          onChangeText={setInput}
          multiline
          maxLength={500}
        />
        <View style={styles.buttonRow}>
          <TouchableOpacity
            style={[styles.voiceButton, isRecording && styles.voiceButtonActive]}
            onPress={handleVoiceInput}
            disabled={isProcessing}
          >
            <Ionicons name="mic" size={24} color="white" />
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.sendButton, isProcessing && styles.sendButtonDisabled]}
            onPress={handleSubmit}
            disabled={isProcessing || !input.trim()}
          >
            {isProcessing ? (
              <ActivityIndicator color="white" />
            ) : (
              <Ionicons name="send" size={24} color="white" />
            )}
          </TouchableOpacity>
        </View>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  content: {
    flex: 1,
    padding: 16,
  },
  card: {
    marginBottom: 16,
    elevation: 2,
  },
  instructionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 8,
    color: '#333',
  },
  instruction: {
    fontSize: 14,
    color: '#666',
    marginBottom: 4,
  },
  responseCard: {
    marginBottom: 16,
    backgroundColor: '#e3f2fd',
    elevation: 2,
  },
  responseTitle: {
    fontSize: 14,
    fontWeight: 'bold',
    marginBottom: 8,
    color: '#1976d2',
  },
  responseText: {
    fontSize: 15,
    color: '#333',
    lineHeight: 22,
  },
  toggleButton: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 12,
    backgroundColor: 'white',
    borderRadius: 8,
    marginBottom: 16,
    elevation: 1,
  },
  toggleButtonText: {
    fontSize: 16,
    color: '#6200ee',
    fontWeight: '600',
  },
  memoriesContainer: {
    marginBottom: 16,
  },
  memoryCard: {
    marginBottom: 12,
    elevation: 2,
  },
  memoryHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 4,
  },
  memoryItem: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
    flex: 1,
  },
  memoryLocation: {
    fontSize: 16,
    color: '#666',
    marginBottom: 4,
  },
  memoryTime: {
    fontSize: 12,
    color: '#999',
    marginBottom: 8,
  },
  navigateButton: {
    marginTop: 8,
  },
  emptyText: {
    fontSize: 14,
    color: '#999',
    textAlign: 'center',
  },
  clearButton: {
    marginBottom: 16,
  },
  inputContainer: {
    padding: 16,
    backgroundColor: 'white',
    borderTopWidth: 1,
    borderTopColor: '#e0e0e0',
  },
  input: {
    backgroundColor: '#f5f5f5',
    borderRadius: 24,
    padding: 12,
    paddingHorizontal: 16,
    fontSize: 16,
    maxHeight: 100,
    marginBottom: 8,
  },
  buttonRow: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: 8,
  },
  voiceButton: {
    backgroundColor: '#03a9f4',
    borderRadius: 28,
    width: 56,
    height: 56,
    justifyContent: 'center',
    alignItems: 'center',
    elevation: 4,
  },
  voiceButtonActive: {
    backgroundColor: '#f44336',
  },
  sendButton: {
    backgroundColor: '#6200ee',
    borderRadius: 28,
    width: 56,
    height: 56,
    justifyContent: 'center',
    alignItems: 'center',
    elevation: 4,
  },
  sendButtonDisabled: {
    backgroundColor: '#ccc',
  },
});

export default LostToFoundScreen;
