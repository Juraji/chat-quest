export interface SseMessageBody {
  source: string
  payload: any
}


export interface SseEvent<T> extends String {
  t?: T
}
