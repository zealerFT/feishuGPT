package dep

import (
	"log"

	"feishu/biz/lark"
	"feishu/config"
	"feishu/service/larksvc"
	"feishu/service/logsvc"
	"feishu/service/lrusvc"
	"feishu/service/snowflakesvc"

	"go.uber.org/dig"
)

var Container *dig.Container

// DI 创建容器，注入全局对象
func DI() *dig.Container {
	return NewContainer(
		WithConfig(),
		WitchNewLogger(),
		WithSnowflake(),
		WithLark(),
		WithLruCache(),
		WithLarkbiz(),
	)
}

func ContainerGet() *dig.Container {
	return Container
}

type Option func(*dig.Container) error

func NewContainer(opts ...Option) *dig.Container {
	container := dig.New()
	for _, opt := range opts {
		if err := opt(container); err != nil {
			log.Fatalf("dig init Container fail: %v", err)
		}
	}
	return container
}

func WithConfig() Option {
	return func(c *dig.Container) error {
		return c.Provide(config.Options)
	}
}

func WitchNewLogger() Option {
	return func(c *dig.Container) error {
		return c.Provide(logsvc.NewLogger)
	}
}

func WithSnowflake() Option {
	return func(c *dig.Container) error {
		return c.Provide(snowflakesvc.MestNewSnowflake)
	}
}

func WithLark() Option {
	return func(c *dig.Container) error {
		return c.Provide(larksvc.MustNewLark)
	}
}

func WithLruCache() Option {
	return func(c *dig.Container) error {
		return c.Provide(lrusvc.NewLruCache)
	}
}

/*业务逻辑*/

func WithLarkbiz() Option {
	return func(c *dig.Container) error {
		return c.Provide(lark.NewLark)
	}
}
