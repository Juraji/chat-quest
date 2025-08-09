package instructions

import (
	"database/sql"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/util"
)

type InstructionTemplate struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Temperature  *float64 `json:"temperature"`
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

func AllInstructionPrompts(db *sql.DB) ([]*InstructionTemplate, error) {
	query := "SELECT * FROM instruction_templates"
	return database.QueryForList(db, query, nil, instructionPromptScanner)
}

func InstructionPromptById(db *sql.DB, id int64) (*InstructionTemplate, error) {
	query := "SELECT * FROM instruction_templates WHERE id = ?"
	args := []any{id}
	return database.QueryForRecord(db, query, args, instructionPromptScanner)
}

func CreateInstructionPrompt(db *sql.DB, prompt *InstructionTemplate) error {
	util.NegFloat64PtrToNil(&prompt.Temperature)

	query := `INSERT INTO instruction_templates (name, type, temperature, system_prompt, instruction)
            VALUES (?, ?, ?, ?, ?) RETURNING id`
	args := []any{prompt.Name, prompt.Type, prompt.Temperature, prompt.SystemPrompt, prompt.Instruction}

	err := database.InsertRecord(db, query, args, &prompt.ID)
	defer util.EmitOnSuccess(InstructionCreatedSignal, prompt, err)

	return err
}

func UpdateInstructionPrompt(db *sql.DB, id int64, prompt *InstructionTemplate) error {
	util.NegFloat64PtrToNil(&prompt.Temperature)

	query := `UPDATE instruction_templates
            SET name = ?,
                type = ?,
                temperature = ?,
                system_prompt = ?,
                instruction = ?
            WHERE id = ?`
	args := []any{prompt.Name, prompt.Type, prompt.Temperature, prompt.SystemPrompt, prompt.Instruction, id}

	err := database.UpdateRecord(db, query, args)
	defer util.EmitOnSuccess(InstructionUpdatedSignal, prompt, err)

	return err
}

func DeleteInstructionPrompt(db *sql.DB, id int64) error {
	query := "DELETE FROM instruction_templates WHERE id = ?"
	args := []any{id}

	err := database.DeleteRecord(db, query, args)
	defer util.EmitOnSuccess(InstructionDeletedSignal, id, err)

	return err
}
