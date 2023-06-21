package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"crypto/subtle"

	"github.com/RustyPunzalan/project/models"
	"github.com/RustyPunzalan/project/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Define your chatbot struct
type Chatbot struct {
	DB *gorm.DB
}

type User struct {
	Fullname string `json:"fullname"`
	Name     string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
}
type Sched struct {
	Date   string `json:"date"`
	Time   string `json:"time"`
	Title  string `json:"title"`
	Reason string `json:"reason"`
}

type Message struct {
	Text string `json:"message"`
}

type Response struct {
	Text string `json:"response"`
}
type Repository struct {
	DB *gorm.DB
}

// Login handles user authentication using Basic Auth
// Login handles user authentication using the request body
func (r *Repository) Login(context *fiber.Ctx) error {
	user := User{}

	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid request payload",
		})
		return err
	}

	// Validate user input
	if user.Name == "" || user.Password == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Missing username or password",
		})
		return nil
	}

	// Query the database to retrieve the user record
	userModel := &models.Users{}
	err = r.DB.Where("name = ?", user.Name).First(userModel).Error
	if err != nil {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Invalid username",
		})
		return nil
	}

	// Compare the provided password with the stored password
	storedPassword := []byte(*userModel.Password)
	if subtle.ConstantTimeCompare([]byte(user.Password), storedPassword) != 1 {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Invalid or password",
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
	if user.Name == "" || user.Fullname == "" || user.Password == "" || user.Address == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Missing required fields",
		})
		return nil
	}

	// Create the user record
	userModel := models.Users{
		Name:     &user.Name,
		Fullname: &user.Fullname,
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
		Date:   &sched.Date,
		Time:   &sched.Time,
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
	if updatedUser.Fullname != "" {
		userModel.Fullname = &updatedUser.Fullname
	}
	if updatedUser.Name != "" {
		userModel.Name = &updatedUser.Name
	}
	if updatedUser.Password != "" {
		userModel.Password = &updatedUser.Password
	}
	if updatedUser.Address != "" {
		userModel.Address = &updatedUser.Address
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
	if updatedSched.Date != "" {
		schedModel.Date = &updatedSched.Date
	}
	if updatedSched.Time != "" {
		schedModel.Time = &updatedSched.Time
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
func (r *Repository) SearchUsers(context *fiber.Ctx) error {
	query := context.Params("query")

	userModels := &[]models.Users{}
	err := r.DB.Where("fullname LIKE ? OR address LIKE ?", "%"+query+"%", "%"+query+"%").Find(userModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not search users",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "users fetched successfully",
		"data":    userModels,
	})
	return nil
}

// HandleChat handles chat interactions with the bot
func (r *Repository) HandleChat(context *fiber.Ctx) error {
	// Parse the chat message from the request body
	message := Message{}
	err := context.BodyParser(&message)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid request payload",
		})
		return err
	}

	// Process the chat message and generate a response
	response := r.GenerateResponse(message.Text)

	// Return the response as JSON
	context.Status(http.StatusOK).JSON(&Response{
		Text: response,
	})
	return nil
}

// GenerateResponse generates a response based on the input text
func (r *Repository) GenerateResponse(input string) string {
	// Add your chatbot logic here to generate a response based on the input
	// For example, you can use a switch statement to handle different input cases

	switch input {
	case "hello", "Hello":
		return "Hello, how can I assist you?"

	case "hi", "Hi":
		return "Welcome! Thank you for choosing our chatbot assistance. \n We are here to help you with any questions, concerns, or information you may need. Feel free to ask us anything, and we'll provide you with the best possible support. Our goal is to make your experience smooth, efficient, and enjoyable. So, go ahead and let us know how we can assist you today!"

		//Tagalog

	case "Kamusta", "kamusta", "Kamusta?", "kamusta?", "Musta", "musta", "Musta?", "musta?":
		return "Kamusta! Ako ang iyong chatbot na handang tumulong sa iyo. Ano ang mga katanungan o tulong na kailangan mo? Sabihin mo lang sa akin at tutulungan kita sa abot ng aking makakaya."

	case "Patulong", "patulong", "Tulong", "tulong", "Patulong po", "patulong po", "tulong po", "Tulong po":
		return "Para makagawa ng appointment, mangyaring ibigay sa app ang mga sumusunod na detalye:\n\n1.	Petsa ng appointment na nais mo.\n2.	Oras ng appointment na nais mo.\n3.	Anumang espesyal na kahilingan o detalye na kailangan kong malaman.\n\n Ipadala lamang ang mga detalye na ito sa aming app at aasikasuhin nito ang iyong appointment. Salamat!"

	case "Paano mag gawa ng appointment?", "paano mag gawa ng appointment?", "pano mag gawa ng appointment?", "Pano mag gawa ng appointment?", "Pano mag gawa ng appointment", "pano mag gawa ng appointment", "Paano mag gawa ng appointment", "paano mag gawa ng appointment":
		return "Tiyak! Ano ang mga katanungan o tulong na kailangan mo? Ipadala lang sa akin ang iyong mga katanungan o mga detalye ng anumang problema, at gagawin ko ang aking makakaya upang tulungan ka."

		//English

	case "help", "Help":
		return "Sure, I can help you with that!"

	case "How can I make an appointment?", "How can i make an appointment?", "how can I make an appointment?", "how can i make an appointment?":
		return "Greetings! Our chatbot is here to assist you. If you'd like to make an appointment, simply provide us with the date, time, and any specific requirements. We'll take care of the rest and confirm the appointment. Feel free to ask any questions or provide further details. We're here to make your experience smooth and convenient."

	case "Steps on how to make an appointment", "steps on how to make an appointment", "Steps on how to make an Appointment", "steps on how to make an Appointment":
		return " Absolutely! We're here to assist you. Please follow these steps:\n\n1.	Choose a preferred date and time.\n2.	Specify any specific requirements.\n3.	Submit your appointment for approval.\n\nWe're ready to help you through the entire process. Feel free to ask any questions along the way!"

	case "Goodbye", "goodbye":
		return "Thank you for chatting with me! I hope I was able to assist you effectively and provide the information you were seeking. Remember, knowledge is a journey, and I'm here to accompany you along the way. If you have any more questions in the future, don't hesitate to reach out. Wishing you continued success, fulfillment, and an abundance of learning experiences. Goodbye for now, and take care!"

	default:
		return "I'm sorry, I didn't understand that. Can you please rephrase?"
	}
}
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/user", r.CreateUser)
	api.Delete("/delete/:id", r.DeleteUser)
	api.Get("/search/:id", r.GetUserByID)
	api.Get("/all_users", r.GetUsers)
	api.Put("/update/:id", r.UpdateUser)
	api.Get("/search_users/:query", r.SearchUsers)

	api.Post("/sched", r.CreateSched)
	api.Delete("/drop_sched/:id", r.DeleteSched)
	api.Get("/search_sched/:id", r.GetSchedByID)
	api.Get("/all_scheds", r.GetScheds)
	api.Put("/update_sched/:id", r.UpdateSched)
	api.Post("/login", r.Login)
	api.Post("/chat", r.HandleChat)
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
	err = models.MigrateScheds(db)
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
