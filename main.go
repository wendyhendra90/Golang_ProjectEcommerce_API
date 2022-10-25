package main

import (
	"student_project/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
)

func main() {
	// session
	store := session.New()

	// load template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// static
	app.Static("/", "./public", fiber.Static{
		Index: "",
	})

	// Middleware to check login
	CheckLogin := func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		val := sess.Get("username")
		if val != nil {
			return c.Next()
		}
		return c.Redirect("/login")
	}

	// controllers
	prodController := controllers.InitProductController(store)
	authController := controllers.InitAuthController(store)
	cartController := controllers.InitCartController(store)
	transaksiController := controllers.InitTransaksiController(store)

	prod := app.Group("/products")
	prod.Get("/", prodController.GetAllProduct)
	prod.Post("/create", CheckLogin, prodController.AddPostedProduct)
	prod.Get("/detail/:id", prodController.DetailProduct)
	prod.Post("/editproduct/:id", CheckLogin, prodController.AddUpdatedProduct)
	prod.Get("/delete/:id", CheckLogin, prodController.DeleteProduct)

	prod.Get("/addtocart/:productid", CheckLogin, cartController.AddToCart)

	cart := app.Group("/cart")
	cart.Get("/:cartid", CheckLogin, cartController.GetCart)

	transaksi := app.Group("/transaction")
	transaksi.Get("/:userid", CheckLogin, transaksiController.AddTransaction)
	transaksi.Get("/ehe/:userid", transaksiController.GetTransaksi)

	app.Post("/login", authController.LoginPosted)
	app.Get("/logout", authController.Logout)
	app.Post("/register", authController.AddRegisteredUser)

	app.Listen(":3000")
}
