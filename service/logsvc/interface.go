package logsvc

import (
	"feishu/service/snowflakesvc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger(snowflake *snowflakesvc.Snowflake) zerolog.Logger {
	return log.With().Str("request_id", snowflake.Generate().String()).Logger()
}
