package pkg

import (
	"sync"
	"time"
)

type TokenData struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	RoleID       uint   `json:"role_id"`
	Provider     string `json:"provider,omitempty"`
	SessionID    string `json:"sid,omitempty"`
	MFACompleted bool   `json:"mfa,omitempty"`
	ExpiresAt    time.Time
}

var refreshTokenStore = make(map[string]TokenData)
var mu sync.RWMutex

func SaveRefreshToken(token string, data TokenData) {
	mu.Lock()
	defer mu.Unlock()
	key := "refresh:" + token
	refreshTokenStore[key] = data
}

func GetRefreshToken(token string) (*TokenData, bool) {
	mu.RLock()
	defer mu.RUnlock()
	key := "refresh:" + token
	data, exists := refreshTokenStore[key]
	if !exists || time.Now().After(data.ExpiresAt) {
		return nil, false
	}
	return &data, true
}

func DeleteRefreshToken(token string) {
	mu.Lock()
	defer mu.Unlock()
	key := "refresh:" + token
	delete(refreshTokenStore, key)
}
