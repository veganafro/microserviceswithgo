package memory

import (
	"context"
	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/model"
	"sync"
)

// Repository defines an in-memory movie metadata repository.
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New creates a new in-memory repository.
func New() *Repository {
	return &Repository{
		data: map[string]*model.Metadata{},
	}
}

// Get retrieves movie metadata using the provided movie id from the in-memory repository.
func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Put updates or creates movie metadata in the in-memory repository.
func (r *Repository) Put(_ context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
