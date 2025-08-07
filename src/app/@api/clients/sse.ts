import {DestroyRef, inject, Injectable, OnDestroy} from '@angular/core';
import {map, Observable, Subject, timer} from 'rxjs';
import {filter, share} from 'rxjs/operators';
import {SseMessageBody} from '@api/model/sse';
import {ChatQuestConfig} from '@api/config';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';

@Injectable({providedIn: 'root'})
export class SSE implements OnDestroy {
  private readonly destroyRef = inject(DestroyRef)
  private readonly config = inject(ChatQuestConfig)
  private readonly events = new Subject<SseMessageBody>();
  private eventSource: EventSource | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = this.config.sse.maxReconnectAttempts;
  private reconnectionDelay = this.config.sse.reconnectionDelay;

  ngOnDestroy() {
    if (this.eventSource) {
      this.eventSource.close();
    }
  }

  on<T>(source: string): Observable<T> {
    return this.events.pipe(
      filter(message => message.source === source),
      map(message => message.payload),
      share(),
      takeUntilDestroyed(this.destroyRef)
    );
  }

  connect() {
    this.reconnect()
  }

  private reconnect() {
    if (this.eventSource) {
      this.eventSource.close();
    }

    // Reset attempts on new connection attempt
    this.reconnectAttempts = 0;

    const s = new EventSource(`${this.config.apiBaseUrl}/sse`);
    this.eventSource = s;

    s.onopen = () => {
      console.log('SSE connection established');
      this.reconnectAttempts = 0;
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
        timer(delay).subscribe(() => this.connect());
      } else {
        console.error('Max reconnection attempts reached');
        s.close();
      }
    };
  }
}
