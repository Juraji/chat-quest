package instructions

import (
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/database"
	p "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/util"
)

type InstructionType string

const (
	ChatInstruction     InstructionType = "CHAT"
	MemoriesInstruction InstructionType = "MEMORIES"
	TitleGeneration     InstructionType = "TITLE_GENERATION"
	CharacterExport     InstructionType = "CHARACTER_EXPORT"
	CharacterBuilder    InstructionType = "CHARACTER_BUILDER"
)

func (i InstructionType) IsValid() bool {
	switch i {
	case ChatInstruction,
		MemoriesInstruction,
		TitleGeneration,
		CharacterExport,
		CharacterBuilder:
		return true
	default:
		return false
	}
}

type Instruction struct {
	ID   int             `json:"id"`
	Name string          `json:"name"`
	Type InstructionType `json:"type"`

	// Model Settings
	Temperature      float32 `json:"temperature"`
	MaxTokens        int     `json:"maxTokens"`
	TopP             float32 `json:"topP"`
	PresencePenalty  float32 `json:"presencePenalty"`
	FrequencyPenalty float32 `json:"frequencyPenalty"`
	Stream           bool    `json:"stream"`
	StopSequences    *string `json:"stopSequences"`
	IncludeReasoning bool    `json:"includeReasoning"`

	// Parsing
	AllowMultiCharacterResponses bool   `json:"allowMultiCharacterResponses"`
	EnableReasoningParsing       bool   `json:"enableReasoningParsing"`
	ReasoningPrefix              string `json:"reasoningPrefix"`
	ReasoningSuffix              string `json:"reasoningSuffix"`
	EnableCharacterMarkers       bool   `json:"enableCharacterMarkers"`
	CharacterIdPrefix            string `json:"characterIdPrefix"`
	CharacterIdSuffix            string `json:"characterIdSuffix"`

	// Prompt Templates
	SystemPrompt *string `json:"systemPrompt"`
	WorldSetup   *string `json:"worldSetup"`
	Instruction  string  `json:"instruction"`
}

// AsLlmParameters converts the instruction's model settings into an LlmParameters struct.
func (i *Instruction) AsLlmParameters() p.LlmParameters {
	return p.LlmParameters{
		MaxTokens:        i.MaxTokens,
		Temperature:      i.Temperature,
		TopP:             i.TopP,
		PresencePenalty:  i.PresencePenalty,
		FrequencyPenalty: i.FrequencyPenalty,
		Stream:           i.Stream,
		StopSequences:    i.StopSequences,
	}
}

// ApplyTemplates processes all string fields in the Instruction that support templating using provided variables.
// It applies text templates to SystemPrompt, WorldSetup, and Instruction fields if they are non-nil,
// replacing placeholders with actual values from the given variables map. Returns an error if template parsing or execution fails.
func (i *Instruction) ApplyTemplates(variables any) error {
	fields := []*string{
		i.SystemPrompt,
		i.WorldSetup,
		&i.Instruction,
	}

	for _, fieldPtr := range fields {
		if fieldPtr == nil {
			continue
		}

		result, err := util.ParseAndApplyTextTemplate(i.Name, *fieldPtr, variables)
		if err != nil {
			return errors.Wrap(err, "Error creating template for instruction template")
		}

		*fieldPtr = strings.TrimSpace(result)
	}

	return nil
}

// CharacterMarkers extracts the first Unicode rune from `CharacterIdPrefix` and returns it along with the full prefix and suffix strings.
// Returns a zero-width space rune if `CharacterIdPrefix` is empty or invalid UTF-8.
func (i *Instruction) CharacterMarkers() (rune, string, string) {
	if !(i.AllowMultiCharacterResponses) {
		return 0, "", ""
	}

	initial, _ := utf8.DecodeRuneInString(i.CharacterIdPrefix)
	return initial, i.CharacterIdPrefix, i.CharacterIdSuffix
}

// ReasoningMarkers Returns the first Unicode rune from `ReasoningPrefix`, along with the full reasoning prefix and suffix strings.
// If `ReasoningPrefix` is empty or invalid UTF-8, returns a zero-width space rune as fallback.
func (i *Instruction) ReasoningMarkers() (rune, string, string) {
	if !i.EnableReasoningParsing {
		return 0, "", ""
	}

	initial, _ := utf8.DecodeRuneInString(i.ReasoningPrefix)
	return initial, i.ReasoningPrefix, i.ReasoningSuffix
}

func instructionPromptScanner(scanner database.RowScanner, dest *Instruction) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Type,
		&dest.Temperature,
		&dest.MaxTokens,
		&dest.TopP,
		&dest.PresencePenalty,
		&dest.FrequencyPenalty,
		&dest.Stream,
		&dest.StopSequences,
		&dest.IncludeReasoning,
		&dest.AllowMultiCharacterResponses,
		&dest.EnableReasoningParsing,
		&dest.ReasoningPrefix,
		&dest.ReasoningSuffix,
		&dest.EnableCharacterMarkers,
		&dest.CharacterIdPrefix,
		&dest.CharacterIdSuffix,
		&dest.SystemPrompt,
		&dest.WorldSetup,
		&dest.Instruction,
	)
}

func AllInstructions() ([]Instruction, error) {
	query := "SELECT * FROM instructions"
	return database.QueryForList(query, nil, instructionPromptScanner)
}

func InstructionById(id int) (*Instruction, error) {
	query := "SELECT * FROM instructions WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(query, args, instructionPromptScanner)
}

func CreateInstruction(inst *Instruction) error {
	es := util.EmptyStrToNil

	query := `INSERT INTO instructions (
                          name,
                          type,
                          temperature,
                          max_tokens,
                          top_p,
                          presence_penalty,
                          frequency_penalty,
                          stream,
                          stop_sequences,
                          include_reasoning,
                          allow_multi_character_responses,
                          enable_reasoning_parsing,
                          reasoning_prefix,
                          reasoning_suffix,
                          enable_character_markers,
                          character_id_prefix,
                          character_id_suffix,
                          system_prompt,
                          world_setup,
                          instruction)
            VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) RETURNING id`
	args := []any{
		inst.Name,
		inst.Type,
		inst.Temperature,
		inst.MaxTokens,
		inst.TopP,
		inst.PresencePenalty,
		inst.FrequencyPenalty,
		inst.Stream,
		es(inst.StopSequences),
		inst.IncludeReasoning,
		inst.AllowMultiCharacterResponses,
		inst.EnableReasoningParsing,
		inst.ReasoningPrefix,
		inst.ReasoningSuffix,
		inst.EnableCharacterMarkers,
		inst.CharacterIdPrefix,
		inst.CharacterIdSuffix,
		es(inst.SystemPrompt),
		es(inst.WorldSetup),
		inst.Instruction,
	}

	err := database.InsertRecord(query, args, &inst.ID)

	if err == nil {
		InstructionCreatedSignal.EmitBG(inst)
	}

	return err
}

func UpdateInstruction(id int, inst *Instruction) error {
	es := util.EmptyStrToNil

	query := `UPDATE instructions
            SET name = ?,
                type = ?,
                temperature = ?,
                max_tokens = ?,
                top_p = ?,
                presence_penalty = ?,
                frequency_penalty = ?,
                stream = ?,
                stop_sequences = ?,
                include_reasoning = ?,
                allow_multi_character_responses = ?,
                enable_reasoning_parsing = ?,
                reasoning_prefix = ?,
                reasoning_suffix = ?,
                enable_character_markers = ?,
                character_id_prefix = ?,
                character_id_suffix = ?,
                system_prompt = ?,
                world_setup = ?,
                instruction = ?
            WHERE id = ?`
	args := []any{
		inst.Name,
		inst.Type,
		inst.Temperature,
		inst.MaxTokens,
		inst.TopP,
		inst.PresencePenalty,
		inst.FrequencyPenalty,
		inst.Stream,
		es(inst.StopSequences),
		inst.IncludeReasoning,
		inst.AllowMultiCharacterResponses,
		inst.EnableReasoningParsing,
		inst.ReasoningPrefix,
		inst.ReasoningSuffix,
		inst.EnableCharacterMarkers,
		inst.CharacterIdPrefix,
		inst.CharacterIdSuffix,
		es(inst.SystemPrompt),
		es(inst.WorldSetup),
		inst.Instruction,
		id,
	}

	err := database.UpdateRecord(query, args)

	if err == nil {
		InstructionUpdatedSignal.EmitBG(inst)
	}

	return err
}

func DeleteInstruction(id int) error {
	query := "DELETE FROM instructions WHERE id = ?"
	args := []any{id}

	_, err := database.DeleteRecord(query, args)

	if err == nil {
		InstructionDeletedSignal.EmitBG(id)
	}

	return err
}
