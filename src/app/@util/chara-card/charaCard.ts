export interface CharaCard {
  spec: CharaCardSpecVersion,
  spec_version: string,
}

export type CharaCardSpecVersion = 'chara_card_v1' | 'chara_card_v2' | 'chara_card_v3';
