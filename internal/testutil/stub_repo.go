package testutil

import (
	"context"
	"github.com/yourusername/resume-app/internal/notes"
)

type StubRepo struct {
	Notes []notes.Note
}

func (s *StubRepo) Create(ctx context.Context, n *notes.Note) error {
	n.ID = int64(len(s.Notes) + 1)
	s.Notes = append(s.Notes, *n)
	return nil
}

func (s *StubRepo) Get(ctx context.Context, id int64) (*notes.Note, error) {
	for _, n := range s.Notes { if n.ID == id { return &n, nil } }
	return nil, nil
}

func (s *StubRepo) List(ctx context.Context) ([]notes.Note, error) {
	return s.Notes, nil
}
