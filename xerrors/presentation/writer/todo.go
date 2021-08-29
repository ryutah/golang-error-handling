package writer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ryutah/golang-error-handling/xerrors/application"
	"github.com/ryutah/golang-error-handling/xerrors/internal/apperrors"
)

type Todo struct {
	application *application.Todo
}

func NewTodo(app *application.Todo) *Todo {
	return &Todo{
		application: app,
	}
}

type addResponse struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Finished bool   `json:"finished,omitempty"`
}

func (t *Todo) Add(ctx context.Context, w io.Writer, title string) {
	newTodo, err := t.application.Add(ctx, application.AddTodoReq{
		Title: title,
	})
	if err != nil {
		renderError(w, err)
		return
	}
	render(w, addResponse{
		ID:       newTodo.ID,
		Title:    newTodo.Title,
		Finished: newTodo.Finished,
	})
}

type getTodoResponse struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Finished bool   `json:"finished,omitempty"`
}

func (t *Todo) Get(ctx context.Context, w io.Writer, id string) {
	todo, err := t.application.Get(ctx, id)
	if err != nil {
		renderError(w, err)
		return
	}
	render(w, getTodoResponse{
		ID:       todo.ID,
		Title:    todo.Title,
		Finished: todo.Finished,
	})
}

func render(w io.Writer, v interface{}) {
	body, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		renderError(w, err)
		return
	}
	fmt.Fprintln(w, string(body))
}

func renderError(w io.Writer, err error) {
	var payload map[string]string
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		payload = map[string]string{
			"error": "指定されたデータが見つかりませんでした。データはすでに削除されているか、指定した内容に誤りがないか確認してください",
		}
	case errors.Is(err, apperrors.ErrInvalidParameter):
		payload = map[string]string{
			"error": "入力内容に誤りがあります。内容を確認し、再実行してください",
		}
	default:
		payload = map[string]string{
			"error": "不明なエラーが発生しました。管理者に連絡をしてください",
		}
	}

	body, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintln(w, string(body))
}
