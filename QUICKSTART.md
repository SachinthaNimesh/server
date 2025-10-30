# Quick Start Guide - ADHD Support App

Get the app running in **5 minutes** with this streamlined guide.

## 1. Install Dependencies (2 minutes)

```bash
npm install
```

This installs all required packages. Wait for it to complete.

## 2. Configure API Key (1 minute)

Choose the **easiest** option for you:

### Option A: OpenAI (Best Quality)
```bash
# Create .env file
echo 'OPENAI_API_KEY=sk-your-key-here' > .env
echo 'LLM_PROVIDER=openai' >> .env

# Get key from: https://platform.openai.com/api-keys
```

### Option B: Google Gemini (Free Tier)
```bash
# Create .env file
echo 'GEMINI_API_KEY=your-key-here' > .env
echo 'LLM_PROVIDER=gemini' >> .env

# Get key from: https://makersuite.google.com/app/apikey
```

### Option C: No API Key (Local Mode Only)
```bash
# Skip AI features, use local search
echo 'LLM_PROVIDER=local' > .env

# Note: This disables AI but app still works
```

## 3. Start the App (1 minute)

```bash
npm start
```

This opens Expo DevTools in your browser.

## 4. Run on Your Phone (1 minute)

### iOS:
1. Install **Expo Go** from App Store
2. Open Camera app
3. Scan the QR code
4. App opens in Expo Go

### Android:
1. Install **Expo Go** from Play Store
2. Open Expo Go app
3. Scan the QR code
4. App opens

## 5. Test the Features (5 minutes)

### Test LostToFound:
```
1. Open LostToFound tab
2. Type: "I'm leaving my keys on the kitchen table"
3. Press send
4. Wait for confirmation
5. Type: "Where are my keys?"
6. Get response!
```

### Test SoundSanctuary:
```
1. Open SoundSanctuary tab
2. Tap "Play Sound"
3. Adjust volume slider
4. Try different sounds
5. Tap "Stop Sound"
```

### Test Settings:
```
1. Open Settings tab
2. Change priority behavior
3. Adjust memory limits
4. View app information
```

## Common Issues

### "Cannot find module 'expo'"
```bash
rm -rf node_modules package-lock.json
npm install
```

### "Network error" when testing
- Check internet connection
- Verify API key is correct
- App will fallback to local mode

### "Permission denied"
- Grant location permission
- Grant microphone permission
- Restart app after granting

### Expo Go crashes
```bash
# Clear cache and restart
npm start -- --clear
```

## Next Steps

- Read [SETUP.md](./SETUP.md) for detailed configuration
- Check [TESTING.md](./TESTING.md) for comprehensive testing
- Review [API_INTEGRATION.md](./API_INTEGRATION.md) for LLM setup

## One-Command Demo

Want to see it without API keys?

```bash
# Run in local mode with sample data
npm install && \
echo 'LLM_PROVIDER=local' > .env && \
npm start
```

Then manually add test data in the app:
1. Go to LostToFound
2. Store 3-4 items
3. Query for them
4. Works without AI!

## Production Deployment

For actual deployment:

1. **Add sound files** to `assets/sounds/`
2. **Create app icons** in `assets/images/`
3. **Get production API keys**
4. **Build with EAS**: 
   ```bash
   npm install -g eas-cli
   eas build --platform all
   ```

See Expo documentation for detailed build instructions:
https://docs.expo.dev/build/introduction/

## Support

Stuck? Check:
- [SETUP.md](./SETUP.md) - Detailed setup
- [Troubleshooting](#common-issues) - Common fixes
- Expo Docs: https://docs.expo.dev/

## Video Tutorial

*Coming soon: Video walkthrough of setup and features*

---

**Time from zero to running app: ~5-10 minutes**

Enjoy testing the ADHD Support App! 🚀
