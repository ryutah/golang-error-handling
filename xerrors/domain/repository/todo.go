package repository

import (
	"context"

	"github.com/google/uuid"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/ryutah/golang-error-handling/xerrors/domain/model"
	"github.com/ryutah/golang-error-handling/xerrors/internal/apperrors"
	"golang.org/x/xerrors"
)

var tablenames = struct {
	todo string
}{
	todo: "todo",
}

var indexNames = struct {
	todo struct {
		id       string
		title    string
		finished string
	}
}{
	todo: struct {
		id       string
		title    string
		finished string
	}{
		id:       "id",
		title:    "title",
		finished: "finished",
	},
}

var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		tablenames.todo: {
			Name: tablenames.todo,
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:   indexNames.todo.id,
					Unique: true,
					Indexer: &memdb.StringFieldIndex{
						Field:     "ID",
						Lowercase: false,
					},
				},
				"title": {
					Name:   indexNames.todo.title,
					Unique: false,
					Indexer: &memdb.StringFieldIndex{
						Field:     "Title",
						Lowercase: false,
					},
				},
				"finished": {
					Name:   indexNames.todo.finished,
					Unique: false,
					Indexer: &memdb.BoolFieldIndex{
						Field: "Finished",
					},
				},
			},
		},
	},
}

type Todo struct {
	store *memdb.MemDB
}

func NewTodoRepsository() (*Todo, error) {
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, xerrors.Errorf("%w: failed to create onmemory database due to %v", apperrors.ErrInternal, err)
	}
	return &Todo{
		store: db,
	}, nil
}

func (t *Todo) NextID(ctx context.Context) (model.TodoID, error) {
	id, err := model.NewTodoID(uuid.NewString())
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	return id, nil
}

func (t *Todo) Store(ctx context.Context, todo model.Todo) error {
	tx := t.store.Txn(true)
	if err := tx.Insert(tablenames.todo, &todo); err != nil {
		return xerrors.Errorf("%w: failed to insert todo(%#v) due to %v", apperrors.ErrInternal, todo, err)
	}

	tx.Commit()
	return nil
}

func (t *Todo) Get(ctx context.Context, id model.TodoID) (*model.Todo, error) {
	tx := t.store.Txn(false)
	raw, err := tx.First(tablenames.todo, indexNames.todo.id, string(id))
	if err != nil {
		return nil, xerrors.Errorf("%w: failed to get todo(%v) due to %v", apperrors.ErrInternal, id, err)
	} else if raw == nil {
		return nil, xerrors.Errorf("%w: todo(id: %v) is not found", apperrors.ErrNotFound, id)
	}
	return raw.(*model.Todo), nil
}
