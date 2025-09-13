package processing

import (
	"sync"

	"github.com/pkg/errors"
	prov "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/util"
	c "juraji.nl/chat-quest/model/characters"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	m "juraji.nl/chat-quest/model/memories"
	p "juraji.nl/chat-quest/model/preferences"
)

// SparseTemplateCharacter

type SparseTemplateCharacter interface {
	ID() int
	CharacterName() string
}

type sparseTemplateCharacterImpl struct {
	id   int
	name string
}

func (t *sparseTemplateCharacterImpl) ID() int {
	return t.id
}
func (t *sparseTemplateCharacterImpl) CharacterName() string {
	return t.name
}

func NewSparseTemplateCharacter(char *c.Character) SparseTemplateCharacter {
	return &sparseTemplateCharacterImpl{
		id:   char.ID,
		name: char.Name,
	}
}

// TemplateCharacter

type TemplateCharacter interface {
	ID() int
	Name() string
	Appearance() (string, error)
	Personality() (string, error)
	History() (string, error)
	DialogueExamples() ([]string, error)
	Memories() ([]string, error)
}

type templateCharacterImpl struct {
	id               int
	name             string
	appearance       func() (string, error)
	personality      func() (string, error)
	history          func() (string, error)
	dialogueExamples func() ([]string, error)
	memories         func() ([]string, error)
}

func (t *templateCharacterImpl) ID() int {
	return t.id
}
func (t *templateCharacterImpl) Name() string {
	return t.name
}
func (t *templateCharacterImpl) Appearance() (string, error) {
	return t.appearance()
}
func (t *templateCharacterImpl) Personality() (string, error) {
	return t.personality()
}
func (t *templateCharacterImpl) History() (string, error) {
	return t.history()
}
func (t *templateCharacterImpl) DialogueExamples() ([]string, error) {
	return t.dialogueExamples()
}
func (t *templateCharacterImpl) Memories() ([]string, error) {
	return t.memories()
}

// NewTemplateCharacter creates a new character object for use in go templates.
// Note that the chatHistory must contain the full history, including the latest user message if applicable.
func NewTemplateCharacter(
	char *c.Character,
	prefs *p.Preferences,
	session *cs.ChatSession,
	chatHistory []cs.ChatMessage,
) TemplateCharacter {
	return &templateCharacterImpl{
		id:   char.ID,
		name: char.Name,

		appearance: sync.OnceValues(func() (string, error) {
			if char.Appearance == nil {
				return "", nil
			}
			charTpl := NewSparseTemplateCharacter(char)
			template, err := util.ParseAndApplyTextTemplate(*char.Appearance, charTpl)
			return template, errors.Wrapf(err, "failed to parse char appearance template for character ID %d", char.ID)
		}),
		personality: sync.OnceValues(func() (string, error) {
			if char.Personality == nil {
				return "", nil
			}
			charTpl := NewSparseTemplateCharacter(char)
			template, err := util.ParseAndApplyTextTemplate(*char.Personality, charTpl)
			return template, errors.Wrapf(err, "failed to parse char personality template for character ID %d", char.ID)
		}),
		history: sync.OnceValues(func() (string, error) {
			if char.History == nil {
				return "", nil
			}
			charTpl := NewSparseTemplateCharacter(char)
			template, err := util.ParseAndApplyTextTemplate(*char.History, charTpl)
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
				template, err := util.ParseAndApplyTextTemplate(example, charTpl)
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

			memories, err := m.GetMemoriesByWorldAndCharacterIdWithEmbeddings(
				session.WorldID, char.ID, *prefs.EmbeddingModelId)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get memories for character ID %d", char.ID)
			}

			// Short circuit: No chat history, just get "AlwaysInclude" memories and return.
			if len(chatHistory) == 0 {
				var staticMemories []string
				for _, mem := range memories {
					if mem.AlwaysInclude {
						staticMemories = append(staticMemories, mem.Content)
					}
				}

				return staticMemories, nil
			}

			// Determine subject, based on the last n message and the trigger message.
			var subject string
			if len(chatHistory) > 0 {
				end := len(chatHistory)
				start := end - prefs.MemoryWindowSize
				if start < 0 {
					start = 0
				}

				for i := start; i < end; i++ {
					subject += chatHistory[i].Content + " "
				}
			}

			// Embed subject text
			embeddingModelInst, err := prov.GetLlmModelInstanceById(*prefs.EmbeddingModelId)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get embedding model while processing memories for character ID %d", char.ID)
			}

			subjectEmbeddings, err := prov.GenerateEmbeddings(embeddingModelInst, subject, true)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to embed subject while processing memories for character ID %d", char.ID)
			}

			// Process memories in batches using goroutines
			const workerCount = 16
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
						}

						similarity := subjectEmbeddings.CosineSimilarity(memory.Embedding)
						if similarity >= minP {
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
