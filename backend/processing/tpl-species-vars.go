package processing

import (
	"juraji.nl/chat-quest/model/species"
)

type TemplateSpecies interface {
	ID() int
	Name() string
	Description() string
}

type templateSpeciesImpl struct {
	id          int
	name        string
	description string
}

func (t *templateSpeciesImpl) ID() int             { return t.id }
func (t *templateSpeciesImpl) Name() string        { return t.name }
func (t *templateSpeciesImpl) Description() string { return t.description }

func NewTemplateSpecies(species *species.Species) TemplateSpecies {
	return &templateSpeciesImpl{
		id:          species.ID,
		name:        species.Name,
		description: species.Description,
	}
}
