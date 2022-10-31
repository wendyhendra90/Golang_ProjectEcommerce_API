package controllers

import (
	"fmt"
	"strconv"
	"student_project/database"
	"student_project/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type ProductController struct {
	// Declare variables
	Db    *gorm.DB
	store *session.Store
}

func InitProductController(s *session.Store) *ProductController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.Product{})

	return &ProductController{Db: db, store: s}
}

// Routing
// GET /products
func (controller *ProductController) GetAllProduct(c *fiber.Ctx) error {
	// Load all Products
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("userId")
	return c.JSON(fiber.Map{
		"Title":    "Daftar Produk",
		"Products": products,
		"UserId":   val,
	})
}



// GET /products/hapus/:id
func (controller *ProductController) DeleteProduct(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intId, errs := strconv.Atoi(params["id"])

	if errs != nil {
		fmt.Println(errs)
	}

	var product models.Product
	err := models.DeleteProductById(controller.Db, &product, intId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
