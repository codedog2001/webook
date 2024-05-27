package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	events2 "xiaoweishu/webook/interactive/events"
	"xiaoweishu/webook/interactive/grpc"
	"xiaoweishu/webook/interactive/ioc"
	"xiaoweishu/webook/interactive/repository"
	"xiaoweishu/webook/interactive/repository/cache"
	"xiaoweishu/webook/interactive/repository/dao"
	"xiaoweishu/webook/interactive/service"
	"xiaoweishu/webook/internal/events"
	"xiaoweishu/webook/pkg/grpcx"
)

type App struct {
	consumers []events.Consumer
	server    *grpcx.Server
}

func main() {
	initViperV1()
	app := InitApp()
	initPrometheus()
	for _, consumer := range app.consumers {
		err := consumer.Start()
		if err != nil {
			panic(err)
		}
	}
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}

}
func InitApp() *App {
	loggerV1 := ioc.InitLogger()
	db := ioc.InitDB(loggerV1)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := ioc.InitRedis()
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	client := ioc.InitSaramaClient()
	interactiveReadEventConsumer := events2.NewInteractiveReadEventConsumer(interactiveRepository, client, loggerV1)
	v := ioc.InitConsumers(interactiveReadEventConsumer)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	server := ioc.NewGrpcxServer(interactiveServiceServer)
	app := &App{
		consumers: v,
		server:    server,
	}
	return app
}

func initPrometheus() {
	go func() {
		// 专门给 prometheus 用的端口
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			return
		}
	}()
}
func initViperV1() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	// 这一步之后，cfile 里面才有值
	pflag.Parse()
	//viper.Set("db.dsn", "localhost:3306")
	// 所有的默认值放好s
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}
