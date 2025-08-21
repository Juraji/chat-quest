package instructions

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
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

func AllInstructions() ([]InstructionTemplate, bool) {
	query := "SELECT * FROM instruction_templates"
	list, err := database.QueryForList(query, nil, instructionPromptScanner)
	if err != nil {
		log.Get().Error("Error fetching instructions", zap.Error(err))
		return nil, false
	}

	return list, true
}

func InstructionById(id int) (*InstructionTemplate, bool) {
	query := "SELECT * FROM instruction_templates WHERE id = ?"
	args := []any{id}
	instruction, err := database.QueryForRecord(query, args, instructionPromptScanner)
	if err != nil {
		log.Get().Error("Error fetching instruction",
			zap.Int("id", id), zap.Error(err))
		return nil, false
	}

	return instruction, true
}

func CreateInstruction(it *InstructionTemplate) bool {
	util.NegFloat32PtrToNil(&it.Temperature)

	query := `INSERT INTO instruction_templates (name, type, temperature, system_prompt, world_setup, instruction)
            VALUES (?, ?, ?, ?, ?, ?) RETURNING id`
	args := []any{it.Name, it.Type, it.Temperature, it.SystemPrompt, it.WorldSetup, it.Instruction}

	err := database.InsertRecord(query, args, &it.ID)
	if err != nil {
		log.Get().Error("Error inserting instruction", zap.Error(err))
		return false
	}

	InstructionCreatedSignal.EmitBG(it)
	return true
}

func UpdateInstruction(id int, it *InstructionTemplate) bool {
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
	if err != nil {
		log.Get().Error("Error updating instruction",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	InstructionUpdatedSignal.EmitBG(it)
	return true
}

func DeleteInstruction(id int) bool {
	query := "DELETE FROM instruction_templates WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(query, args)
	if err != nil {
		log.Get().Error("Error deleting instruction",
			zap.Int("id", id), zap.Error(err))
		return false
	}

	InstructionDeletedSignal.EmitBG(id)
	return true
}
