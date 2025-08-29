package instructions

import (
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
)

type InstructionType string

const (
	ChatInstructionType     InstructionType = "CHAT"
	MemoriesInstructionType InstructionType = "MEMORIES"
)

func (i InstructionType) IsValid() bool {
	switch i {
	case ChatInstructionType:
		return true
	case MemoriesInstructionType:
		return true
	default:
		return false
	}
}

type InstructionTemplate struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Type         InstructionType `json:"type"`
	Temperature  *float32        `json:"temperature"`
	SystemPrompt string          `json:"systemPrompt"`
	WorldSetup   string          `json:"worldSetup"`
	Instruction  string          `json:"instruction"`
}

func instructionPromptScanner(scanner database.RowScanner, dest *InstructionTemplate) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Type,
		&dest.Temperature,
		&dest.SystemPrompt,
		&dest.WorldSetup,
		&dest.Instruction,
	)
}

func AllInstructions() ([]InstructionTemplate, error) {
	query := "SELECT * FROM instruction_templates"
	return database.QueryForList(query, nil, instructionPromptScanner)
}

func InstructionById(id int) (*InstructionTemplate, error) {
	query := "SELECT * FROM instruction_templates WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(query, args, instructionPromptScanner)
}

func CreateInstruction(it *InstructionTemplate) error {
	util.NegFloat32PtrToNil(&it.Temperature)

	query := `INSERT INTO instruction_templates (name, type, temperature, system_prompt, world_setup, instruction)
            VALUES (?, ?, ?, ?, ?, ?) RETURNING id`
	args := []any{it.Name, it.Type, it.Temperature, it.SystemPrompt, it.WorldSetup, it.Instruction}

	err := database.InsertRecord(query, args, &it.ID)

	if err == nil {
		InstructionCreatedSignal.EmitBG(it)
	}

	return err
}

func UpdateInstruction(id int, it *InstructionTemplate) error {
	util.NegFloat32PtrToNil(&it.Temperature)

	query := `UPDATE instruction_templates
            SET name = ?,
                type = ?,
                temperature = ?,
                system_prompt = ?,
                world_setup = ?,
                instruction = ?
            WHERE id = ?`
	args := []any{it.Name, it.Type, it.Temperature,
		it.SystemPrompt, it.WorldSetup, it.Instruction, id}

	err := database.UpdateRecord(query, args)

	if err == nil {
		InstructionUpdatedSignal.EmitBG(it)
	}

	return err
}

func DeleteInstruction(id int) error {
	query := "DELETE FROM instruction_templates WHERE id = ?"
	args := []any{id}

	_, err := database.DeleteRecord(query, args)

	if err == nil {
		InstructionDeletedSignal.EmitBG(id)
	}

	return err
}
