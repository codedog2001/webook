package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"xiaoweishu/webook/internal/domain"
	"xiaoweishu/webook/internal/service"
)

type UserHandLer struct {
	svc         *service.UserService //其实就是gorm,db
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

func NewUserHandLer(svc *service.UserService) *UserHandLer {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandLer{
		emailExp:    emailExp,
		passwordExp: passwordExp,
		svc:         svc,
	}
}

func (u *UserHandLer) RegisterUsersRoutes(server *gin.Engine) {
	//设置分组路由
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/profile1", u.Profile1)
	ug.POST("/edit", u.Edit)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/profile", u.ProfileJWT)
}
func (u *UserHandLer) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	//binb方法会根据content -Type 来解析你的数据到req里面
	//解析错了，就会直接写回一个400错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//const只能放编译期就能确定的东西，emailExp是运行才能确定的东西，是不能直接放到这里面的

	//在参数校验的时候，一般只有超时了才会出现err ，timeout

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		//这里不能把详细的错误返回给前端，但是可以记录到日志中
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		//这里不能把详细的错误返回给前端，但是可以记录到日志中
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的密码必须大于8位，包含数字，特殊字符")
		return
	}
	fmt.Println(req)
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	return
	//接下来就到了数据库操作

}

//	func (u UserHandLer) Login(ctx *gin.Context) {
//		type LoginReq struct {
//			Email    string `json:"email"`
//			Password string `json:"password"`
//		}
//		var req LoginReq
//		if err := ctx.Bind(&req); err != nil {
//			return
//		} //拿到参数之后，就要进入下一层进行其他的逻辑处理
//		err := u.svc.Login(ctx, req.Email, req.Password)
//		if errors.Is(err, service.ErrInvalidUserOrPassword) {
//			ctx.String(http.StatusOK, "用户名或密码不对")
//			return
//		}
//		if err != nil {
//			ctx.String(http.StatusOK, "系统错误")
//			return
//		}
//		ctx.String(http.StatusOK, "登录成功")
//		return
//	}
func (u UserHandLer) Edit(ctx *gin.Context) {

}
func (u UserHandLer) Profile1(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是你的profile")
}
func (u UserHandLer) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	//必然有claims
	if !ok {
		//奇怪的错误，监控住这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	claims, ok := c.(*UserClaims) //类型断言
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "这是你的profile")
	println(claims.Uid)
	//这边补充p'r
}
func (h UserHandLer) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		uc := UserClaims{
			Uid:       u.Id,
			UserAgent: ctx.GetHeader("User-Agent"),
			RegisteredClaims: jwt.RegisteredClaims{
				// 1 分钟过期
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		tokenStr, err := token.SignedString(JWTKey)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
		}
		ctx.Header("x-jwt-token", tokenStr)
		ctx.String(http.StatusOK, "登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

// UserClaims 声明你自己要放进去token里面的数据
type UserClaims struct {
	jwt.RegisteredClaims //这个字段是包里面实现好的，直接组合起来使用就行了
	Uid                  int64
	UserAgent            string
}
