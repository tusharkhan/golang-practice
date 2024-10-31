package models

import (
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
	DB            *sql.DB
	BytesPerToken int
	Duration      time.Duration
}

func (pr *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)
	var userId int
	err := pr.DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userId)
	fmt.Println(userId)
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

	var sqlString string = "INSERT INTO password_resets (user_id, token_hash, expires_at) VALUES($1, $2, $3) ON CONFLICT (user_id) SET token_hash=$3 RETURNING id"
	err = pr.DB.QueryRow(sqlString, psReste.UserID, psReste.TokenHash, psReste.ExpiresAt).Scan(&psReste.ID)
	if err != nil {
		return nil, fmt.Errorf("Error in inserting token")
	}

	return &psReste, nil
}

func (pr *PasswordResetService) Consume(token, newPassword string) (*User, error) {
	return nil, fmt.Errorf("Not implemented")
}
