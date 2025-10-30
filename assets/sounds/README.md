# Sound Assets

This directory should contain soothing sound files for the SoundSanctuary feature.

## Required Files

- `rain.mp3` - Rain sound effect
- `ocean.mp3` - Ocean waves sound
- `forest.mp3` - Forest ambience
- `white-noise.mp3` - White noise

## Where to Get Free Sounds

### Recommended Sources

1. **Freesound.org**
   - https://freesound.org/
   - Free sounds under Creative Commons
   - Search for: "rain loop", "ocean waves", "forest ambience", "white noise"

2. **Pixabay**
   - https://pixabay.com/sound-effects/
   - Completely free, no attribution required
   - High quality ambient sounds

3. **YouTube Audio Library**
   - https://www.youtube.com/audiolibrary
   - Free to use
   - Download and convert to MP3

4. **Free Music Archive**
   - https://freemusicarchive.org/
   - Various Creative Commons licenses
   - Ambient and sound effects sections

## File Requirements

- **Format**: MP3 (preferred) or WAV
- **Length**: 1-5 minutes (will loop)
- **Sample Rate**: 44.1 kHz or 48 kHz
- **Bit Rate**: 128-320 kbps
- **Size**: Keep under 5MB per file for app performance

## Example Downloads

### Rain Sound
Search for: "gentle rain loop"
- Duration: 2-3 minutes
- Type: Steady rain without thunder
- Volume: Medium, consistent

### Ocean Waves
Search for: "ocean waves calm"
- Duration: 2-3 minutes
- Type: Gentle waves, no seagulls
- Volume: Medium, rhythmic

### Forest
Search for: "forest ambience birds"
- Duration: 3-5 minutes
- Type: Birds chirping, leaves rustling
- Volume: Medium, varied

### White Noise
Search for: "white noise sleep"
- Duration: 1-2 minutes (simple loop)
- Type: Pure white noise
- Volume: Medium, constant

## Quick Setup Script

```bash
# Download free samples from Pixabay
cd assets/sounds

# Rain (example - replace with actual URLs)
curl -o rain.mp3 "https://example.com/rain.mp3"

# Ocean
curl -o ocean.mp3 "https://example.com/ocean.mp3"

# Forest
curl -o forest.mp3 "https://example.com/forest.mp3"

# White noise
curl -o white-noise.mp3 "https://example.com/whitenoise.mp3"
```

## Converting Files

If you have WAV files, convert to MP3:

```bash
# Using ffmpeg
ffmpeg -i input.wav -codec:a libmp3lame -b:a 192k output.mp3

# Or use online converter
# https://cloudconvert.com/wav-to-mp3
```

## Testing Sounds

After adding files, test them:

1. Open the app
2. Go to SoundSanctuary tab
3. Select each sound type
4. Play and verify quality
5. Check looping is smooth
6. Adjust volume levels

## App Integration

Once files are added, update `src/context/AppContext.js` line ~208:

```javascript
// Replace this:
{ uri: 'https://example.com/placeholder.mp3' }

// With this:
require(`../../assets/sounds/${soundName}.mp3`)
```

## Attribution

If using Creative Commons sounds, add attributions here:

- Rain: [Artist Name] - [Source URL] - [License]
- Ocean: [Artist Name] - [Source URL] - [License]
- Forest: [Artist Name] - [Source URL] - [License]
- White Noise: [Artist Name] - [Source URL] - [License]

## Notes

- The app currently uses placeholder audio URLs
- Real sound files are needed for full functionality
- Keep files optimized for mobile devices
- Test on actual devices for quality and performance
- Consider adding more sound options in future versions
