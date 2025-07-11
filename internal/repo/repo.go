package repo

import (
	"context"
	"fmt"
	"github.com/AramLab/module3_prac1/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	createTaskQuery       = `INSERT INTO tasks (user_id, title, description, status) VALUES ($1, $2, $3, $4) RETURNING id;`
	getTaskByIdQuery      = `SELECT id, user_id, title, description, status, created_at FROM tasks WHERE id = $1;`
	getTasksQuery         = `SELECT id, user_id, title, description, status, created_at FROM tasks;`
	getTaskByUserIDQuery  = `SELECT id, user_id, title, description, status, created_at FROM tasks WHERE user_id = $1 LIMIT 1;`
	getTasksByUserIDQuery = `SELECT id, user_id, title, description, status, created_at FROM tasks WHERE user_id = $1;`
	getUserByIDQuery      = `SELECT id, username, password FROM users WHERE id = $1;`
	updateTaskStatusQuery = `UPDATE tasks SET status = $1 WHERE id = $2;`
	deleteTaskQuery       = `DELETE FROM tasks WHERE id = $1;`
	createUserQuery       = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id;`
)

type repository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	CreateTask(ctx context.Context, task Task) (int, error)
	CreateUser(ctx context.Context, user User) (int, error)

	GetTaskById(ctx context.Context, id int) (*Task, error)
	GetTasks(ctx context.Context) ([]Task, error)
	GetTaskByUserID(ctx context.Context, userID int) (*Task, error)
	GetTasksByUserID(ctx context.Context, userID int) ([]Task, error)
	GetUserByID(ctx context.Context, id int) (*User, error)

	UpdateTaskStatus(ctx context.Context, id int, status string) error
	DeleteTask(ctx context.Context, id int) error
}

func NewRepository(ctx context.Context, cfg config.PostgreSQL) (Repository, error) {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse connection string")
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pool")
	}

	return &repository{pool}, nil
}

func (r *repository) CreateTask(ctx context.Context, task Task) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, createTaskQuery, task.UserID, task.Title, task.Description, task.Status).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create task")
	}
	return id, nil
}

func (r *repository) CreateUser(ctx context.Context, user User) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, createUserQuery, user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create user")
	}
	return id, nil
}

func (r *repository) GetTaskById(ctx context.Context, id int) (*Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, getTaskByIdQuery, id).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get task by id")
	}
	return &task, nil
}

func (r *repository) GetTasks(ctx context.Context) ([]Task, error) {
	rows, err := r.pool.Query(ctx, getTasksQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tasks")
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan task")
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *repository) GetTaskByUserID(ctx context.Context, userID int) (*Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, getTaskByUserIDQuery, userID).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get task by userID")
	}
	return &task, nil
}

func (r *repository) GetTasksByUserID(ctx context.Context, userID int) ([]Task, error) {
	rows, err := r.pool.Query(ctx, getTasksByUserIDQuery, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tasks by userID")
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan task")
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *repository) GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx, getUserByIDQuery, id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // пользователь не найден — это не ошибка, просто nil
		}
		return nil, errors.Wrap(err, "failed to get user by id")
	}
	return &user, nil
}

func (r *repository) UpdateTaskStatus(ctx context.Context, id int, status string) error {
	_, err := r.pool.Exec(ctx, updateTaskStatusQuery, status, id)
	if err != nil {
		return errors.Wrap(err, "failed to update task status")
	}
	return nil
}

func (r *repository) DeleteTask(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, deleteTaskQuery, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete task")
	}
	return nil
}
