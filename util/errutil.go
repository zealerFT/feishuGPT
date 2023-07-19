package util

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ApiSeg(ctx *gin.Context) func() {
	return func() {
		if err := recover(); err != nil {
			// 记录日志
			errs := errors.New(fmt.Sprintf("stack:%+v\n", err))
			log.Err(errs).Msgf("ApiSeg end")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.New("服务器异常～"))
		}
	}
}
