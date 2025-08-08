import {Component, inject} from '@angular/core';
import {SSE} from '@api/sse';

@Component({
  selector: 'sse-connection-status',
  imports: [],
  templateUrl: './sse-connection-status.html',
})
export class SseConnectionStatus {
  private readonly sse = inject(SSE)

  readonly status = this.sse.connectionState
}
