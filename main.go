package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tabish-bp/go_db_app/models"
	"github.com/tabish-bp/go_db_app/storage"
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

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := Book{}

	err := c.BodyParser(&book)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"status": "error", "message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"status": "error", "message": "book not created"})

		return err
	}
	c.Status(http.StatusCreated).JSON(
		&fiber.Map{"status": "success", "message": "book created"})

	return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	bookModels := &[]models.Book{}

	err := r.DB.Find(&bookModels).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"status": "error", "message": "could not get books"})
		return err
	}
	c.Status(http.StatusOK).JSON(
		&fiber.Map{"status": "success", "message": "books retrieved", "data": bookModels})
	return nil
}

func (r *Repository) GetBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	bookModel := &models.Book{}

	if id == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "could not find book"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book retrieved", "data": bookModel})
	return nil
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	bookModel := &[]models.Book{}
	id := c.Params("id")

	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
	}

	err := r.DB.Delete(bookModel, id)
	if err.Error != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete book"})
		return err.Error
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book deleted"})
	return nil
}

// func (r *Repository) UpdateBook(c *fiber.Ctx) error {
// 	book := Book{}
// 	bookModel := &[]models.Book{}
// 	id := c.Params("id")

// 	if id == "" {
// 		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "id cannot be empty"})
// 	}

// 	if err := c.BodyParser(&book); err != nil {
// 		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
// 	}

// 	if err := r.DB.Model(bookModel).Where("id = ?", id).Updates(book).Error; err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not update book"})
// 	}

// 	return c.Status(http.StatusOK).JSON(&fiber.Map{"message": "book updated"})
// }

func (r *Repository) UpdateBook(c *fiber.Ctx) error {
	book := Book{}
	bookModel := &[]models.Book{}
	id := c.Params("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "id cannot be empty"})
	}

	if err := c.BodyParser(&book); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
	}

	if err := r.DB.First(bookModel, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "could not find book"})
	}

	if err := r.DB.Model(bookModel).Updates(book).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not update book"})
	} else {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "book updated."})
	}
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_book", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Put("/update_book/:id", r.UpdateBook)
	api.Get("/get_book/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func AutoMigrate(db *gorm.DB) error {
	err := models.MigrateBooks(db)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = models.MigrateAuthors(db)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal(err)
	}

	err = AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()

	r.SetupRoutes(app)
	app.Listen(":8080")

}
