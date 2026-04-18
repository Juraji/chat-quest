import {SseEvent} from '@api/sse';

export type LogMessageLevel = "DEBUG" | "INFO" | "WARN" | "ERROR"

export interface LogMessage {
  level: LogMessageLevel
  time: string
  message: string
  fields: Record<string, any>
}

export const LogMessages: SseEvent<LogMessage> = 'LogMessages'
