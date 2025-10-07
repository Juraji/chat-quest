package instructions

import (
	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/database"
	p "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/util"
)

type InstructionType string

const (
	ChatInstruction     InstructionType = "CHAT"
	MemoriesInstruction InstructionType = "MEMORIES"
)

func (i InstructionType) IsValid() bool {
	switch i {
	case ChatInstruction:
		return true
	case MemoriesInstruction:
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
	ReasoningPrefix   string `json:"reasoningPrefix"`
	ReasoningSuffix   string `json:"reasoningSuffix"`
	CharacterIdPrefix string `json:"characterIdPrefix"`
	CharacterIdSuffix string `json:"characterIdSuffix"`

	// Prompt Templates
	SystemPrompt string `json:"systemPrompt"`
	WorldSetup   string `json:"worldSetup"`
	Instruction  string `json:"instruction"`
}

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

func (i *Instruction) ApplyTemplates(variables any) error {
	fields := []*string{
		&i.SystemPrompt,
		&i.WorldSetup,
		&i.Instruction,
	}

	for _, fieldPtr := range fields {
		result, err := util.ParseAndApplyTextTemplate(*fieldPtr, variables)
		if err != nil {
			return errors.Wrap(err, "Error creating template for instruction template")
		}

		*fieldPtr = result
	}

	return nil
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
	zf := util.ZeroFloat32ToNil
	zi := util.ZeroIntToNil

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
            VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?.?,?) RETURNING id`
	args := []any{
		inst.Name,
		inst.Type,
		zf(inst.Temperature),
		zi(inst.MaxTokens),
		zf(inst.TopP),
		zf(inst.PresencePenalty),
		zf(inst.FrequencyPenalty),
		inst.Stream,
		es(inst.StopSequences),
		inst.IncludeReasoning,
		inst.ReasoningPrefix,
		inst.ReasoningSuffix,
		inst.CharacterIdPrefix,
		inst.CharacterIdSuffix,
		inst.SystemPrompt,
		inst.WorldSetup,
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
	zf := util.ZeroFloat32ToNil
	zi := util.ZeroIntToNil

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
		zf(inst.Temperature),
		zi(inst.MaxTokens),
		zf(inst.TopP),
		zf(inst.PresencePenalty),
		zf(inst.FrequencyPenalty),
		inst.Stream,
		es(inst.StopSequences),
		inst.IncludeReasoning,
		inst.ReasoningPrefix,
		inst.ReasoningSuffix,
		inst.CharacterIdPrefix,
		inst.CharacterIdSuffix,
		inst.SystemPrompt,
		inst.WorldSetup,
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
