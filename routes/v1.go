package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hibakun/nova-store/handlers"
	"github.com/hibakun/nova-store/middleware"
)

func V1(app *fiber.App) {
	app.Static("/public", "./public")
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/login", handlers.Login)
	v1.Post("/register", handlers.CreateUser)
	v1.Get("/logout", middleware.Protected, handlers.Logout)
	v1.Get("/validate", handlers.Validate)
	v1.Get("/user", middleware.Protected, handlers.GetUser)
	v1.Get("/user/transaction", middleware.Protected, handlers.GetUserTransactions)

	game := v1.Group("/game")
	game.Post("/create", middleware.Protected, handlers.CreateGame)
	game.Get("/", handlers.GetAllGames)
	game.Get("/:id", handlers.GetGameById)

	genre := v1.Group("/genre", middleware.Protected)
	genre.Post("/create", handlers.CreateGenre)
	genre.Get("/", handlers.GetAllGenres)
	genre.Get("/:id", handlers.GetGenreById)

	item := v1.Group("/item", middleware.Protected)
	item.Post("/create", handlers.CreateItem)
	item.Get("/:id", handlers.GetItemById)
	item.Get("/game/:id", handlers.GetItemsByGameId)

	payment := v1.Group("/payment")
	payment.Post("/create", middleware.Protected, handlers.CreatePayment)
	payment.Get("/", handlers.GetAllPayments)
	payment.Get("/:id", handlers.GetPaymentById)

	transaction := v1.Group("/transaction")
	transaction.Post("/create", middleware.Protected, handlers.CreateTransaction)
	transaction.Get("/all", handlers.GetAllTransactions)
	transaction.Get("/user", middleware.Protected, handlers.GetTransactionByUserId)
	transaction.Get("/:id", handlers.GetTransactionByUuid)
}
