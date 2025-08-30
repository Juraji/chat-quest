package chat_sessions

import (
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"time"
)

type ChatSession struct {
	ID                      int        `json:"id"`
	WorldID                 int        `json:"worldId"`
	CreatedAt               *time.Time `json:"createdAt"`
	Name                    string     `json:"name"`
	ScenarioID              *int       `json:"scenarioId"`
	EnableMemories          bool       `json:"enableMemories"`
	PauseAutomaticResponses bool       `json:"pauseAutomaticResponses"`
}

func chatSessionScanner(scanner database.RowScanner, dest *ChatSession) error {
	return scanner.Scan(
		&dest.ID,
		&dest.WorldID,
		&dest.CreatedAt,
		&dest.Name,
		&dest.ScenarioID,
		&dest.EnableMemories,
		&dest.PauseAutomaticResponses,
	)
}

func GetAllByWorldId(worldId int) ([]ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE world_id=?"
	args := []any{worldId}

	return database.QueryForList(query, args, chatSessionScanner)
}

func GetByWorldIdAndId(worldId int, id int) (*ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE world_id=? AND id=?"
	args := []any{worldId, id}
	return database.QueryForRecord(query, args, chatSessionScanner)
}

func GetById(id int) (*ChatSession, error) {
	query := "SELECT * FROM chat_sessions WHERE id=?"
	args := []any{id}
	return database.QueryForRecord(query, args, chatSessionScanner)
}

func Create(worldId int, session *ChatSession, characterIds []int) error {
	session.WorldID = worldId
	session.CreatedAt = nil

	var addedParticipants []*ChatParticipant
	err := database.Transactional(func(ctx *database.TxContext) error {
		query := `INSERT INTO chat_sessions (world_id, name, scenario_id, enable_memories, pause_automatic_responses)
            VALUES (?, ?, ?, ?, ?) RETURNING id, created_at`
		args := []any{
			session.WorldID,
			session.Name,
			session.ScenarioID,
			session.EnableMemories,
			session.PauseAutomaticResponses,
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

	if err == nil {
		ChatSessionCreatedSignal.EmitBG(session)
		ChatParticipantAddedSignal.EmitAllBG(addedParticipants)
	}

	return err
}

func Update(worldId int, id int, session *ChatSession) error {
	query := `UPDATE chat_sessions
            SET name = ?,
                scenario_id = ?,
                enable_memories = ?,
                pause_automatic_responses = ?
            WHERE world_id = ?
              AND id = ?`
	args := []any{
		session.Name,
		session.ScenarioID,
		session.EnableMemories,
		session.PauseAutomaticResponses,
		worldId,
		id,
	}

	err := database.UpdateRecord(query, args)

	if err == nil {
		ChatSessionUpdatedSignal.EmitBG(session)
	}

	return err
}

func Delete(worldId int, id int) error {
	query := "DELETE FROM chat_sessions WHERE world_id=? AND id=?"
	args := []any{worldId, id}

	_, err := database.DeleteRecord(query, args)

	if err == nil {
		ChatSessionDeletedSignal.EmitBG(worldId)
	}

	return err
}
