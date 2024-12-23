package main

import (
	"fmt"
	"goFiberPostgres/models"
	"goFiberPostgres/storage"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2" // Ensure Fiber v2 is imported
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}
type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := Book{}

	err := c.BodyParser(&book)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request failed"})
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not create book"})
	}

	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Book has been added"})
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	err := r.DB.Find(bookModels).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Books couldn't be fetched"})
	}
	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Books fetched successfully", "data": bookModels})
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	bookModel := models.Books{}
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "ID cannot be empty"})
	}
	err := r.DB.Delete(&bookModel, id).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not delete book"})
	}
	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Deletion successful"})
}

func (r *Repository) GetBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "ID cannot be empty"})
	}
	fmt.Println("The ID is", id) // Optional for debugging

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get the book"})
	}
	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Book fetched successfully", "data": bookModel})
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load the database")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New() // Initialize Fiber app with v2 syntax
	r.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
