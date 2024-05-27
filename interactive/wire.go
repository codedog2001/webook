//go:build wireinject

package main

import (
	"github.com/google/wire"
	"xiaoweishu/webook/interactive/events"
	"xiaoweishu/webook/interactive/grpc"
	"xiaoweishu/webook/interactive/ioc"
	repository2 "xiaoweishu/webook/interactive/repository"
	cache2 "xiaoweishu/webook/interactive/repository/cache"
	dao2 "xiaoweishu/webook/interactive/repository/dao"
	service2 "xiaoweishu/webook/interactive/service"
)

var thirdPartySet = wire.NewSet(ioc.InitDB,
	ioc.InitLogger,
	ioc.InitSaramaClient,
	ioc.InitRedis)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService,
)

func InitApp() *App {
	wire.Build(thirdPartySet,
		interactiveSvcSet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitConsumers,
		ioc.NewGrpcxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
