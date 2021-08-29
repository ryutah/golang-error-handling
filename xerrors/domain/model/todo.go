package model

import (
	"github.com/ryutah/golang-error-handling/xerrors/internal/apperrors"
	"golang.org/x/xerrors"
)

type TodoID string

func NewTodoID(s string) (TodoID, error) {
	if s == "" {
		return "", xerrors.Errorf("%w: id must not be blank", apperrors.ErrInvalidParameter)
	}
	return TodoID(s), nil
}

type TodoStatus string

const (
	TodoStatusTodo       TodoStatus = "Todo"
	TodoStatusInprogress TodoStatus = "In Progress"
	TodoStatusDone       TodoStatus = "Done"
)

type Todo struct {
	ID       TodoID
	Title    string
	Finished bool
}

func NewTodo(id TodoID, title string) (*Todo, error) {
	newTodo := &Todo{
		ID: id,
	}
	if err := newTodo.SetTitle(title); err != nil {
		return nil, xerrors.Errorf("failed to create todo: %w", err)
	}
	return newTodo, nil
}

func (t *Todo) SetTitle(title string) error {
	if title == "" {
		return xerrors.Errorf("%w: title should not be blank", apperrors.ErrInvalidParameter)
	}
	t.Title = title
	return nil
}
