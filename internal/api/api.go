package api

import (
	"github.com/AramLab/module3_prac1/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Routers - структура для хранения зависимостей роутов
type Routers struct {
	Service service.Service
}

// NewRouters - конструктор для настройки API
func NewRouters(r *Routers) *fiber.App {
	app := fiber.New()

	// Настройка CORS (разрешенные методы, заголовки, авторизация)
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET, POST, PUT, DELETE",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))

	// Группа маршрутов с авторизацией
	apiGroup := app.Group("/v1")

	// Роут для создания задачи
	apiGroup.Post("/tasks", r.Service.CreateTask)
	apiGroup.Get("/tasks", r.Service.GetTasks)
	apiGroup.Get("/tasks/:id", r.Service.GetTaskById)
	apiGroup.Put("/tasks/:id", r.Service.UpdateTask)
	apiGroup.Delete("/tasks/:id", r.Service.DeleteTask)

	return app
}
