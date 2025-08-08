import {inject} from '@angular/core';
import {SSE} from '@api/sse';

export function sseInitializer() {
  const sse = inject(SSE)
  sse.connect()

  window.addEventListener("beforeunload", () => sse.disconnect())
}
