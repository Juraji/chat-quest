package processing

import (
	"sync"

	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/util"
	c "juraji.nl/chat-quest/model/characters"
	sp "juraji.nl/chat-quest/model/species"
	w "juraji.nl/chat-quest/model/worlds"
)

type CharacterBuilderVars interface {
	CharacterName() string
	Age() *int
	Pronouns() string
	Species() (string, error)
	CurrentAppearance() (string, error)
	CurrentPersonality() (string, error)
	CurrentHistory() (string, error)
	World() (string, error)
	UserInput() string
}

type characterBuilderVarsImpl struct {
	name        string
	age         *int
	pronouns    *string
	species     func() (string, error)
	appearance  func() (string, error)
	personality func() (string, error)
	history     func() (string, error)
	world       func() (string, error)
	userInput   string
}

func (c *characterBuilderVarsImpl) CharacterName() string               { return c.name }
func (c *characterBuilderVarsImpl) Age() *int                           { return c.age }
func (c *characterBuilderVarsImpl) Pronouns() string                    { return util.StrPtrOrDefault(c.pronouns, "") }
func (c *characterBuilderVarsImpl) Species() (string, error)            { return c.species() }
func (c *characterBuilderVarsImpl) CurrentAppearance() (string, error)  { return c.appearance() }
func (c *characterBuilderVarsImpl) CurrentPersonality() (string, error) { return c.personality() }
func (c *characterBuilderVarsImpl) CurrentHistory() (string, error)     { return c.history() }
func (c *characterBuilderVarsImpl) World() (string, error)              { return c.world() }
func (c *characterBuilderVarsImpl) UserInput() string                   { return c.userInput }

func NewCharacterBuilderVars(char *c.Character, world *w.World, userInput string) CharacterBuilderVars {
	return &characterBuilderVarsImpl{
		name:     char.Name,
		age:      char.Age,
		pronouns: char.Pronouns,

		species: sync.OnceValues(func() (string, error) {
			if char.SpeciesID == nil {
				return "", nil
			}
			species, err := sp.SpeciesByID(*char.SpeciesID)
			if err != nil || species == nil {
				return "", err
			}

			return species.Description, nil
		}),
		appearance: sync.OnceValues(func() (string, error) {
			if char.Appearance == nil {
				return "", nil
			}
			charTpl := NewSparseTemplateCharacter(char)
			template, err := util.ParseAndApplyTextTemplate("Appeance for "+char.Name, *char.Appearance, charTpl)
			return template, errors.Wrapf(err, "failed to parse char appearance template for character ID %d", char.ID)
		}),
		personality: sync.OnceValues(func() (string, error) {
			if char.Personality == nil {
				return "", nil
			}
			charTpl := NewSparseTemplateCharacter(char)
			template, err := util.ParseAndApplyTextTemplate("Personality for "+char.Name, *char.Personality, charTpl)
			return template, errors.Wrapf(err, "failed to parse char personality template for character ID %d", char.ID)
		}),
		history: sync.OnceValues(func() (string, error) {
			if char.History == nil {
				return "", nil
			}
			charTpl := NewSparseTemplateCharacter(char)
			template, err := util.ParseAndApplyTextTemplate("History for "+char.Name, *char.History, charTpl)
			return template, errors.Wrapf(err, "failed to parse char history template for character ID %d", char.ID)
		}),
		world: sync.OnceValues(func() (string, error) {
			if world == nil || world.Description == nil {
				return "", nil
			}

			return *world.Description, nil
		}),
		userInput: userInput,
	}
}
