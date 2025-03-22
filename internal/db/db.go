package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"ai-chat-service-go/internal/config"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pressly/goose/v3"
)

// Store provides access to all database operations
type Store struct {
	db *sql.DB
}

// Queries is used temporarily until sqlc generates the real code
type Queries struct {
	db DBTX
}

// DBTX is an interface for database operations
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// NewStore creates a new database store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// NewQueries creates a new queries instance
func NewQueries(db DBTX) *Queries {
	return &Queries{
		db: db,
	}
}

// InitDB initializes the database connection
func InitDB(cfg config.DatabaseConfig) (*Store, error) {
	connStr := cfg.GetConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return NewStore(db), nil
}

// RunMigrations runs database migrations using goose
func RunMigrations(cfg config.DatabaseConfig) error {
	db, err := sql.Open("postgres", cfg.GetConnectionString())
	if err != nil {
		return fmt.Errorf("failed to open database connection for migrations: %w", err)
	}
	defer db.Close()

	// Set goose verbosity
	if os.Getenv("VERBOSE") == "true" {
		goose.SetVerbose(true)
	}

	// Set goose dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Determine migration directory
	migrationDir := cfg.MigrationDir
	if !filepath.IsAbs(migrationDir) {
		// If it's not an absolute path, make it relative to the working directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}
		migrationDir = filepath.Join(cwd, migrationDir)
	}

	// Verify migration directory exists
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		return fmt.Errorf("migration directory does not exist: %s", migrationDir)
	}

	// Run migrations
	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Get current version
	version, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get DB version: %w", err)
	}

	log.Printf("Successfully ran migrations, current DB version: %d", version)
	return nil
}

// ExecTx executes a function within a database transaction
func (s *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := NewQueries(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// Temporary implementations of the database operations
// These will be replaced by the generated code from sqlc

// Chat represents a chat in the database
type Chat struct {
	ID             string
	Title          string
	UserEmail      string
	LastActiveDate time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Message represents a message in the database
type Message struct {
	ID         string
	Content    string
	SenderType string
	ChatID     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CreateChatParams holds parameters for creating a chat
type CreateChatParams struct {
	ID             string
	Title          string
	UserEmail      string
	LastActiveDate time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// CreateChat creates a new chat
func (q *Queries) CreateChat(ctx context.Context, arg CreateChatParams) (Chat, error) {
	// This is a placeholder implementation until sqlc generates the real code
	return Chat{
		ID:             arg.ID,
		Title:          arg.Title,
		UserEmail:      arg.UserEmail,
		LastActiveDate: arg.LastActiveDate,
		CreatedAt:      arg.CreatedAt,
		UpdatedAt:      arg.UpdatedAt,
	}, nil
}

// GetChat gets a chat by ID
func (s *Store) GetChat(ctx context.Context, id string) (Chat, error) {
	// This is a placeholder implementation until sqlc generates the real code
	return Chat{}, sql.ErrNoRows
}

// GetChatsByUserEmail gets all chats for a user
func (s *Store) GetChatsByUserEmail(ctx context.Context, email string) ([]Chat, error) {
	// This is a placeholder implementation until sqlc generates the real code
	return []Chat{}, nil
}

// UpdateChatLastActiveParams holds parameters for updating a chat's last active date
type UpdateChatLastActiveParams struct {
	ID             string
	LastActiveDate time.Time
	UpdatedAt      time.Time
}

// UpdateChatLastActive updates a chat's last active date
func (q *Queries) UpdateChatLastActive(ctx context.Context, arg UpdateChatLastActiveParams) error {
	// This is a placeholder implementation until sqlc generates the real code
	return nil
}

// CreateMessageParams holds parameters for creating a message
type CreateMessageParams struct {
	ID         string
	Content    string
	SenderType string
	ChatID     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CreateMessage creates a new message
func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	// This is a placeholder implementation until sqlc generates the real code
	return Message{
		ID:         arg.ID,
		Content:    arg.Content,
		SenderType: arg.SenderType,
		ChatID:     arg.ChatID,
		CreatedAt:  arg.CreatedAt,
		UpdatedAt:  arg.UpdatedAt,
	}, nil
}

// GetMessage gets a message by ID
func (s *Store) GetMessage(ctx context.Context, id string) (Message, error) {
	// This is a placeholder implementation until sqlc generates the real code
	return Message{}, sql.ErrNoRows
}

// GetMessagesByChatID gets all messages for a chat
func (s *Store) GetMessagesByChatID(ctx context.Context, chatID string) ([]Message, error) {
	// This is a placeholder implementation until sqlc generates the real code
	return []Message{}, nil
}