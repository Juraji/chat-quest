import {Injectable, Signal, signal, WritableSignal} from '@angular/core';
import {arrayRemove} from '@util/array';
import {finalize, from, Observable, ObservableInput} from 'rxjs';

type ToastType = "INFO" | "WARNING" | "DANGER"

interface Toast {
  id: number
  createdAt: Date,
  type: ToastType
  message: string,
  timerHandle: number | null,
}

const DEFAULT_TOAST_TIMEOUT = 5000

@Injectable({
  providedIn: 'root'
})
export class Notifications {
  private lastId: number = 1

  protected readonly _toasts: WritableSignal<Toast[]> = signal([])
  readonly toasts: Signal<Toast[]> = this._toasts

  constructor() {
  }

  /**
   * Displays a toast notification with the given message and optional styling.
   *
   * @param {string} message - The text content to display in the toast notification.
   * @param {ToastType} [type='INFO'] - The visual type of the toast, determining its styling and severity level. Defaults to 'INFO'.
   * @param {number} [timeout=DEFAULT_TOAST_TIMEOUT] - The duration in milliseconds before automatically removing the toast. Set to 0 or negative to prevent auto-dismissal.
   *
   * @return {number} A unique identifier for the created toast, used to reference and manage it later (e.g., removal).
   */
  toast(
    message: string,
    type: ToastType = 'INFO',
    timeout: number = DEFAULT_TOAST_TIMEOUT
  ): number {
    const toastId = ++this.lastId
    const timerHandle = timeout <= 0 ? null : window.setTimeout(() => this.removeToast(toastId), timeout)

    const toast: Toast = {
      id: toastId,
      createdAt: new Date(),
      type, message,
      timerHandle,
    }

    this._toasts.update(toasts => [toast, ...toasts])

    return toastId
  }

  run<R>(message: string, type: ToastType = 'INFO', action: () => ObservableInput<R>): Observable<R> {
    const toastId = this.toast(message, type, -1)
    return from(action()).pipe(finalize(() => this.removeToast(toastId)))
  }

  removeToast(toastId: number) {
    this._toasts.update(toasts => {
      // Find the toast with matching ID
      const index = toasts.findIndex(t => t.id === toastId)

      if (index !== -1) {
        // Get the toast object
        const toast = toasts[index]

        // Clear the timer if it exists
        if (toast.timerHandle) {
          clearTimeout(toast.timerHandle)
        }

        // Remove the toast from the array and return the updated array
        return arrayRemove(toasts, index)
      }

      return toasts
    })
  }
}
