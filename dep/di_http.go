package dep

import (
	larkBiz "feishu/biz/lark"
	"feishu/config"
	"feishu/service/snowflakesvc"

	lru "github.com/hashicorp/golang-lru"
	"go.uber.org/dig"
)

func DIHttpDependency() (out *HttpDependency) {
	container := DI()
	if err := container.Invoke(func(dep HttpDependency) { out = &dep }); err != nil {
		panic(err)
	}

	return
}

type HttpDependency struct {
	dig.In

	Config    *config.AppConfig
	Snowflake *snowflakesvc.Snowflake
	Lru       *lru.Cache
	LarkHub   larkBiz.Hub
}
