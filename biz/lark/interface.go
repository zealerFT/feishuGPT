package lark

import (
	"context"

	"feishu/config"
	"feishu/model"
	"feishu/service/snowflakesvc"

	lru "github.com/hashicorp/golang-lru"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/rs/zerolog"
	"go.uber.org/dig"
)

type Hub struct {
	dig.In

	Lark Lark
}

//go:generate mockgen -destination mock/lark.go -package larkmock feishu/biz/Lark Lark
type Lark interface {
	GroupInfo(ctx context.Context, chatId string) error
	GroupList(ctx context.Context) ([]*larkim.ListChat, error)
	GroupMessageCreate(ctx context.Context, chatId string, body string) error
	SendInteractiveMonitorMsg(ctx context.Context, receiveId, userIdType string) error
	SendInteractiveMsg(ctx context.Context, receiveId, receiveIdType string) error
	SendArgocdMsg(ctx context.Context, argocdBody model.ArgocdBody, receiveId, receiveIdType string) error
	SendJarvisMsg(ctx context.Context, jarvis model.Jarvis, receiveId, receiveIdType string) error
	SendImagesSyncMsg(ctx context.Context, message, receiveId, receiveIdType string) error
	OpenaiReply(ctx context.Context, content, receiveId, messageId, receiveIdType string) error
	ReactionReply(ctx context.Context, messageId string) error
}

type Larkbiz struct {
	dig.In

	Cfg *config.AppConfig // 配置

	Snowflake *snowflakesvc.Snowflake
	Lark      *lark.Client
	Lru       *lru.Cache
	Log       zerolog.Logger
}

// Object 当invoke的时候会实例化Larkbiz对象，并使用go的interface实现方式实现接口，所以这个变量在初始化的时候已是一个全局变量
var Object Lark

func NewLark(entry Larkbiz) Lark {
	Object = &entry
	return Object
}
