import {ChatQuestModel} from '@api/common';

export type ProviderType = "OPEN_AI"
export type LlmModelType =
  "UNKNOWN" |
  "CHAT_MODEL" |
  "EMBEDDING_MODEL"

export interface ConnectionProfile extends ChatQuestModel {
  name: string;
  providerType: ProviderType
  baseUrl: string
  apiKey: string
}

export interface LlmModel extends ChatQuestModel {
  profileId: number
  modelId: string
  modelType: LlmModelType
  disabled: boolean
}

export interface LlmModelView {
  id: number
  modelId: string
  profileId: number
  profileName: string
}

export interface AiProviders {
  offline: ConnectionProfile[],
  online: ConnectionProfile[],
}
