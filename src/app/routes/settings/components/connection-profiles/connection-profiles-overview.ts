import {Component, input, InputSignal} from '@angular/core';
import {RouterLink} from '@angular/router';
import {EmptyPipe} from '@components/empty.pipe';
import {ConnectionProfile} from '@api/providers';

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
