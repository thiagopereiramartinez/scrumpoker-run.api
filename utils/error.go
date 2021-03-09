package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/models"
)

func SendError(c *fiber.Ctx, statusCode int, err error) error {
	if err := c.JSON(models.Error{
		Code:    statusCode,
		Message: err.Error(),
	}); err != nil {
		return err
	}
	if err := c.SendStatus(statusCode); err != nil {
		return err
	}

	return nil
}
