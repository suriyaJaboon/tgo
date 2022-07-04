package service

import (
	"tgo/store"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type lc struct {
	rlc store.LC
}

type lcDto struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func NewLC(f *fiber.App, rlc store.LC) *lc {
	slc := lc{rlc: rlc}

	g := f.Group("/lc")
	g.Get("", slc.LC)
	g.Get("/:id", slc.LCByID)

	return &lc{}
}

func (l lc) LC(c *fiber.Ctx) error {
	lcx, err := l.rlc.LCX()
	if err != nil {
		err = fiber.NewError(fiber.StatusInternalServerError, err.Error())
		return c.Status(500).JSON(err)
	}

	return c.JSON(lcx)
}

func (l lc) LCByID(c *fiber.Ctx) error {
	idx, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		err = fiber.NewError(fiber.StatusNotAcceptable, err.Error())
		return c.Status(fiber.StatusNotAcceptable).JSON(err)
	}

	lcx, err := l.rlc.LCXByID(idx)
	if err != nil {
		err = fiber.NewError(fiber.StatusInternalServerError, err.Error())
		return c.Status(500).JSON(err)
	}

	return c.JSON(lcx)
}
