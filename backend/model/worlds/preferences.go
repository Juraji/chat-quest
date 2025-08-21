package worlds

import (
	"errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
)

type ChatPreferences struct {
	ChatModelID       *int `json:"chatModelId"`
	ChatInstructionID *int `json:"chatInstructionId"`
}

func (p *ChatPreferences) Validate() error {
	if p == nil {
		return errors.New("chat preferences is nil")
	}
	if p.ChatModelID == nil {
		return errors.New("ChatModelId is nil")
	}
	if p.ChatInstructionID == nil {
		return errors.New("ChatInstructionId is nil")
	}

	return nil
}

func chatPreferencesScanner(scanner database.RowScanner, dest *ChatPreferences) error {
	return scanner.Scan(
		&dest.ChatModelID,
		&dest.ChatInstructionID,
	)
}

func GetChatPreferences() (*ChatPreferences, bool) {
	query := "SELECT chat_model_id, chat_instruction_id FROM chat_preferences WHERE id = 0"
	prefs, err := database.QueryForRecord(query, nil, chatPreferencesScanner)
	if err != nil {
		log.Get().Error("Error fetching chat_preferences", zap.Error(err))
		return nil, false
	}

	return prefs, true
}

func UpdateChatPreferences(prefs *ChatPreferences) bool {
	query := `UPDATE chat_preferences
            SET chat_model_id = ?,
                chat_instruction_id = ?
            WHERE id = 0`
	args := []any{
		prefs.ChatModelID,
		prefs.ChatInstructionID,
	}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating chat_preferences", zap.Error(err))
		return false
	}

	ChatPreferencesUpdatedSignal.EmitBG(prefs)
	return true
}
