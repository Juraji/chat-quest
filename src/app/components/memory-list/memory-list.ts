import {
  booleanAttribute,
  Component,
  computed,
  effect,
  inject,
  input,
  InputSignal,
  Signal,
  signal,
  WritableSignal
} from '@angular/core';
import {Memories, Memory, MemoryCreated, MemoryDeleted, MemoryUpdated} from '@api/memories';
import {MemoryListItem} from './memory-list-item';
import {arrayAdd, arrayRemove, arrayReplace} from '@util/array';
import {NEW_ID} from '@api/common';
import {defer, map} from 'rxjs';
import {booleanSignal} from '@util/ng';
import {SSE} from '@api/sse';
import {Characters} from '@api/characters';

interface PerCharacterCount {
  characterName: string
  total: number
  always: number
}

@Component({
  selector: 'memory-list',
  imports: [
    MemoryListItem
  ],
  templateUrl: './memory-list.html',
  styleUrls: ['./memory-list.scss'],
})
export class MemoryList {
  private readonly memoriesService = inject(Memories)
  private readonly charactersService = inject(Characters)
  private readonly sse = inject(SSE)

  readonly worldId: InputSignal<number> = input.required()
  readonly characterId: InputSignal<Nullable<number>> = input()
  readonly disabled = input(false, {transform: booleanAttribute})
  readonly collapsed = input(true, {transform: booleanAttribute})
  protected readonly isCollapsed = booleanSignal(this.collapsed)

  protected readonly memories: WritableSignal<Memory[]> = signal([])
  protected readonly newMemories: WritableSignal<Memory[]> = signal([])

  protected readonly showCounts = booleanSignal(false)
  protected readonly counts: Signal<PerCharacterCount[]> = computed(() => {
    const characters = this.charactersService.all()
    const memories = this.memories()

    const counts: PerCharacterCount[] = []
    for (const char of characters) {
      const charMemories = memories.filter(m => m.characterId === char.id)
      if (charMemories.length === 0) continue
      counts.push({
        characterName: char.name,
        total: charMemories.length,
        always: charMemories.filter(m => m.alwaysInclude).length
      })
    }

    return counts
  })

  constructor() {
    effect(() => {
      const worldId = this.worldId()
      const characterId = this.characterId()

      const memories$ = defer(() => {
        if (characterId === undefined) {
          return this.memoriesService.getAll(worldId)
        } else if (!!characterId) {
          return this.memoriesService.getAllByCharacter(worldId, characterId)
        } else {
          return []
        }
      })

      memories$
        .pipe(map(memories => memories
          .sort((a, b) => b.createdAt!.localeCompare(a.createdAt!))))
        .subscribe(memories => this.memories.set(memories))
    });

    const eventFilter: (m: Memory) => boolean = m => {
      if (m.worldId !== this.worldId()) return false;

      const cId = this.characterId()
      if (cId === undefined) return true;
      return cId === m.characterId

    }
    this.sse
      .on(MemoryCreated, eventFilter)
      .subscribe(m => this.memories.update(memories =>
        arrayAdd(memories, m, true)))
    this.sse
      .on(MemoryUpdated, eventFilter)
      .subscribe(m => this.memories.update(memories =>
        arrayReplace(memories, m, m2 => m2.id === m.id)))
    this.sse
      .on(MemoryDeleted)
      .subscribe(id => this.memories.update(memories =>
        arrayRemove(memories, m => m.id === id)))
  }

  onAddMemory() {
    const newMemory: Memory = {
      id: NEW_ID,
      worldId: this.worldId(),
      characterId: this.characterId(),
      createdAt: null,
      content: '',
      alwaysInclude: false
    }

    this.newMemories.update(memories => [newMemory, ...memories])
  }

  onSaveMemory(memory: Memory, newMemoryIndex: Nullable<number>) {
    this.memoriesService
      .save(this.worldId(), memory)
      .subscribe(updated => {
        this.memories.update(memories =>
          arrayReplace(memories, updated, m => m.id === memory.id));
        if (typeof newMemoryIndex === 'number') {
          this.onDeleteNewMemory(newMemoryIndex)
        }
      })
  }

  onDeleteNewMemory(atIndex: number) {
    this.newMemories.update(memories =>
      arrayRemove(memories, atIndex))
  }

  onDeleteMemory(id: number) {
    const doDelete = confirm('Are you sure you want to delete this memory?');

    if (doDelete) {
      this.memoriesService
        .delete(this.worldId(), id)
        .subscribe(() => this.memories.update(memories =>
          arrayRemove(memories, m => m.id === id)))
    }
  }
}
