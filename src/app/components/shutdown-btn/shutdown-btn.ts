import {Component, inject} from '@angular/core';
import {System} from '@api/system';

@Component({
  selector: 'app-shutdown-btn',
  imports: [],
  templateUrl: './shutdown-btn.html',
})
export class ShutdownBtn {
  private readonly system = inject(System);


  onShutdown() {
    const doShutdown = confirm("Do you want to shutdown? This will shutdown the ChatQuest backend.");

    if (doShutdown) {
      this.system
        .shutdown()
        .subscribe()
    }
  }
}
