package startup

import (
	"xiaoweishu/webook/internal/pkg/logger"
	"xiaoweishu/webook/internal/service/oauth2/wechat"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", l)
}
