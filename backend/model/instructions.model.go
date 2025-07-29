package model

import (
	"database/sql"
	"juraji.nl/chat-quest/util"
)

type InstructionPrompt struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Temperature  *float64 `json:"temperature"`
	SystemPrompt string   `json:"systemPrompt"`
	Instruction  string   `json:"instruction"`
}

func InstructionPromptScanner(scanner RowScanner, dest *InstructionPrompt) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Type,
		&dest.Temperature,
		&dest.SystemPrompt,
		&dest.Instruction,
	)
}

func AllInstructionPrompts(db *sql.DB) ([]*InstructionPrompt, error) {
	query := "SELECT * FROM instruction_prompts"
	return queryForList(db, query, nil, InstructionPromptScanner)
}

func InstructionPromptById(db *sql.DB, id int64) (*InstructionPrompt, error) {
	query := "SELECT * FROM instruction_prompts WHERE id = $1"
	args := []any{id}
	return queryForRecord(db, query, args, InstructionPromptScanner)
}

func CreateInstructionPrompt(db *sql.DB, prompt *InstructionPrompt) error {
	util.NegFloat64PtrToNil(&prompt.Temperature)

	query := `INSERT INTO instruction_prompts (name, type, temperature, system_prompt, instruction)
            VALUES ($1, $2, $3, $4, $5) RETURNING id`
	args := []any{prompt.Name, prompt.Type, prompt.Temperature, prompt.SystemPrompt, prompt.Instruction}
	scanFunc := func(scanner RowScanner) error {
		return scanner.Scan(&prompt.ID)
	}

	return insertRecord(db, query, args, scanFunc)
}

func UpdateInstructionPrompt(db *sql.DB, id int64, prompt *InstructionPrompt) error {
	util.NegFloat64PtrToNil(&prompt.Temperature)

	query := `UPDATE instruction_prompts
            SET name = $1,
                type = $2,
                temperature = $3,
                system_prompt = $4,
                instruction = $5
            WHERE id = $6`
	args := []any{prompt.Name, prompt.Type, prompt.Temperature, prompt.SystemPrompt, prompt.Instruction, id}

	return updateRecord(db, query, args)
}

func DeleteInstructionPrompt(db *sql.DB, id int64) error {
	query := "DELETE FROM instruction_prompts WHERE id = $1"
	args := []any{id}

	return deleteRecord(db, query, args)
}
