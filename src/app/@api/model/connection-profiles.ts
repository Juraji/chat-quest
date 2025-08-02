import {ChatQuestModel} from './model';

export type ProviderType = "OPEN_AI"

export interface ConnectionProfile extends ChatQuestModel {
  name: string;
  providerType: ProviderType
  baseUrl: string
  apiKey: string
}

export interface LlmModel extends ChatQuestModel{
  connectionProfileId: number
  modelId: string
  temperature: number
  maxTokens: number
  topP: number
  stream: boolean
  stopSequences: string
  disabled: boolean
}

export interface AiProviders {
  offline: ConnectionProfile[],
  online: ConnectionProfile[],
}
