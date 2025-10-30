import React, { useState, useEffect, useRef } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  Alert,
  Switch,
} from 'react-native';
import { Button, Card, Slider, Chip } from 'react-native-paper';
import { Audio } from 'expo-av';
import { Ionicons } from '@expo/vector-icons';
import { useAppContext } from '../context/AppContext';

/**
 * SoundSanctuary Screen - Noise-calming system
 * Monitors ambient noise and plays soothing sounds when threshold is exceeded
 */
const SoundSanctuaryScreen = () => {
  const {
    isPlaying,
    currentSound,
    volume,
    playSoothingSound,
    pauseSound,
    resumeSound,
    stopSound,
    setPlaybackVolume,
    settings,
    updateSettings,
  } = useAppContext();

  const [isMonitoring, setIsMonitoring] = useState(false);
  const [currentNoiseLevel, setCurrentNoiseLevel] = useState(0);
  const [recording, setRecording] = useState(null);
  const [hasAudioPermission, setHasAudioPermission] = useState(false);
  
  const monitoringInterval = useRef(null);

  // Available soothing sounds
  const availableSounds = [
    { id: 'rain', name: 'Rain', icon: '🌧️' },
    { id: 'ocean', name: 'Ocean Waves', icon: '🌊' },
    { id: 'forest', name: 'Forest', icon: '🌲' },
    { id: 'white-noise', name: 'White Noise', icon: '⚪' },
  ];

  useEffect(() => {
    requestAudioPermissions();
    return () => {
      stopMonitoring();
    };
  }, []);

  /**
   * Request audio recording permission
   */
  const requestAudioPermissions = async () => {
    try {
      const { status } = await Audio.requestPermissionsAsync();
      setHasAudioPermission(status === 'granted');
      
      if (status !== 'granted') {
        Alert.alert(
          'Permission Required',
          'Microphone access is needed to monitor ambient noise levels.',
          [{ text: 'OK' }]
        );
      }

      // Configure audio mode
      await Audio.setAudioModeAsync({
        allowsRecordingIOS: true,
        playsInSilentModeIOS: true,
        staysActiveInBackground: false,
      });
    } catch (error) {
      console.error('Error requesting audio permissions:', error);
    }
  };

  /**
   * Start monitoring ambient noise
   */
  const startMonitoring = async () => {
    if (!hasAudioPermission) {
      await requestAudioPermissions();
      return;
    }

    try {
      setIsMonitoring(true);
      
      // Start periodic monitoring (every 5 seconds for battery efficiency)
      monitoringInterval.current = setInterval(async () => {
        await measureNoiseLevel();
      }, 5000);

      Alert.alert(
        'Monitoring Started',
        'The app will now monitor ambient noise and play soothing sounds when needed.',
        [{ text: 'OK' }]
      );
    } catch (error) {
      console.error('Error starting monitoring:', error);
      Alert.alert('Error', 'Failed to start monitoring');
      setIsMonitoring(false);
    }
  };

  /**
   * Stop monitoring ambient noise
   */
  const stopMonitoring = async () => {
    setIsMonitoring(false);
    
    if (monitoringInterval.current) {
      clearInterval(monitoringInterval.current);
      monitoringInterval.current = null;
    }

    if (recording) {
      try {
        await recording.stopAndUnloadAsync();
        setRecording(null);
      } catch (error) {
        console.error('Error stopping recording:', error);
      }
    }

    setCurrentNoiseLevel(0);
  };

  /**
   * Measure ambient noise level
   * Note: This is a simplified implementation
   * In a real app, you would use expo-av's Audio.Recording with metering
   */
  const measureNoiseLevel = async () => {
    try {
      // Create a short recording to measure noise
      const { recording: newRecording } = await Audio.Recording.createAsync(
        Audio.RecordingOptionsPresets.HIGH_QUALITY
      );
      
      setRecording(newRecording);

      // Record for 1 second
      await new Promise(resolve => setTimeout(resolve, 1000));

      // Stop recording
      await newRecording.stopAndUnloadAsync();
      
      // In a real implementation, you would analyze the recording's metering data
      // For this demo, we simulate a noise level (0-100)
      const simulatedNoiseLevel = Math.floor(Math.random() * 100);
      setCurrentNoiseLevel(simulatedNoiseLevel);

      // Check if noise exceeds threshold
      if (simulatedNoiseLevel > settings.soundThreshold && !isPlaying) {
        await playSoothingSound(settings.selectedSound);
      } else if (simulatedNoiseLevel <= settings.soundThreshold && isPlaying) {
        // Optional: Auto-stop when noise level drops
        // await stopSound();
      }

      setRecording(null);
    } catch (error) {
      console.error('Error measuring noise level:', error);
      // Continue monitoring despite errors
    }
  };

  /**
   * Handle sound selection
   */
  const handleSoundSelect = async (soundId) => {
    await updateSettings({ selectedSound: soundId });
    if (isPlaying) {
      await playSoothingSound(soundId);
    }
  };

  /**
   * Handle volume change
   */
  const handleVolumeChange = async (value) => {
    await setPlaybackVolume(value);
  };

  /**
   * Handle threshold change
   */
  const handleThresholdChange = async (value) => {
    await updateSettings({ soundThreshold: Math.round(value) });
  };

  /**
   * Toggle playback
   */
  const handleTogglePlayback = async () => {
    if (isPlaying) {
      await pauseSound();
    } else {
      await playSoothingSound(settings.selectedSound);
    }
  };

  /**
   * Show threshold tutorial
   */
  const showThresholdTutorial = () => {
    Alert.alert(
      'Finding Your Threshold',
      'To find the right noise threshold:\n\n' +
      '1. Enable monitoring\n' +
      '2. Observe the current noise level in different environments\n' +
      '3. Adjust the threshold slider based on when you want sounds to play\n' +
      '4. Typical values:\n' +
      '   • Quiet room: 30-40\n' +
      '   • Normal conversation: 50-60\n' +
      '   • Busy street: 70-80\n' +
      '   • Loud environment: 80-90',
      [{ text: 'Got it!' }]
    );
  };

  return (
    <ScrollView style={styles.container}>
      <View style={styles.content}>
        {/* Status Card */}
        <Card style={styles.card}>
          <Card.Content>
            <View style={styles.statusRow}>
              <Text style={styles.statusLabel}>Monitoring:</Text>
              <Chip
                mode="flat"
                style={[
                  styles.statusChip,
                  isMonitoring && styles.statusChipActive,
                ]}
              >
                {isMonitoring ? 'Active' : 'Inactive'}
              </Chip>
            </View>
            <View style={styles.statusRow}>
              <Text style={styles.statusLabel}>Playback:</Text>
              <Chip
                mode="flat"
                style={[
                  styles.statusChip,
                  isPlaying && styles.statusChipActive,
                ]}
              >
                {isPlaying ? 'Playing' : 'Stopped'}
              </Chip>
            </View>
          </Card.Content>
        </Card>

        {/* Current Noise Level */}
        <Card style={styles.card}>
          <Card.Content>
            <View style={styles.noiseLevelContainer}>
              <Text style={styles.sectionTitle}>Current Noise Level</Text>
              <View style={styles.noiseMeter}>
                <Text style={styles.noiseLevelText}>{currentNoiseLevel}</Text>
                <Text style={styles.noiseLevelUnit}>dB</Text>
              </View>
              <View style={styles.noiseLevelBar}>
                <View
                  style={[
                    styles.noiseLevelFill,
                    {
                      width: `${currentNoiseLevel}%`,
                      backgroundColor:
                        currentNoiseLevel > settings.soundThreshold
                          ? '#f44336'
                          : '#4caf50',
                    },
                  ]}
                />
              </View>
            </View>
          </Card.Content>
        </Card>

        {/* Noise Threshold */}
        <Card style={styles.card}>
          <Card.Content>
            <View style={styles.sliderHeader}>
              <Text style={styles.sectionTitle}>Noise Threshold</Text>
              <Button mode="text" onPress={showThresholdTutorial}>
                Tutorial
              </Button>
            </View>
            <Text style={styles.sliderValue}>{settings.soundThreshold} dB</Text>
            <Slider
              style={styles.slider}
              minimumValue={0}
              maximumValue={100}
              value={settings.soundThreshold}
              onValueChange={handleThresholdChange}
              minimumTrackTintColor="#6200ee"
              maximumTrackTintColor="#e0e0e0"
              thumbTintColor="#6200ee"
            />
            <Text style={styles.sliderDescription}>
              Soothing sounds will play when noise exceeds this level
            </Text>
          </Card.Content>
        </Card>

        {/* Sound Selection */}
        <Card style={styles.card}>
          <Card.Content>
            <Text style={styles.sectionTitle}>Choose Your Sound</Text>
            <View style={styles.soundsGrid}>
              {availableSounds.map((sound) => (
                <Chip
                  key={sound.id}
                  mode={settings.selectedSound === sound.id ? 'flat' : 'outlined'}
                  selected={settings.selectedSound === sound.id}
                  onPress={() => handleSoundSelect(sound.id)}
                  style={[
                    styles.soundChip,
                    settings.selectedSound === sound.id && styles.soundChipSelected,
                  ]}
                >
                  {sound.icon} {sound.name}
                </Chip>
              ))}
            </View>
          </Card.Content>
        </Card>

        {/* Volume Control */}
        <Card style={styles.card}>
          <Card.Content>
            <Text style={styles.sectionTitle}>Volume</Text>
            <Text style={styles.sliderValue}>{Math.round(volume * 100)}%</Text>
            <Slider
              style={styles.slider}
              minimumValue={0}
              maximumValue={1}
              value={volume}
              onValueChange={handleVolumeChange}
              minimumTrackTintColor="#6200ee"
              maximumTrackTintColor="#e0e0e0"
              thumbTintColor="#6200ee"
            />
          </Card.Content>
        </Card>

        {/* Control Buttons */}
        <View style={styles.controlButtons}>
          <Button
            mode="contained"
            onPress={isMonitoring ? stopMonitoring : startMonitoring}
            style={[
              styles.controlButton,
              isMonitoring && styles.controlButtonActive,
            ]}
            icon={isMonitoring ? 'stop' : 'play'}
          >
            {isMonitoring ? 'Stop Monitoring' : 'Start Monitoring'}
          </Button>

          <Button
            mode="contained"
            onPress={handleTogglePlayback}
            style={styles.controlButton}
            icon={isPlaying ? 'pause' : 'play'}
          >
            {isPlaying ? 'Pause Sound' : 'Play Sound'}
          </Button>

          {isPlaying && (
            <Button
              mode="outlined"
              onPress={stopSound}
              style={styles.controlButton}
              icon="stop"
            >
              Stop Sound
            </Button>
          )}
        </View>

        {/* Information Card */}
        <Card style={styles.infoCard}>
          <Card.Content>
            <Text style={styles.infoTitle}>💡 How it works</Text>
            <Text style={styles.infoText}>
              • Enable monitoring to automatically detect ambient noise
            </Text>
            <Text style={styles.infoText}>
              • When noise exceeds your threshold, soothing sounds play automatically
            </Text>
            <Text style={styles.infoText}>
              • Adjust volume and choose your preferred sound
            </Text>
            <Text style={styles.infoText}>
              • The app checks noise levels every 5 seconds to save battery
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
  },
  card: {
    marginBottom: 16,
    elevation: 2,
  },
  statusRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  statusLabel: {
    fontSize: 16,
    color: '#333',
    fontWeight: '600',
  },
  statusChip: {
    backgroundColor: '#e0e0e0',
  },
  statusChipActive: {
    backgroundColor: '#c8e6c9',
  },
  noiseLevelContainer: {
    alignItems: 'center',
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 12,
  },
  noiseMeter: {
    flexDirection: 'row',
    alignItems: 'baseline',
    marginVertical: 16,
  },
  noiseLevelText: {
    fontSize: 48,
    fontWeight: 'bold',
    color: '#6200ee',
  },
  noiseLevelUnit: {
    fontSize: 24,
    color: '#999',
    marginLeft: 8,
  },
  noiseLevelBar: {
    width: '100%',
    height: 12,
    backgroundColor: '#e0e0e0',
    borderRadius: 6,
    overflow: 'hidden',
  },
  noiseLevelFill: {
    height: '100%',
    borderRadius: 6,
  },
  sliderHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  sliderValue: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#6200ee',
    textAlign: 'center',
    marginBottom: 8,
  },
  slider: {
    width: '100%',
    height: 40,
  },
  sliderDescription: {
    fontSize: 12,
    color: '#666',
    textAlign: 'center',
    marginTop: 4,
  },
  soundsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  soundChip: {
    marginRight: 8,
    marginBottom: 8,
  },
  soundChipSelected: {
    backgroundColor: '#e1bee7',
  },
  controlButtons: {
    marginBottom: 16,
  },
  controlButton: {
    marginBottom: 12,
  },
  controlButtonActive: {
    backgroundColor: '#f44336',
  },
  infoCard: {
    backgroundColor: '#e3f2fd',
    marginBottom: 16,
    elevation: 2,
  },
  infoTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 8,
    color: '#1976d2',
  },
  infoText: {
    fontSize: 14,
    color: '#333',
    marginBottom: 4,
    lineHeight: 20,
  },
});

export default SoundSanctuaryScreen;
