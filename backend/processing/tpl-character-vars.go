package processing

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	prov "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/util"
	c "juraji.nl/chat-quest/model/characters"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	m "juraji.nl/chat-quest/model/memories"
	p "juraji.nl/chat-quest/model/preferences"
	sp "juraji.nl/chat-quest/model/species"
)

// SparseTemplateCharacter

type SparseTemplateCharacter interface {
	ID() int
	CharacterName() string
	Age() *int
	Pronouns() string
	Species() (string, error)
}

type sparseTemplateCharacterImpl struct {
	id       int
	name     string
	age      *int
	pronouns *string
	species  func() (string, error)
}

func (t *sparseTemplateCharacterImpl) ID() int                  { return t.id }
func (t *sparseTemplateCharacterImpl) CharacterName() string    { return t.name }
func (t *sparseTemplateCharacterImpl) Age() *int                { return t.age }
func (t *sparseTemplateCharacterImpl) Pronouns() string         { return util.StrPtrOrDefault(t.pronouns, "") }
func (t *sparseTemplateCharacterImpl) Species() (string, error) { return t.species() }

func NewSparseTemplateCharacter(char *c.Character) SparseTemplateCharacter {
	return &sparseTemplateCharacterImpl{
		id:       char.ID,
		name:     char.Name,
		age:      char.Age,
		pronouns: char.Pronouns,
		species: sync.OnceValues(func() (string, error) {
			if char.SpeciesID == nil {
				return "", nil
			}
			species, err := sp.SpeciesByID(*char.SpeciesID)
			if err != nil {
				return "", err
			}
			if species == nil {
				return "", nil
			}

			return species.Name, nil
		}),
	}
}

type TemplateCharacter interface {
	ID() int
	Name() string
	Appearance() (string, error)
	Personality() (string, error)
	History() (string, error)
	DialogueExamples() ([]string, error)
	Memories() ([]string, error)
	Age() *int
	Pronouns() string
	Species() (string, error)
}

type templateCharacterImpl struct {
	id               int
	name             string
	age              *int
	pronouns         *string
	species          func() (string, error)
	appearance       func() (string, error)
	personality      func() (string, error)
	history          func() (string, error)
	dialogueExamples func() ([]string, error)
	memories         func() ([]string, error)
}

func (t templateCharacterImpl) ID() int                             { return t.id }
func (t templateCharacterImpl) Name() string                        { return t.name }
func (t templateCharacterImpl) Appearance() (string, error)         { return t.appearance() }
func (t templateCharacterImpl) Personality() (string, error)        { return t.personality() }
func (t templateCharacterImpl) History() (string, error)            { return t.history() }
func (t templateCharacterImpl) DialogueExamples() ([]string, error) { return t.dialogueExamples() }
func (t templateCharacterImpl) Memories() ([]string, error)         { return t.memories() }
func (t templateCharacterImpl) Age() *int                           { return t.age }
func (t templateCharacterImpl) Pronouns() string                    { return util.StrPtrOrDefault(t.pronouns, "") }
func (t templateCharacterImpl) Species() (string, error)            { return t.species() }

// NewTemplateCharacter creates a new character object for use in go templates.
func NewTemplateCharacter(
	char *c.Character,
	prefs *p.Preferences,
	session *cs.ChatSession,
) TemplateCharacter {
	logger := log.Get().With(
		zap.Int("characterId", char.ID))

	return &templateCharacterImpl{
		id:       char.ID,
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

			return species.Name, nil
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
		dialogueExamples: sync.OnceValues(func() ([]string, error) {
			examples, err := c.DialogueExamplesByCharacterId(char.ID)
			if err != nil {
				return nil, err
			}
			if len(examples) == 0 {
				return nil, nil
			}

			charTpl := NewSparseTemplateCharacter(char)
			for i, example := range examples {
				template, err := util.ParseAndApplyTextTemplate("Dialogue Example for "+char.Name, example, charTpl)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to parse char dialogue example template for character ID %d: %s", char.ID, example)
				}
				examples[i] = template
			}
			return examples, nil
		}),
		memories: sync.OnceValues(func() ([]string, error) {
			if session == nil {
				// Not in a chat session context
				return nil, nil
			}
			if !session.UseMemories {
				// Memories are disabled
				return nil, nil
			}

			chatHistory, err := cs.GetTailChatMessages(session.ID, prefs.MemoryIncludeChatSize)
			if err != nil {
				return nil, errors.Wrap(err, "failed to fetch chat messages for memory scan")
			}

			memories, err := m.GetMemoriesByWorldAndCharacterIdWithEmbeddings(
				session.WorldID, char.ID, *prefs.EmbeddingModelId)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get memories for character ID %d", char.ID)
			}

			// Short circuit: Character has no memories
			if len(memories) == 0 {
				return nil, nil
			}

			// Short circuit: No chat history, just get "AlwaysInclude" memories and return.
			if len(chatHistory) == 0 && session.ChatNotes == nil {
				var staticMemories []string
				for _, mem := range memories {
					if mem.AlwaysInclude {
						staticMemories = append(staticMemories, mem.Content)
					}
				}

				return staticMemories, nil
			}

			// The subject is made up of
			// - Current participant names.
			// - Optionally the Chat notes.
			// - The messages to scan.
			var subjectBuffer strings.Builder
			{
				participants, err := cs.GetAllParticipantsAsCharacters(session.ID)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get session participants")
				}
				for _, participant := range participants {
					subjectBuffer.WriteString(participant.Name)
					subjectBuffer.WriteRune(' ')
				}

				if prefs.MemoryIncludeChatNotes && session.ChatNotes != nil {
					subjectBuffer.WriteString(*session.ChatNotes)
					subjectBuffer.WriteRune(' ')
				}

				if len(chatHistory) > 0 {
					for _, msg := range chatHistory {
						subjectBuffer.WriteString(msg.Content)
						subjectBuffer.WriteRune(' ')
					}
				}
			}

			// Embed subject text
			embeddingModelInst, err := prov.GetLlmModelInstanceById(*prefs.EmbeddingModelId)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get embedding model while processing memories for character ID %d", char.ID)
			}

			subjectEmbeddings, err := prov.GenerateEmbeddings(embeddingModelInst, subjectBuffer.String())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to embed subject while processing memories for character ID %d", char.ID)
			}

			// Process memories in batches using goroutines
			const workerCount = 8
			const minChunkSize = 10
			minP := prefs.MemoryMinP

			chunkSize := (len(memories) + workerCount - 1) / workerCount
			if chunkSize < minChunkSize {
				chunkSize = minChunkSize
			}

			relevantMemoriesChan := make(chan string, len(memories))
			var wg sync.WaitGroup

			for start := 0; start < len(memories); start = start + chunkSize {
				end := start + chunkSize
				if end > len(memories) {
					end = len(memories)
				}

				wg.Add(1)
				go func(start, end int) {
					defer wg.Done()

					for _, memory := range memories[start:end] {
						if memory.AlwaysInclude {
							relevantMemoriesChan <- memory.Content
							continue
						}

						similarity := subjectEmbeddings.CosineSimilarity(memory.Embedding)
						if similarity >= minP {
							logger.Debug("Including memory",
								zap.Float64("similarity", similarity),
								zap.String("memory", memory.Content))
							relevantMemoriesChan <- memory.Content
						}
					}
				}(start, end)
			}

			// Wait for all workers to finish
			go func() {
				wg.Wait()
				close(relevantMemoriesChan)
			}()

			var relevantMemories []string
			for memory := range relevantMemoriesChan {
				relevantMemories = append(relevantMemories, memory)
			}

			return relevantMemories, nil
		}),
	}
}
