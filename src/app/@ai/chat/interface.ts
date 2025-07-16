export interface OpenAiSettings {
  baseUri: string
  apiKey: string

  temperature: number
  maxTokens: number
  topP: number
  stream: boolean
  stop: string
}

export interface OpenAiModel {
  id: string;
  object: string;
  created: number;
  root?: string;
}

export interface OpenAiChatCompletionRequest {
  model: string;
  messages: Array<{ role: string; content: string }>;
  temperature?: number;
  max_tokens?: number;
  top_p?: number;
  n?: number;
  stream?: boolean;
  stop?: string | string[];
}

export interface OpenAiChatCompletionChoice {
  message: { role: string; content: string };
  finish_reason: string;
  index: number;
}

export interface OpenAiChatCompletionUsage {
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
}

export interface OpenAiChatCompletionResponse {
  id: string;
  object: string;
  created: number;
  model: string;
  choices: OpenAiChatCompletionChoice[];
  usage?: OpenAiChatCompletionUsage;
}

export interface OpenAIListResponse<T> {
  object: 'list',
  data: T[],
}
