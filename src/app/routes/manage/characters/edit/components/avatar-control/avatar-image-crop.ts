import {
  Component,
  effect,
  input,
  InputSignal,
  output,
  OutputEmitterRef,
  Signal,
  signal,
  viewChild,
  WritableSignal
} from '@angular/core';
import {AngularCropperjsModule, CropperComponent} from 'angular-cropperjs';
import {readBlobAsDataUrl} from '@util/blobs';

const aspectRatio = 1;
const maxWidth = 350
const maxHeight = 350
const exportType = "image/jpeg";
const exportQuality = 1;

@Component({
  selector: 'app-avatar-image-crop',
  imports: [AngularCropperjsModule],
  templateUrl: './avatar-image-crop.html',
  styleUrl: './avatar-image-crop.scss',
})
export class AvatarImageCrop {
  readonly imageBlob: InputSignal<Blob> = input.required()

  readonly onCropComplete: OutputEmitterRef<Blob> = output()
  readonly onCropCanceled: OutputEmitterRef<void> = output()

  protected readonly b64ImageUrl: WritableSignal<string> = signal('')
  protected readonly cropperCmp: Signal<CropperComponent | undefined> = viewChild('angularCropper')

  readonly cropperOptions = {
    aspectRatio,
    dragMode: 'move',
    background: true,
    movable: true,
    rotatable: true,
    scalable: true,
    zoomable: false,
    viewMode: 1,
  };

  constructor() {
    effect(() => {
      const file = this.imageBlob()
      readBlobAsDataUrl(file).then(data => this.b64ImageUrl.set(data))
    });
  }

  onAcceptCrop() {
    this.cropperCmp()?.cropper
      .getCroppedCanvas({maxWidth, maxHeight})
      .toBlob(b => this.onCropComplete.emit(b!), exportType, exportQuality);
  }

  onCancel() {
    this.onCropCanceled.emit()
    this.b64ImageUrl.set('')
  }
}
