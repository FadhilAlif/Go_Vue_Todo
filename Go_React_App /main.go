package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	// "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"go-react-app/database"
	"go-react-app/logger"
	"go-react-app/models"
)

func loggerFiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		logger.Log.WithFields(logrus.Fields{
			"method":  c.Method(),
			"path":    c.Path(),
			"status":  c.Response().StatusCode(),
			"latency": c.Response().Header.Peek("X-Response-Time"),
		}).Info("HTTP Request")

		return err
	}
}


func main() {
	fmt.Println("Starting Go-React-App...")

	// Init logrus
	logger.InitLogger()

	// Log startup info
	logger.Log.Info("📦 Inisialisasi aplikasi...")

	err := godotenv.Load(".env")
	if err != nil {
		logger.Log.Fatal("❌ Gagal memuat .env file")
	}

	database.ConnectDB()
	logger.Log.Info("✅ Koneksi database berhasil")

	err = database.DB.AutoMigrate(&models.Todo{})
	if err != nil {
		logger.Log.Fatalf("❌ AutoMigrate error: %v", err)
	}

	app := fiber.New()
	app.Use(loggerFiberMiddleware()) 

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Put("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	logger.Log.Infof("🚀 Listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}


// GET /api/todos
func getTodos(c *fiber.Ctx) error {
	var todos []models.Todo
	if err := database.DB.Find(&todos).Error; err != nil {
		logger.Log.Errorf("❌ Gagal mengambil data: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch todos"})
	}
	logger.Log.Info("✅ Todos berhasil diambil")
	return c.Status(200).JSON(todos)
}


// POST /api/todos
func createTodo(c *fiber.Ctx) error {
	todo := new(models.Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Body cannot be empty"})
	}

	if err := database.DB.Create(&todo).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create todo"})
	}

	return c.Status(201).JSON(todo)
}

// PUT /api/todos/:id
func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var todo models.Todo
	if err := database.DB.First(&todo, todoID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	}

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := database.DB.Save(&todo).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update todo"})
	}

	return c.Status(200).JSON(todo)
}

// DELETE /api/todos/:id
func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := database.DB.Delete(&models.Todo{}, todoID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete todo"})
	}

	return c.SendStatus(204)
}
