package models

import (
	"course/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
)

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(use_id int) (*Session, error) {
	token, tokenError := rand.SessionToken()

	if tokenError != nil {
		return nil, fmt.Errorf("error in creating token %w", tokenError)
	}

	var session Session = Session{
		UserID:    use_id,
		Token:     token,
		TokenHash: ss.hashToken(token),
	}

	// database, databaseError := helper.ConnectDatabase()

	// if databaseError != nil {
	// 	panic(databaseError)
	// }

	// ss.DB = database

	row := ss.DB.QueryRow(`INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2) RETURNING id;`, session.UserID, session.TokenHash)
	sqlError := row.Scan(&session.ID)

	if sqlError != nil {
		return nil, fmt.Errorf("Errror in creating session %w", sqlError)
	}

	return &session, nil
}

func (ss *SessionService) Update() (*Session, error) {
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}

func (ss *SessionService) hashToken(token string) string {
	tok := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tok[:])
}
