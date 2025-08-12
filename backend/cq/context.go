package cq

import (
	"context"
	"database/sql"
	"log"
)

type ChatQuestContext struct {
	ctx context.Context
	db  *sql.DB
	log *log.Logger
}

func NewChatQuestContext(
	rootCtx context.Context,
	db *sql.DB,
	log *log.Logger,
) *ChatQuestContext {
	return &ChatQuestContext{
		ctx: rootCtx,
		db:  db,
		log: log,
	}
}

func (cq *ChatQuestContext) Context() context.Context {
	return cq.ctx
}

func (cq *ChatQuestContext) DB() *sql.DB {
	return cq.db
}

func (cq *ChatQuestContext) Logger() *log.Logger {
	return cq.log
}

func (cq *ChatQuestContext) WithContext(newContext context.Context) *ChatQuestContext {
	return &ChatQuestContext{
		ctx: newContext,
		db:  cq.db,
		log: cq.log,
	}
}
