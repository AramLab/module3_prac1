package service

import (
	"github.com/AramLab/module3_prac1/internal/repo"
	"github.com/AramLab/module3_prac1/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
)

type Service interface {
	CreateUser(ctx *fiber.Ctx) error
	CreateTask(ctx *fiber.Ctx) error
	GetTaskById(ctx *fiber.Ctx) error
	GetTasks(ctx *fiber.Ctx) error
	GetTaskByUserId(ctx *fiber.Ctx) error
	GetTasksByUserId(ctx *fiber.Ctx) error
	UpdateTaskStatus(ctx *fiber.Ctx) error
	DeleteTask(ctx *fiber.Ctx) error
}

type service struct {
	repo repo.Repository
	log  *zap.SugaredLogger
}

func NewService(repo repo.Repository, log *zap.SugaredLogger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) CreateTask(ctx *fiber.Ctx) error {
	var request TaskRequest
	if err := ctx.BodyParser(&request); err != nil {
		s.log.Errorf("CreateTask: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := validator.Validate(ctx.Context(), request); err != nil {
		s.log.Errorf("CreateTask validation: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 🔍 Проверка существования пользователя по ID
	user, err := s.repo.GetUserByID(ctx.Context(), request.UserID)
	if err != nil {
		s.log.Errorf("CreateTask: user not found: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	// Лог (необязательно): кто создаёт задачу
	s.log.Infof("Creating task for user: %s", user.Username)

	task := repo.Task{
		UserID:      request.UserID,
		Title:       request.Title,
		Description: request.Description,
		Status:      request.Status,
	}

	id, err := s.repo.CreateTask(ctx.Context(), task)
	if err != nil {
		s.log.Errorf("CreateTask: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create task"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"id": id})
}

func (s *service) CreateUser(ctx *fiber.Ctx) error {
	var request UserRequest
	if err := ctx.BodyParser(&request); err != nil {
		s.log.Errorf("CreateUser: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := validator.Validate(ctx.Context(), request); err != nil {
		s.log.Errorf("CreateUser validation: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user := repo.User{
		Username: request.Username,
		Password: request.Password, // ⚠️ В проде здесь должен быть хеш
	}

	id, err := s.repo.CreateUser(ctx.Context(), user)
	if err != nil {
		s.log.Errorf("CreateUser: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"id": id})
}

func (s *service) GetTaskById(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	task, err := s.repo.GetTaskById(ctx.Context(), id)
	if err != nil {
		s.log.Errorf("GetTaskById: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(task)
}

func (s *service) GetTasks(ctx *fiber.Ctx) error {
	tasks, err := s.repo.GetTasks(ctx.Context())
	if err != nil {
		s.log.Errorf("GetTasks: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get tasks"})
	}
	return ctx.Status(fiber.StatusOK).JSON(tasks)
}

func (s *service) GetTaskByUserId(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Params("userID"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid userID"})
	}

	task, err := s.repo.GetTaskByUserID(ctx.Context(), userID)
	if err != nil {
		s.log.Errorf("GetTaskByUserId: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found for user"})
	}

	return ctx.Status(fiber.StatusOK).JSON(task)
}

func (s *service) GetTasksByUserId(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Params("userID"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid userID"})
	}

	tasks, err := s.repo.GetTasksByUserID(ctx.Context(), userID)
	if err != nil {
		s.log.Errorf("GetTasksByUserId: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get tasks for user"})
	}

	return ctx.Status(fiber.StatusOK).JSON(tasks)
}

func (s *service) UpdateTaskStatus(ctx *fiber.Ctx) error {
	var req TaskStatusUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := validator.Validate(ctx.Context(), req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.repo.UpdateTaskStatus(ctx.Context(), req.ID, req.Status); err != nil {
		s.log.Errorf("UpdateTaskStatus: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update task status"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "task status updated"})
}

func (s *service) DeleteTask(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := s.repo.DeleteTask(ctx.Context(), id); err != nil {
		s.log.Errorf("DeleteTask: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete task"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "task deleted"})
}
