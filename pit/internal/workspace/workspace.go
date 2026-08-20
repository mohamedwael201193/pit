package workspace

import (
	"fmt"
	"sync"
	"time"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

type Workspace struct {
	ID        string
	EVM       identity.Address
	Network   config.Network
	CreatedAt time.Time
}

type MemoryBlob struct {
	Kind    string
	Payload []byte
}

type Store struct {
	mu   sync.RWMutex
	byID map[string]Workspace
	mem  map[string]map[string]MemoryBlob // ws -> key -> blob
}

func NewStore() *Store {
	return &Store{
		byID: make(map[string]Workspace),
		mem:  make(map[string]map[string]MemoryBlob),
	}
}

func (s *Store) Create(evm identity.Address, net config.Network) (Workspace, error) {
	if evm == "" {
		return Workspace{}, fmt.Errorf("missing wallet")
	}
	if err := config.RejectGlobalUser(string(evm)); err != nil {
		return Workspace{}, err
	}
	ws := Workspace{
		ID:        identity.NewWorkspaceID(),
		EVM:       evm,
		Network:   net,
		CreatedAt: time.Now().UTC(),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID[ws.ID] = ws
	s.mem[ws.ID] = make(map[string]MemoryBlob)
	return ws, nil
}

func (s *Store) Get(id string) (Workspace, error) {
	parsed, err := identity.ParseWorkspaceID(id)
	if err != nil {
		return Workspace{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	ws, ok := s.byID[parsed]
	if !ok {
		return Workspace{}, fmt.Errorf("not found")
	}
	return ws, nil
}

func memKey(kind, name string) string {
	return kind + "/" + name
}

func (s *Store) PutMemory(wsID, kind, name string, payload []byte) error {
	if _, err := s.Get(wsID); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.mem[wsID] == nil {
		s.mem[wsID] = make(map[string]MemoryBlob)
	}
	s.mem[wsID][memKey(kind, name)] = MemoryBlob{Kind: kind, Payload: payload}
	return nil
}

func (s *Store) GetMemory(wsID, kind, name string) (MemoryBlob, error) {
	if _, err := s.Get(wsID); err != nil {
		return MemoryBlob{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.mem[wsID][memKey(kind, name)]
	if !ok {
		return MemoryBlob{}, fmt.Errorf("not found")
	}
	out := make([]byte, len(b.Payload))
	copy(out, b.Payload)
	return MemoryBlob{Kind: b.Kind, Payload: out}, nil
}

func (s *Store) ListMemory(wsID string) ([]string, error) {
	if _, err := s.Get(wsID); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.mem[wsID]))
	for k := range s.mem[wsID] {
		keys = append(keys, k)
	}
	return keys, nil
}
