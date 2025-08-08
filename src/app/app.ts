import {Component} from '@angular/core';
import {RouterLink, RouterLinkActive, RouterOutlet} from '@angular/router';
import {NotificationsDisplay} from '@components/notifications';
import {ShutdownBtn} from '@components/shutdown-btn';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, RouterLink, RouterLinkActive, NotificationsDisplay, ShutdownBtn],
  templateUrl: './app.html',
})
export class App {
}
