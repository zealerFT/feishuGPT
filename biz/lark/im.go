package lark

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"feishu/model"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/rs/zerolog/log"
	goopenai "github.com/sashabaranov/go-openai"
)

func (l *Larkbiz) GroupMessageCreate(ctx context.Context, chatId string, body string) error {
	request := larkim.NewCreateMessageReqBuilder().ReceiveIdType("chat_id").Body(larkim.NewCreateMessageReqBodyBuilder().
		MsgType(larkim.MsgTypeText).
		ReceiveId(chatId).
		Content("{\"text\":\"firstline \\n second line  \"}").
		Build())
	resp, err := l.Lark.Im.Message.Create(ctx, request.Build())
	// 处理错误
	if err != nil {
		log.Err(err).Msgf("GroupMessageCreate is fail: %v", resp.Err.Details)
		return err
	}
	// 服务端错误处理
	if !resp.Success() {
		log.Error().Msgf("GroupMessageCreate is fail: %v - %v -%v", resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

// SendInteractiveMonitorMsg 运维报警通知
// https://open.feishu.cn/tool/cardbuilder?from=cotentmodule
func (l *Larkbiz) SendInteractiveMonitorMsg(ctx context.Context, receiveId, receiveIdType string) error {
	if receiveIdType == "" {
		receiveIdType = larkcontact.UserIdTypeOpenId
	}

	// config
	config := larkcard.NewMessageCardConfig().
		EnableForward(true).
		UpdateMulti(true).
		Build()

	// header
	header := larkcard.NewMessageCardHeader().
		Template(larkcard.TemplateRed).
		Title(larkcard.NewMessageCardPlainText().
			Content("1 级报警 - 数据平台").
			Build()).
		Build()

	// Elements
	divElement1 := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("**🕐 时间：**2021-02-23 20:17:51").
					Build()).
				IsShort(true).
				Build(),
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("**🔢 事件 ID：：**336720").
					Build()).
				IsShort(true).
				Build(),
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("").
					Build()).
				IsShort(false).
				Build(),

			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("**📋 项目：**\nQA 7").
					Build()).
				IsShort(true).
				Build(),
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("**👤 一级值班：**\n<at id=ou_c245b0a7dff2725cfa2fb104f8b48b9d>加多</at>").
					Build()).
				IsShort(true).
				Build(),

			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("").
					Build()).
				IsShort(false).
				Build(),
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("**👤 二级值班：**\n<at id=ou_c245b0a7dff2725cfa2fb104f8b48b9d>加多</at>").
					Build()).
				IsShort(true).
				Build()}).
		Build()

	divElement3 := larkcard.NewMessageCardNote().
		Elements([]larkcard.MessageCardNoteElement{larkcard.NewMessageCardPlainText().
			Content("🔴 支付失败数  🔵 支付成功数").
			Build()}).
		Build()

	divElement4 := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{larkcard.NewMessageCardEmbedButton().
			Type(larkcard.MessageCardButtonTypePrimary).
			Value(map[string]interface{}{"key1": "value1"}).
			Text(larkcard.NewMessageCardPlainText().
				Content("跟进处理").
				Build()),
			larkcard.NewMessageCardEmbedSelectMenuStatic().
				MessageCardEmbedSelectMenuStatic(larkcard.NewMessageCardEmbedSelectMenuBase().
					Options([]*larkcard.MessageCardEmbedSelectOption{larkcard.NewMessageCardEmbedSelectOption().
						Value("1").
						Text(larkcard.NewMessageCardPlainText().
							Content("屏蔽10分钟").
							Build()),
						larkcard.NewMessageCardEmbedSelectOption().
							Value("2").
							Text(larkcard.NewMessageCardPlainText().
								Content("屏蔽30分钟").
								Build()),
						larkcard.NewMessageCardEmbedSelectOption().
							Value("3").
							Text(larkcard.NewMessageCardPlainText().
								Content("屏蔽1小时").
								Build()),
						larkcard.NewMessageCardEmbedSelectOption().
							Value("4").
							Text(larkcard.NewMessageCardPlainText().
								Content("屏蔽24小时").
								Build()),
					}).
					Placeholder(larkcard.NewMessageCardPlainText().
						Content("暂时屏蔽报警").
						Build()).
					Value(map[string]interface{}{"key": "value"}).
					Build()).
				Build()}).
		Build()

	divElement5 := larkcard.NewMessageCardHr().Build()

	divElement6 := larkcard.NewMessageCardDiv().
		Text(larkcard.NewMessageCardLarkMd().
			Content("🙋🏼 [我要反馈误报](https://open.feishu.cn/) | 📝 [录入报警处理过程](https://open.feishu.cn/)").
			Build()).
		Build()

	// CardUrl
	cardLink := larkcard.NewMessageCardURL().
		PcUrl("http://www.baidu.com").
		IoSUrl("http://www.google.com").
		Url("http://open.feishu.com").
		AndroidUrl("http://www.jianshu.com").
		Build()

	low := "low"
	priority := larkcard.NewMessageCardMarkdown().
		Content(fmt.Sprintf(`**Priority**: (~~*%s*~~)  **%s**`, low, "high")).
		Build()
	fmt.Println(priority)
	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{divElement1, divElement3, divElement4, divElement5, divElement6, priority}).
		CardLink(cardLink).
		String()
	if err != nil {
		fmt.Println(err)
		return err
	}

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(receiveId).
			Content(cardContent).
			Build()).
		Build()
	resp, err := l.Lark.Im.Message.Create(ctx, req)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
	return nil
}

func (l *Larkbiz) SendInteractiveMsg(ctx context.Context, receiveId, receiveIdType string) error {
	if receiveIdType == "" {
		receiveIdType = larkcontact.UserIdTypeOpenId
	}
	// config
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(true).
		UpdateMulti(false).
		Build()

	// CardUrl
	cardLink := larkcard.NewMessageCardURL().
		PcUrl("http://www.baidu.com").
		IoSUrl("http://www.google.com").
		Url("http://open.feishu.com").
		AndroidUrl("http://www.jianshu.com").
		Build()

	// header
	header := larkcard.NewMessageCardHeader().
		Template(larkcard.TemplateGreen).
		Title(larkcard.NewMessageCardPlainText().
			Content("上线通知").
			Build()).
		Build()

	// Elements
	divElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content("**🕐 时间：**\\n2021-02-23 20:17:51").
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 谁处理了问题
	content := "✅ " + "name" + "已处理了此告警"
	processPersonElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(content).
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{divElement, processPersonElement}).
		CardLink(cardLink).
		String()
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := l.Lark.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(receiveId).
			Content(cardContent).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return err
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}

	return nil
}

func (l *Larkbiz) SendArgocdMsg(ctx context.Context, argocdBody model.ArgocdBody, receiveId, receiveIdType string) error {
	if receiveIdType == "" {
		receiveIdType = larkcontact.UserIdTypeOpenId
	}

	// "author": "{{(call .repo.GetCommitMetadata .app.status.sync.revision).Author}}",
	// "message": "{{(call .repo.GetCommitMetadata .app.status.operationState.operation.sync.revision).Message}}"
	// oncePer: app.status.operationState.syncResult.revision
	// config
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(true).
		UpdateMulti(false).
		Build()

	// CardUrl
	newURL := strings.Replace(argocdBody.TitleLink, "https://localhost:4000", l.Cfg.LarkSetting.ArgocdBaseUrl, 1)
	cardLink := larkcard.NewMessageCardURL().
		PcUrl(newURL).
		IoSUrl(newURL).
		Url(newURL).
		AndroidUrl(newURL).
		Build()

	// header
	cardColor := larkcard.TemplateGreen
	if argocdBody.State == "success" {
		cardColor = larkcard.TemplateGreen
	} else if argocdBody.State == "pending" {
		cardColor = larkcard.TemplateOrange
	} else {
		cardColor = larkcard.TemplateRed
	}
	header := larkcard.NewMessageCardHeader().
		Template(cardColor).
		Title(larkcard.NewMessageCardPlainText().
			Content("argocd 部署通知").
			Build()).
		Build()

	// Elements
	location, _ := time.LoadLocation("Asia/Shanghai")
	currentTime := time.Now().In(location).Format("2006-01-02 15:04:05")
	contentOne := fmt.Sprintf("**🕐 时间：** %s", currentTime)
	contentTwo := fmt.Sprintf("**⭐️ app: [部署详情](%s)", newURL)
	divElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content(contentOne).
					Build()).
				IsShort(true).
				Build(),
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content(contentTwo).
					Build()).
				IsShort(true).
				Build(),
		}).
		Build()

	// 归属
	content := "✅ " + argocdBody.Author + ": " + argocdBody.Message
	processPersonElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(content).
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{divElement, processPersonElement}).
		CardLink(cardLink).
		String()
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := l.Lark.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(receiveId).
			Content(cardContent).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return err
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
	return nil
}

func (l *Larkbiz) SendJarvisMsg(ctx context.Context, jarvis model.Jarvis, receiveId, receiveIdType string) error {
	if receiveIdType == "" {
		receiveIdType = larkcontact.UserIdTypeOpenId
	}

	// config
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(true).
		UpdateMulti(false).
		Build()

	// header
	header := larkcard.NewMessageCardHeader().
		Template(larkcard.TemplateBlue).
		Title(larkcard.NewMessageCardPlainText().
			Content("jarvis auto deploy通知").
			Build()).
		Build()

	// Elements
	location, _ := time.LoadLocation("Asia/Shanghai")
	currentTime := time.Now().In(location).Format("2006-01-02 15:04:05")
	contentOne := fmt.Sprintf("**🕐 时间：** %s", currentTime)
	divElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content(contentOne).
					Build()).
				IsShort(true).
				Build(),
		}).
		Build()

	// 归属
	content := "✅ " + jarvis.Author + ": " + jarvis.CommitMessage
	processPersonElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(content).
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{divElement, processPersonElement}).
		// CardLink(cardLink).
		String()
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := l.Lark.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(receiveId).
			Content(cardContent).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return err
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}

	return nil
}

func (l *Larkbiz) SendImagesSyncMsg(ctx context.Context, message, receiveId, receiveIdType string) error {
	if receiveIdType == "" {
		receiveIdType = larkcontact.UserIdTypeOpenId
	}

	// config
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(true).
		UpdateMulti(false).
		Build()

	// header
	header := larkcard.NewMessageCardHeader().
		Template(larkcard.TemplateBlue).
		Title(larkcard.NewMessageCardPlainText().
			Content("images 同步通知").
			Build()).
		Build()

	// Elements
	location, _ := time.LoadLocation("Asia/Shanghai")
	currentTime := time.Now().In(location).Format("2006-01-02 15:04:05")
	contentOne := fmt.Sprintf("**🕐 时间：** %s", currentTime)
	divElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content(contentOne).
					Build()).
				IsShort(true).
				Build(),
		}).
		Build()

	// 归属
	content := "✅ 镜像及版本: " + message
	processPersonElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(content).
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{divElement, processPersonElement}).
		// CardLink(cardLink).
		String()
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := l.Lark.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(receiveId).
			Content(cardContent).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return err
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}

	return nil
}
func (l *Larkbiz) OpenaiReply(ctx context.Context, content, receiveId, messageId, receiveIdType string) error {
	// 幂等
	_, ok := l.Lru.Get(messageId)
	if ok {
		l.Log.Info().Msgf("repeat message~")
		return nil
	} else {
		l.Lru.Add(messageId, 1)
	}
	openaiConfig := goopenai.DefaultConfig(l.Cfg.OpenaiSetting.Token)
	openaiConfig.BaseURL = l.Cfg.OpenaiSetting.OpenaiAddr
	client := goopenai.NewClientWithConfig(openaiConfig)
	var req []goopenai.ChatCompletionMessage

	ctxParent, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	// 处理streaming
	StreamChatCompletion := func(ctx context.Context, stream *goopenai.ChatCompletionStream) string {
		log.Log().Msgf("begin new chat stream")
		var content string
		defer func() {
			log.Log().Msgf("current chat stream end :%s", content)
		}()
		for {
			select {
			case <-ctx.Done():
				return ""
			default:
				response, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						// 完成读取
						// 自动上下文
						// req = append(req, goopenai.ChatCompletionMessage{
						// 	Role:    goopenai.ChatMessageRoleUser,
						// 	Content: content,
						// })
						return content
					} else {
						// 出现意外错误的情况
						log.Err(err).Msgf("StreamChatCompletion An error occurred midway %v", err)
					}
					return ""
				} else {
					if response.Choices[0].Delta.Content == "" {
						// 特殊处理
						continue
					}
					content = content + response.Choices[0].Delta.Content
				}
			}
		}
	}

	var result string
	wait := goAndWait(ctxParent,
		func() {
			defer func() {
				log.Log().Msg("StreamChatCompletion close!")
			}()
			req = append(req, goopenai.ChatCompletionMessage{
				Role:    goopenai.ChatMessageRoleUser,
				Content: content,
			})

			stream, err := client.CreateChatCompletionStream(ctxParent, goopenai.ChatCompletionRequest{
				Model:    goopenai.GPT3Dot5Turbo,
				Messages: req,
				Stream:   true,
			})
			if err != nil {
				log.Err(err).Msgf("StreamChatCompletion fail %v", err)
				return
			}
			result = StreamChatCompletion(ctxParent, stream)
		},
	)

	wait()

	// config
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(true).
		UpdateMulti(false).
		Build()

	// 消息主题
	processPersonElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(result).
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		// Header(header).
		Elements([]larkcard.MessageCardElement{processPersonElement}).
		// CardLink(cardLink).
		String()
	if err != nil {
		l.Log.Err(err).Msgf("cardContent Message.Create fail: %v", messageId)
		return err
	}

	resp, err := l.Lark.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(receiveId).
			Content(cardContent).
			Build()).
		Build())

	if err != nil {
		l.Log.Err(err).Msgf("OpenaiReply Message.Create fail: %v", messageId)
		return err
	}

	if !resp.Success() {
		l.Log.Error().Msgf("OpenaiReply Message.Create fail: %v, %v, %v", resp.Code, resp.Msg, resp.RequestId())
		return err
	}

	return nil

}

func (l *Larkbiz) ReactionReply(ctx context.Context, messageId string) error {
	emojis := []string{"SMILE", "OK", "THUMBSUP", "APPLAUSE", "JIAYI", "SMIRK", "PARTY", "FIREWORKS", "Trophy", "RAINBOWPUKE", "Hundred", "CheckMark"}
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	emoji := emojis[r.Intn(len(emojis))]
	resp, err := l.Lark.Im.MessageReaction.Create(ctx, // "SMILE"
		larkim.NewCreateMessageReactionReqBuilder().MessageId(messageId).Body(&larkim.CreateMessageReactionReqBody{
			ReactionType: &larkim.Emoji{
				EmojiType: &emoji,
			},
		}).Build(),
	)
	if err != nil {
		l.Log.Err(err).Msgf("MessageReaction.Create fail: %v", messageId)
		return err
	}

	if !resp.Success() {
		l.Log.Error().Msgf("MessageReaction.Create fail: %v, %v, %v", resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

func goAndWait(ctx context.Context, fns ...func()) func() {
	var wg sync.WaitGroup
	wg.Add(len(fns))
	for _, fn := range fns {
		go func(fn func()) {
			defer func() {
				if err := recover(); err != nil {
					// 记录日志
					log.Err(err.(error)).Msgf("stack: %s", debug.Stack())
				}
			}()
			defer wg.Done()
			fn()
		}(fn)
	}
	return wg.Wait
}
