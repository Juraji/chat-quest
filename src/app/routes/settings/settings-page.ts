import {Component, inject, Signal} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {ConnectionProfilesOverview} from './components/connection-profiles';
import {InstructionTemplatesOverview} from './components/instruction-templates';
import {ChatSettings} from './components/chat-settings';
import {ActivatedRoute} from '@angular/router';
import {ChatPreferences, ConnectionProfile, InstructionTemplate, MemoryPreferences} from '@api/model';
import {routeDataSignal} from '@util/ng';
import {MemorySettings} from './components/memory-settings';

@Component({
  selector: 'app-settings-page',
  imports: [
    PageHeader,
    ConnectionProfilesOverview,
    InstructionTemplatesOverview,
    ChatSettings,
    MemorySettings
  ],
  templateUrl: './settings-page.html'
})
export class SettingsPage {
  private readonly activatedRoute = inject(ActivatedRoute);

  readonly profiles: Signal<ConnectionProfile[]> = routeDataSignal(this.activatedRoute, 'profiles')
  readonly templates: Signal<InstructionTemplate[]> = routeDataSignal(this.activatedRoute, 'templates')
  readonly chatPreferences: Signal<ChatPreferences> = routeDataSignal(this.activatedRoute, 'chatPreferences');
  readonly memoryPreferences: Signal<MemoryPreferences> = routeDataSignal(this.activatedRoute, 'memoryPreferences');

}
