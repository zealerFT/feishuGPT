package lark

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func (l *Larkbiz) GroupInfo(ctx context.Context, chatId string) error {
	request := larkim.NewGetChatReqBuilder().UserIdType("chat_id").ChatId(chatId).UserIdType("user_id")
	resp, err := l.Lark.Im.Chat.Get(ctx, request.Build())
	// 处理错误
	if err != nil {
		log.Err(err).Msgf("GroupInfo is fail: %v", resp.Err.Details)
		return err
	}
	// 服务端错误处理
	if !resp.Success() {
		log.Error().Msgf("GroupInfo is fail: %v - %v -%v", resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

func (l *Larkbiz) GroupList(ctx context.Context) ([]*larkim.ListChat, error) {
	request := larkim.NewListChatReqBuilder().PageSize(1).Limit(100)
	list, err := l.Lark.Im.Chat.List(ctx, request.Build())
	if err != nil {
		log.Err(err).Msgf("GroupList is fail: %v", list.Err.Details)
		return nil, err
	}
	if !list.Success() {
		log.Error().Msgf("GroupList is fail: %v - %v -%v", list.Code, list.Msg, list.RequestId())
		return nil, list
	}
	fmt.Println("data:", list.Data)

	return list.Data.Items, nil
}
