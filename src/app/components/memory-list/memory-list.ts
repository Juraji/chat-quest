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
import {booleanSignal, controlValueSignal, formControl} from '@util/ng';
import {SSE} from '@api/sse';
import {ReactiveFormsModule} from '@angular/forms';
import {Characters} from '@api/characters';

type SearchType = "CHARACTER" | "CONTENT"
const PAGE_SIZE = 10

@Component({
  selector: 'memory-list',
  imports: [
    MemoryListItem,
    ReactiveFormsModule
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

  protected readonly addMemoryActive = booleanSignal(false)
  protected readonly newMemoryTpl = computed(() => ({
    id: NEW_ID,
    worldId: this.worldId(),
    characterId: this.characterId(),
    createdAt: null,
    content: '',
    alwaysInclude: false
  }))

  protected readonly memories: WritableSignal<Memory[]> = signal([])

  // Stage 1: Text search
  protected readonly searchControl = formControl('')
  protected readonly searchControlValue = controlValueSignal(this.searchControl)
  protected readonly searchType: WritableSignal<SearchType> = signal("CONTENT")
  protected readonly searchResult: Signal<Memory[]> = computed(() => {
    const val = this.searchControlValue().toLowerCase()
    const memories = this.memories()

    if (val.length === 0) return memories

    if (this.searchType() === "CHARACTER") {
      const characterIds = this.charactersService.all()
        .filter(c => c.name.toLowerCase().includes(val))
        .map(c => c.id)
      return memories.filter(m => !!m.characterId && characterIds.includes(m.characterId))
    } else {
      return memories.filter(m => m.content.toLowerCase().includes(val))
    }
  })

  // Stage 2: Pagination
  protected readonly offsetCount = computed(() => Math.floor(this.searchResult().length / PAGE_SIZE))
  protected readonly currentOffset = signal(0)
  protected readonly paginated = computed(() => {
    const searchResult = this.searchResult()
    const offset = PAGE_SIZE * this.currentOffset()

    return searchResult.slice(offset, offset + PAGE_SIZE)
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

  onToggleSearchType() {
    this.searchType.update(current => {
      switch (current) {
        case "CHARACTER":
          return "CONTENT"
        case "CONTENT":
          return "CHARACTER"
      }
    })
  }

  onPreviousPage() {
    this.currentOffset.update(c => c - 1 || 0)
  }

  onNextPage() {
    this.currentOffset.update(c => c + 1)
  }

  onSaveMemory(memory: Memory) {
    this.memoriesService
      .save(this.worldId(), memory)
      .subscribe(updated => {
        this.memories.update(memories =>
          arrayReplace(memories, updated, m => m.id === memory.id));
        this.addMemoryActive.setFalse()
      })
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
