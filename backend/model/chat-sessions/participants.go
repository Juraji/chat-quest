package chat_sessions

import (
	"errors"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
	"juraji.nl/chat-quest/model/characters"
	"math/rand"
	"time"
)

type ChatParticipant struct {
	ChatSessionID int        `json:"chatSessionId"`
	CharacterID   int        `json:"characterId"`
	AddedOn       *time.Time `json:"addedOn"`
	RemovedOn     *time.Time `json:"removedOn"`
}

func chatParticipantScanner(scanner database.RowScanner, dest *ChatParticipant) error {
	return scanner.Scan(
		&dest.ChatSessionID,
		&dest.CharacterID,
		&dest.AddedOn,
		&dest.RemovedOn,
	)
}

func GetCurrentParticipants(sessionId int) ([]characters.Character, error) {
	query := `SELECT c.* FROM chat_participants p
                JOIN characters c ON p.character_id = c.id
            WHERE p.chat_session_id = ?
              AND p.removed_on IS NULL`
	args := []any{sessionId}
	return database.QueryForList(query, args, characters.CharacterScanner)
}

func GetParticipantsBefore(sessionId int, before time.Time) ([]characters.Character, error) {
	query := `SELECT c.* FROM chat_participants p
                JOIN characters c ON p.character_id = c.id
            WHERE p.chat_session_id = ?
              AND p.added_on <= ?
              AND (p.removed_on IS NULL OR p.removed_on > ?)`
	args := []any{sessionId, before, before}
	return database.QueryForList(query, args, characters.CharacterScanner)
}

func GetParticipant(sessionId int, characterId int) (*ChatParticipant, error) {
	query := `SELECT * FROM chat_participants
         	  WHERE chat_session_id = ?
         	    AND character_id = ?
         	    AND removed_on IS NULL`
	args := []any{sessionId, characterId}
	return database.QueryForRecord(query, args, chatParticipantScanner)
}

func IsGroupSession(sessionId int) (*bool, error) {
	query := `SELECT COUNT(*) > 1 FROM chat_participants WHERE chat_session_id = ?`
	args := []any{sessionId}
	return database.QueryForRecord(query, args, database.BoolScanner)
}

// RandomParticipantId selects a participant from a chat session with weighted randomness based on talkativeness.
// The function returns a pointer to the selected participant's ID or nil if no selection could be made.
func RandomParticipantId(sessionId int) (*int, error) {
	const scale float32 = 20
	const minT float32 = 0.05
	type choice struct {
		cId           int
		talkativeness float32
	}

	query := `SELECT
                c.id AS participant_id,
                c.group_talkativeness AS group_talkativeness
            FROM chat_participants p
                JOIN characters c ON p.character_id = c.id
            WHERE chat_session_id = ?
			  AND removed_on IS NULL;`

	args := []any{sessionId}
	scanFunc := func(scanner database.RowScanner, dest *choice) error {
		err := scanner.Scan(&dest.cId, &dest.talkativeness)
		if err != nil {
			return err
		}

		dest.talkativeness = util.MaxFloat32(dest.talkativeness, minT)
		return nil
	}

	choices, err := database.QueryForList(query, args, scanFunc)
	if err != nil {
		return nil, err
	}

	if len(choices) == 0 {
		// No participants
		return nil, nil
	}
	if len(choices) == 1 {
		// Single participant
		return &choices[0].cId, nil
	}

	// Calculate total weight (scaled by scale)
	totalWeight := 0
	for _, p := range choices {
		totalWeight += int(p.talkativeness * scale)
	}

	if totalWeight == 0 {
		return nil, nil
	}

	// Generate random number between 0 and totalWeight-1
	randomValue := rand.Intn(totalWeight)

	// Find the participant that contains our random value
	accumulatedWeight := 0
	for _, p := range choices {
		weight := int(p.talkativeness * scale)
		accumulatedWeight += weight
		if randomValue < accumulatedWeight {
			return &p.cId, nil
		}
	}

	// This line should theoretically never be reached if weights are correct
	return nil, errors.New("incorrect weights")
}

func AddParticipant(sessionId int, characterId int) error {
	participant := ChatParticipant{
		ChatSessionID: sessionId,
		CharacterID:   characterId,
		AddedOn:       nil,
		RemovedOn:     nil,
	}

	query := `INSERT INTO chat_participants (chat_session_id, character_id)
			  VALUES (?, ?)
			  ON CONFLICT (chat_session_id, character_id)
			      DO UPDATE SET removed_on = NULL
              RETURNING added_on, removed_on; `
	args := []any{sessionId, characterId}

	err := database.InsertRecord(query, args, &participant.AddedOn, &participant.RemovedOn)

	if err == nil {
		ChatParticipantAddedSignal.EmitBG(&participant)
	}

	return err
}

func RemoveParticipant(sessionId int, characterId int) error {
	participant := ChatParticipant{
		ChatSessionID: sessionId,
		CharacterID:   characterId,
		AddedOn:       nil,
		RemovedOn:     nil,
	}

	query := `UPDATE chat_participants
			  SET removed_on = CURRENT_TIMESTAMP
              WHERE chat_session_id = ? AND character_id = ?
              RETURNING added_on, removed_on;`
	args := []any{sessionId, characterId}

	err := database.InsertRecord(query, args, &participant.AddedOn, &participant.RemovedOn)

	if err == nil {
		ChatParticipantRemovedSignal.EmitBG(&participant)
	}

	return err
}
