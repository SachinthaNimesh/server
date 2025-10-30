import axios from 'axios';
import {
  OPENAI_API_KEY,
  OPENAI_MODEL,
  GEMINI_API_KEY,
  GEMINI_MODEL,
  HUGGINGFACE_API_KEY,
  HUGGINGFACE_MODEL,
  LLM_PROVIDER,
} from '@env';

/**
 * LLM Service for interacting with different AI providers
 * Supports OpenAI, Google Gemini, and Hugging Face
 */

/**
 * Call OpenAI API
 */
const callOpenAI = async (messages) => {
  try {
    const response = await axios.post(
      'https://api.openai.com/v1/chat/completions',
      {
        model: OPENAI_MODEL || 'gpt-3.5-turbo',
        messages: messages,
        temperature: 0.7,
        max_tokens: 500,
      },
      {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${OPENAI_API_KEY}`,
        },
      }
    );

    return {
      success: true,
      content: response.data.choices[0].message.content,
      provider: 'openai',
    };
  } catch (error) {
    console.error('OpenAI API Error:', error.response?.data || error.message);
    return {
      success: false,
      error: error.response?.data?.error?.message || error.message,
      provider: 'openai',
    };
  }
};

/**
 * Call Google Gemini API
 */
const callGemini = async (messages) => {
  try {
    // Convert messages to Gemini format
    const prompt = messages.map(msg => `${msg.role}: ${msg.content}`).join('\n\n');

    const response = await axios.post(
      `https://generativelanguage.googleapis.com/v1beta/models/${GEMINI_MODEL || 'gemini-1.5-flash'}:generateContent?key=${GEMINI_API_KEY}`,
      {
        contents: [{
          parts: [{
            text: prompt
          }]
        }]
      },
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    return {
      success: true,
      content: response.data.candidates[0].content.parts[0].text,
      provider: 'gemini',
    };
  } catch (error) {
    console.error('Gemini API Error:', error.response?.data || error.message);
    return {
      success: false,
      error: error.response?.data?.error?.message || error.message,
      provider: 'gemini',
    };
  }
};

/**
 * Call Hugging Face Inference API
 */
const callHuggingFace = async (messages) => {
  try {
    // Convert messages to a single prompt
    const prompt = messages.map(msg => `${msg.role}: ${msg.content}`).join('\n\n');

    const response = await axios.post(
      `https://api-inference.huggingface.co/models/${HUGGINGFACE_MODEL || 'mistralai/Mistral-7B-Instruct-v0.1'}`,
      {
        inputs: prompt,
        parameters: {
          max_new_tokens: 500,
          temperature: 0.7,
          return_full_text: false,
        },
      },
      {
        headers: {
          'Authorization': `Bearer ${HUGGINGFACE_API_KEY}`,
          'Content-Type': 'application/json',
        },
      }
    );

    return {
      success: true,
      content: response.data[0].generated_text,
      provider: 'huggingface',
    };
  } catch (error) {
    console.error('Hugging Face API Error:', error.response?.data || error.message);
    return {
      success: false,
      error: error.response?.data?.error || error.message,
      provider: 'huggingface',
    };
  }
};

/**
 * Main function to call LLM based on provider
 */
export const callLLM = async (messages, provider = LLM_PROVIDER) => {
  const selectedProvider = provider || 'openai';

  switch (selectedProvider.toLowerCase()) {
    case 'openai':
      return await callOpenAI(messages);
    case 'gemini':
      return await callGemini(messages);
    case 'huggingface':
      return await callHuggingFace(messages);
    default:
      return {
        success: false,
        error: `Unknown provider: ${selectedProvider}`,
      };
  }
};

/**
 * Extract item and location from user input using LLM
 */
export const extractMemoryInfo = async (userInput, conversationHistory = []) => {
  const systemPrompt = `You are a helpful assistant that extracts information about items and their locations from user statements.
When a user tells you where they placed an item, extract:
1. Item name (what they placed)
2. Location description (where they placed it)

Respond in JSON format:
{
  "item": "item name",
  "location": "location description"
}

If the user is asking about an item (not storing), respond with:
{
  "action": "query",
  "item": "item name"
}`;

  const messages = [
    { role: 'system', content: systemPrompt },
    ...conversationHistory.slice(-10).map(h => ({
      role: h.role === 'user' ? 'user' : 'assistant',
      content: h.content,
    })),
    { role: 'user', content: userInput },
  ];

  const response = await callLLM(messages);

  if (response.success) {
    try {
      // Try to parse JSON from the response
      const jsonMatch = response.content.match(/\{[\s\S]*\}/);
      if (jsonMatch) {
        return {
          success: true,
          data: JSON.parse(jsonMatch[0]),
          rawResponse: response.content,
        };
      }
    } catch (parseError) {
      console.error('Error parsing LLM response:', parseError);
    }
  }

  return {
    success: false,
    error: response.error || 'Failed to parse response',
    rawResponse: response.content,
  };
};

/**
 * Query memories using LLM
 */
export const queryMemories = async (query, memories, conversationHistory = []) => {
  // Create context from memories
  const memoryContext = memories.map((m, idx) => 
    `${idx + 1}. ${m.item} - ${m.location} (stored on ${new Date(m.timestamp).toLocaleString()})`
  ).join('\n');

  const systemPrompt = `You are a helpful memory assistant. The user has stored the following items:

${memoryContext}

Help the user find what they're looking for based on this information. Be conversational and helpful.`;

  const messages = [
    { role: 'system', content: systemPrompt },
    ...conversationHistory.slice(-10).map(h => ({
      role: h.role === 'user' ? 'user' : 'assistant',
      content: h.content,
    })),
    { role: 'user', content: query },
  ];

  return await callLLM(messages);
};

/**
 * Fallback local search when LLM is unavailable
 */
export const localMemorySearch = (query, memories) => {
  const queryLower = query.toLowerCase();
  const matches = memories.filter(m => 
    m.item.toLowerCase().includes(queryLower) ||
    m.location.toLowerCase().includes(queryLower)
  );

  if (matches.length === 0) {
    return {
      success: true,
      content: "I couldn't find any matching items in your memory. Try being more specific or check if you've stored that item.",
      matches: [],
    };
  }

  const responseText = matches.map((m, idx) => 
    `${idx + 1}. ${m.item} is at ${m.location} (stored on ${new Date(m.timestamp).toLocaleString()})`
  ).join('\n');

  return {
    success: true,
    content: `I found these items:\n\n${responseText}`,
    matches,
  };
};
