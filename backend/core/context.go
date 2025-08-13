package core

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
)

type ChatQuestContext struct {
	ctx    context.Context
	db     *sql.DB
	logger *zap.Logger
}

func NewChatQuestContext(
	rootCtx context.Context,
	db *sql.DB,
	log *zap.Logger,
) *ChatQuestContext {
	return &ChatQuestContext{
		ctx:    rootCtx,
		db:     db,
		logger: log,
	}
}

func (cq *ChatQuestContext) Context() context.Context {
	return cq.ctx
}

func (cq *ChatQuestContext) DB() *sql.DB {
	return cq.db
}

func (cq *ChatQuestContext) Logger() *zap.Logger {
	return cq.logger
}

func (cq *ChatQuestContext) WithContext(newContext context.Context) *ChatQuestContext {
	return &ChatQuestContext{
		ctx:    newContext,
		db:     cq.db,
		logger: cq.logger,
	}
}
