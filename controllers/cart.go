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

type CartController struct {
	// Declare variables
	Db    *gorm.DB
	store *session.Store
}

func InitCartController(s *session.Store) *CartController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.Cart{})

	return &CartController{Db: db, store: s}
}

// GET /addtocart/products/:productid
func (controller *CartController) AddToCart(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"
	//get id user session
	sess, errs := controller.store.Get(c)
	if errs != nil {
		panic(errs)
	}
	val := sess.Get("userId")
	//change val to int
	// intCartId, _ := strconv.Atoi(params["cartid"])
	//print user id
	str := fmt.Sprintf("%v", val)
	intCartId, _ := strconv.Atoi(str)
	intProductId, _ := strconv.Atoi(params["productid"])

	var cart models.Cart
	var product models.Product

	// Find the product first,
	err := models.ReadProductById(controller.Db, &product, intProductId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// Then find the cart
	errss := models.ReadCartById(controller.Db, &cart, intCartId)
	if errss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	cart.Quantity += 1
	// Finally, insert the product to cart
	errsss := models.InsertProductToCart(controller.Db, &cart, &product)
	if errsss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(
		fiber.Map{
			"Title":    "Detail Product",
			"Products": cart.Products,
		},
	)
}

// GET /shoppingcart/:cartid
func (controller *CartController) GetCart(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intCartId, _ := strconv.Atoi(params["cartid"])

	var cart models.Cart
	err := models.ReadAllProductsInCart(controller.Db, &cart, intCartId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("userId")

	return c.JSON(
		fiber.Map{
			"Title":    "Shopping Cart",
			"Products": cart.Products,
			"UserId":   val,
		},
	)
}
