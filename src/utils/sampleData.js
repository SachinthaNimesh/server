/**
 * Sample data for testing and demonstration
 * This can be used to quickly populate the app with example memories
 */

export const sampleMemories = [
  {
    id: '1',
    item: 'car keys',
    location: 'kitchen table',
    coordinates: {
      latitude: 37.7749,
      longitude: -122.4194,
    },
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(), // 2 hours ago
  },
  {
    id: '2',
    item: 'wallet',
    location: 'bedroom drawer',
    coordinates: {
      latitude: 37.7750,
      longitude: -122.4195,
    },
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 4).toISOString(), // 4 hours ago
  },
  {
    id: '3',
    item: 'phone charger',
    location: 'desk',
    coordinates: {
      latitude: 37.7751,
      longitude: -122.4196,
    },
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 6).toISOString(), // 6 hours ago
  },
  {
    id: '4',
    item: 'sunglasses',
    location: 'car glove box',
    coordinates: {
      latitude: 37.7752,
      longitude: -122.4197,
    },
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(), // 1 day ago
  },
  {
    id: '5',
    item: 'laptop',
    location: 'living room coffee table',
    coordinates: {
      latitude: 37.7753,
      longitude: -122.4198,
    },
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 2).toISOString(), // 2 days ago
  },
  {
    id: '6',
    item: 'headphones',
    location: 'backpack',
    coordinates: {
      latitude: 37.7754,
      longitude: -122.4199,
    },
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 3).toISOString(), // 3 days ago
  },
  {
    id: '7',
    item: 'water bottle',
    location: 'gym bag',
    coordinates: {
      latitude: 37.7755,
      longitude: -122.4200,
    },
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 4).toISOString(), // 4 days ago
  },
  {
    id: '8',
    item: 'medication',
    location: 'bathroom cabinet',
    coordinates: {
      latitude: 37.7756,
      longitude: -122.4201,
    },
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 5).toISOString(), // 5 days ago
  },
  {
    id: '9',
    item: 'passport',
    location: 'bedroom safe',
    coordinates: {
      latitude: 37.7757,
      longitude: -122.4202,
    },
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 10).toISOString(), // 10 days ago
  },
  {
    id: '10',
    item: 'book',
    location: 'nightstand',
    coordinates: {
      latitude: 37.7758,
      longitude: -122.4203,
    },
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 15).toISOString(), // 15 days ago
  },
];

export const sampleConversationHistory = [
  {
    role: 'user',
    content: "I'm leaving my car keys on the kitchen table",
    timestamp: Date.now() - 1000 * 60 * 60 * 2,
  },
  {
    role: 'assistant',
    content: "Got it! I've remembered that your car keys are at kitchen table.",
    timestamp: Date.now() - 1000 * 60 * 60 * 2 + 1000,
  },
  {
    role: 'user',
    content: 'Where are my car keys?',
    timestamp: Date.now() - 1000 * 60 * 30,
  },
  {
    role: 'assistant',
    content: 'Your car keys are at kitchen table (stored on today at 10:30 AM)',
    timestamp: Date.now() - 1000 * 60 * 30 + 1000,
  },
];

export const testQueries = [
  'Where are my car keys?',
  'Where did I put my wallet?',
  "I can't find my phone charger",
  'Where are my sunglasses?',
  'Where is my laptop?',
  'Have you seen my headphones?',
  'Where did I leave my water bottle?',
  'Where is my medication?',
  'Where did I put my passport?',
  'Where is my book?',
];

export const testStorageStatements = [
  "I'm putting my car keys on the kitchen table",
  'My wallet is in the bedroom drawer',
  'Phone charger is on my desk',
  'Sunglasses are in the car glove box',
  'Laptop is on the living room coffee table',
  'Headphones are in my backpack',
  'Water bottle is in my gym bag',
  'Medication is in the bathroom cabinet',
  'Passport is in the bedroom safe',
  'Book is on the nightstand',
  'Remote control is on the couch',
  'Umbrella is by the front door',
  'Shoes are in the closet',
  'Jacket is hanging in the hallway',
  'Watch is on the dresser',
];

/**
 * Helper function to load sample data into the app
 * Use this for demo purposes or testing
 */
export async function loadSampleData(appContext) {
  try {
    // Load sample memories
    for (const memory of sampleMemories) {
      await appContext.addMemory({
        item: memory.item,
        location: memory.location,
        coordinates: memory.coordinates,
      });
    }

    // Load sample conversation history
    for (const message of sampleConversationHistory) {
      await appContext.addToConversationHistory(message.role, message.content);
    }

    console.log('Sample data loaded successfully');
    return true;
  } catch (error) {
    console.error('Error loading sample data:', error);
    return false;
  }
}

/**
 * Helper function to clear all data
 */
export async function clearAllData(appContext) {
  try {
    await appContext.clearConversationHistory();
    // Note: You would need to add a method to clear all memories
    console.log('All data cleared');
    return true;
  } catch (error) {
    console.error('Error clearing data:', error);
    return false;
  }
}
