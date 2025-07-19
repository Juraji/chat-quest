import {Component} from '@angular/core';
import {ChatSettingsOpenAi} from './chat-settings-open-ai/chat-settings-open-ai';
import {ChatSettingsSystemPrompts} from './chat-settings-system-prompts/chat-settings-system-prompts';

@Component({
  selector: 'app-chat-settings-page',
  imports: [
    ChatSettingsOpenAi,
    ChatSettingsSystemPrompts
  ],
  templateUrl: './chat-settings-page.html'
})
export class ChatSettingsPage {

}
