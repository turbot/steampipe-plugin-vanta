package rest_api

import "sync"

// StaticTokenStore implements TokenStore for static token authentication
type StaticTokenStore struct {
	mutex     sync.RWMutex
	tokenType string
	token     string
}

// NewStaticTokenStore creates a new static token store
func NewStaticTokenStore(tokenType, token string) *StaticTokenStore {
	return &StaticTokenStore{
		tokenType: tokenType,
		token:     token,
	}
}

// GetToken returns the stored token
func (s *StaticTokenStore) GetToken() (string, string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.tokenType, s.token
}

// SetToken updates the stored token
func (s *StaticTokenStore) SetToken(tokenType, token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.tokenType = tokenType
	s.token = token
}

// BearerTokenStore is a convenience wrapper for Bearer token authentication
type BearerTokenStore struct {
	*StaticTokenStore
}

// NewBearerTokenStore creates a new Bearer token store
func NewBearerTokenStore(token string) *BearerTokenStore {
	return &BearerTokenStore{
		StaticTokenStore: NewStaticTokenStore("Bearer", token),
	}
}

// SetBearerToken updates the Bearer token
func (b *BearerTokenStore) SetBearerToken(token string) {
	b.SetToken("Bearer", token)
}
