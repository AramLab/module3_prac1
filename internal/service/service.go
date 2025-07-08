package service

import (
	"github.com/AramLab/module3_prac1/internal/repo"
	"github.com/AramLab/module3_prac1/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strconv"
)

type Service interface {
	CreateTask(ctx *fiber.Ctx) error
	GetTaskById(ctx *fiber.Ctx) error
	GetTasks(ctx *fiber.Ctx) error
	UpdateTask(ctx *fiber.Ctx) error
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := validator.Validate(ctx.Context(), request); err != nil {
		s.log.Errorf("CreateTask validation: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	task := repo.Task{
		Title:       request.Title,
		Description: request.Description,
		Status:      request.Status,
	}

	id, err := s.repo.CreateTask(task)
	if err != nil {
		s.log.Errorf("CreateTask: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.Wrap(err, "failed to create task").Error(),
		})
	}
	s.log.Debugf("CreateTask: %v", task)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": id,
	})
}

func (s *service) GetTaskById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.log.Errorf("GetTaskById: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid task id",
		})
	}

	req := TaskRequestById{ID: id}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		s.log.Errorf("GetTaskById validation: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	task, err := s.repo.GetTaskById(id)
	if err != nil {
		s.log.Errorf("GetTaskById: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "task not found",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(task)
}

func (s *service) GetTasks(ctx *fiber.Ctx) error {
	tasks, err := s.repo.GetTasks()
	if err != nil {
		s.log.Errorf("GetTasks: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve tasks",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(tasks)
}

func (s *service) UpdateTask(ctx *fiber.Ctx) error {
	var request TaskRequestUpdate
	if err := ctx.BodyParser(&request); err != nil {
		s.log.Errorf("UpdateTask: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := validator.Validate(ctx.Context(), request); err != nil {
		s.log.Errorf("UpdateTask validation: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	task, err := s.repo.GetTaskById(request.ID)
	if err != nil {
		s.log.Errorf("UpdateTask: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "task not found",
		})
	}

	// Для частичного обновления
	//if request.Title != "" {
	//	task.Title = request.Title
	//}
	//if request.Description != "" {
	//	task.Description = request.Description
	//}
	//if request.Status != "" {
	//	task.Status = request.Status
	//}

	s.repo.UpdateTask(*task)

	return ctx.Status(fiber.StatusOK).JSON(task)
}

func (s *service) DeleteTask(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.log.Errorf("DeleteTask: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid task id",
		})
	}

	req := TaskRequestById{ID: id}
	if err := validator.Validate(ctx.Context(), req); err != nil {
		s.log.Errorf("DeleteTask validation: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	s.repo.DeleteTask(id)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "task deleted successfully",
	})
}
