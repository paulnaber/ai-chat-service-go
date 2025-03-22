// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createChatStmt, err = db.PrepareContext(ctx, createChat); err != nil {
		return nil, fmt.Errorf("error preparing query CreateChat: %w", err)
	}
	if q.createMessageStmt, err = db.PrepareContext(ctx, createMessage); err != nil {
		return nil, fmt.Errorf("error preparing query CreateMessage: %w", err)
	}
	if q.getChatStmt, err = db.PrepareContext(ctx, getChat); err != nil {
		return nil, fmt.Errorf("error preparing query GetChat: %w", err)
	}
	if q.getChatsByUserEmailStmt, err = db.PrepareContext(ctx, getChatsByUserEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetChatsByUserEmail: %w", err)
	}
	if q.getMessageStmt, err = db.PrepareContext(ctx, getMessage); err != nil {
		return nil, fmt.Errorf("error preparing query GetMessage: %w", err)
	}
	if q.getMessagesByChatIDStmt, err = db.PrepareContext(ctx, getMessagesByChatID); err != nil {
		return nil, fmt.Errorf("error preparing query GetMessagesByChatID: %w", err)
	}
	if q.updateChatLastActiveStmt, err = db.PrepareContext(ctx, updateChatLastActive); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateChatLastActive: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createChatStmt != nil {
		if cerr := q.createChatStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createChatStmt: %w", cerr)
		}
	}
	if q.createMessageStmt != nil {
		if cerr := q.createMessageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createMessageStmt: %w", cerr)
		}
	}
	if q.getChatStmt != nil {
		if cerr := q.getChatStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getChatStmt: %w", cerr)
		}
	}
	if q.getChatsByUserEmailStmt != nil {
		if cerr := q.getChatsByUserEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getChatsByUserEmailStmt: %w", cerr)
		}
	}
	if q.getMessageStmt != nil {
		if cerr := q.getMessageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getMessageStmt: %w", cerr)
		}
	}
	if q.getMessagesByChatIDStmt != nil {
		if cerr := q.getMessagesByChatIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getMessagesByChatIDStmt: %w", cerr)
		}
	}
	if q.updateChatLastActiveStmt != nil {
		if cerr := q.updateChatLastActiveStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateChatLastActiveStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                       DBTX
	tx                       *sql.Tx
	createChatStmt           *sql.Stmt
	createMessageStmt        *sql.Stmt
	getChatStmt              *sql.Stmt
	getChatsByUserEmailStmt  *sql.Stmt
	getMessageStmt           *sql.Stmt
	getMessagesByChatIDStmt  *sql.Stmt
	updateChatLastActiveStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                       tx,
		tx:                       tx,
		createChatStmt:           q.createChatStmt,
		createMessageStmt:        q.createMessageStmt,
		getChatStmt:              q.getChatStmt,
		getChatsByUserEmailStmt:  q.getChatsByUserEmailStmt,
		getMessageStmt:           q.getMessageStmt,
		getMessagesByChatIDStmt:  q.getMessagesByChatIDStmt,
		updateChatLastActiveStmt: q.updateChatLastActiveStmt,
	}
}
