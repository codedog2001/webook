//go:build wireinject

package main

import (
	"github.com/google/wire"
	"xiaoweishu/webook/interactive/events"
	repository2 "xiaoweishu/webook/interactive/repository"
	cache2 "xiaoweishu/webook/interactive/repository/cache"
	dao2 "xiaoweishu/webook/interactive/repository/dao"
	service2 "xiaoweishu/webook/interactive/service"
	"xiaoweishu/webook/internal/events/article"
	"xiaoweishu/webook/internal/repository"
	"xiaoweishu/webook/internal/repository/cache"
	"xiaoweishu/webook/internal/repository/dao"
	"xiaoweishu/webook/internal/service"
	"xiaoweishu/webook/internal/web"
	ijwt "xiaoweishu/webook/internal/web/jwt"
	"xiaoweishu/webook/ioc"
)

var interactiveSVCSet = wire.NewSet(
	dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService)

func InitWebServer() *App {
	wire.Build(
		//最低层的第三方依赖
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLogger,
		//初始化sarama客户端，生产者一般用自带的
		ioc.InitSaramaClient,
		ioc.InitSyncProducer,
		ioc.InitIntrClient,
		ioc.
			interactiveSVCSet,
		//dao层初始化
		dao.NewUserDAO,
		dao.NewArticleGORMDAO,

		events.NewInteractiveReadEventConsumer,
		article.NewSaramaSyncProducer,
		//消费者一般需要自己自定义
		ioc.InitConsumers,
		//cache初始化
		cache.NewCodeCache,
		cache.NewUserCache,
		cache.NewArticleRedisCache,
		//repository
		repository.NewCodeRepository,
		repository.NewUserRepository,
		repository.NewCachedArticleRepository,
		//service部分
		service.NewCodeService,
		service.NewUserService,
		service.NewArticleService,
		ioc.InitSMSService,
		ioc.InitWechatService,

		//这里必须使用函数返回sms.service ,直接使用mem.service会报错
		//handler部分
		web.NewArticleHandler,
		web.NewUserHandLer,
		ioc.InitGinMiddlewares,
		ijwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitWebServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)

}
