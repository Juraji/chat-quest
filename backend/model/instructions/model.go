package instructions

import (
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
)

type InstructionTemplate struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Temperature  *float32 `json:"temperature"`
	SystemPrompt string   `json:"systemPrompt"`
	Instruction  string   `json:"instruction"`
}

func instructionPromptScanner(scanner database.RowScanner, dest *InstructionTemplate) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Type,
		&dest.Temperature,
		&dest.SystemPrompt,
		&dest.Instruction,
	)
}

func AllInstructionPrompts() ([]InstructionTemplate, error) {
	query := "SELECT * FROM instruction_templates"
	return database.QueryForList(database.GetDB(), query, nil, instructionPromptScanner)
}

func InstructionPromptById(id int) (*InstructionTemplate, error) {
	query := "SELECT * FROM instruction_templates WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(database.GetDB(), query, args, instructionPromptScanner)
}

func CreateInstructionPrompt(prompt *InstructionTemplate) error {
	util.NegFloat32PtrToNil(&prompt.Temperature)

	query := `INSERT INTO instruction_templates (name, type, temperature, system_prompt, instruction)
            VALUES (?, ?, ?, ?, ?) RETURNING id`
	args := []any{prompt.Name, prompt.Type, prompt.Temperature, prompt.SystemPrompt, prompt.Instruction}

	err := database.InsertRecord(database.GetDB(), query, args, &prompt.ID)
	if err != nil {
		return err
	}

	util.Emit(InstructionCreatedSignal, prompt)
	return nil
}

func UpdateInstructionPrompt(id int, prompt *InstructionTemplate) error {
	util.NegFloat32PtrToNil(&prompt.Temperature)

	query := `UPDATE instruction_templates
            SET name = ?,
                type = ?,
                temperature = ?,
                system_prompt = ?,
                instruction = ?
            WHERE id = ?`
	args := []any{prompt.Name, prompt.Type, prompt.Temperature, prompt.SystemPrompt, prompt.Instruction, id}

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(InstructionUpdatedSignal, prompt)
	return nil
}

func DeleteInstructionPrompt(id int) error {
	query := "DELETE FROM instruction_templates WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(InstructionDeletedSignal, id)
	return nil
}
