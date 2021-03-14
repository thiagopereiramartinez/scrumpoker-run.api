package test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

type MockRouter struct {
	mock.Mock
}

func (m *MockRouter) Use(a ...interface{}) fiber.Router {
	args := m.Called(a)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Get(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Head(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Post(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Put(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Delete(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Connect(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Options(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Trace(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Patch(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Add(method, path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Static(prefix, root string, config ...fiber.Static) fiber.Router {
	args := m.Called(prefix, root, config)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) All(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Group(prefix string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(prefix, handlers)
	return args.Get(0).(fiber.Router)
}

func (m *MockRouter) Mount(prefix string, f *fiber.App) fiber.Router {
	args := m.Called(prefix, f)
	return args.Get(0).(fiber.Router)
}