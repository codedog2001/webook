// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"xiaoweishu/webook/internal/repository"
	"xiaoweishu/webook/internal/repository/cache"
	"xiaoweishu/webook/internal/repository/dao"
	"xiaoweishu/webook/internal/service"
	"xiaoweishu/webook/internal/web"
	"xiaoweishu/webook/internal/web/jwt"
	"xiaoweishu/webook/ioc"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	loggerV1 := InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, handler, loggerV1)
	db := InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	codeSerVice := service.NewCodeService(codeRepository, smsService)
	userHandLer := web.NewUserHandLer(userService, codeSerVice, handler)
	wechatService := InitWechatService(loggerV1)
	oAuth2WechatHandLer := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	articleDAO := dao.NewArticleGORMDAO(db)
	articleCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(articleDAO, userRepository, articleCache)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(loggerV1, articleService)
	engine := ioc.InitWebServer(v, userHandLer, oAuth2WechatHandLer, articleHandler)
	return engine
}

func InitArticleHandler(dao2 dao.ArticleDAO) *web.ArticleHandler {
	loggerV1 := InitLogger()
	db := InitDB()
	userDAO := dao.NewUserDAO(db)
	cmdable := InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	articleCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(dao2, userRepository, articleCache)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(loggerV1, articleService)
	return articleHandler
}

// wire.go:

var thirdPartySet = wire.NewSet(
	InitRedis, InitDB,
	InitLogger)

var userSvcProvider = wire.NewSet(dao.NewUserDAO, cache.NewUserCache, repository.NewUserRepository, service.NewUserService)

var articlSvcProvider = wire.NewSet(repository.NewCachedArticleRepository, cache.NewArticleRedisCache, dao.NewArticleGORMDAO, service.NewArticleService)