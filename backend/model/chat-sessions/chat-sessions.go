package chat_sessions

import (
	"time"

	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
)

type ChatSession struct {
	ID                      int        `json:"id"`
	WorldID                 int        `json:"worldId"`
	CreatedAt               *time.Time `json:"createdAt"`
	Name                    string     `json:"name"`
	ScenarioID              *int       `json:"scenarioId"`
	GenerateMemories        bool       `json:"generateMemories"`
	UseMemories             bool       `json:"useMemories"`
	AutoArchiveMessages     bool       `json:"autoArchiveMessages"`
	PauseAutomaticResponses bool       `json:"pauseAutomaticResponses"`
}

func chatSessionScanner(scanner database.RowScanner, dest *ChatSession) error {
	return scanner.Scan(
		&dest.ID,
		&dest.WorldID,
		&dest.CreatedAt,
		&dest.Name,
		&dest.ScenarioID,
		&dest.GenerateMemories,
		&dest.UseMemories,
		&dest.AutoArchiveMessages,
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
		query := `INSERT INTO chat_sessions (world_id, name, scenario_id, generate_memories, use_memories,
                           auto_archive_messages, pause_automatic_responses)
            VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id, created_at`
		args := []any{
			session.WorldID,
			session.Name,
			session.ScenarioID,
			session.GenerateMemories,
			session.UseMemories,
			session.AutoArchiveMessages,
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
			participant := ChatParticipant{
				ChatSessionID: session.ID,
				CharacterID:   characterId,
				AddedOn:       nil,
				RemovedOn:     nil,
			}

			query = `INSERT INTO chat_participants (chat_session_id, character_id, muted)
					 VALUES (?, ?, FALSE)
					 RETURNING added_on`
			args = []any{sessionId, characterId}
			err = ctx.InsertRecord(query, args, &participant.AddedOn)
			if err != nil {
				log.Get().Error("Error adding chat participant",
					zap.Int("sessionId", sessionId),
					zap.Int("characterId", characterId),
					zap.Error(err))
				return err
			}
			addedParticipants = append(addedParticipants, &participant)
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
                generate_memories = ?,
                use_memories = ?,
                auto_archive_messages = ?,
                pause_automatic_responses = ?
            WHERE world_id = ?
              AND id = ?`
	args := []any{
		session.Name,
		session.ScenarioID,
		session.GenerateMemories,
		session.UseMemories,
		session.AutoArchiveMessages,
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

func ForkChatSession(sessionId int, messageId int) (*ChatSession, error) {
	var newSession *ChatSession

	txErr := database.Transactional(func(ctx *database.TxContext) error {
		// Copy chat session
		var err error
		query := `INSERT INTO chat_sessions (world_id, name, scenario_id, generate_memories, use_memories,
                           auto_archive_messages, pause_automatic_responses)
				  SELECT world_id,
				         name || ' (copy)',
				         scenario_id,
				         generate_memories,
				         use_memories,
				         auto_archive_messages,
				         pause_automatic_responses
				  FROM chat_sessions
				  WHERE id = ?
				  RETURNING id, world_id, created_at, name, scenario_id, generate_memories, use_memories,
				      auto_archive_messages, pause_automatic_responses;`
		args := []any{sessionId}
		if newSession, err = database.QueryForRecord(query, args, chatSessionScanner); err != nil {
			return err
		}

		newSessionId := newSession.ID

		// Copy participants
		query = `INSERT INTO chat_participants (chat_session_id, character_id, added_on, removed_on, muted)
				 SELECT ?, character_id, added_on, removed_on, muted
				 FROM chat_participants
				 WHERE chat_session_id = ?`
		args = []any{newSessionId, sessionId}
		if err = database.UpdateRecord(query, args); err != nil {
			return err
		}

		// Copy messages up to messageId
		query = `INSERT INTO chat_messages (chat_session_id, created_at, is_user, is_generating, is_archived, character_id, content)
				 SELECT ?, created_at, is_user, is_generating, is_archived, character_id, content
				 FROM chat_messages
				 WHERE chat_session_id = ? AND id <= ?;`
		args = []any{newSessionId, sessionId, messageId}
		if err = database.UpdateRecord(query, args); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	ChatSessionCreatedSignal.EmitBG(newSession)
	return newSession, nil
}
