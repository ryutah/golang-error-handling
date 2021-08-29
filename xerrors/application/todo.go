package application

import (
	"context"
	"log"

	"github.com/ryutah/golang-error-handling/xerrors/domain/model"
	"github.com/ryutah/golang-error-handling/xerrors/domain/repository"
	"golang.org/x/xerrors"
)

type Todo struct {
	repository *repository.Todo
}

func NewTodo(repo *repository.Todo) *Todo {
	return &Todo{
		repository: repo,
	}
}

type (
	AddTodoReq struct {
		Title string
	}
	AddTodoResp struct {
		ID       string
		Title    string
		Finished bool
	}
)

func (t *Todo) Add(ctx context.Context, req AddTodoReq) (resp *AddTodoResp, err error) {
	defer func() { logError(err) }()

	id, err := t.repository.NextID(ctx)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	newTodo, err := model.NewTodo(id, req.Title)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	if err := t.repository.Store(ctx, *newTodo); err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return &AddTodoResp{
		ID:       string(newTodo.ID),
		Title:    newTodo.Title,
		Finished: newTodo.Finished,
	}, nil
}

type GetTodoResp struct {
	ID       string
	Title    string
	Finished bool
}

func (t *Todo) Get(ctx context.Context, id string) (resp *GetTodoResp, err error) {
	defer func() { logError(err) }()

	todoID, err := model.NewTodoID(id)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	todo, err := t.repository.Get(ctx, todoID)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return &GetTodoResp{
		ID:       string(todo.ID),
		Title:    todo.Title,
		Finished: todo.Finished,
	}, nil
}

func logError(err error) {
	if err == nil {
		return
	}
	log.Printf("%+v\n", err)
}
