package larksvc

import (
	"fmt"
	"net/http"
	"time"

	"feishu/config"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

func MustNewLark(config *config.AppConfig) *lark.Client {
	fmt.Println("config.LarkSetting.AppID:", config.LarkSetting.AppID)
	fmt.Println("config.LarkSetting.AppSecret:", config.LarkSetting.AppSecret)
	return lark.NewClient(config.LarkSetting.AppID, config.LarkSetting.AppSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithReqTimeout(3*time.Second),
		lark.WithEnableTokenCache(true),
		lark.WithHelpdeskCredential("id", "token"),
		lark.WithHttpClient(http.DefaultClient),
		lark.WithLogReqAtDebug(true),
	)
}
