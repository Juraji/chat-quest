import {CharaCard} from './charaCard';

/** @see https://github.com/malfoyslastname/character-card-spec-v2/blob/main/spec_v1.md */
export interface CharaCardV1 extends CharaCard {
  name: string
  description: string
  personality: string
  scenario: string
  first_mes: string
  mes_example: string
}
