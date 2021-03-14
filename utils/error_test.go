package utils

import (
	"errors"
	Assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/models"
	"testing"
)

type MockCtx struct {
	mock.Mock
}

func (m *MockCtx) JSON(i interface{}) error {
	args := m.Called(i)
	return args.Error(0)
}

func (m *MockCtx) SendStatus(statusCode int) error {
	args := m.Called(statusCode)
	return args.Error(0)
}

func TestSendErrorValid(t *testing.T) {

	assert := Assert.New(t)

	m := new(MockCtx)
	m.On("JSON", models.Error{
		Code:    500,
		Message: "test error",
	}).Return(nil)
	m.On("SendStatus", 500).Return(nil)

	err := SendError(m, 500, errors.New("test error"))
	assert.NoError(err)

	m.AssertExpectations(t)

}

func TestSendErrorNil(t *testing.T) {

	assert := Assert.New(t)

	m := new(MockCtx)

	err := SendError(m, 500, nil)
	assert.Error(err)
	assert.Equal("error property cannot be nil", err.Error())

}

func TestSendSenderNil(t *testing.T) {

	assert := Assert.New(t)

	err := SendError(nil, 500, errors.New("new error"))
	assert.Error(err)
	assert.Equal("sender property cannot be nil", err.Error())

}

func TestSendErrorFailInJSON(t *testing.T) {

	assert := Assert.New(t)

	m := new(MockCtx)
	m.On("JSON", models.Error{
		Code:    500,
		Message: "test error",
	}).Return(errors.New("error to generate JSON"))

	err := SendError(m, 500, errors.New("test error"))
	assert.Error(err)
	assert.Equal("error to generate JSON", err.Error())

	m.AssertExpectations(t)

}

func TestSendErrorFailInSendStatus(t *testing.T) {

	assert := Assert.New(t)

	m := new(MockCtx)
	m.On("JSON", models.Error{
		Code:    500,
		Message: "test error",
	}).Return(nil)
	m.On("SendStatus", 500).Return(errors.New("error to send status code"))

	err := SendError(m, 500, errors.New("test error"))
	assert.Error(err)
	assert.Equal("error to send status code", err.Error())

	m.AssertExpectations(t)

}
