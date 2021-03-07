package utils

import "github.com/gofiber/fiber/v2"

func SendError(c *fiber.Ctx, statusCode int, err error) error {
	if err := c.JSON(fiber.Error{
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
