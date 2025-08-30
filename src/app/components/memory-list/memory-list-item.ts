import {
  booleanAttribute,
  Component,
  computed,
  effect,
  inject,
  input,
  InputSignal,
  output,
  OutputEmitterRef
} from '@angular/core';
import {Memory} from '@api/memories';
import {FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {Characters} from '@api/characters';
import {booleanSignal, BooleanSignal, formControl, formGroup, readOnlyControl} from '@util/ng';
import {isNew} from '@api/common';
import {TimeAgoPipe} from '@components/time-ago.pipe';

@Component({
  selector: 'memory-list-item',
  imports: [
    FormsModule,
    ReactiveFormsModule,
    TimeAgoPipe
  ],
  templateUrl: './memory-list-item.html'
})
export class MemoryListItem {
  private readonly characters = inject(Characters)

  readonly memory: InputSignal<Memory> = input.required()
  readonly characterId: InputSignal<Nullable<number>> = input()
  readonly disabled = input(false, {transform: booleanAttribute})
  readonly memoryChanged: OutputEmitterRef<Memory> = output()
  readonly deleteRequested: OutputEmitterRef<void> = output()

  protected readonly isNew = computed(() => isNew(this.memory()))
  protected readonly editMode: BooleanSignal = booleanSignal(false)
  protected readonly allCharacters = this.characters.all
  protected readonly character = this.characters.listViewBy(() => this.memory().characterId)

  readonly formGroup = formGroup<Memory>({
    id: readOnlyControl(),
    worldId: readOnlyControl(),
    characterId: formControl(null),
    createdAt: readOnlyControl(),
    content: formControl('', [Validators.required]),
  })

  constructor() {
    effect(() => {
      const m = this.memory()
      this.formGroup.reset(m)
      this.editMode.set(isNew(m))
    });
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const update: Memory = {
      ...this.memory(),
      ...this.formGroup.value,
    }

    this.memoryChanged.emit(update)

    if (!this.isNew()) {
      this.editMode.setFalse()
    }
  }

  onDelete(): void {
    this.deleteRequested.emit()
  }
}
