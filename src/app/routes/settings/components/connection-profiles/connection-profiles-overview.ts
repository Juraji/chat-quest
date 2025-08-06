import {Component, input, InputSignal} from '@angular/core';
import {RouterLink} from '@angular/router';
import {ConnectionProfile} from '@api/model';
import {EmptyPipe} from '@components/empty-pipe';

@Component({
  selector: 'connection-profiles-overview',
  imports: [
    RouterLink,
    EmptyPipe
  ],
  templateUrl: './connection-profiles-overview.html'
})
export class ConnectionProfilesOverview {
  readonly profiles: InputSignal<ConnectionProfile[]> = input.required()
}
