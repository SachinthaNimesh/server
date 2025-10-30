# Image Assets

This directory should contain app icons and images.

## Required Files

For Expo apps, you need:

- `icon.png` - App icon (1024x1024 px)
- `splash.png` - Splash screen (1242x2436 px for iOS, 1080x1920 px for Android)
- `adaptive-icon.png` - Android adaptive icon (1024x1024 px)
- `favicon.png` - Web favicon (48x48 px)

## Placeholder Creation

Until you have custom icons, you can use these online tools:

### Icon Generators
- https://www.figma.com/community - Free templates
- https://www.canva.com/ - Easy design tool
- https://makeappicon.com/ - Generates all sizes

### Recommended Design

For an ADHD support app:
- **Colors**: Calming blues, purples, greens
- **Icon**: Brain symbol, helper hand, or sound wave
- **Style**: Minimal, modern, accessible

### Quick Setup (Temporary)

Create solid color placeholders:

```bash
# Using ImageMagick
convert -size 1024x1024 xc:#6200EE icon.png
convert -size 1242x2436 xc:#FFFFFF splash.png
convert -size 1024x1024 xc:#6200EE adaptive-icon.png
convert -size 48x48 xc:#6200EE favicon.png
```

Or download from:
- https://placeholder.com/
- https://via.placeholder.com/1024x1024

## Current Setup

The app is configured in `app.json` to use these files. If they don't exist, Expo will use defaults.

## Future Enhancements

Consider adding:
- Logo variations (light/dark mode)
- Tutorial images
- Feature screenshots
- Empty state illustrations
- Loading animations
