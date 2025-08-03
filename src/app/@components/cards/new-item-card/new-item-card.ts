import {Component} from '@angular/core';

@Component({
  selector: 'app-new-item-card',
  imports: [],
  templateUrl: './new-item-card.html',
  styleUrl: './new-item-card.scss',
  host: {
    '[class.chat-quest-card]': 'true',
  }
})
export class NewItemCard {

}
