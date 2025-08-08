import {DestroyRef, inject, Injectable, Signal, signal} from '@angular/core';
import {map, Observable, Subject, timer} from 'rxjs';
import {filter, share} from 'rxjs/operators';
import {SseMessageBody} from './sse.model';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import {ChatQuestConfig} from '@config/config';

@Injectable({providedIn: 'root'})
export class SSE {
  private readonly destroyRef = inject(DestroyRef)
  private readonly config = inject(ChatQuestConfig)
  private readonly events = new Subject<SseMessageBody>();
  private eventSource: EventSource | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = this.config.sse.maxReconnectAttempts;
  private reconnectionDelay = this.config.sse.reconnectionDelay;

  private readonly _connectionState = signal<number>(EventSource.CLOSED);
  readonly connectionState: Signal<number> = this._connectionState;

  on<T>(source: string): Observable<T> {
    return this.events.pipe(
      filter(message => message.source === source),
      map(message => message.payload),
      share(),
      takeUntilDestroyed(this.destroyRef)
    );
  }

  connect() {
    this.reconnectAttempts = 0;
    this.reconnect()
  }

  disconnect() {
    // Prevent reconnect on destroy
    this.reconnectAttempts = this.maxReconnectAttempts;
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
  }

  private reconnect() {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }

    const s = new EventSource(`${this.config.apiBaseUrl}/sse`);
    this.eventSource = s;
    this._connectionState.set(s.CONNECTING);

    s.onopen = () => {
      console.log('SSE connection established');
      this.reconnectAttempts = 0;
      this._connectionState.set(s.OPEN);
    };

    s.addEventListener('message', e => {
      try {
        const data = JSON.parse(e.data);
        this.events.next(data);
      } catch (e) {
        console.error('Error parsing SSE message:', e);
      }
    })

    s.addEventListener('ping', e => {
      console.log(`Received SSE ping with timestamp: ${e.data}`);
    })

    s.onerror = () => {
      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        const delay = this.reconnectionDelay * (this.reconnectAttempts * this.reconnectAttempts);

        console.log(`SSE connection error. Reconnecting in ${delay}ms...`);
        timer(delay).subscribe(() => this.reconnect());
      } else {
        console.error('Max reconnection attempts reached');
        s.close();
        this._connectionState.set(s.CLOSED);
      }
    };
  }
}
