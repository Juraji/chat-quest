import {inject} from '@angular/core';
import {SSE} from '@api/clients';

export function sseInitializer() {
    const sse = inject(SSE)
    sse.connect()
}
