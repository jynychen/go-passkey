package main

import (
	"crypto/rand"
	"encoding/base64"
	"sync"

	"github.com/go-webauthn/webauthn/webauthn"
)

type InMem struct {
	mu       sync.RWMutex // protects users and sessions
	users    map[string]PasskeyUser
	sessions map[string]*webauthn.SessionData

	log Logger
}

func (i *InMem) GenSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil

}

func NewInMem(log Logger) *InMem {
	return &InMem{
		users:    make(map[string]PasskeyUser),
		sessions: make(map[string]*webauthn.SessionData),
		log:      log,
	}
}

func (i *InMem) GetSession(token string) (*webauthn.SessionData, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	i.log.Printf("[DEBUG] GetSession: %v", i.sessions[token])
	val, ok := i.sessions[token]

	return val, ok
}

func (i *InMem) SaveSession(token string, data webauthn.SessionData) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.log.Printf("[DEBUG] SaveSession: %s - %v", token, data)
	i.sessions[token] = &data
}

func (i *InMem) DeleteSession(token string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.log.Printf("[DEBUG] DeleteSession: %v", token)
	delete(i.sessions, token)
}

func (i *InMem) GetOrCreateUser(userName string) PasskeyUser {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.log.Printf("[DEBUG] GetOrCreateUser: %v", userName)
	if _, ok := i.users[userName]; !ok {
		i.log.Printf("[DEBUG] GetOrCreateUser: creating new user: %v", userName)
		i.users[userName] = &User{
			ID:          []byte(userName),
			DisplayName: userName,
			Name:        userName,
		}
	}

	return i.users[userName]
}

func (i *InMem) SaveUser(user PasskeyUser) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.log.Printf("[DEBUG] SaveUser: %v", user.WebAuthnName())
	i.log.Printf("[DEBUG] SaveUser: %v", user)
	i.users[user.WebAuthnName()] = user
}
