import {Component} from '@angular/core';

@Component({
  selector: 'app-new-item-card',
  imports: [],
  templateUrl: './new-item-card.html',
  styleUrl: './new-item-card.scss',
  host: {
    '[class.item-card]': 'true',
  }
})
export class NewItemCard {

}
