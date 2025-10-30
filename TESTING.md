# ADHD Support App - Testing Guide

## Testing Checklist

This guide helps you systematically test all features of the ADHD Support App on both Android and iOS devices.

## Prerequisites

- App installed and running on device/emulator
- LLM API key configured in `.env`
- Location and microphone permissions granted
- Test in both quiet and noisy environments

---

## 1. Initial Setup Tests

### ✅ App Launch
- [ ] App launches without errors
- [ ] Bottom navigation displays three tabs
- [ ] All screens are accessible
- [ ] No console errors in Expo logs

### ✅ Permissions
- [ ] Location permission requested on first use
- [ ] Microphone permission requested when accessing features
- [ ] Permission denial handled gracefully
- [ ] Settings link provided for denied permissions

---

## 2. LostToFound Feature Tests

### ✅ Basic Memory Storage

**Test 1: Text Input - Store Item**
1. Go to LostToFound tab
2. Type: "I'm leaving my car keys on the kitchen table"
3. Press send button
4. Expected: 
   - Success message appears
   - Text-to-speech reads confirmation
   - Item appears in stored memories list

**Test 2: Multiple Items**
1. Store 5 different items in different locations:
   - "My wallet is in the bedroom drawer"
   - "Phone charger is on my desk"
   - "Sunglasses are in the car"
   - "Laptop is in the living room"
   - "Headphones are in my backpack"
2. Expected: All 5 items appear in memories list

**Test 3: Memory Persistence**
1. Store an item
2. Close and restart the app
3. Expected: Item still appears in memories list

### ✅ Querying Memories

**Test 4: Simple Query**
1. After storing items, type: "Where are my car keys?"
2. Expected: 
   - AI responds with location
   - Response mentions "kitchen table"
   - Text-to-speech reads response

**Test 5: Natural Language Query**
1. Type: "Where did I put my wallet?"
2. Expected: AI responds with correct location

**Test 6: Query Non-existent Item**
1. Type: "Where is my laptop?"
2. Expected: AI politely indicates item not found

**Test 7: Conversation Context**
1. Store: "My book is on the nightstand"
2. Query: "Where is my book?"
3. Follow-up: "When did I put it there?"
4. Expected: AI maintains context and answers appropriately

### ✅ Location Features

**Test 8: GPS Coordinates**
1. Store an item
2. View memory in list
3. Expected: Navigate button appears if location available

**Test 9: Google Maps Integration**
1. Tap Navigate button on a memory
2. Expected: Google Maps opens with coordinates

### ✅ Memory Management

**Test 10: Delete Memory**
1. Tap delete icon on a memory
2. Confirm deletion
3. Expected: Memory removed from list

**Test 11: Rolling History (30 items)**
1. Store 35 items
2. Expected: Only last 30 items retained

**Test 12: Clear All Memories**
1. Go to Settings
2. Tap "Clear All Memories"
3. Confirm
4. Return to LostToFound
5. Expected: All memories removed

### ✅ LLM Integration

**Test 13: OpenAI Provider**
1. Configure OpenAI in .env
2. Store and query items
3. Expected: Intelligent responses

**Test 14: Gemini Provider**
1. Configure Gemini in .env
2. Restart app
3. Test storage and queries
4. Expected: Intelligent responses

**Test 15: Hugging Face Provider**
1. Configure Hugging Face in .env
2. Restart app
3. Test storage and queries
4. Expected: Intelligent responses

**Test 16: LLM Fallback**
1. Use invalid API key
2. Try to store/query items
3. Expected: Local search fallback works

---

## 3. SoundSanctuary Feature Tests

### ✅ Basic Playback

**Test 17: Manual Sound Playback**
1. Go to SoundSanctuary tab
2. Tap "Play Sound"
3. Expected: 
   - Status changes to "Playing"
   - Default rain sound plays (if audio files available)

**Test 18: Pause/Resume**
1. Play a sound
2. Tap "Pause Sound"
3. Tap "Play Sound" again
4. Expected: Sound pauses and resumes correctly

**Test 19: Stop Playback**
1. Play a sound
2. Tap "Stop Sound"
3. Expected: Sound stops completely

### ✅ Sound Selection

**Test 20: Change Sound Type**
1. Select "Ocean Waves"
2. Play sound
3. Expected: Ocean sound plays
4. Repeat for all available sounds

### ✅ Volume Control

**Test 21: Adjust Volume**
1. Play a sound
2. Move volume slider
3. Expected: Volume changes in real-time

**Test 22: Volume at 0%**
1. Set volume to minimum
2. Play sound
3. Expected: Sound plays but inaudible

**Test 23: Volume at 100%**
1. Set volume to maximum
2. Play sound
3. Expected: Sound plays at full volume

### ✅ Noise Monitoring

**Test 24: Start Monitoring**
1. Tap "Start Monitoring"
2. Expected:
   - Status changes to "Active"
   - Current noise level updates periodically
   - Alert confirms monitoring started

**Test 25: Threshold Detection (Simulated)**
1. Set threshold to 50 dB
2. Start monitoring
3. Wait for noise level to exceed threshold
4. Expected: Sound plays automatically

**Test 26: Stop Monitoring**
1. Stop monitoring
2. Expected: Noise level resets to 0

### ✅ Threshold Configuration

**Test 27: Adjust Threshold**
1. Move threshold slider to different values
2. Expected: Value updates in real-time

**Test 28: Threshold Tutorial**
1. Tap "Tutorial" button
2. Expected: Help dialog explains how to find threshold

---

## 4. Settings Tests

### ✅ Priority Management

**Test 29: Auto-Pause Behavior**
1. Go to Settings
2. Select "Auto-pause"
3. Go to SoundSanctuary, start playing sound
4. Go to LostToFound, send a message
5. Expected: Sound pauses during processing, resumes after

**Test 30: Duck Volume Behavior**
1. Select "Duck volume" in Settings
2. Play sound
3. Send LostToFound message
4. Expected: Volume reduces during processing

**Test 31: Ignore Behavior**
1. Select "Ignore" in Settings
2. Play sound
3. Send LostToFound message
4. Expected: Sound continues at normal volume

### ✅ Memory Settings

**Test 32: Change Max Memories**
1. Change maximum memories to 20
2. Store 25 items
3. Expected: Only last 20 items retained

**Test 33: Change Conversation History**
1. Set to 5 interactions
2. Have 10 back-and-forth conversations
3. Expected: Only last 5 retained

### ✅ Data Management

**Test 34: Clear Conversation History**
1. Have several conversations
2. Tap "Clear Conversation History"
3. Expected: History cleared, new queries start fresh

**Test 35: Reset Settings**
1. Change multiple settings
2. Tap "Reset Settings"
3. Confirm
4. Expected: All settings return to defaults

### ✅ Information Screens

**Test 36: App Information**
1. Tap "App Information"
2. Expected: Version and description displayed

**Test 37: LLM Setup Guide**
1. Tap "LLM Setup Guide"
2. Expected: Configuration instructions displayed

---

## 5. Cross-Feature Integration Tests

### ✅ Priority Coexistence

**Test 38: Sound Playing + Memory Query**
1. Start SoundSanctuary playback
2. Use LostToFound to query an item
3. Expected: Behavior follows settings (pause/duck/ignore)

**Test 39: Monitoring + Voice Command**
1. Enable noise monitoring
2. Query items on LostToFound
3. Expected: Monitoring continues, priority handled

### ✅ State Persistence

**Test 40: App Background**
1. Configure settings
2. Send app to background
3. Return to app
4. Expected: State preserved

**Test 41: App Restart**
1. Store memories, configure settings
2. Force quit and restart app
3. Expected: All data persisted

---

## 6. Error Handling Tests

### ✅ Network Issues

**Test 42: No Internet Connection**
1. Disable internet
2. Try to query items
3. Expected: Fallback to local search works

**Test 43: Invalid API Key**
1. Set wrong API key in .env
2. Try to store/query
3. Expected: Graceful error, fallback works

### ✅ Permission Issues

**Test 44: Location Denied**
1. Deny location permission
2. Store items
3. Expected: Items stored without coordinates, text input still works

**Test 45: Microphone Denied**
1. Deny microphone permission
2. Try voice input
3. Expected: Appropriate message, text input available

### ✅ Edge Cases

**Test 46: Empty Input**
1. Try to submit empty text
2. Expected: Error message or button disabled

**Test 47: Very Long Input**
1. Type 1000 characters
2. Expected: Input limited or handles gracefully

**Test 48: Special Characters**
1. Type: "I left my $ @ # in the ^ & *"
2. Expected: Handles without crashing

---

## 7. Platform-Specific Tests

### ✅ iOS Specific

**Test 49: iOS Permissions UI**
- [ ] Permission dialogs follow iOS design
- [ ] Settings deep-link works

**Test 50: iOS Maps Integration**
- [ ] Opens Apple Maps correctly

### ✅ Android Specific

**Test 51: Android Permissions UI**
- [ ] Permission dialogs follow Material Design
- [ ] Settings deep-link works

**Test 52: Android Maps Integration**
- [ ] Opens Google Maps correctly

**Test 53: Back Button Behavior**
- [ ] Android back button works as expected

---

## 8. Performance Tests

### ✅ Responsiveness

**Test 54: Large Memory List**
1. Store 30 items
2. Scroll through list
3. Expected: Smooth scrolling

**Test 55: Quick Actions**
1. Rapidly switch between tabs
2. Expected: No lag or crashes

### ✅ Battery Usage

**Test 56: Monitoring Battery Impact**
1. Enable monitoring for 30 minutes
2. Check battery usage
3. Expected: Reasonable battery consumption

---

## 9. Accessibility Tests

**Test 57: Screen Reader**
1. Enable TalkBack (Android) or VoiceOver (iOS)
2. Navigate through app
3. Expected: All elements properly labeled

**Test 58: Large Text**
1. Enable large text in device settings
2. Expected: App scales appropriately

---

## 10. Example Data Tests

### Sample Memories for Testing

Use these to quickly populate the app:

```
1. "My house keys are on the hook by the front door"
2. "Wallet is in my jacket pocket"
3. "Phone charger is plugged in by my bed"
4. "Laptop is on the dining table"
5. "Sunglasses are in the car glove box"
6. "Water bottle is in my gym bag"
7. "Medication is in the bathroom cabinet"
8. "Passport is in the safe"
9. "Headphones are in my backpack"
10. "Watch is on the nightstand"
```

### Sample Queries

```
1. "Where are my keys?"
2. "Where did I put my wallet?"
3. "I can't find my phone charger"
4. "Have you seen my laptop?"
5. "Where's my water bottle?"
```

---

## Test Results Template

| Test # | Test Name | iOS Pass | Android Pass | Notes |
|--------|-----------|----------|--------------|-------|
| 1      | App Launch | ✅ | ✅ | |
| 2      | Text Storage | ✅ | ✅ | |
| ...    | ...       | ... | ... | ... |

---

## Reporting Issues

When reporting bugs, include:
1. Test number and name
2. Device/emulator details
3. OS version
4. Steps to reproduce
5. Expected vs actual behavior
6. Screenshots/logs if applicable
7. Expo and React Native versions

---

## Automated Testing (Future)

Consider adding:
- Jest unit tests for utilities
- React Native Testing Library for components
- Detox for E2E tests
- CI/CD pipeline integration
