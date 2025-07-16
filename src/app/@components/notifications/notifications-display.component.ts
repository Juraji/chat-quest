import {Component, inject} from '@angular/core';
import {Notifications} from '@components/notifications/notifications';
import {TimeAgoPipe} from '@components/notifications/time-ago.pipe';
import {AsyncPipe} from '@angular/common';

@Component({
  selector: 'app-notifications',
  imports: [
    TimeAgoPipe,
    AsyncPipe
  ],
  templateUrl: './notifications-display.component.html',
  styleUrl: './notifications-display.component.scss'
})
export class NotificationsDisplay {
  readonly notifications = inject(Notifications)
}
