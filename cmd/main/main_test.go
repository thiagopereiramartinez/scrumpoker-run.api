package main

import (
	Assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/test"
	"testing"
)

func TestSetupRouter(t *testing.T) {

	_ = Assert.New(t)
	router := new(test.MockRouter)

	router.On("Use", mock.Anything).Return(router)
	router.On("Get", mock.Anything, mock.Anything).Return(router)
	router.On("Put", mock.Anything, mock.Anything).Return(router)
	router.On("Post", mock.Anything, mock.Anything).Return(router)
	router.On("Patch", mock.Anything, mock.Anything).Return(router)
	router.On("Group", mock.Anything, mock.Anything).Return(router)
	router.On("Delete", mock.Anything, mock.Anything).Return(router)

	SetupRouter(router)

	router.AssertExpectations(t)
	Assert.True(t, len(router.Mock.Calls) > 0)
}
