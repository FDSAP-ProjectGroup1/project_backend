package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"crypto/subtle"
	"encoding/base64"
	"strings"

	"github.com/RustyPunzalan/project/models"
	"github.com/RustyPunzalan/project/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type User struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
}
type Sched struct {
	Title  string `json:"title"`
	Reason string `json:"reason"`
}

type Repository struct {
	DB *gorm.DB
}

// Login handles user authentication using Basic Auth
func (r *Repository) Login(context *fiber.Ctx) error {
	auth := context.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Authorization header missing or invalid",
		})
		return nil
	}

	payload, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Invalid authorization payload",
		})
		return nil
	}

	credentials := strings.SplitN(string(payload), ":", 2)
	if len(credentials) != 2 {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Invalid username or password",
		})
		return nil
	}

	username := credentials[0]
	password := credentials[1]

	// Query the database to retrieve the user record
	userModel := &models.Users{}
	err = r.DB.Where("name = ?", username).First(userModel).Error
	if err != nil {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Invalid username or password",
		})
		return nil
	}

	// Compare the provided password with the stored password
	storedPassword := []byte(*userModel.Password)
	if subtle.ConstantTimeCompare([]byte(password), storedPassword) != 1 {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Invalid username or password",
		})
		return nil
	}

	// Authentication successful
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Login successful",
	})
	return nil
}

func (r *Repository) CreateUser(context *fiber.Ctx) error {
	user := User{}

	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid request payload",
		})
		return err
	}

	// Validate user input
	if user.Name == "" || user.UserName == "" || user.Password == "" || user.Address == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Missing required fields",
		})
		return nil
	}

	// Create the user record
	userModel := models.Users{
		Name:     &user.Name,
		UserName: &user.UserName,
		Password: &user.Password,
		Address:  &user.Address,
	}

	err = r.DB.Create(&userModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not create user",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "User has been added",
	})
	return nil
}
func (r *Repository) CreateSched(context *fiber.Ctx) error {
	sched := Sched{}

	err := context.BodyParser(&sched)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid request payload",
		})
		return err
	}

	// Validate user input
	if sched.Title == "" || sched.Reason == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Missing required fields",
		})
		return nil
	}

	// Create the sched record
	schedModel := models.Scheds{
		Title:  &sched.Title,
		Reason: &sched.Reason,
	}

	err = r.DB.Create(&schedModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not create schedule",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Schedule has been added",
	})
	return nil
}

func (r *Repository) DeleteUser(context *fiber.Ctx) error {
	userModel := models.Users{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(userModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete user",
		})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user deleted successfully",
	})
	return nil
}
func (r *Repository) DeleteSched(context *fiber.Ctx) error {
	schedModel := models.Scheds{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(schedModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete schedule",
		})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "schedule deleted successfully",
	})
	return nil
}

func (r *Repository) GetUsers(context *fiber.Ctx) error {
	userModels := &[]models.Users{}

	err := r.DB.Find(userModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get users"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "users fetched successfully",
		"data":    userModels,
	})
	return nil
}
func (r *Repository) GetScheds(context *fiber.Ctx) error {
	schedModels := &[]models.Scheds{}

	err := r.DB.Find(schedModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get schedules"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "scheduls fetched successfully",
		"data":    schedModels,
	})
	return nil
}

func (r *Repository) GetUserByID(context *fiber.Ctx) error {
	id := context.Params("id")
	userModel := &models.Users{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(userModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the user"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user id fetched successfully",
		"data":    userModel,
	})
	return nil
}
func (r *Repository) GetSchedByID(context *fiber.Ctx) error {
	id := context.Params("id")
	schedModel := &models.Scheds{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(schedModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the user"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user id fetched successfully",
		"data":    schedModel,
	})
	return nil
}

func (r *Repository) UpdateUser(context *fiber.Ctx) error {
	userModel := models.Users{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.First(&userModel, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not retrieve the user",
		})
		return err
	}
	updatedUser := User{}
	err = context.BodyParser(&updatedUser)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	// Update the user fields
	if updatedUser.Name != "" {
		userModel.Name = &updatedUser.Name
	}
	if updatedUser.Password != "" {
		userModel.Password = &updatedUser.Password
	}

	err = r.DB.Save(&userModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not update the user"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user updated successfully",
		"data":    userModel,
	})
	return nil
}
func (r *Repository) UpdateSched(context *fiber.Ctx) error {
	schedModel := models.Scheds{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.First(&schedModel, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not retrieve the schedule",
		})
		return err
	}
	updatedSched := Sched{}
	err = context.BodyParser(&updatedSched)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	// Update the sched fields
	if updatedSched.Title != "" {
		schedModel.Title = &updatedSched.Title
	}
	if updatedSched.Reason != "" {
		schedModel.Reason = &updatedSched.Reason
	}

	err = r.DB.Save(&schedModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not update the schedule"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "schedule updated successfully",
		"data":    schedModel,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_users", r.CreateUser)
	api.Delete("/delete_user/:id", r.DeleteUser)
	api.Get("/get_users/:id", r.GetUserByID)
	api.Get("/users", r.GetUsers)
	api.Put("/update_user/:id", r.UpdateUser)

	api.Post("/create_scheds", r.CreateSched)
	api.Delete("/delete_sched/:id", r.DeleteSched)
	api.Get("/get_sched/:id", r.GetSchedByID)
	api.Get("/scheds", r.GetScheds)
	api.Put("/update_sched/:id", r.UpdateSched)
	api.Post("/login", r.Login)
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
		log.Fatal("could not load the database")
	}
	err = models.MigrateUsers(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}