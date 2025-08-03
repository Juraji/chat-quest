CREATE TABLE scenarios
(
  id                  INTEGER PRIMARY KEY AUTOINCREMENT,
  name                VARCHAR(100) NOT NULL,
  description         TEXT         NOT NULL,
  avatar_url          TEXT DEFAULT NULL,
  linked_character_id INTEGER      REFERENCES characters (id) ON DELETE SET NULL
);

INSERT INTO scenarios (name, description, linked_character_id)
VALUES ('Hearthstone Haven',
        'The warm glow of flickering candles and the central fire pit cast dancing shadows across the worn wooden tables as
a bard in the corner plucks his lute, singing tales of old to an audience of elves, dwarves, and humans enjoying
their hard-earned rest.

The pattering rain against shuttered windows creates a soothing rhythm that blends with
murmured conversations about blacksmithing forges and farming harvests over steaming bowls of stew simmering atop
small stoves built into the walls. Smoke from extinguished candles curls upwards, mingling with the aroma of roasted
meats and ale''s bitter tang that lingers in the air—creating a comforting haven where weary travelers can rest their
heads on soft cushions or share laughter over half-eaten meals scattered across tables. Outside, rain continues to fall
steadily as wind occasionally sneaks under gaps in doors bringing raindrop scents indoors while warmth radiates from
the fire pit toward cooler patches near exits keeping patrons cozy despite stormy weather outside.',
        NULL);
