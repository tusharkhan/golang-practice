package models

import (
	"course/helper"
	"course/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

const (
	DefaultTimeDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB             *sql.DB
	BytesPerToken  int
	Duration       time.Duration
	SessionService *SessionService
}

func (pr *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)
	var userId int
	err := pr.DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userId)

	if err != nil {
		return nil, fmt.Errorf("User fot found", err)
	}

	bytesPertoken := pr.BytesPerToken
	if bytesPertoken < rand.SessionTokenBytes {
		bytesPertoken = rand.SessionTokenBytes
	}

	token, tokenError := rand.String(bytesPertoken)
	if tokenError != nil {
		return nil, fmt.Errorf("Error in creating token", tokenError)
	}

	duration := pr.Duration
	if duration <= 0 {
		duration = DefaultTimeDuration
	}

	tok := sha256.Sum256([]byte(token))
	hashToken := base64.URLEncoding.EncodeToString(tok[:])

	psReste := PasswordReset{
		UserID:    userId,
		Token:     token,
		TokenHash: hashToken,
		ExpiresAt: time.Now().Add(duration),
	}

	var sqlString string = `INSERT INTO password_resets (user_id, token_hash, expires_at) 
			VALUES($1, $2, $3) 
			ON CONFLICT (user_id) 
			DO UPDATE SET token_hash = $2, expires_at = $3 
			RETURNING id`
	err = pr.DB.QueryRow(sqlString, psReste.UserID, psReste.TokenHash, psReste.ExpiresAt).Scan(&psReste.ID)

	if err != nil {
		return nil, fmt.Errorf("Error in inserting token")
	}

	return &psReste, nil
}

func (pr *PasswordResetService) Consume(newPassword string, user_id int) (*User, error) {

	var selectUser string = "SELECT * FROM users WHERE id = $1"
	var user User

	selectUserQueryError := pr.DB.QueryRow(selectUser, user_id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_at)

	if selectUserQueryError != nil {
		return nil, fmt.Errorf("Selected User not found", selectUserQueryError)
	}

	var updateUser string = "UPDATE users SET password = $1 WHERE id = $2 RETURNING id"

	pass, hashError := helper.HashString(newPassword)

	if hashError != nil {
		return nil, hashError
	}

	updateError := pr.DB.QueryRow(updateUser, pass, user_id).Scan(&user.ID)

	if updateError != nil {
		return nil, fmt.Errorf("Error in updating user ", updateError)
	}

	var deleteSessionQuery string = "DELETE FROM sessions WHERE user_id = $1"

	_, deleteSessionError := pr.DB.Exec(deleteSessionQuery, user_id)

	if deleteSessionError != nil {
		return nil, fmt.Errorf("Error in deleting session ", updateError)
	}

	_, passwordResetsDelete := pr.DB.Exec("DELETE FROM password_resets WHERE user_id = $1", user_id)

	if passwordResetsDelete != nil {
		return nil, fmt.Errorf("Error in deleting session ", updateError)
	}

	return &user, nil
}
