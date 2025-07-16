import {Injectable, Signal, signal, WritableSignal} from '@angular/core';

type ToastType = "INFO" | "WARNING" | "DANGER"

interface Toast {
  id: number
  createdAt: Date,
  type: ToastType
  label: string,
  message: string,
  timerHandle: number | null,
}

const DEFAULT_TOAST_TIMEOUT = 1e4 // 10 seconds

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
   * Add a toast.
   *
   * @param label The label for the toast. (Generally, this is the type of item spawned this message, like "Models" or "Scenes".)
   * @param message The message to display (HTML is not supported and will throw a sanitization error!)
   * @param type One of ToastType.
   * @param timeout When not <=0 the message will be dismissed after the given amount of milliseconds approximately.
   */
  toast(
    label: string,
    message: string,
    type: ToastType = 'INFO',
    timeout: number = DEFAULT_TOAST_TIMEOUT
  ) {
    const toastId = ++this.lastId
    const timerHandle = timeout <= 0 ? null : setTimeout(() => this.removeToast(toastId), timeout)

    const toast: Toast = {
      id: toastId,
      createdAt: new Date(),
      type, label, message,
      timerHandle,
    }

    this._toasts.update(toasts => [toast, ...toasts])
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
        return [...toasts.slice(0, index), ...toasts.slice(index + 1)]
      }

      return toasts
    })
  }
}
