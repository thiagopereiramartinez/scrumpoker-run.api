package utils

import (
	"errors"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/models"
)

type SenderContext interface {
	JSON(interface{}) error
	SendStatus(statusCode int) error
}

func SendError(c SenderContext, statusCode int, err error) error {
	if err == nil {
		return errors.New("error property cannot be nil")
	}
	if c == nil {
		return errors.New("sender property cannot be nil")
	}

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
