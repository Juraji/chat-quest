import {
  Component,
  computed,
  effect,
  forwardRef,
  inject,
  input,
  InputSignal,
  linkedSignal,
  Signal,
  signal,
  WritableSignal
} from '@angular/core';
import {ControlValueAccessor, FormsModule, NG_VALUE_ACCESSOR} from '@angular/forms';
import {filter, iif, map, mergeMap, of, toArray} from 'rxjs';
import {DropdownContainer, DropdownMenu, DropdownToggle} from '../dropdown';
import {isNew, NEW_ID} from '@api/common';
import {arrayAdd, arrayMerge} from '@util/array';
import {Tag, Tags} from '@api/characters';

@Component({
  selector: 'app-tags-control',
  templateUrl: './tags-control.html',
  styleUrls: ['./tags-control.scss'],
  imports: [
    FormsModule,
    DropdownContainer,
    DropdownMenu,
    DropdownToggle,
  ],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => TagsControl),
      multi: true
    }
  ]
})
export class TagsControl implements ControlValueAccessor {
  private readonly tags = inject(Tags)

  private onChange: (value: Tag[]) => void = () => null;
  private onTouched: () => void = () => null;

  private readonly availableTags: Signal<Tag[]> = this.tags.all

  readonly tagsInput: InputSignal<Tag[]> = input<Tag[]>([], {alias: 'tags'})
  readonly disabled: InputSignal<boolean> = input(false)

  protected readonly inputText: WritableSignal<string> = signal('')
  protected readonly isDisabled: WritableSignal<boolean> = linkedSignal(() => this.disabled())
  protected readonly currentTags: WritableSignal<Tag[]> = linkedSignal(() => this.tagsInput())
  protected readonly currentTagIds: Signal<number[]> = computed(() => this.currentTags().map(t => t.id))

  protected readonly missingTags: Signal<Tag[]> = computed(() => {
    const availableTags = this.availableTags()
    const currentTagIds = this.currentTagIds()
    return availableTags.filter(t => !currentTagIds.includes(t.id))
  })

  constructor() {
    effect(() => {
      // When inputText changes we are touched
      this.inputText()
      this.onTouched()
    });
  }

  writeValue(obj: Tag[]): void {
    if (!Array.isArray(obj)) {
      throw new Error('writeValue expects an array of tags');
    }

    this.currentTags.set(obj)
  }

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this.onTouched = fn;
  }

  setDisabledState?(isDisabled: boolean): void {
    this.isDisabled.set(isDisabled);
  }

  onAddFromMenu(t: Tag) {
    this.currentTags.update(tags => arrayAdd(tags, t))
    this.onChange(this.currentTags())
  }

  addTag() {
    this.onTouched()
    const inputText = this.inputText().trim()
    if (!inputText) return

    const available = this.availableTags()
    const currentTagIds = this.currentTagIds()

    of(inputText)
      .pipe(
        mergeMap(names => names
          .split(',')
          .map(tag => tag.trim())),
        map((label): Tag => {
          const lc = label.toLowerCase()
          const existing = available.find(t => t.lowercase === lc)
          return !!existing ? existing : {id: NEW_ID, label, lowercase: lc}
        }),
        filter(t => isNew(t) || currentTagIds.includes(t.id)),
        mergeMap(tag => iif(() => isNew(tag), this.tags.save(tag), [tag])),
        toArray()
      )
      .subscribe(newTags => {
        this.currentTags.update(tags => arrayMerge(tags, newTags, (a, b) => a.id === b.id))
        this.onChange(this.currentTags())
        this.inputText.set('')
      })
  }

  removeTag(tagId: number): void {
    this.onTouched()
    this.currentTags.update(tags => tags.filter(t => t.id !== tagId))
    this.onChange(this.currentTags())
  }

  onInputKeyDown(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      event.stopPropagation()
      event.preventDefault()
      this.addTag()
    }
  }
}
