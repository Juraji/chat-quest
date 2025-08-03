import {Component, input, InputSignal, output, OutputEmitterRef, Signal, viewChild} from '@angular/core';
import {AngularCropperjsModule, CropperComponent} from 'angular-cropperjs';

const CROPPER_OPTIONS = {
  aspectRatio: 1,
  dragMode: 'move',
  background: true,
  movable: true,
  rotatable: true,
  scalable: true,
  zoomable: false,
  viewMode: 1,
}
const maxWidth = 500
const maxHeight = 500
const exportType = "image/jpeg";
const exportQuality = 1;

@Component({
  selector: 'app-avatar-image-crop',
  imports: [AngularCropperjsModule],
  templateUrl: './avatar-image-crop.html',
  styleUrl: './avatar-image-crop.scss',
})
export class AvatarImageCrop {
  readonly imageDataUrl: InputSignal<string> = input.required()

  readonly onCropComplete: OutputEmitterRef<string> = output()
  readonly onCropCanceled: OutputEmitterRef<void> = output()

  protected readonly cropperCmp: Signal<CropperComponent | undefined> = viewChild('angularCropper')

  readonly cropperOptions = CROPPER_OPTIONS

  onAcceptCrop() {
    const dataUrl = this.cropperCmp()!.cropper
      .getCroppedCanvas({maxWidth, maxHeight})
      .toDataURL(exportType, exportQuality)
    this.onCropComplete.emit(dataUrl)
  }

  onCancel() {
    this.onCropCanceled.emit()
  }
}
