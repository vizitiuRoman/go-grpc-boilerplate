package repo

import (
	"context"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/db"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/log"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/model"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/domain/repo"
	"go.uber.org/zap"
)

type todoRepoFactory struct {
	logger log.Logger
}

func NewTodoRepoFactory(logger log.Logger) repo.TodoRepoFactory {
	return &todoRepoFactory{
		logger: logger,
	}
}

func (f *todoRepoFactory) Create(ctx context.Context, db db.DB) repo.TodoRepo {
	return &todoRepo{
		logger: f.logger.WithComponent(ctx, "TodoRepo"),
		db:     db,
	}
}

type todoRepo struct {
	logger log.Logger
	db     db.DB
}

func (r *todoRepo) Find(ctx context.Context, id int64) (*model.Todo, error) {
	const q = `SELECT id, name, description FROM todo WHERE id = $1`

	var todo model.Todo
	row := r.db.QueryRowContext(ctx, q, id)
	err := wrapErrNoRows(row.Scan(&todo.ID, &todo.Name, &todo.Description))
	if err != nil {
		r.logger.WithMethod(ctx, "Find").Error("query error", zap.Error(err))
		return nil, err
	}

	return &todo, err
}

func (r *todoRepo) FindAll(ctx context.Context) ([]*model.Todo, error) {
	logger := r.logger.WithMethod(ctx, "FindAll")

	const q = `SELECT id, name, description FROM todo`

	todos := make([]*model.Todo, 0)

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return nil, err
	}
	if err = rows.Err(); err != nil {
		logger.Error("reading error", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		todo := model.Todo{}
		if err = rows.Scan(&todo.ID, &todo.Name, &todo.Description); err != nil {
			logger.Error("scan error", zap.Error(err))
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, err
}

func (r *todoRepo) Create(ctx context.Context, input *model.Todo) (*model.Todo, error) {
	const q = `INSERT INTO todo (name, description) VALUES ($1, $2) RETURNING id, name, description`

	var todo model.Todo

	row := r.db.QueryRowContext(ctx, q, input.Name, input.Description)
	if err := wrapUniqueViolation(row.Scan(&todo.ID, &todo.Name, &todo.Description)); err != nil {
		r.logger.WithMethod(ctx, "Create").Error("insert todo", zap.Error(err))
		return nil, err
	}

	return &todo, nil
}

func (r *todoRepo) Update(ctx context.Context, input *model.Todo) (*model.Todo, error) {
	const q = `
		UPDATE todo
		SET name = $1,
		    description = $2
		WHERE id = $3
		RETURNING id, name, description
	`

	var todo model.Todo

	row := r.db.QueryRowContext(ctx, q, input.Name, input.Description, input.ID)
	if err := wrapErrNoRows(row.Scan(&todo.ID, &todo.Name, &todo.Description)); err != nil {
		r.logger.WithMethod(ctx, "Update").Error("update todo", zap.Error(err))
		return nil, err
	}

	return &todo, nil
}

func (r *todoRepo) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM todo WHERE id = $1`

	rows, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		r.logger.WithMethod(ctx, "Delete").Error("query error", zap.Error(err))
		return err
	}

	count, _ := rows.RowsAffected()
	if count == 0 {
		r.logger.WithMethod(ctx, "Delete").Error("no rows affected", zap.Error(domain.ErrNotFound))
		return domain.ErrNotFound
	}

	return err
}
