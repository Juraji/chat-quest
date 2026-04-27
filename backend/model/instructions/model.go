package instructions

import (
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
)

func (i InstructionType) IsValid() bool {
	switch i {
	case ChatInstruction,
		MemoriesInstruction,
		TitleGeneration,
		CharacterExport:
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
	ReasoningPrefix   *string `json:"reasoningPrefix"`
	ReasoningSuffix   *string `json:"reasoningSuffix"`
	CharacterIdPrefix *string `json:"characterIdPrefix"`
	CharacterIdSuffix *string `json:"characterIdSuffix"`

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

		result, err := util.ParseAndApplyTextTemplate(*fieldPtr, variables)
		if err != nil {
			return errors.Wrap(err, "Error creating template for instruction template")
		}

		*fieldPtr = result
	}

	return nil
}

// CharacterMarkerEnabled checks if character markers are enabled for this instruction by verifying
// that both the prefix and suffix fields for character IDs are non-nil.
func (i *Instruction) CharacterMarkerEnabled() bool {
	return i.CharacterIdPrefix != nil && i.CharacterIdSuffix != nil
}

// CharacterMarkers retrieves the character identifier markers if enabled, returning the first rune of the prefix,
// the full prefix string, and the suffix string for use in text formatting. If markers are disabled,
// all returned values will be zero-value (empty string or Unicode replacement character).
func (i *Instruction) CharacterMarkers() (rune, string, string) {
	var initial rune
	var prefix string
	var suffix string

	if i.CharacterMarkerEnabled() {
		initial, _ = utf8.DecodeRuneInString(*i.CharacterIdPrefix)
		prefix = *i.CharacterIdPrefix
		suffix = *i.CharacterIdSuffix
	}

	return initial, prefix, suffix
}

// ReasoningMarkerEnabled checks whether reasoning markers are enabled for this instruction by verifying
// that both the prefix and suffix fields for reasoning content are non-nil.
func (i *Instruction) ReasoningMarkerEnabled() bool {
	return i.ReasoningPrefix != nil && i.ReasoningSuffix != nil
}

// ReasoningMarkers returns the reasoning content markers if enabled, extracting the first rune of the prefix,
// the full prefix string, and the suffix string for use in text formatting. If disabled or fields are nil,
// all returned values will be zero-value (empty string or Unicode replacement character).
func (i *Instruction) ReasoningMarkers() (rune, string, string) {
	var initial rune
	var prefix string
	var suffix string

	if i.ReasoningMarkerEnabled() {
		initial, _ = utf8.DecodeRuneInString(*i.ReasoningPrefix)
		prefix = *i.ReasoningPrefix
		suffix = *i.ReasoningSuffix
	}

	return initial, prefix, suffix
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
		&dest.ReasoningPrefix,
		&dest.ReasoningSuffix,
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
                          reasoning_prefix,
                          reasoning_suffix,
                          character_id_prefix,
                          character_id_suffix,
                          system_prompt,
                          world_setup,
                          instruction)
            VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) RETURNING id`
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
		es(inst.ReasoningPrefix),
		es(inst.ReasoningSuffix),
		es(inst.CharacterIdPrefix),
		es(inst.CharacterIdSuffix),
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
                reasoning_prefix = ?,
                reasoning_suffix = ?,
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
		es(inst.ReasoningPrefix),
		es(inst.ReasoningSuffix),
		es(inst.CharacterIdPrefix),
		es(inst.CharacterIdSuffix),
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
