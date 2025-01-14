package memory

import (
	"context"
	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

// Repository defines a rating repository.
type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

// New creates a new in-memory repository.
func New() *Repository {
	return &Repository{
		data: map[model.RecordType]map[model.RecordID][]model.Rating{},
	}
}

// Get retrieves all ratings for a given record.
func (r *Repository) Get(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
) ([]model.Rating, error) {
	records, ok := r.data[recordType]
	if !ok {
		return nil, repository.ErrNotFound
	}
	ratings, ok := records[recordID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return ratings, nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
	rating *model.Rating,
) error {
	_, ok := r.data[recordType]
	if !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}
	r.data[recordType][recordID] = append(r.data[recordType][recordID], *rating)
	return nil
}
