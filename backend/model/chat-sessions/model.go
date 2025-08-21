package chat_sessions

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
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

func chatParticipantScanner(scanner database.RowScanner, dest *ChatParticipant) error {
	return scanner.Scan(
		&dest.ChatSessionID,
		&dest.CharacterID,
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

func GetAllByWorldId(worldId int) ([]ChatSession, bool) {
	query := "SELECT * FROM chat_sessions WHERE world_id=?"
	args := []any{worldId}

	list, err := database.QueryForList(query, args, chatSessionScanner)
	if err != nil {
		log.Get().Error("Error fetching chat sessions",
			zap.Int("worldId", worldId), zap.Error(err))
		return nil, false
	}

	return list, true
}

func GetByWorldIdAndId(worldId int, id int) (*ChatSession, bool) {
	query := "SELECT * FROM chat_sessions WHERE world_id=? AND id=?"
	args := []any{worldId, id}
	session, err := database.QueryForRecord(query, args, chatSessionScanner)
	if err != nil {
		log.Get().Error("Error fetching chat session",
			zap.Int("worldId", worldId),
			zap.Int("id", id),
			zap.Error(err),
		)
		return nil, false
	}

	return session, true
}

func GetById(id int) (*ChatSession, bool) {
	query := "SELECT * FROM chat_sessions WHERE id=?"
	args := []any{id}
	session, err := database.QueryForRecord(query, args, chatSessionScanner)
	if err != nil {
		log.Get().Error("Error fetching chat session",
			zap.Int("id", id),
			zap.Error(err),
		)
		return nil, false
	}

	return session, true
}

func Create(worldId int, session *ChatSession, characterIds []int) bool {
	session.WorldID = worldId
	session.CreatedAt = nil

	var addedParticipants []*ChatParticipant
	err := database.Transactional(func(ctx *database.TxContext) error {
		query := `INSERT INTO chat_sessions (world_id, name, scenario_id, enable_memories)
            VALUES (?, ?, ?, ?) RETURNING id, created_at`
		args := []any{
			session.WorldID,
			session.Name,
			session.ScenarioID,
			session.EnableMemories,
		}

		err := ctx.InsertRecord(query, args, &session.ID, &session.CreatedAt)
		if err != nil {
			log.Get().Error("Error creating chat session",
				zap.Int("worldId", worldId),
				zap.Error(err),
			)
			return err
		}

		if len(characterIds) == 0 {
			return nil
		}

		sessionId := session.ID
		for _, characterId := range characterIds {
			query = `INSERT INTO chat_participants (chat_session_id, character_id) VALUES (?, ?)`
			args = []any{sessionId, characterId}
			err := ctx.InsertRecord(query, args)
			if err != nil {
				log.Get().Error("Error adding chat participant",
					zap.Int("sessionId", sessionId),
					zap.Int("characterId", characterId),
					zap.Error(err))
				return err
			}
			addedParticipants = append(addedParticipants, &ChatParticipant{session.ID, characterId})
		}

		return nil
	})

	ChatSessionCreatedSignal.EmitBG(session)
	ChatParticipantAddedSignal.EmitAllBG(addedParticipants)
	return err == nil
}

func Update(worldId int, id int, session *ChatSession) bool {
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

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating chat session",
			zap.Int("worldId", worldId),
			zap.Int("id", id),
			zap.Error(err))
		return false
	}

	ChatSessionUpdatedSignal.EmitBG(session)
	return true
}

func Delete(worldId int, id int) bool {
	query := "DELETE FROM chat_sessions WHERE world_id=? AND id=?"
	args := []any{worldId, id}

	err := database.DeleteRecord(query, args)
	if err != nil {
		log.Get().Error("Error deleting chat session",
			zap.Int("worldId", worldId),
			zap.Int("id", id),
			zap.Error(err))
		return false
	}

	ChatSessionDeletedSignal.EmitBG(worldId)
	return true
}

func GetChatMessages(sessionId int) ([]ChatMessage, bool) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=?"
	args := []any{sessionId}
	list, err := database.QueryForList(query, args, chatMessageScanner)
	if err != nil {
		log.Get().Error("Error fetching chat session messages",
			zap.Int("sessionId", sessionId),
			zap.Error(err))
		return nil, false
	}

	return list, true
}

func GetChatMessagesPreceding(sessionId int, messageId int) ([]ChatMessage, bool) {
	query := "SELECT * FROM chat_messages WHERE chat_session_id=? and id < ?"
	args := []any{sessionId, messageId}
	list, err := database.QueryForList(query, args, chatMessageScanner)
	if err != nil {
		log.Get().Error("Error fetching preceding chat session messages",
			zap.Int("sessionId", sessionId),
			zap.Int("precedingId", messageId),
			zap.Error(err))
		return nil, false
	}

	return list, true
}

func CreateChatMessage(sessionId int, chatMessage *ChatMessage) bool {
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

	err := database.InsertRecord(query, args, &chatMessage.ID, &chatMessage.CreatedAt)
	if err != nil {
		log.Get().Error("Error creating chat session message",
			zap.Int("sessionId", sessionId),
			zap.Error(err))
		return false
	}

	ChatMessageCreatedSignal.EmitBG(chatMessage)
	return true
}

func UpdateChatMessage(sessionId int, id int, chatMessage *ChatMessage) bool {
	query := `UPDATE chat_messages
            SET content = ?,
                is_generating = ?
            WHERE chat_session_id = ?
              AND id = ?`
	args := []any{chatMessage.Content, chatMessage.IsGenerating, sessionId, id}

	err := database.UpdateRecord(query, args)
	if err != nil {
		log.Get().Error("Error updating chat session message",
			zap.Int("sessionId", sessionId),
			zap.Int("id", id),
			zap.Error(err))
		return false
	}

	ChatMessageUpdatedSignal.EmitBG(chatMessage)
	return true
}

func DeleteChatMessagesFrom(sessionId int, id int) bool {
	//language=SQL
	query := `DELETE
            FROM chat_messages
            WHERE chat_session_id = ?
              AND id >= ?
            RETURNING id`
	args := []any{sessionId, id}

	deletedIds, err := database.QueryForList(query, args, func(scanner database.RowScanner, dest *int) error {
		return scanner.Scan(dest)
	})
	if err != nil {
		log.Get().Error("Error deleting chat session message",
			zap.Int("sessionId", sessionId),
			zap.Int("startingAtId", id),
			zap.Error(err))
		return false
	}

	ChatMessageDeletedSignal.EmitAllBG(deletedIds)
	return true
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
