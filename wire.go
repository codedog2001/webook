//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"xiaoweishu/webook/internal/repository"
	"xiaoweishu/webook/internal/repository/cache"
	"xiaoweishu/webook/internal/repository/dao"
	"xiaoweishu/webook/internal/service"
	"xiaoweishu/webook/internal/web"
	ijwt "xiaoweishu/webook/internal/web/jwt"
	"xiaoweishu/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//最低层的第三方依赖
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLogger,
		//dao层初始化
		dao.NewUserDAO,
		//cache初始化
		cache.NewCodeCache, cache.NewUserCache,
		//repository
		repository.NewCodeRepository, repository.NewUserRepository,
		//service部分
		service.NewCodeService, service.NewUserService, ioc.InitSMSService, ioc.InitWechatService,
		//这里必须使用函数返回sms.service ,直接使用mem.service会报错
		//handler部分
		web.NewUserHandLer,
		ioc.InitGinMiddlewares,
		ijwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitWebServer,
	)
	return gin.Default()

}
