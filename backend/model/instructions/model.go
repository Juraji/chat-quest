package instructions

import (
	"juraji.nl/chat-quest/core"
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

func AllInstructionPrompts(cq *core.ChatQuestContext) ([]*InstructionTemplate, error) {
	query := "SELECT * FROM instruction_templates"
	return database.QueryForList(cq.DB(), query, nil, instructionPromptScanner)
}

func InstructionPromptById(cq *core.ChatQuestContext, id int) (*InstructionTemplate, error) {
	query := "SELECT * FROM instruction_templates WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(cq.DB(), query, args, instructionPromptScanner)
}

func CreateInstructionPrompt(cq *core.ChatQuestContext, prompt *InstructionTemplate) error {
	util.NegFloat32PtrToNil(&prompt.Temperature)

	query := `INSERT INTO instruction_templates (name, type, temperature, system_prompt, instruction)
            VALUES (?, ?, ?, ?, ?) RETURNING id`
	args := []any{prompt.Name, prompt.Type, prompt.Temperature, prompt.SystemPrompt, prompt.Instruction}

	err := database.InsertRecord(cq.DB(), query, args, &prompt.ID)
	if err != nil {
		return err
	}

	InstructionCreatedSignal.Emit(cq.Context(), prompt)
	return nil
}

func UpdateInstructionPrompt(cq *core.ChatQuestContext, id int, prompt *InstructionTemplate) error {
	util.NegFloat32PtrToNil(&prompt.Temperature)

	query := `UPDATE instruction_templates
            SET name = ?,
                type = ?,
                temperature = ?,
                system_prompt = ?,
                instruction = ?
            WHERE id = ?`
	args := []any{prompt.Name, prompt.Type, prompt.Temperature, prompt.SystemPrompt, prompt.Instruction, id}

	err := database.UpdateRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	InstructionUpdatedSignal.Emit(cq.Context(), prompt)
	return nil
}

func DeleteInstructionPrompt(cq *core.ChatQuestContext, id int) error {
	query := "DELETE FROM instruction_templates WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(cq.DB(), query, args)
	if err != nil {
		return err
	}

	InstructionDeletedSignal.Emit(cq.Context(), id)
	return nil
}
