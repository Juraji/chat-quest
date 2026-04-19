import {inject} from '@angular/core';
import {SSE} from '@api/sse';
import {Notifications} from '@components/notifications';
import {LogMessages} from '@api/log';

export function sseInitializer() {
  const sse = inject(SSE)
  const notifications = inject(Notifications)

  sse.connect()

  // Display backend error messages
  sse
    .on(LogMessages, m => m.level == "ERROR")
    .subscribe(logMessage => notifications.toast(`Backend error: ${logMessage.message}`, "DANGER"))

  window.addEventListener("beforeunload", () => sse.disconnect())
}
