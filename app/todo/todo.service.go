package todo

import (
	"errors"
	"numtostr/gotodo/app/dal"
	"numtostr/gotodo/utils"

	"github.com/gofiber/fiber"
	"gorm.io/gorm"
)

// CreateTodo is responsible for create todo
func CreateTodo(c *fiber.Ctx) {
	b := new(CreateDTO)

	if err := utils.ParseBodyAndValidate(c, b); err != nil {
		c.Next(err)
		return
	}

	d := &dal.Todo{
		Task: b.Task,
		User: utils.GetUser(c),
	}

	if err := dal.CreateTodo(d).Error; err != nil {
		c.Next(fiber.NewError(fiber.StatusConflict, err.Error()))
		return
	}

	c.JSON(&CreateRes{
		Todo: &Response{
			ID:        d.ID,
			Task:      d.Task,
			Completed: d.Completed,
		},
	})
}

// GetTodos returns the todos list
func GetTodos(c *fiber.Ctx) {
	d := &[]Response{}

	err := dal.FindTodosByUser(d, utils.GetUser(c)).Error
	if err != nil {
		c.Next(fiber.NewError(fiber.StatusConflict, err.Error()))
		return
	}

	c.JSON(&ListRes{
		Todos: d,
	})
}

// GetTodo return a single todo
func GetTodo(c *fiber.Ctx) {
	todoID := c.Params("todoID")

	if todoID == "" {
		c.Next(fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid todoID"))
		return
	}

	d := &Response{}

	err := dal.FindTodoByUser(d, todoID, utils.GetUser(c)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(&CreateRes{})
		return
	}

	c.JSON(&CreateRes{
		Todo: d,
	})
}

// DeleteTodo deletes a single todo
func DeleteTodo(c *fiber.Ctx) {
	todoID := c.Params("todoID")

	if todoID == "" {
		c.Next(fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid todoID"))
		return
	}

	res := dal.DeleteTodo(todoID, utils.GetUser(c))
	if res.RowsAffected == 0 {
		c.Next(fiber.NewError(fiber.StatusConflict, "Unable to delete todo"))
		return
	}

	err := res.Error
	if err != nil {
		c.Next(fiber.NewError(fiber.StatusConflict, err.Error()))
		return
	}

	c.JSON(&CreateRes{})
}
