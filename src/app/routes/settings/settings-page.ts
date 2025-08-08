import {Component, inject, Signal} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {ConnectionProfilesOverview} from './components/connection-profiles';
import {InstructionOverview} from './components/instructions';
import {ChatSettings} from './components/chat-settings';
import {ActivatedRoute} from '@angular/router';
import {ChatPreferences} from '@api/worlds';
import {routeDataSignal} from '@util/ng';
import {MemorySettings} from './components/memory-settings';
import {ConnectionProfile, LlmModelView} from '@api/providers';
import {Instruction} from '@api/instructions';
import {MemoryPreferences} from '@api/memories';

@Component({
  selector: 'app-settings-page',
  imports: [
    PageHeader,
    ConnectionProfilesOverview,
    InstructionOverview,
    ChatSettings,
    MemorySettings
  ],
  templateUrl: './settings-page.html'
})
export class SettingsPage {
  private readonly activatedRoute = inject(ActivatedRoute);

  readonly profiles: Signal<ConnectionProfile[]> = routeDataSignal(this.activatedRoute, 'profiles')
  readonly templates: Signal<Instruction[]> = routeDataSignal(this.activatedRoute, 'templates')
  readonly chatPreferences: Signal<ChatPreferences> = routeDataSignal(this.activatedRoute, 'chatPreferences');
  readonly memoryPreferences: Signal<MemoryPreferences> = routeDataSignal(this.activatedRoute, 'memoryPreferences');
  readonly llmModelViews: Signal<LlmModelView[]> = routeDataSignal(this.activatedRoute, 'llmModelViews');

}
