package store

import (
	"database/sql"
	"time"

	"github.com/Rizzwaan/workoutVerse/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

// TokenStore interface defines methods for token operations
type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int64, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensByUserID(userID int64, scope string) error
}

func (s *PostgresTokenStore) CreateNewToken(userID int64, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = s.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `INSERT INTO tokens (user_id, token, hash, expiry, scope) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(query, token.UserID, token.PlainText, token.Hash, token.Expiry, token.Scope)
	return err
}

func (s *PostgresTokenStore) DeleteAllTokensByUserID(userID int64, scope string) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`
	_, err := s.db.Exec(query, userID, scope)
	return err
}
