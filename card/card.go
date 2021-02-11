package card

import (
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	"kanban/database"
	"log"
	"strconv"
)

type Swimlane string

const (
	Todo  Swimlane = "TODO"
	Doing Swimlane = "DOING"
)

type Board struct {
	Todo  []Card `json:"todo"`
	Doing []Card `json:"doing"`
	//Done  uint   `json:"done"`
}

type Card struct {
	gorm.Model
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Swimlane Swimlane `json:"swimlane"`
}

func GetCards(c *fiber.Ctx) {
	queryShowArchived := c.Query("showArchived")
	var showArchived bool
	if queryShowArchived != "" {
		parseBool, err := strconv.ParseBool(c.Query("showArchived"))
		if err != nil {
			c.Status(503).Send(err)
			return
		}
		showArchived = parseBool
	}

	db := database.DBConn
	var cards []Card
	if showArchived {
		db.Find(&cards)
	} else {
		db.Not("swimlane = ?", "DONE").Find(&cards)
	}
	board := new(Board)
	cardsListToBoard(board, cards)
	c.JSON(board)
}

func cardsListToBoard(out *Board, cards []Card) {
	for _, c := range cards {
		switch c.Swimlane {
		case Todo:
			out.Todo = append(out.Todo, c)
		case Doing:
			out.Doing = append(out.Doing, c)
		}
	}
}

func GetCard(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn
	var card Card
	db.Find(&card, id)
	c.JSON(card)
}

func UpdateCard(c *fiber.Ctx) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Fatal("id not valid")
		return
	}
	db := database.DBConn

	card := new(Card)
	if err := c.BodyParser(card); err != nil {
		c.Status(503).Send(err)
		return
	}
	card.ID = uint(id)
	db.Save(&card)
	c.JSON(card)
}

type MoveCardRequest struct {
	Swimlane Swimlane `json:"swimlane"`
}

func MoveCard(c *fiber.Ctx) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Fatal("id not valid")
		return
	}
	moveCardRequest := new(MoveCardRequest)
	if err := c.BodyParser(moveCardRequest); err != nil {
		c.Status(503).Send(err)
		return
	}
	db := database.DBConn
	var card Card
	find := db.Find(&card, id)
	if find.Error != nil {
		c.Status(503).Send(find.Error)
		return
	}
	card.Swimlane = moveCardRequest.Swimlane
	db.Save(&card)
	c.JSON(card)
}

func AddCard(c *fiber.Ctx) {
	db := database.DBConn
	card := new(Card)
	if err := c.BodyParser(card); err != nil {
		c.Status(503).Send(err)
		return
	}
	card.Swimlane = Todo
	db.Create(&card)
	c.JSON(card)
}

func DeleteCard(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn
	var card Card
	db.First(&card, id)
	if card.Title == "" {
		c.Status(500).Send("No card found with ID: " + id)
		return
	}
	db.Delete(&card)
	c.Status(204)
}
