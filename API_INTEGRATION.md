# LLM API Integration Guide

## Overview

This app supports three LLM providers: OpenAI, Google Gemini, and Hugging Face. This guide helps you set up each provider and understand how the integration works.

## Quick Setup

### Step 1: Choose Your Provider

| Provider | Cost | Speed | Quality | Free Tier |
|----------|------|-------|---------|-----------|
| **OpenAI** | $$$ | Fast | Excellent | $5 credit |
| **Gemini** | $ | Fast | Excellent | 60 req/min |
| **Hugging Face** | Free | Medium | Good | Yes |

### Step 2: Get API Key

#### OpenAI (gpt-3.5-turbo)

1. Go to https://platform.openai.com/signup
2. Create account or sign in
3. Navigate to https://platform.openai.com/api-keys
4. Click "Create new secret key"
5. Name it "ADHD Support App"
6. Copy the key (starts with `sk-`)

**Pricing**: $0.50 per 1M input tokens, $1.50 per 1M output tokens  
**Free Tier**: $5 credit for new accounts

#### Google Gemini (gemini-1.5-flash)

1. Go to https://makersuite.google.com/
2. Sign in with Google account
3. Navigate to https://makersuite.google.com/app/apikey
4. Click "Create API key"
5. Select project or create new one
6. Copy the key

**Pricing**: Free tier - 15 requests per minute, 1500 per day  
**Free Tier**: Yes, generous limits

#### Hugging Face (Free Inference API)

1. Go to https://huggingface.co/join
2. Create account or sign in
3. Navigate to https://huggingface.co/settings/tokens
4. Click "New token"
5. Name it "ADHD Support App"
6. Select "Read" permission
7. Copy the token

**Pricing**: Free for community inference API  
**Note**: May have rate limits or cold start delays

### Step 3: Configure .env File

Create a `.env` file in the project root:

```bash
cp .env.example .env
```

Then edit `.env` and add your chosen provider:

**For OpenAI:**
```env
OPENAI_API_KEY=sk-proj-xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
OPENAI_MODEL=gpt-3.5-turbo
LLM_PROVIDER=openai
```

**For Gemini:**
```env
GEMINI_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
GEMINI_MODEL=gemini-1.5-flash
LLM_PROVIDER=gemini
```

**For Hugging Face:**
```env
HUGGINGFACE_API_KEY=hf_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
HUGGINGFACE_MODEL=mistralai/Mistral-7B-Instruct-v0.1
LLM_PROVIDER=huggingface
```

### Step 4: Restart App

```bash
# Stop the current Expo server
# Press Ctrl+C

# Start again
npm start
```

---

## Testing Your Setup

### Test 1: Store an Item

```
Input: "I'm putting my laptop on the dining table"
Expected: Success message confirming storage
```

### Test 2: Query an Item

```
Input: "Where is my laptop?"
Expected: Response mentioning "dining table"
```

### Test 3: Check Logs

Look for these in your Expo console:

**Success:**
```
✓ LLM API call successful (openai)
✓ Extracted: {item: "laptop", location: "dining table"}
```

**Failure:**
```
✗ LLM API Error: Invalid API key
→ Falling back to local search
```

---

## API Implementation Details

### How It Works

```
User Input
    ↓
Voice/Text Recognition
    ↓
LLM API Call (with conversation history)
    ↓
Response Parsing
    ↓
Memory Storage / Query Response
    ↓
Text-to-Speech Output
```

### Message Format

All LLM providers receive messages in this format:

```javascript
[
  {
    role: "system",
    content: "You are a helpful assistant that extracts..."
  },
  {
    role: "user", 
    content: "I'm leaving my keys on the kitchen table"
  },
  {
    role: "assistant",
    content: '{"item": "keys", "location": "kitchen table"}'
  },
  {
    role: "user",
    content: "Where are my keys?"
  }
]
```

### Expected Responses

**For Storage:**
```json
{
  "item": "car keys",
  "location": "kitchen table"
}
```

**For Queries:**
```json
{
  "action": "query",
  "item": "car keys"
}
```

**Free-form response** (also supported):
```
Your car keys are on the kitchen table. You stored them there on January 15 at 10:30 AM.
```

---

## Provider-Specific Details

### OpenAI

**Endpoint:**
```
POST https://api.openai.com/v1/chat/completions
```

**Request Format:**
```javascript
{
  model: "gpt-3.5-turbo",
  messages: [...],
  temperature: 0.7,
  max_tokens: 500
}
```

**Response Format:**
```javascript
{
  choices: [{
    message: {
      role: "assistant",
      content: "..."
    }
  }]
}
```

**Rate Limits:**
- 3 requests per minute (free tier)
- 200 requests per minute (paid tier)

**Common Errors:**
- `401`: Invalid API key
- `429`: Rate limit exceeded
- `500`: Server error

### Google Gemini

**Endpoint:**
```
POST https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key={API_KEY}
```

**Request Format:**
```javascript
{
  contents: [{
    parts: [{
      text: "Combined prompt from messages..."
    }]
  }]
}
```

**Response Format:**
```javascript
{
  candidates: [{
    content: {
      parts: [{
        text: "..."
      }]
    }
  }]
}
```

**Rate Limits:**
- 15 requests per minute (free tier)
- 1500 requests per day (free tier)

**Common Errors:**
- `400`: Invalid request format
- `403`: API key issue
- `429`: Quota exceeded

### Hugging Face

**Endpoint:**
```
POST https://api-inference.huggingface.co/models/{MODEL_NAME}
```

**Request Format:**
```javascript
{
  inputs: "Combined prompt from messages...",
  parameters: {
    max_new_tokens: 500,
    temperature: 0.7,
    return_full_text: false
  }
}
```

**Response Format:**
```javascript
[{
  generated_text: "..."
}]
```

**Rate Limits:**
- Varies by model
- May experience cold starts (~20 seconds)

**Common Errors:**
- `503`: Model loading (wait and retry)
- `401`: Invalid token
- `429`: Rate limited

**Recommended Models:**
- `mistralai/Mistral-7B-Instruct-v0.1` (good quality)
- `google/flan-t5-xxl` (faster, smaller)
- `bigscience/bloom-560m` (very fast, lower quality)

---

## Troubleshooting

### "Invalid API Key" Error

**OpenAI:**
```bash
# Test your key
curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY"

# Should return list of models
```

**Gemini:**
```bash
# Test your key
curl "https://generativelanguage.googleapis.com/v1beta/models?key=$GEMINI_API_KEY"

# Should return list of models
```

**Hugging Face:**
```bash
# Test your token
curl https://huggingface.co/api/whoami \
  -H "Authorization: Bearer $HUGGINGFACE_API_KEY"

# Should return your username
```

### "Rate Limit Exceeded"

1. Wait a minute and try again
2. Upgrade to paid tier (OpenAI/Gemini)
3. Switch to different provider
4. Use local fallback (automatic in app)

### "Unexpected Response Format"

Check logs for raw API response:
```javascript
console.log('LLM Response:', response.data);
```

The app should handle most formats, but you can modify parsing in `src/utils/llmService.js`:

```javascript
// Around line 180
try {
  const jsonMatch = response.content.match(/\{[\s\S]*\}/);
  if (jsonMatch) {
    return JSON.parse(jsonMatch[0]);
  }
} catch (parseError) {
  // Add custom parsing here
}
```

### Local Fallback Not Working

The app automatically falls back to local search when LLM fails. Test it:

1. Set invalid API key
2. Store items
3. Query items
4. Should get results from local memory

---

## Switching Providers

To switch providers, just update `.env`:

```bash
# Change from OpenAI to Gemini
LLM_PROVIDER=gemini

# Or Hugging Face
LLM_PROVIDER=huggingface
```

Then restart the app. No code changes needed!

---

## Cost Optimization

### Reduce Token Usage

Edit `src/utils/llmService.js` to reduce conversation history:

```javascript
// Line 155 - reduce from 10 to 5
...conversationHistory.slice(-5).map(h => ({
```

### Use Cheaper Models

**OpenAI:**
```env
OPENAI_MODEL=gpt-3.5-turbo-0125  # Newer, cheaper
```

**Hugging Face:**
```env
HUGGINGFACE_MODEL=google/flan-t5-base  # Smaller, faster
```

### Cache Responses

Add caching to avoid redundant API calls:

```javascript
const responseCache = new Map();

async function callLLMWithCache(messages) {
  const key = JSON.stringify(messages);
  if (responseCache.has(key)) {
    return responseCache.get(key);
  }
  const response = await callLLM(messages);
  responseCache.set(key, response);
  return response;
}
```

---

## Advanced Configuration

### Custom System Prompts

Edit `src/utils/llmService.js` line 143:

```javascript
const systemPrompt = `You are a helpful assistant that extracts information about items and their locations from user statements.

Your custom instructions here...

Respond in JSON format:
{
  "item": "item name",
  "location": "location description"
}`;
```

### Timeout Configuration

Add timeout to API calls:

```javascript
const response = await axios.post(url, data, {
  headers: { ... },
  timeout: 10000  // 10 seconds
});
```

### Retry Logic

Add automatic retries on failure:

```javascript
async function callLLMWithRetry(messages, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await callLLM(messages);
    } catch (error) {
      if (i === maxRetries - 1) throw error;
      await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1)));
    }
  }
}
```

---

## Security Best Practices

### ✅ DO

- Store API keys in `.env` (not in code)
- Add `.env` to `.gitignore`
- Use environment-specific keys (dev/prod)
- Rotate keys periodically
- Monitor API usage

### ❌ DON'T

- Commit API keys to git
- Share keys publicly
- Use production keys in development
- Hardcode keys in source files
- Expose keys in client-side code

### Production Considerations

For production apps, use a backend proxy:

```
Mobile App → Your Backend → LLM API
```

This keeps API keys secure on the server.

---

## Monitoring and Debugging

### Enable Debug Logging

Add to `src/utils/llmService.js`:

```javascript
const DEBUG = true;

if (DEBUG) {
  console.log('Request:', messages);
  console.log('Response:', response);
}
```

### Track API Usage

```javascript
let apiCallCount = 0;
let totalTokens = 0;

// After each call
apiCallCount++;
totalTokens += response.usage?.total_tokens || 0;

console.log(`API Calls: ${apiCallCount}, Tokens: ${totalTokens}`);
```

### Error Tracking

Integrate error tracking service:

```javascript
import * as Sentry from '@sentry/react-native';

try {
  await callLLM(messages);
} catch (error) {
  Sentry.captureException(error);
  throw error;
}
```

---

## FAQ

**Q: Which provider is best?**  
A: OpenAI (gpt-3.5-turbo) for best quality. Gemini for free tier. Hugging Face for fully free.

**Q: Can I use multiple providers?**  
A: Yes! Set up all keys in `.env` and switch `LLM_PROVIDER` as needed.

**Q: What if all providers fail?**  
A: The app automatically uses local search as fallback.

**Q: How much will it cost?**  
A: For typical use (~50 queries/day), less than $1/month on OpenAI.

**Q: Can I run without any LLM?**  
A: Yes! Leave API keys empty and the app uses local search only.

**Q: How do I add a new provider?**  
A: Add a new function in `llmService.js` following the existing pattern.

---

## Support

- OpenAI: https://help.openai.com/
- Gemini: https://ai.google.dev/docs
- Hugging Face: https://discuss.huggingface.co/

For app-specific issues, see main README.md troubleshooting section.
