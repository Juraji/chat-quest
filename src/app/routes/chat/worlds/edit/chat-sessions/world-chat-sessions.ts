import {Component} from '@angular/core';
import {ReactiveFormsModule} from '@angular/forms';
import {NewChatSession} from './components/new-chat-session-form/new-chat-session';

@Component({
  selector: 'world-chat-sessions',
  imports: [
    ReactiveFormsModule,
    NewChatSession
  ],
  templateUrl: './world-chat-sessions.html'
})
export class WorldChatSessions {

}
