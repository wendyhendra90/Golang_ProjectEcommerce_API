package controllers

import (
	"strconv"
	"student_project/database"
	"student_project/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type TransaksiController struct {
	// Declare variables
	Db    *gorm.DB
	store *session.Store
}

func InitTransaksiController(s *session.Store) *TransaksiController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.Transaksi{})

	return &TransaksiController{Db: db, store: s}
}

// GET /checkout/:userid
func (controller *TransaksiController) AddTransaction(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intUserId, _ := strconv.Atoi(params["userid"])

	var transaksi models.Transaksi
	var cart models.Cart

	// Find the product first,
	err := models.ReadAllProductsInCart(controller.Db, &cart, intUserId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	errs := models.CreateTransaksi(controller.Db, &transaksi, uint(intUserId), cart.Products)
	if errs != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// Delete products in cart
	errss := models.UpdateCart(controller.Db, cart.Products, &cart, uint(intUserId))
	if errss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(
		fiber.Map{
			"Title":    "Checkout",
			"Products": cart.Products,
		},
	)
}

// GET /historytransaksi/:userid
func (controller *TransaksiController) GetTransaksi(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intUserId, _ := strconv.Atoi(params["userid"])

	var transaksis []models.Transaksi
	err := models.ReadTransaksiById(controller.Db, &transaksis, intUserId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(
		fiber.Map{
			"Title":     "History Transaksi",
			"Transaksi": transaksis,
		},
	)

}

// GET /history/detail/:transaksiid
func (controller *TransaksiController) DetailTransaksi(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intTransaksiId, _ := strconv.Atoi(params["transaksiid"])

	var transaksi models.Transaksi
	err := models.ReadAllProductsInTransaksi(controller.Db, &transaksi, intTransaksiId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(
		fiber.Map{
			"Title":    "Detail Transaksi",
			"Products": transaksi.Products,
		},
	)
}
