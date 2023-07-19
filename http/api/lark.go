package api

import (
	"context"
	"net/http"

	larkBiz "feishu/biz/lark"
	"feishu/config"
	"feishu/http/middleware"
	"feishu/model"
	"feishu/util"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/rs/zerolog/log"
)

func GroupInfo(c *gin.Context) {
	defer util.ApiSeg(c)()

	var body struct {
		ChatId string `json:"chat_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Err(err).Msg("failed to bind body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "malformed request"})
		return
	}

	err := middleware.Dependency(c).LarkHub.Lark.GroupInfo(c, body.ChatId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get GroupInfo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func GroupList(c *gin.Context) {
	defer util.ApiSeg(c)()

	res, err := middleware.Dependency(c).LarkHub.Lark.GroupList(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get GroupList"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func GroupMessageCreate(c *gin.Context) {
	defer util.ApiSeg(c)()

	var body struct {
		ChatId string `json:"chat_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Err(err).Msg("failed to bind body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "malformed request"})
		return
	}

	err := middleware.Dependency(c).LarkHub.Lark.GroupMessageCreate(c, body.ChatId, "")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get GroupMessageCreate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func SendInteractiveMonitorMsg(c *gin.Context) {
	defer util.ApiSeg(c)()

	var body struct {
		ReceiveId     string `json:"receive_id"`
		ReceiveIdType string `json:"receive_id_type"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Err(err).Msg("failed to bind body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "malformed request"})
		return
	}

	err := middleware.Dependency(c).LarkHub.Lark.SendInteractiveMonitorMsg(c, body.ReceiveId, body.ReceiveIdType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get SendInteractiveMonitorMsg"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func SendInteractiveMsg(c *gin.Context) {
	defer util.ApiSeg(c)()

	var body struct {
		ReceiveId     string `json:"receive_id"`
		ReceiveIdType string `json:"receive_id_type"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Err(err).Msg("failed to bind body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "malformed request"})
		return
	}

	err := middleware.Dependency(c).LarkHub.Lark.SendInteractiveMsg(c, body.ReceiveId, body.ReceiveIdType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get SendInteractiveMsg"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func SendArgocdMsg(c *gin.Context) {
	defer util.ApiSeg(c)()

	log.Info().Msgf("c.Request: %v", c.Request.Body)
	body := model.ArgocdBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Err(err).Msg("failed to bind body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "malformed request"})
		return
	}

	log.Info().Msgf("SendArgocdMsg body: %v", body)

	err := middleware.Dependency(c).LarkHub.Lark.SendArgocdMsg(c, body, config.Options().LarkSetting.ArgocdChatId, "chat_id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get SendArgocdMsg"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ReactionReply(ctx context.Context, event *larkim.P2MessageReactionCreatedV1) error {
	log.Log().Msgf("ReactionReply event %v", *event.Event.MessageId)
	err := larkBiz.Object.ReactionReply(ctx, *event.Event.MessageId)
	if err != nil {
		log.Err(err).Msgf("failed to OpenaiReply %v", event.Event.MessageId)
		return err
	}
	return nil
}

func ReactionReplyHttp(c *gin.Context) {
	defer util.ApiSeg(c)()

	var body struct {
		MessageId string `json:"message_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Err(err).Msg("failed to bind body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ReactionReplyHttp malformed request"})
		return
	}

	err := larkBiz.Object.ReactionReply(c, body.MessageId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get OpenaiReplyHttp"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func OpenaiReply(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	log.Log().Msgf("OpenaiReply event content:%v - chatId:%v - senderId:%v - messageId:%v", *event.Event.Message.Content, *event.Event.Message.ChatId, *event.Event.Sender.SenderId, *event.Event.Message.MessageId)
	var p struct {
		Text string `json:"text"`
	}
	err := json.Unmarshal([]byte(*event.Event.Message.Content), &p)
	if err != nil {
		log.Err(err).Msgf("failed to Unmarshal %v", event.Event.Message.Content)
		return err
	}
	err = larkBiz.Object.OpenaiReply(ctx, p.Text, *event.Event.Message.ChatId, *event.Event.Message.MessageId, "chat_id")
	if err != nil {
		log.Err(err).Msgf("failed to OpenaiReply %v", event.Event.Message.Content)
		return err
	}

	return nil
}

func OpenaiReplyHttp(c *gin.Context) {
	defer util.ApiSeg(c)()

	var body struct {
		ChatId    string `json:"chat_id"`
		Content   string `json:"content"`
		MessageId string `json:"message_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Err(err).Msg("failed to bind body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "OpenaiReplyHttp malformed request"})
		return
	}

	err := larkBiz.Object.OpenaiReply(c, body.Content, body.ChatId, body.MessageId, "chat_id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to get OpenaiReplyHttp"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}
