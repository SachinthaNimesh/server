import React, { createContext, useContext, useState, useEffect } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { Audio } from 'expo-av';

/**
 * App Context for managing global state
 * Handles memory storage, audio state, and settings
 */
const AppContext = createContext();

export const useAppContext = () => {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error('useAppContext must be used within AppProvider');
  }
  return context;
};

export const AppProvider = ({ children }) => {
  // Memory state for LostToFound
  const [memories, setMemories] = useState([]);
  const [conversationHistory, setConversationHistory] = useState([]);

  // Audio state for SoundSanctuary
  const [isPlaying, setIsPlaying] = useState(false);
  const [soundObject, setSoundObject] = useState(null);
  const [currentSound, setCurrentSound] = useState('rain');
  const [volume, setVolume] = useState(0.5);

  // Settings state
  const [settings, setSettings] = useState({
    soundThreshold: 70,
    selectedSound: 'rain',
    priorityBehavior: 'pause', // 'pause', 'duck', or 'ignore'
    maxMemories: 30,
    maxConversationHistory: 10,
  });

  // Load data from AsyncStorage on mount
  useEffect(() => {
    loadMemories();
    loadSettings();
    loadConversationHistory();
  }, []);

  /**
   * Load memories from AsyncStorage
   */
  const loadMemories = async () => {
    try {
      const stored = await AsyncStorage.getItem('memories');
      if (stored) {
        setMemories(JSON.parse(stored));
      }
    } catch (error) {
      console.error('Error loading memories:', error);
    }
  };

  /**
   * Save memories to AsyncStorage
   */
  const saveMemories = async (newMemories) => {
    try {
      await AsyncStorage.setItem('memories', JSON.stringify(newMemories));
      setMemories(newMemories);
    } catch (error) {
      console.error('Error saving memories:', error);
    }
  };

  /**
   * Add a new memory entry
   */
  const addMemory = async (memory) => {
    const newMemory = {
      id: Date.now().toString(),
      timestamp: new Date().toISOString(),
      ...memory,
    };

    let updatedMemories = [newMemory, ...memories];
    
    // Keep only the last N memories (rolling history)
    if (updatedMemories.length > settings.maxMemories) {
      updatedMemories = updatedMemories.slice(0, settings.maxMemories);
    }

    await saveMemories(updatedMemories);
    return newMemory;
  };

  /**
   * Delete a memory by ID
   */
  const deleteMemory = async (id) => {
    const updatedMemories = memories.filter(m => m.id !== id);
    await saveMemories(updatedMemories);
  };

  /**
   * Update a memory by ID
   */
  const updateMemory = async (id, updates) => {
    const updatedMemories = memories.map(m => 
      m.id === id ? { ...m, ...updates } : m
    );
    await saveMemories(updatedMemories);
  };

  /**
   * Load conversation history from AsyncStorage
   */
  const loadConversationHistory = async () => {
    try {
      const stored = await AsyncStorage.getItem('conversationHistory');
      if (stored) {
        setConversationHistory(JSON.parse(stored));
      }
    } catch (error) {
      console.error('Error loading conversation history:', error);
    }
  };

  /**
   * Save conversation history to AsyncStorage
   */
  const saveConversationHistory = async (history) => {
    try {
      await AsyncStorage.setItem('conversationHistory', JSON.stringify(history));
      setConversationHistory(history);
    } catch (error) {
      console.error('Error saving conversation history:', error);
    }
  };

  /**
   * Add to conversation history
   */
  const addToConversationHistory = async (role, content) => {
    let updated = [...conversationHistory, { role, content, timestamp: Date.now() }];
    
    // Keep only last N conversations
    if (updated.length > settings.maxConversationHistory * 2) {
      updated = updated.slice(-settings.maxConversationHistory * 2);
    }

    await saveConversationHistory(updated);
  };

  /**
   * Clear conversation history
   */
  const clearConversationHistory = async () => {
    await saveConversationHistory([]);
  };

  /**
   * Load settings from AsyncStorage
   */
  const loadSettings = async () => {
    try {
      const stored = await AsyncStorage.getItem('settings');
      if (stored) {
        setSettings(JSON.parse(stored));
      }
    } catch (error) {
      console.error('Error loading settings:', error);
    }
  };

  /**
   * Save settings to AsyncStorage
   */
  const saveSettings = async (newSettings) => {
    try {
      await AsyncStorage.setItem('settings', JSON.stringify(newSettings));
      setSettings(newSettings);
    } catch (error) {
      console.error('Error saving settings:', error);
    }
  };

  /**
   * Update specific settings
   */
  const updateSettings = async (updates) => {
    const newSettings = { ...settings, ...updates };
    await saveSettings(newSettings);
  };

  /**
   * Play soothing sound
   */
  const playSoothingSound = async (soundName = currentSound) => {
    try {
      // Stop current sound if playing
      if (soundObject) {
        await soundObject.stopAsync();
        await soundObject.unloadAsync();
      }

      // Note: In a real app, you would have actual sound files
      // For this example, we're setting up the structure
      const { sound } = await Audio.Sound.createAsync(
        // This would be: require(`../../assets/sounds/${soundName}.mp3`),
        // For now, using a placeholder
        { uri: 'https://example.com/placeholder.mp3' },
        { shouldPlay: true, isLooping: true, volume }
      );

      setSoundObject(sound);
      setIsPlaying(true);
      setCurrentSound(soundName);
    } catch (error) {
      console.error('Error playing sound:', error);
    }
  };

  /**
   * Pause sound playback
   */
  const pauseSound = async () => {
    if (soundObject) {
      await soundObject.pauseAsync();
      setIsPlaying(false);
    }
  };

  /**
   * Resume sound playback
   */
  const resumeSound = async () => {
    if (soundObject) {
      await soundObject.playAsync();
      setIsPlaying(true);
    }
  };

  /**
   * Stop sound playback
   */
  const stopSound = async () => {
    if (soundObject) {
      await soundObject.stopAsync();
      await soundObject.unloadAsync();
      setSoundObject(null);
      setIsPlaying(false);
    }
  };

  /**
   * Set volume for sound playback
   */
  const setPlaybackVolume = async (newVolume) => {
    setVolume(newVolume);
    if (soundObject) {
      await soundObject.setVolumeAsync(newVolume);
    }
  };

  /**
   * Duck volume (reduce for voice commands)
   */
  const duckVolume = async () => {
    if (soundObject && isPlaying) {
      await soundObject.setVolumeAsync(volume * 0.3);
    }
  };

  /**
   * Restore volume after ducking
   */
  const restoreVolume = async () => {
    if (soundObject && isPlaying) {
      await soundObject.setVolumeAsync(volume);
    }
  };

  const value = {
    // Memory state and methods
    memories,
    addMemory,
    deleteMemory,
    updateMemory,
    conversationHistory,
    addToConversationHistory,
    clearConversationHistory,

    // Audio state and methods
    isPlaying,
    currentSound,
    volume,
    playSoothingSound,
    pauseSound,
    resumeSound,
    stopSound,
    setPlaybackVolume,
    duckVolume,
    restoreVolume,

    // Settings
    settings,
    updateSettings,
  };

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
};
