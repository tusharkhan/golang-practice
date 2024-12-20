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

	checkIfexists, existsError := ss.checkSeeeionExists(use_id)

	if existsError != nil {
		return nil, fmt.Errorf("errror in checking session %w", existsError)
	}

	if checkIfexists {
		return ss.updateToken(&session)
	} else {
		return ss.CreateNewToken(&session)
	}
}

func (ss *SessionService) User(token string) (*User, error) {
	var hashToken string = ss.hashToken(token)
	var getUser User = User{}
	var query string = `SELECT users.id, users.name, users.email, users.created_at FROM sessions
		JOIN users ON users.id = sessions.user_id WHERE sessions.token_hash = $1`
	getUserQueryError := ss.DB.QueryRow(query, hashToken).Scan(&getUser.ID, &getUser.Name, &getUser.Email, &getUser.Created_at)

	if getUserQueryError != nil {
		return nil, getUserQueryError
	}

	return &getUser, nil
}

func (ss *SessionService) hashToken(token string) string {
	tok := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tok[:])
}

func (ss *SessionService) checkSeeeionExists(id int) (bool, error) {
	var exists bool
	var query string = `SELECT EXISTS(SELECT 1 FROM sessions WHERE user_id = $1)`

	queryError := ss.DB.QueryRow(query, id).Scan(&exists)

	if queryError != nil {
		return false, queryError
	}

	return exists, nil
}

func (ss *SessionService) updateToken(updatedSession *Session) (*Session, error) {
	row := ss.DB.QueryRow(`UPDATE sessions SET token_hash = $1 WHERE user_id = $2 RETURNING id;`, updatedSession.TokenHash, updatedSession.UserID)
	fmt.Println(updatedSession)
	scanError := row.Scan(&updatedSession.ID)

	if scanError != nil {
		return nil, scanError
	}

	return updatedSession, nil
}

func (ss *SessionService) CreateNewToken(session *Session) (*Session, error) {

	row := ss.DB.QueryRow(`INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2) RETURNING id;`, session.UserID, session.TokenHash)

	sqlError := row.Scan(&session.ID)

	if sqlError != nil {
		return nil, fmt.Errorf("Errror in creating session %w", sqlError)
	}

	return session, nil
}

func (ss *SessionService) DestroySession(hashToken string) bool {
	var tokenString string = ss.hashToken(hashToken)

	var deleteQuery string = "DELETE FROM sessions WHERE token_hash = $1"

	_, deleteError := ss.DB.Exec(deleteQuery, tokenString)

	return (deleteError == nil)
}
