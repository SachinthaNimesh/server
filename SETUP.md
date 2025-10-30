# ADHD Support App - Setup Guide

## Project Overview

This is a React Native mobile application built with Expo framework, designed to help people with ADHD manage sensory overload and forgetfulness. The app includes two main features:

1. **LostToFound** - AI-powered memory assistant using voice or text
2. **SoundSanctuary** - Noise-calming system with automatic ambient sound playback

## Prerequisites

- Node.js (v16 or higher)
- npm or yarn
- Expo Go app on your mobile device (for testing)
- An LLM API key (OpenAI, Google Gemini, or Hugging Face)

## Installation Steps

### 1. Install Dependencies

```bash
npm install
```

This will install all required packages including:
- Expo SDK 51
- React Navigation
- Expo Audio (expo-av)
- Expo Location (expo-location)
- Expo Speech (expo-speech)
- AsyncStorage
- React Native Paper
- Axios

### 2. Configure Environment Variables

Create a `.env` file in the project root by copying the example:

```bash
cp .env.example .env
```

Edit `.env` and add your LLM API key. Choose ONE of the following providers:

#### Option A: OpenAI (Recommended for best results)

```env
OPENAI_API_KEY=sk-your_actual_openai_key_here
OPENAI_MODEL=gpt-3.5-turbo
LLM_PROVIDER=openai
```

Get your API key from: https://platform.openai.com/api-keys

#### Option B: Google Gemini (Free tier available)

```env
GEMINI_API_KEY=your_gemini_api_key_here
GEMINI_MODEL=gemini-1.5-flash
LLM_PROVIDER=gemini
```

Get your API key from: https://makersuite.google.com/app/apikey

#### Option C: Hugging Face (Free inference API)

```env
HUGGINGFACE_API_KEY=your_huggingface_token_here
HUGGINGFACE_MODEL=mistralai/Mistral-7B-Instruct-v0.1
LLM_PROVIDER=huggingface
```

Get your token from: https://huggingface.co/settings/tokens

### 3. Add Sound Assets (Optional)

For the SoundSanctuary feature to work with real audio files, you need to add sound files to `assets/sounds/`:

```
assets/sounds/
  ├── rain.mp3
  ├── ocean.mp3
  ├── forest.mp3
  └── white-noise.mp3
```

You can download free ambient sounds from:
- https://freesound.org/
- https://pixabay.com/sound-effects/

**Note:** The app is currently configured with placeholder audio. To use real sounds, update the audio source in `src/context/AppContext.js` around line 208.

### 4. Start the Development Server

```bash
npm start
```

This will start the Expo development server and display a QR code.

### 5. Run on Your Device

#### Using Expo Go (Easiest Method)

1. Install Expo Go from your app store:
   - iOS: https://apps.apple.com/app/expo-go/id982107779
   - Android: https://play.google.com/store/apps/details?id=host.exp.exponent

2. Scan the QR code displayed in your terminal:
   - iOS: Use the Camera app
   - Android: Use the Expo Go app

#### Using iOS Simulator (Mac only)

```bash
npm run ios
```

#### Using Android Emulator

```bash
npm run android
```

Make sure you have Android Studio installed with an emulator configured.

## Project Structure

```
adhd-support-app/
├── App.js                          # Main app entry point
├── app.json                        # Expo configuration
├── babel.config.js                 # Babel configuration
├── package.json                    # Dependencies
├── .env.example                    # Environment variables template
├── assets/                         # Static assets
│   ├── sounds/                     # Soothing sounds for SoundSanctuary
│   └── images/                     # App icons and images
└── src/
    ├── screens/                    # Screen components
    │   ├── LostToFoundScreen.js    # Memory assistant screen
    │   ├── SoundSanctuaryScreen.js # Noise-calming screen
    │   └── SettingsScreen.js       # Settings screen
    ├── context/                    # State management
    │   └── AppContext.js           # Global app context
    └── utils/                      # Utility functions
        └── llmService.js           # LLM API integration
```

## Key Features Explained

### LostToFound (Memory Assistant)

- **Voice/Text Input**: Users can record where they placed items
- **AI Processing**: LLM extracts item name and location
- **Location Tracking**: Stores GPS coordinates automatically
- **Smart Queries**: Ask where items are using natural language
- **Navigation**: Opens Google Maps to stored locations
- **Rolling History**: Keeps last 30 entries (configurable)

### SoundSanctuary (Noise Calming)

- **Ambient Monitoring**: Periodically checks noise levels
- **Automatic Playback**: Plays soothing sounds when threshold exceeded
- **Sound Selection**: Choose from rain, ocean, forest, or white noise
- **Volume Control**: Adjust playback volume
- **Battery Efficient**: Checks every 5 seconds instead of continuous monitoring

### Priority Management

The app intelligently manages conflicts between features:
- **Pause**: Stops sound playback during voice commands
- **Duck**: Reduces volume during voice commands
- **Ignore**: Continues playback normally

## Troubleshooting

### Issue: "LLM API Error"

- Verify your API key is correct in `.env`
- Check that you have credits/quota remaining
- Test API key directly using curl or Postman
- The app will fallback to local search if LLM fails

### Issue: "Location permission denied"

- Go to device Settings → Privacy → Location
- Enable location access for Expo Go
- Restart the app

### Issue: "Microphone permission denied"

- Go to device Settings → Privacy → Microphone
- Enable microphone access for Expo Go
- Restart the app

### Issue: "Sound doesn't play"

- Ensure device volume is not muted
- Add actual sound files to `assets/sounds/`
- Update audio source paths in AppContext.js

### Issue: "App crashes on startup"

- Clear Expo cache: `expo start -c`
- Reinstall dependencies: `rm -rf node_modules && npm install`
- Check for Node.js version compatibility

## Testing Checklist

See [TESTING.md](./TESTING.md) for a comprehensive testing guide.

## Next Steps

1. Add actual ambient sound files
2. Implement proper voice recognition (using @react-native-voice/voice)
3. Add proper audio metering for noise level detection
4. Implement user authentication (optional)
5. Add data export/import functionality
6. Optimize battery usage
7. Add unit and integration tests

## Converting to Wearable Device

See [WEARABLE_GUIDE.md](./WEARABLE_GUIDE.md) for suggestions on converting this mobile prototype into a wearable device.

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review Expo documentation: https://docs.expo.dev/
3. Check React Native documentation: https://reactnative.dev/

## License

This project is for demonstration purposes.
