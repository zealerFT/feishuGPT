package route

import (
	"context"
	"fmt"

	"feishu/http/api"
	"feishu/http/middleware"

	"feishu/config"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func AppAPI(opts ...func(engine *gin.Engine)) func(s *gin.Engine) {
	return func(s *gin.Engine) {
		for _, opt := range opts {
			opt(s)
		}

		// Not Found
		s.NoRoute(api.Handle404)
		// Health Check
		s.GET("/check", api.Health)

		// feishu
		route := authRouteGroup(s, "/api/v1/feishu")
		// 已包含初始化event校验，无需专门写
		route.POST("/webhook/event", sdkginext.NewEventHandlerFunc(LarkEvent()))
		// biz
		route.GET("/group/info", api.GroupInfo)
		route.GET("/group/list", api.GroupList)
		route.POST("/message/group/create", api.GroupMessageCreate)
		route.POST("/infra", api.SendInteractiveMonitorMsg)
		route.POST("/infra2", api.SendInteractiveMsg)
		route.POST("/argocd", api.SendArgocdMsg)
		route.POST("/openai", api.OpenaiReplyHttp)
		route.POST("/reaction", api.ReactionReplyHttp)

	}
}

// 需要jwt鉴权的group
func authRouteGroup(s *gin.Engine, relativePath string) *gin.RouterGroup {
	group := s.Group(relativePath)
	group.Use(middleware.Auth())
	return group
}

func LarkEvent() *dispatcher.EventDispatcher {
	handler := dispatcher.NewEventDispatcher(config.Options().LarkSetting.VerificationToken, config.Options().LarkSetting.EncryptKey)
	return handler.OnP2MessageReceiveV1(api.OpenaiReply).OnP1MessageReceiveV1(func(ctx context.Context, event *larkim.P1MessageReceiveV1) error {
		// 当用户发送消息给机器人或在群聊中@机器人时触发此事件。(当群里只有机器人时会触发)
		return nil
	}).OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2UserCreatedV3(func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2MessageReactionCreatedV1(api.ReactionReply).OnCustomizedEvent("custom_event_type", func(ctx context.Context, event *larkevent.EventReq) error {
		// 原生消息体
		fmt.Println(string(event.Body))
		fmt.Println(larkcore.Prettify(event.Header))
		fmt.Println(larkcore.Prettify(event.RequestURI))
		fmt.Println(event.RequestId())

		// 处理消息
		cipherEventJsonStr, err := handler.ParseReq(ctx, event)
		if err != nil {
			//  错误处理
			return err
		}

		plainEventJsonStr, err := handler.DecryptEvent(ctx, cipherEventJsonStr)
		if err != nil {
			//  错误处理
			return err
		}

		// 处理解密后的 消息体
		fmt.Println(plainEventJsonStr)

		return nil
	})
}
