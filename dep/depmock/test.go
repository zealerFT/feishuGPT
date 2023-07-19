package depmock

import (
	"testing"

	"feishu/dep"

	larkmock "feishu/biz/lark/mock"

	"github.com/golang/mock/gomock"
	"go.uber.org/dig"
)

func GetMock(t *testing.T, f interface{}) {
	container := dep.NewContainer(WithGomockControlle(t), WithMockLark())
	if err := container.Invoke(f); err != nil {
		panic(err)
	}

	return
}

func Mocks(t *testing.T) *Mock {
	currMock := &Mock{}
	GetMock(t, func(m Mock) { currMock = &m })
	return currMock
}

func WithGomockControlle(t *testing.T) dep.Option {
	ctrl := gomock.NewController(t)
	return func(c *dig.Container) error {
		// c.Provide need func to return ctrl, not just ctrl
		return c.Provide(func() *gomock.Controller { return ctrl })
	}
}

func WithMockLark() dep.Option {
	return func(c *dig.Container) error {
		return c.Provide(larkmock.NewMockLark)
	}
}

type Mock struct {
	dig.In
	Ctrl    *gomock.Controller
	LarkHub *larkmock.MockLark
}
