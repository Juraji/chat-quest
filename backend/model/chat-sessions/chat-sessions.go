package chat_sessions

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
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
