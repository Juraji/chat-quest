package model

import "database/sql"

type SystemPrompt struct {
	DbEntity
	Name   string `json:"name"`
	Type   string `json:"type"`
	Prompt string `json:"prompt"`
}

func systemPromptScanner(scanner RowScanner, dest *SystemPrompt) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Type,
		&dest.Prompt,
	)
}

func AllSystemPrompts(db *sql.DB) ([]*SystemPrompt, error) {
	query := "SELECT * FROM system_prompts"
	return queryForList(db, query, nil, systemPromptScanner)
}

func SystemPromptById(db *sql.DB, id int64) (*SystemPrompt, error) {
	query := "SELECT * FROM system_prompts WHERE id = $1"
	args := []any{id}
	return queryForRecord(db, query, args, systemPromptScanner)
}

func CreateSystemPrompt(db *sql.DB, prompt *SystemPrompt) error {
	query := "INSERT INTO system_prompts (name, type, prompt) VALUES ($1, $2, $3)"
	args := []any{prompt.Name, prompt.Type, prompt.Prompt}
	scanFunc := func(scanner RowScanner) error {
		return scanner.Scan(&prompt.ID)
	}

	return insertRecord(db, query, args, scanFunc)
}

func UpdateSystemPrompt(db *sql.DB, id int64, prompt *SystemPrompt) error {
	query := "UPDATE system_prompts SET name = $1, type = $2, prompt = $3 WHERE id = $4"
	args := []any{prompt.Name, prompt.Type, prompt.Prompt, id}

	return updateRecord(db, query, args)
}

func DeleteSystemPrompt(db *sql.DB, id int64) error {
	query := "DELETE FROM system_prompts WHERE id = $1"
	args := []any{id}

	return deleteRecord(db, query, args)
}
