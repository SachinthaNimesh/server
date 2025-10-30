import React, { useEffect } from 'react';
import { StatusBar } from 'expo-status-bar';
import { NavigationContainer } from '@react-navigation/native';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { Provider as PaperProvider } from 'react-native-paper';
import { Ionicons } from '@expo/vector-icons';

// Context Providers
import { AppProvider } from './src/context/AppContext';

// Screens
import LostToFoundScreen from './src/screens/LostToFoundScreen';
import SoundSanctuaryScreen from './src/screens/SoundSanctuaryScreen';
import SettingsScreen from './src/screens/SettingsScreen';

const Tab = createBottomTabNavigator();

/**
 * Main App Component
 * Sets up navigation, providers, and app structure
 */
export default function App() {
  return (
    <PaperProvider>
      <AppProvider>
        <NavigationContainer>
          <StatusBar style="auto" />
          <Tab.Navigator
            screenOptions={({ route }) => ({
              tabBarIcon: ({ focused, color, size }) => {
                let iconName;

                if (route.name === 'LostToFound') {
                  iconName = focused ? 'search' : 'search-outline';
                } else if (route.name === 'SoundSanctuary') {
                  iconName = focused ? 'volume-high' : 'volume-high-outline';
                } else if (route.name === 'Settings') {
                  iconName = focused ? 'settings' : 'settings-outline';
                }

                return <Ionicons name={iconName} size={size} color={color} />;
              },
              tabBarActiveTintColor: '#6200ee',
              tabBarInactiveTintColor: 'gray',
              headerStyle: {
                backgroundColor: '#6200ee',
              },
              headerTintColor: '#fff',
              headerTitleStyle: {
                fontWeight: 'bold',
              },
            })}
          >
            <Tab.Screen 
              name="LostToFound" 
              component={LostToFoundScreen}
              options={{ title: 'Memory Assistant' }}
            />
            <Tab.Screen 
              name="SoundSanctuary" 
              component={SoundSanctuaryScreen}
              options={{ title: 'Sound Sanctuary' }}
            />
            <Tab.Screen 
              name="Settings" 
              component={SettingsScreen}
              options={{ title: 'Settings' }}
            />
          </Tab.Navigator>
        </NavigationContainer>
      </AppProvider>
    </PaperProvider>
  );
}
