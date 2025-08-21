package chat_sessions

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/util"
	"juraji.nl/chat-quest/model/characters"
	"math/rand"
)

type ChatParticipant struct {
	ChatSessionID int `json:"chatSessionId"`
	CharacterID   int `json:"characterId"`
}

func chatParticipantScanner(scanner database.RowScanner, dest *ChatParticipant) error {
	return scanner.Scan(
		&dest.ChatSessionID,
		&dest.CharacterID,
	)
}

func GetParticipants(sessionId int) ([]characters.Character, bool) {
	query := `SELECT c.* FROM chat_participants cp
                JOIN characters c ON cp.character_id = c.id
            WHERE cp.chat_session_id = ?`
	args := []any{sessionId}
	list, err := database.QueryForList(query, args, characters.CharacterScanner)
	if err != nil {
		log.Get().Error("Error fetching chat participants",
			zap.Int("sessionId", sessionId),
			zap.Error(err))
		return nil, false
	}

	return list, true
}

func GetParticipant(sessionId int, characterId int) (*ChatParticipant, bool) {
	query := `SELECT * FROM chat_participants WHERE chat_session_id = ? and character_id = ?`
	args := []any{sessionId, characterId}
	participant, err := database.QueryForRecord(query, args, chatParticipantScanner)
	if err != nil {
		log.Get().Error("Error fetching chat participant",
			zap.Int("sessionId", sessionId),
			zap.Int("characterId", characterId),
			zap.Error(err))
		return nil, false
	}

	return participant, true
}

func IsGroupSession(sessionId int) (bool, bool) {
	query := `SELECT COUNT(*) > 1 FROM chat_participants WHERE chat_session_id = ?`
	args := []any{sessionId}
	isGroupChat, err := database.QueryForRecord(query, args, database.BoolScanner)
	if err != nil || isGroupChat == nil {
		log.Get().Error("Error fetching session group status",
			zap.Int("sessionId", sessionId),
			zap.Error(err))
		return false, false
	}
	return *isGroupChat, true
}

// RandomParticipantId selects a participant from a chat session with weighted randomness based on talkativeness.
// The function returns a pointer to the selected participant's ID or nil if no selection could be made.
func RandomParticipantId(sessionId int) (*int, bool) {
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
            WHERE chat_session_id = ?;`

	args := []any{sessionId}
	scanFunc := func(scanner database.RowScanner, dest *choice) error {
		var t float32

		err := scanner.Scan(
			&dest.cId,
			&t)
		if err != nil {
			return err
		}

		dest.talkativeness = util.MaxFloat32(t, minT)
		return nil
	}

	choices, err := database.QueryForList(query, args, scanFunc)
	if err != nil {
		log.Get().Error("Error fetching chat participants",
			zap.Int("sessionId", sessionId),
			zap.Error(err))
		return nil, false
	}

	if len(choices) == 0 {
		// No participants
		return nil, true
	}
	if len(choices) == 1 {
		// Single participant
		return &choices[0].cId, true
	}

	// Calculate total weight (scaled by scale)
	totalWeight := 0
	for _, p := range choices {
		totalWeight += int(p.talkativeness * scale)
	}

	if totalWeight == 0 {
		return nil, true
	}

	// Generate random number between 0 and totalWeight-1
	randomValue := rand.Intn(totalWeight)

	// Find the participant that contains our random value
	accumulatedWeight := 0
	for _, p := range choices {
		weight := int(p.talkativeness * scale)
		accumulatedWeight += weight
		if randomValue < accumulatedWeight {
			return &p.cId, true
		}
	}

	// This line should theoretically never be reached if weights are correct
	return nil, false
}

func AddParticipant(sessionId int, characterId int) bool {
	query := `INSERT INTO chat_participants (chat_session_id, character_id) VALUES (?, ?)`
	args := []any{sessionId, characterId}
	err := database.InsertRecord(query, args)
	if err != nil {
		log.Get().Error("Error adding chat participant",
			zap.Int("sessionId", sessionId),
			zap.Int("characterId", characterId),
			zap.Error(err))
		return false
	}

	participant := ChatParticipant{sessionId, characterId}
	ChatParticipantAddedSignal.EmitBG(&participant)
	return true
}

func RemoveParticipant(sessionId int, characterId int) bool {
	query := `DELETE FROM chat_participants WHERE chat_session_id = ? AND character_id = ?`
	args := []any{sessionId, characterId}

	err := database.DeleteRecord(query, args)
	if err != nil {
		log.Get().Error("Error removing chat participant",
			zap.Int("sessionId", sessionId),
			zap.Int("characterId", characterId),
			zap.Error(err))
		return false
	}

	participant := ChatParticipant{sessionId, characterId}
	ChatParticipantRemovedSignal.EmitBG(&participant)
	return true
}
