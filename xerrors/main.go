package main

import (
	"context"
	"os"

	"github.com/ryutah/golang-error-handling/xerrors/application"
	"github.com/ryutah/golang-error-handling/xerrors/domain/repository"
	"github.com/ryutah/golang-error-handling/xerrors/presentation/writer"
)

func main() {
	repo, err := repository.NewTodoRepsository()
	if err != nil {
		panic(err)
	}
	pres := writer.NewTodo(application.NewTodo(repo))

	ctx := context.Background()
	pres.Add(ctx, os.Stdout, "new todo!")
	pres.Get(ctx, os.Stdout, "notfound")
}
