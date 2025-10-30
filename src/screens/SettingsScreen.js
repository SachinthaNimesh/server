import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  Switch,
  Alert,
} from 'react-native';
import { Card, List, Button, Divider, RadioButton } from 'react-native-paper';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { useAppContext } from '../context/AppContext';

/**
 * Settings Screen
 * Allows users to configure app preferences and behavior
 */
const SettingsScreen = () => {
  const { settings, updateSettings, memories, clearConversationHistory } = useAppContext();
  const [localSettings, setLocalSettings] = useState(settings);

  /**
   * Handle priority behavior change
   */
  const handlePriorityBehaviorChange = async (value) => {
    const newSettings = { ...localSettings, priorityBehavior: value };
    setLocalSettings(newSettings);
    await updateSettings({ priorityBehavior: value });
  };

  /**
   * Handle max memories change
   */
  const handleMaxMemoriesChange = async (value) => {
    const newSettings = { ...localSettings, maxMemories: value };
    setLocalSettings(newSettings);
    await updateSettings({ maxMemories: value });
  };

  /**
   * Handle max conversation history change
   */
  const handleMaxConversationChange = async (value) => {
    const newSettings = { ...localSettings, maxConversationHistory: value };
    setLocalSettings(newSettings);
    await updateSettings({ maxConversationHistory: value });
  };

  /**
   * Clear all memories
   */
  const handleClearAllMemories = () => {
    Alert.alert(
      'Clear All Memories',
      'Are you sure you want to delete all stored memories? This action cannot be undone.',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Clear All',
          style: 'destructive',
          onPress: async () => {
            try {
              await AsyncStorage.setItem('memories', JSON.stringify([]));
              Alert.alert('Success', 'All memories have been cleared');
            } catch (error) {
              Alert.alert('Error', 'Failed to clear memories');
            }
          },
        },
      ]
    );
  };

  /**
   * Clear conversation history
   */
  const handleClearConversation = () => {
    Alert.alert(
      'Clear Conversation History',
      'Are you sure you want to clear the conversation history?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Clear',
          style: 'destructive',
          onPress: async () => {
            await clearConversationHistory();
            Alert.alert('Success', 'Conversation history cleared');
          },
        },
      ]
    );
  };

  /**
   * Reset all settings to default
   */
  const handleResetSettings = () => {
    Alert.alert(
      'Reset Settings',
      'Are you sure you want to reset all settings to default values?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Reset',
          style: 'destructive',
          onPress: async () => {
            const defaultSettings = {
              soundThreshold: 70,
              selectedSound: 'rain',
              priorityBehavior: 'pause',
              maxMemories: 30,
              maxConversationHistory: 10,
            };
            setLocalSettings(defaultSettings);
            await updateSettings(defaultSettings);
            Alert.alert('Success', 'Settings have been reset to defaults');
          },
        },
      ]
    );
  };

  /**
   * Show app information
   */
  const showAppInfo = () => {
    Alert.alert(
      'ADHD Support App',
      'Version 1.0.0\n\n' +
      'An app designed to help people with ADHD manage sensory overload and forgetfulness.\n\n' +
      'Features:\n' +
      '• LostToFound: AI-powered memory assistant\n' +
      '• SoundSanctuary: Noise-calming system\n\n' +
      'Built with React Native and Expo',
      [{ text: 'OK' }]
    );
  };

  /**
   * Show LLM configuration help
   */
  const showLLMHelp = () => {
    Alert.alert(
      'LLM Configuration',
      'To use the AI features, you need to configure an LLM API key:\n\n' +
      '1. Create a .env file in the project root\n' +
      '2. Copy the contents from .env.example\n' +
      '3. Add your API key for OpenAI, Gemini, or Hugging Face\n' +
      '4. Restart the app\n\n' +
      'Supported providers:\n' +
      '• OpenAI (gpt-3.5-turbo)\n' +
      '• Google Gemini (gemini-1.5-flash)\n' +
      '• Hugging Face (free inference API)',
      [{ text: 'OK' }]
    );
  };

  return (
    <ScrollView style={styles.container}>
      <View style={styles.content}>
        {/* Priority Management */}
        <Card style={styles.card}>
          <Card.Content>
            <Text style={styles.sectionTitle}>Priority Management</Text>
            <Text style={styles.sectionDescription}>
              Choose how SoundSanctuary behaves when you use LostToFound voice commands
            </Text>
          </Card.Content>
          <Divider />
          <RadioButton.Group
            onValueChange={handlePriorityBehaviorChange}
            value={localSettings.priorityBehavior}
          >
            <List.Item
              title="Auto-pause"
              description="Pause sounds completely during voice commands"
              left={() => <RadioButton value="pause" />}
            />
            <List.Item
              title="Duck volume"
              description="Lower volume during voice commands"
              left={() => <RadioButton value="duck" />}
            />
            <List.Item
              title="Ignore"
              description="Keep playing at normal volume"
              left={() => <RadioButton value="ignore" />}
            />
          </RadioButton.Group>
        </Card>

        {/* Memory Settings */}
        <Card style={styles.card}>
          <Card.Content>
            <Text style={styles.sectionTitle}>Memory Settings</Text>
            <Text style={styles.sectionDescription}>
              Configure how LostToFound stores and manages your memories
            </Text>
          </Card.Content>
          <Divider />
          <List.Item
            title="Maximum Stored Memories"
            description={`Current: ${localSettings.maxMemories} entries`}
            left={(props) => <List.Icon {...props} icon="database" />}
          />
          <View style={styles.buttonRow}>
            <Button
              mode="outlined"
              onPress={() => handleMaxMemoriesChange(20)}
              style={styles.optionButton}
            >
              20
            </Button>
            <Button
              mode="outlined"
              onPress={() => handleMaxMemoriesChange(30)}
              style={styles.optionButton}
            >
              30
            </Button>
            <Button
              mode="outlined"
              onPress={() => handleMaxMemoriesChange(50)}
              style={styles.optionButton}
            >
              50
            </Button>
          </View>
          <Divider style={styles.divider} />
          <List.Item
            title="Conversation History Length"
            description={`Current: ${localSettings.maxConversationHistory} interactions`}
            left={(props) => <List.Icon {...props} icon="message-text" />}
          />
          <View style={styles.buttonRow}>
            <Button
              mode="outlined"
              onPress={() => handleMaxConversationChange(5)}
              style={styles.optionButton}
            >
              5
            </Button>
            <Button
              mode="outlined"
              onPress={() => handleMaxConversationChange(10)}
              style={styles.optionButton}
            >
              10
            </Button>
            <Button
              mode="outlined"
              onPress={() => handleMaxConversationChange(20)}
              style={styles.optionButton}
            >
              20
            </Button>
          </View>
        </Card>

        {/* Data Management */}
        <Card style={styles.card}>
          <Card.Content>
            <Text style={styles.sectionTitle}>Data Management</Text>
          </Card.Content>
          <Divider />
          <List.Item
            title="Stored Memories"
            description={`${memories.length} items currently stored`}
            left={(props) => <List.Icon {...props} icon="folder" />}
          />
          <View style={styles.buttonContainer}>
            <Button
              mode="outlined"
              onPress={handleClearAllMemories}
              style={styles.actionButton}
              icon="delete-sweep"
            >
              Clear All Memories
            </Button>
            <Button
              mode="outlined"
              onPress={handleClearConversation}
              style={styles.actionButton}
              icon="message-off"
            >
              Clear Conversation History
            </Button>
          </View>
        </Card>

        {/* LLM Configuration */}
        <Card style={styles.card}>
          <Card.Content>
            <Text style={styles.sectionTitle}>AI Configuration</Text>
            <Text style={styles.sectionDescription}>
              Configure your LLM provider for AI-powered features
            </Text>
          </Card.Content>
          <Divider />
          <List.Item
            title="LLM Setup Guide"
            description="Learn how to configure your AI provider"
            left={(props) => <List.Icon {...props} icon="help-circle" />}
            onPress={showLLMHelp}
          />
        </Card>

        {/* About */}
        <Card style={styles.card}>
          <Card.Content>
            <Text style={styles.sectionTitle}>About</Text>
          </Card.Content>
          <Divider />
          <List.Item
            title="App Information"
            description="Version 1.0.0"
            left={(props) => <List.Icon {...props} icon="information" />}
            onPress={showAppInfo}
          />
          <List.Item
            title="Reset Settings"
            description="Restore default settings"
            left={(props) => <List.Icon {...props} icon="restore" />}
            onPress={handleResetSettings}
          />
        </Card>

        {/* Development Note */}
        <Card style={styles.infoCard}>
          <Card.Content>
            <Text style={styles.infoTitle}>📝 Note for Development</Text>
            <Text style={styles.infoText}>
              This app requires proper LLM API configuration to work fully. 
              Make sure to set up your .env file with valid API keys before testing AI features.
            </Text>
          </Card.Content>
        </Card>
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  content: {
    padding: 16,
    paddingBottom: 32,
  },
  card: {
    marginBottom: 16,
    elevation: 2,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 4,
  },
  sectionDescription: {
    fontSize: 14,
    color: '#666',
    marginTop: 4,
  },
  divider: {
    marginVertical: 8,
  },
  buttonRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    paddingHorizontal: 16,
    paddingVertical: 8,
  },
  optionButton: {
    flex: 1,
    marginHorizontal: 4,
  },
  buttonContainer: {
    padding: 16,
  },
  actionButton: {
    marginBottom: 12,
  },
  infoCard: {
    backgroundColor: '#fff3e0',
    marginBottom: 16,
    elevation: 2,
  },
  infoTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 8,
    color: '#e65100',
  },
  infoText: {
    fontSize: 14,
    color: '#333',
    lineHeight: 20,
  },
});

export default SettingsScreen;
