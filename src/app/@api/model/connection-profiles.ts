import {ChatQuestModel} from './model';

export type ProviderType = "OPEN_AI"

export interface ConnectionProfile extends ChatQuestModel {
  name: string;
  providerType: ProviderType
  baseUrl: string
  apiKey: string
}

export interface LlmModel {
  id: number
  connectionProfileId: number
  modelId: string
  temperature: number
  maxTokens: number
  topP: number
  stream: boolean
  stop: string[]
}
