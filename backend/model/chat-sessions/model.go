package chat_sessions

import (
	"fmt"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/util"
	"juraji.nl/chat-quest/model/characters"
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
		&dest.CharacterID,
		&dest.Content,
		&dest.MemoryID,
	)
}

func NewChatMessage(sessionId int, isUser bool, characterId *int, content string) *ChatMessage {
	return &ChatMessage{
		ID:            0,
		ChatSessionID: sessionId,
		CreatedAt:     nil,
		IsUser:        isUser,
		CharacterID:   characterId,
		Content:       content,
		MemoryID:      nil,
	}
}

func GetAllChatSessionsByWorldId(worldId int) ([]ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE world_id=?"
	args := []any{worldId}

	return database.QueryForList(database.GetDB(), query, args, chatSessionScanner)
}

func GetChatSessionById(worldId int, id int) (*ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE world_id=? AND id=?"
	args := []any{worldId, id}
	return database.QueryForRecord(database.GetDB(), query, args, chatSessionScanner)
}

func CreateChatSession(worldId int, chatSession *ChatSession, characterIds []int) error {
	chatSession.WorldID = worldId
	chatSession.CreatedAt = nil

	tx, err := database.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer database.RollBackOnErr(tx, err)

	query := `INSERT INTO chat_sessions (world_id, name, scenario_id, enable_memories)
            VALUES (?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		chatSession.WorldID,
		chatSession.Name,
		chatSession.ScenarioID,
		chatSession.EnableMemories,
	}

	err = database.InsertRecord(tx, query, args, &chatSession.ID, &chatSession.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert chat session: %w", err)
	}

	for _, characterId := range characterIds {
		err = addChatSessionParticipant(tx, chatSession.ID, characterId)
		if err != nil {
			return fmt.Errorf("failed to insert chat participant (%d -> %d):  %w", chatSession.ID, characterId, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	util.Emit(ChatSessionCreatedSignal, chatSession)
	for _, characterId := range characterIds {
		participant := ChatParticipant{chatSession.ID, characterId}
		util.Emit(ChatParticipantAddedSignal, &participant)
	}

	return nil
}

func UpdateChatSession(worldId int, id int, chatSession *ChatSession) error {
	query := `UPDATE chat_sessions
            SET name = ?,
                scenario_id = ?,
                enable_memories = ?
            WHERE world_id = ?
              AND id = ?`
	args := []any{
		chatSession.Name,
		chatSession.ScenarioID,
		chatSession.EnableMemories,
		worldId,
		id,
	}

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(ChatSessionUpdatedSignal, chatSession)
	return nil
}

func DeleteChatSessionById(worldId int, id int) error {
	query := "DELETE FROM chat_sessions WHERE world_id=? AND id=?"
	args := []any{worldId, id}

	err := database.DeleteRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(ChatSessionDeletedSignal, worldId)
	return nil
}

func GetChatMessages(sessionId int) ([]ChatMessage, error) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=?"
	args := []any{sessionId}
	return database.QueryForList(database.GetDB(), query, args, chatMessageScanner)
}

func CreateChatMessage(sessionId int, chatMessage *ChatMessage) error {
	chatMessage.ChatSessionID = sessionId
	chatMessage.CreatedAt = nil

	query := `INSERT INTO chat_messages (chat_session_id, is_user, character_id, content)
            VALUES (?, ?, ?, ?) RETURNING id, created_at`
	args := []any{
		chatMessage.ChatSessionID,
		chatMessage.IsUser,
		chatMessage.CharacterID,
		chatMessage.Content,
	}

	err := database.InsertRecord(database.GetDB(), query, args, &chatMessage.ID, &chatMessage.CreatedAt)
	if err != nil {
		return err
	}

	util.Emit(ChatMessageCreatedSignal, chatMessage)

	return nil
}

func UpdateChatMessage(sessionId int, id int, chatMessage *ChatMessage) error {
	query := `UPDATE chat_messages
            SET content = ?
            WHERE chat_session_id = ?
              AND id = ?`
	args := []any{chatMessage.Content, sessionId, id}

	err := database.UpdateRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	util.Emit(ChatMessageUpdatedSignal, chatMessage)
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

	util.EmitAll(ChatMessageDeletedSignal, deletedIds)
	return nil
}

func GetChatSessionParticipants(chatSessionId int) ([]characters.Character, error) {
	query := `SELECT c.* FROM chat_participants cp
                JOIN characters c ON cp.character_id = c.id
            WHERE cp.chat_session_id = ?`
	args := []any{chatSessionId}
	return database.QueryForList(database.GetDB(), query, args, characters.CharacterScanner)
}

func IsGroupChatSession(chatSessionId int) (bool, error) {
	query := `SELECT COUNT(*) > 1 FROM chat_participants WHERE chat_session_id = ?`
	args := []any{chatSessionId}
	isGroupChat, err := database.QueryForRecord(database.GetDB(), query, args, database.BoolScanner)
	if err != nil || isGroupChat == nil {
		return false, err
	}
	return *isGroupChat, nil
}

// RandomChatSessionParticipantId selects a random character ID from a chat session,
// biased by each character's group_talkativeness value. The query creates a weighted
// probability distribution and uses it to select a participant with higher talkativeness
// characters being more likely to be chosen.
// Note that a character is always chosen, if there are any. Even if all are really not talkative.
func RandomChatSessionParticipantId(chatSessionId int) (*int, error) {
	// language=sqlite
	query := `WITH ranked_characters AS (
              SELECT
                cp.character_id,
                cd.group_talkativeness,
                SUM(cd.group_talkativeness) OVER (ORDER BY RANDOM()) as running_sum,
                SUM(cd.group_talkativeness) OVER () as total_sum
              FROM chat_participants cp
              JOIN character_details cd ON cp.character_id = cd.character_id
              WHERE cp.chat_session_id = ?
            )
            SELECT character_id
            FROM ranked_characters
            WHERE RANDOM() * total_sum <= running_sum
            LIMIT 1;`
	args := []any{chatSessionId}

	return database.QueryForRecord(database.GetDB(), query, args, database.IntScanner)
}

func AddChatSessionParticipant(chatSessionId int, characterId int) error {
	err := addChatSessionParticipant(database.GetDB(), chatSessionId, characterId)
	if err != nil {
		return err
	}

	participant := ChatParticipant{chatSessionId, characterId}
	util.Emit(ChatParticipantAddedSignal, &participant)
	return nil
}

func addChatSessionParticipant(db database.QueryExecutor, chatSessionId int, characterId int) error {
	query := `INSERT INTO chat_participants (chat_session_id, character_id) VALUES (?, ?)`
	args := []any{chatSessionId, characterId}
	return database.InsertRecord(db, query, args)
}

func RemoveChatSessionParticipant(chatSessionId int, characterId int) error {
	query := `DELETE FROM chat_participants WHERE chat_session_id = ? AND character_id = ?`
	args := []any{chatSessionId, characterId}

	err := database.DeleteRecord(database.GetDB(), query, args)
	if err != nil {
		return err
	}

	participant := ChatParticipant{chatSessionId, characterId}
	util.Emit(ChatParticipantRemovedSignal, &participant)
	return nil
}
