package chat_sessions

import (
	"fmt"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
	"juraji.nl/chat-quest/model/characters"
	"math/rand"
	"time"
)

type ChatSession struct {
	ID             int        `json:"id"`
	WorldID        int        `json:"worldId"`
	CreatedAt      *time.Time `json:"createdAt"`
	Name           string     `json:"name"`
	ScenarioID     *int       `json:"scenarioId"`
	EnableMemories bool       `json:"enableMemories"`
}

type ChatMessage struct {
	ID            int        `json:"id"`
	ChatSessionID int        `json:"chatSessionId"`
	CreatedAt     *time.Time `json:"createdAt"`
	IsUser        bool       `json:"isUser"`
	IsSystem      bool       `json:"isSystem"`
	IsGenerating  bool       `json:"isGenerating"`
	CharacterID   *int       `json:"characterId"`
	Content       string     `json:"content"`

	// Managed by memories
	MemoryID *int `json:"memoryId"`
}

type ChatParticipant struct {
	ChatSessionID int `json:"chatSessionId"`
	CharacterID   int `json:"characterId"`
}

func chatSessionScanner(scanner database.RowScanner, dest *ChatSession) error {
	return scanner.Scan(
		&dest.ID,
		&dest.WorldID,
		&dest.CreatedAt,
		&dest.Name,
		&dest.ScenarioID,
		&dest.EnableMemories,
	)
}

func chatMessageScanner(scanner database.RowScanner, dest *ChatMessage) error {
	return scanner.Scan(
		&dest.ID,
		&dest.ChatSessionID,
		&dest.CreatedAt,
		&dest.IsUser,
		&dest.IsSystem,
		&dest.IsGenerating,
		&dest.CharacterID,
		&dest.Content,
		&dest.MemoryID,
	)
}

func NewChatMessage(isUser bool, isSystem bool, isGenerating bool, characterId *int, content string) *ChatMessage {
	return &ChatMessage{
		ID:            0,
		ChatSessionID: 0,
		CreatedAt:     nil,
		IsUser:        isUser,
		IsSystem:      isSystem,
		IsGenerating:  isGenerating,
		CharacterID:   characterId,
		Content:       content,
		MemoryID:      nil,
	}
}

func GetAllByWorldId(worldId int) ([]ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE world_id=?"
	args := []any{worldId}

	return database.QueryForList(database.GetDB(), query, args, chatSessionScanner)
}

func GetByWorldIdAndId(worldId int, id int) (*ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE world_id=? AND id=?"
	args := []any{worldId, id}
	return database.QueryForRecord(database.GetDB(), query, args, chatSessionScanner)
}

func GetById(id int) (*ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE id=?"
	args := []any{id}
	return database.QueryForRecord(database.GetDB(), query, args, chatSessionScanner)
}

func Create(worldId int, session *ChatSession, characterIds []int) error {
	session.WorldID = worldId
	session.CreatedAt = nil

	tx, err := database.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(tx, err)

	query := `INSERT INTO chat_sessions (world_id, name, scenario_id, enable_memories)
            VALUES (?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		session.WorldID,
		session.Name,
		session.ScenarioID,
		session.EnableMemories,
	}

	err = database.InsertRecord(tx, query, args, &session.ID, &session.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert chat session: %w", err)
	}

	var addedParticipants []*ChatParticipant
	for _, characterId := range characterIds {
		err = addParticipant(tx, session.ID, characterId)
		addedParticipants = append(addedParticipants, &ChatParticipant{session.ID, characterId})
		if err != nil {
			return fmt.Errorf("failed to insert chat participant (%d -> %d):  %w", session.ID, characterId, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	ChatSessionCreatedSignal.EmitBG(session)
	ChatParticipantAddedSignal.EmitAllBG(addedParticipants)
	return nil
}

func Update(worldId int, id int, session *ChatSession) error {
	query := `UPDATE chat_sessions
            SET name = ?,
                scenario_id = ?,
                enable_memories = ?
            WHERE world_id = ?
              AND id = ?`
	args := []any{
		session.Name,
		session.ScenarioID,
		session.EnableMemories,
		worldId,
		id,
	}

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	ChatSessionUpdatedSignal.EmitBG(session)
	return nil
}

func Delete(worldId int, id int) error {
	query := "DELETE FROM chat_sessions WHERE world_id=? AND id=?"
	args := []any{worldId, id}

	err := database.DeleteRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	ChatSessionDeletedSignal.EmitBG(worldId)
	return nil
}

func GetChatMessages(sessionId int) ([]ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=?"
	args := []any{sessionId}
	return database.QueryForList(database.GetDB(), query, args, chatMessageScanner)
}

func GetChatMessagesPreceding(sessionId int, messageId int) ([]ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=? and id < ?"
	args := []any{sessionId, messageId}
	return database.QueryForList(database.GetDB(), query, args, chatMessageScanner)
}

func CreateChatMessage(sessionId int, chatMessage *ChatMessage) error {
	chatMessage.ChatSessionID = sessionId
	chatMessage.CreatedAt = nil

	query := `INSERT INTO chat_messages (chat_session_id, is_user, is_system, is_generating, character_id, content)
            VALUES (?, ?, ?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		chatMessage.ChatSessionID,
		chatMessage.IsUser,
		chatMessage.IsSystem,
		chatMessage.IsGenerating,
		chatMessage.CharacterID,
		chatMessage.Content,
	}

	err := database.InsertRecord(database.GetDB(), query, args, &chatMessage.ID, &chatMessage.CreatedAt)
	if err != nil {
		return err
	}

	ChatMessageCreatedSignal.EmitBG(chatMessage)

	return nil
}

func UpdateChatMessage(sessionId int, id int, chatMessage *ChatMessage) error {
	query := `UPDATE chat_messages
            SET content = ?,
                is_generating = ?
            WHERE chat_session_id = ?
              AND id = ?`
	args := []any{chatMessage.Content, chatMessage.IsGenerating, sessionId, id}

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	ChatMessageUpdatedSignal.EmitBG(chatMessage)
	return nil
}

func DeleteChatMessagesFrom(sessionId int, id int) error {
	//language=SQL
	query := `DELETE
            FROM chat_messages
            WHERE chat_session_id = ?
              AND id >= ?
            RETURNING id`
	args := []any{sessionId, id}

	deletedIds, err := database.QueryForList(database.GetDB(), query, args, func(scanner database.RowScanner, dest *int) error {
		return scanner.Scan(dest)
	})
	if err != nil {
		return err
	}

	ChatMessageDeletedSignal.EmitAllBG(deletedIds)
	return nil
}

func GetParticipants(sessionId int) ([]characters.Character, error) {
	query := `SELECT c.* FROM chat_participants cp
                JOIN characters c ON cp.character_id = c.id
            WHERE cp.chat_session_id = ?`
	args := []any{sessionId}
	return database.QueryForList(database.GetDB(), query, args, characters.CharacterScanner)
}

func IsGroupSession(sessionId int) (bool, error) {
	query := `SELECT COUNT(*) > 1 FROM chat_participants WHERE chat_session_id = ?`
	args := []any{sessionId}
	isGroupChat, err := database.QueryForRecord(database.GetDB(), query, args, database.BoolScanner)
	if err != nil || isGroupChat == nil {
		return false, err
	}
	return *isGroupChat, nil
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

	choices, err := database.QueryForList(database.GetDB(), query, args, scanFunc)
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
	return nil, nil
}

func AddParticipant(sessionId int, characterId int) error {
	err := addParticipant(database.GetDB(), sessionId, characterId)
	if err != nil {
		return err
	}

	participant := ChatParticipant{sessionId, characterId}
	ChatParticipantAddedSignal.EmitBG(&participant)
	return nil
}

func addParticipant(db database.QueryExecutor, sessionId int, characterId int) error {
	query := `INSERT INTO chat_participants (chat_session_id, character_id) VALUES (?, ?)`
	args := []any{sessionId, characterId}
	return database.InsertRecord(db, query, args)
}

func RemoveParticipant(sessionId int, characterId int) error {
	query := `DELETE FROM chat_participants WHERE chat_session_id = ? AND character_id = ?`
	args := []any{sessionId, characterId}

	err := database.DeleteRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	participant := ChatParticipant{sessionId, characterId}
	ChatParticipantRemovedSignal.EmitBG(&participant)
	return nil
}
