package session

import (
	"sync"

	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
)

type Session struct {
	sessionData map[string][]schema.WordStat
	dataMutex   sync.RWMutex
}

func New() *Session {
	return &Session{
		sessionData: make(map[string][]schema.WordStat),
		dataMutex: sync.RWMutex{},
	}
}

func(s *Session) GetState(sessionID string) ([]schema.WordStat, bool) {
	s.dataMutex.RLock()
	defer s.dataMutex.RUnlock()
	stat, exists := s.sessionData[sessionID]
	return stat, exists
}

func (s *Session) SetState(sessionID string, stats []schema.WordStat) error {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()
	s.sessionData[sessionID] = stats
	return nil
}
