package startup

import "xiaoweishu/webook/internal/pkg/logger"

func InitLogger() logger.LoggerV1 {
	return logger.NewNopLogger() //也就是暂时先不用日志，测试用不到日志
}
