package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/repository"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrTooMany      = errors.New("too many attempts")
	ErrInvalidCode  = errors.New("invalid code")
	ErrExpired      = errors.New("expired")
	ErrUnauthorized = errors.New("unauthorized")
)

type BotSender interface {
	SendMessage(chatID int64, msg string) error
}

type Service struct {
	repo      *repository.InMemory
	bot       BotSender
	jwtSecret []byte
}

func New(repo *repository.InMemory, bot BotSender, jwtSecret []byte) *Service {
	return &Service{repo: repo, bot: bot, jwtSecret: jwtSecret}
}

func (s *Service) RegisterUser(tgID int64, username string) error {
	_, err := s.repo.CreateUser(tgID, username)
	if errors.Is(err, repository.ErrUserExists) {
		return ErrConflict
	}
	return err
}

func genCode() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	n := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	return fmt.Sprintf("%06d", n), nil
}

func (s *Service) SendCode(username string) (time.Time, error) {
	user, ok := s.repo.GetUserByUsername(username)
	if !ok {
		return time.Time{}, ErrNotFound
	}
	if c, ok := s.repo.GetCode(username); ok {
		if time.Since(c.SentAt) < time.Minute {
			return time.Time{}, ErrTooMany
		}
	}
	code, err := genCode()
	if err != nil {
		return time.Time{}, err
	}
	expires := time.Now().Add(5 * time.Minute)
	s.repo.SaveCode(username, code, expires)
	if s.bot != nil {
		err = s.bot.SendMessage(user.TGID, fmt.Sprintf("Ваш код: %s", code))
		if err != nil {
			slog.Warn("failed to send code via bot", "username", username, "error", err)
		}
	}
	return expires, nil
}

func (s *Service) VerifyCode(username, code string) (*entities.TokenPair, error) {
	user, ok := s.repo.GetUserByUsername(username)
	if !ok {
		return nil, ErrNotFound
	}
	c, ok := s.repo.GetCode(username)
	if !ok {
		return nil, ErrInvalidCode
	}
	if time.Now().After(c.ExpiresAt) {
		s.repo.DeleteCode(username)
		return nil, ErrExpired
	}
	if c.Attempts >= 3 {
		return nil, ErrTooMany
	}
	if c.Code != code {
		c.Attempts++
		s.repo.UpdateCode(username, c)
		return nil, ErrUnauthorized
	}
	s.repo.DeleteCode(username)
	pair, err := s.generateTokens(user.ID, user.Roles)
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func (s *Service) RefreshToken(refreshToken string) (*entities.TokenPair, error) {
	payload, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, ErrUnauthorized
	}
	jti, ok := payload["jti"].(string)
	if !ok {
		return nil, errors.New("invalid token payload: missing jti")
	}
	rt, ok := s.repo.GetRefreshToken(jti)
	if !ok || rt.Revoked || rt.ExpiresAt.Before(time.Now()) {
		return nil, ErrUnauthorized
	}
	userID, ok := payload["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token payload: missing sub")
	}
	pair, err := s.generateTokens(userID, nil)
	if err != nil {
		return nil, err
	}
	s.repo.RevokeRefreshToken(rt.JTI)
	return pair, nil
}

func (s *Service) Logout(accessToken, refreshToken string) error {
	if accessToken != "" {
		payload, err := s.parseToken(accessToken)
		if err != nil {
			return ErrUnauthorized
		}
		if payload != nil {
			if jti, ok := payload["jti"].(string); ok {
				s.repo.BlacklistToken(jti)
			}
		}
	}

	if refreshToken != "" {
		payload, err := s.parseToken(refreshToken)
		if err != nil {
			return ErrUnauthorized
		}
		if payload != nil {
			if jti, ok := payload["jti"].(string); ok {
				s.repo.RevokeRefreshToken(jti)
			}
		}
	}

	return nil
}

func (s *Service) generateTokens(userID string, roles []string) (*entities.TokenPair, error) {
	accessJTI := uuid.NewString()
	refreshJTI := uuid.NewString()
	now := time.Now()

	accessPayload := map[string]interface{}{
		"sub":   userID,
		"roles": roles,
		"iat":   now.Unix(),
		"exp":   now.Add(30 * time.Minute).Unix(),
		"jti":   accessJTI,
	}
	refreshPayload := map[string]interface{}{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(30 * 24 * time.Hour).Unix(),
		"jti": refreshJTI,
	}
	access, err := s.sign(accessPayload)
	if err != nil {
		return nil, err
	}
	refresh, err := s.sign(refreshPayload)
	if err != nil {
		return nil, err
	}
	s.repo.StoreRefreshToken(entities.RefreshToken{
		JTI:       refreshJTI,
		UserID:    userID,
		IssuedAt:  now,
		ExpiresAt: now.Add(30 * 24 * time.Hour),
		Revoked:   false,
	})
	return &entities.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *Service) sign(payload map[string]interface{}) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}
	payloadEnc := base64.RawURLEncoding.EncodeToString(bodyBytes)
	signingInput := fmt.Sprintf("%s.%s", header, payloadEnc)
	mac := hmac.New(sha256.New, s.jwtSecret)
	mac.Write([]byte(signingInput))
	sig := mac.Sum(nil)
	signature := base64.RawURLEncoding.EncodeToString(sig)
	return signingInput + "." + signature, nil
}

func (s *Service) parseToken(token string) (map[string]interface{}, error) {
	parts := split(token, '.')
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}
	sigInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, s.jwtSecret)
	mac.Write([]byte(sigInput))
	if !hmac.Equal(mac.Sum(nil), mustDecode(parts[2])) {
		return nil, errors.New("bad signature")
	}
	payloadData, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(payloadData, &payload); err != nil {
		return nil, err
	}
	if exp, ok := payload["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, ErrExpired
		}
	}
	return payload, nil
}

func split(s string, sep byte) []string {
	var parts []string
	last := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			parts = append(parts, s[last:i])
			last = i + 1
		}
	}
	parts = append(parts, s[last:])
	return parts
}

func mustDecode(s string) []byte {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		slog.Error("failed to decode base64", "error", err)
	}
	return b
}
