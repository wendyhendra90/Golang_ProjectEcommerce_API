package controllers

import (
	"student_project/database"
	"student_project/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginForm struct {
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

type AuthController struct {
	// Declare variables
	Db    *gorm.DB
	store *session.Store
}

func InitAuthController(s *session.Store) *AuthController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.User{})

	return &AuthController{Db: db, store: s}
}

// post /login
func (controller *AuthController) LoginPosted(c *fiber.Ctx) error {
	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	var user models.User
	var myform LoginForm

	if err := c.BodyParser(&myform); err != nil {
		return c.SendStatus(400) // Bad Request, RegisterForm is not complete
	}

	// Find user
	errs := models.FindUserByUsername(controller.Db, &user, myform.Username)
	if errs != nil {
		return c.JSON(
			fiber.Map{
				"message": "Username not found",
			},
		) // Unsuccessful login (cannot find user)
	}

	// Compare password
	compare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(myform.Password))
	if compare == nil { // compare == nil artinya hasil compare di atas true
		sess.Set("username", user.Username)
		sess.Set("userId", user.ID)
		sess.Save()

		return c.JSON(
			fiber.Map{
				"message": "Login successful",
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"message": "Login failed",
		},
	)
}

// POST /register
func (controller *AuthController) AddRegisteredUser(c *fiber.Ctx) error {
	var user models.User
	var cart models.Cart

	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(400) // Bad Request, RegisterForm is not complete
	}

	// Hash password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	sHash := string(bytes)

	// Simpan hashing, bukan plain passwordnya
	user.Password = sHash

	// save user
	err := models.CreateUser(controller.Db, &user)
	if err != nil {
		return c.SendStatus(500) // Server error, gagal menyimpan user
	}

	// Find user
	errs := models.FindUserByUsername(controller.Db, &user, user.Username)
	if errs != nil {
		return c.SendStatus(500) // Server error, gagal menyimpan user
	}

	// also create cart
	errCart := models.CreateCart(controller.Db, &cart, user.ID)
	if errCart != nil {
		return c.SendStatus(500) // Server error, gagal menyimpan user
	}

	// if succeed
	return c.JSON(
		fiber.Map{
			"message": "Register successful",
		},
	)
}

// /logout
func (controller *AuthController) Logout(c *fiber.Ctx) error {

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Destroy()
	return c.JSON(
		fiber.Map{
			"message": "Logout successful",
		},
	)
}
