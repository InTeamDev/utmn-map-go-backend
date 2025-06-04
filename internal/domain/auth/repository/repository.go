package repository

import (
	"sync"
	"time"

	"errors"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/entities"
	"github.com/google/uuid"
)

type InMemory struct {
	mu        sync.RWMutex
	users     map[string]entities.User         // username -> user
	tgIndex   map[int64]string                 // tgID -> username
	codes     map[string]entities.AuthCode     // username -> code
	refresh   map[string]entities.RefreshToken // jti -> refresh token
	blacklist map[string]time.Time             // jti -> time
}

func NewInMemory() *InMemory {
	return &InMemory{
		users:     make(map[string]entities.User),
		tgIndex:   make(map[int64]string),
		codes:     make(map[string]entities.AuthCode),
		refresh:   make(map[string]entities.RefreshToken),
		blacklist: make(map[string]time.Time),
	}
}

func (r *InMemory) CreateUser(tgID int64, username string) (entities.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[username]; ok {
		return entities.User{}, ErrUserExists
	}
	if _, ok := r.tgIndex[tgID]; ok {
		return entities.User{}, ErrUserExists
	}
	user := entities.User{ID: uuid.NewString(), TGID: tgID, Username: username, Roles: []string{"user"}, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	r.users[username] = user
	r.tgIndex[tgID] = username
	return user, nil
}

func (r *InMemory) GetUserByUsername(username string) (entities.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[username]
	return u, ok
}

func (r *InMemory) SaveCode(username, code string, exp time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.codes[username] = entities.AuthCode{Code: code, ExpiresAt: exp, Attempts: 0, SentAt: time.Now()}
}

func (r *InMemory) GetCode(username string) (entities.AuthCode, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.codes[username]
	return c, ok
}

func (r *InMemory) UpdateCode(username string, c entities.AuthCode) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.codes[username] = c
}

func (r *InMemory) DeleteCode(username string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.codes, username)
}

func (r *InMemory) StoreRefreshToken(t entities.RefreshToken) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.refresh[t.JTI] = t
}

func (r *InMemory) GetRefreshToken(jti string) (entities.RefreshToken, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.refresh[jti]
	return t, ok
}

func (r *InMemory) RevokeRefreshToken(jti string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t, ok := r.refresh[jti]; ok {
		t.Revoked = true
		r.refresh[jti] = t
	}
}

func (r *InMemory) BlacklistToken(jti string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blacklist[jti] = time.Now()
}

func (r *InMemory) IsBlacklisted(jti string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.blacklist[jti]
	return ok
}

var ErrUserExists = errors.New("user exists")
