package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"xiaoweishu/webook/internal/repository"
	"xiaoweishu/webook/internal/repository/cache"
	"xiaoweishu/webook/internal/repository/dao"
	"xiaoweishu/webook/internal/service"
	"xiaoweishu/webook/internal/web"
	"xiaoweishu/webook/internal/web/jwt"
	"xiaoweishu/webook/ioc"
)

func main() {
	initLogger()
	initViper()
	server := InitWebServerv1()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello，启动成功了！")
	})
	err := server.Run(":8080")
	if err != nil {
		return
	}
}
func InitWebServerv1() *gin.Engine {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	loggerV1 := ioc.InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, handler, loggerV1)
	db := ioc.InitDB(loggerV1)
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	codeSerVice := service.NewCodeService(codeRepository, smsService)
	userHandLer := web.NewUserHandLer(userService, codeSerVice, handler)
	wechatService := ioc.InitWechatService(loggerV1)
	oAuth2WechatHandLer := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	engine := ioc.InitWebServer(v, userHandLer, oAuth2WechatHandLer)
	return engine
}
func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
func initViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	// 当前工作目录的 config 子目录
	viper.AddConfigPath("config")
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

// 用viper和远程控制中心etcd配合使用
func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3",
		"http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("远程配置中心发生变更")
	})
	go func() {
		for {
			err = viper.WatchRemoteConfig()
			if err != nil {
				panic(err)
			}
			log.Println("watch", viper.GetString("test.key"))
			//time.Sleep(time.Second)
		}
	}()
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

//func useSession(server *gin.Engine) {
//	login := &middleware.LoginMiddlewareBuilder{}
//	// 存储数据的，也就是你 userId 存哪里
//	// 直接存 cookie
//	store := cookie.NewStore([]byte("secret"))
//	// 基于内存的实现
//	//store := memstore.NewStore([]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
//	//	[]byte("eF1`yQ9>yT1`tH1,sJ0.zD8;mZ9~nC6("))
//	//store, err := redis.NewStore(16, "tcp",
//	//	"localhost:6379", "",
//	//	[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
//	//	[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgA"))
//	//if err != nil {
//	//	panic(err)
//	//}
//	var jwt:=web.NewUserHandLer()
//	server.Use(sessions.Sessions("ssid", store), web.NewUserHandLer)
