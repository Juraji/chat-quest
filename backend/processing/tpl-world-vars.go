package processing

import (
	"sync"

	"juraji.nl/chat-quest/model/characters"
	"juraji.nl/chat-quest/model/chat-sessions"
)

type GreetingVars interface {
	CharacterName() (string, error)
	PersonaName() (string, error)
}

type greetingVarsImpl struct {
	characterName func() (string, error)
	personaName   func() (string, error)
}

func (w *greetingVarsImpl) CharacterName() (string, error) {
	return w.characterName()
}
func (w *greetingVarsImpl) PersonaName() (string, error) {
	return w.personaName()
}

func NewGreetingVars(sessionId int, char *characters.Character) GreetingVars {
	return &greetingVarsImpl{
		characterName: func() (string, error) {
			return char.Name, nil
		},
		personaName: sync.OnceValues(func() (string, error) {
			session, err := chat_sessions.GetById(sessionId)
			if err != nil {
				return "", err
			}

			if session.PersonaID == nil {
				return "User", nil
			}

			character, err := characters.CharacterById(*session.PersonaID)
			if err != nil {
				return "", err
			}

			return character.Name, nil
		}),
	}
}
