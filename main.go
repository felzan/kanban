package main

import (
	"fmt"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	"kanban/card"
	"kanban/database"
	"log"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	initDatabase()

	setupRoutes(app)
	err := app.Listen(8080)
	if err != nil {
		log.Fatal(err)
	}

	defer database.DBConn.Close()
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "kanban.db")
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection opened to database")
	database.DBConn.AutoMigrate(&card.Card{})
}

func setupRoutes(app *fiber.App) {
	app.Get("/", getRoot)
	v1 := app.Group("/api/v1/")
	v1.Get("/cards", card.GetCards)
	v1.Get("/cards/:id", card.GetCard)
	v1.Put("/cards/:id", card.UpdateCard)
	v1.Patch("/cards/:id", card.MoveCard)
	v1.Post("/cards", card.AddCard)
	v1.Delete("/cards/:id", card.DeleteCard)
}

func getRoot(c *fiber.Ctx) {
	c.Send("Root")
}
