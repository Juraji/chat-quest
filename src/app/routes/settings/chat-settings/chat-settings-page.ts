import { Component } from '@angular/core';
import {ChatSettingsOpenAi} from './chat-settings-open-ai/chat-settings-open-ai';

@Component({
  selector: 'app-chat-settings-page',
  imports: [
    ChatSettingsOpenAi
  ],
  templateUrl: './chat-settings-page.html'
})
export class ChatSettingsPage {

}
