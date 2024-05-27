package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"xiaoweishu/webook/internal/domain"
	"xiaoweishu/webook/pkg/logger"
)

type Result struct {
	AccessToken string `json:"access_token"`
	// access_token接口调用凭证超时时间，单位（秒）
	ExpiresIn int64 `json:"expires_in"`
	// 用户刷新access_token
	RefreshToken string `json:"refresh_token"`
	// 授权用户唯一标识
	OpenId string `json:"openid"`
	// 用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`
	// 当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段。
	UnionId string `json:"unionid"`

	// 错误返回
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// 回调链接
var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx *gin.Context, state string) (string, error)
	VerifyCode(ctx *gin.Context, code string, state string) (domain.WechatInfo, error)
}
type service struct {
	appID     string
	appSecret string
	client    *http.Client
}

// service必须要实现Service的所有方法，才能将他作为接口的值返出去
func NewService(appID string, appSecret string, l logger.LoggerV1) Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
		//严格来说，这里是要依赖注入的，但是这里并没有别的地方用，所以直接用默认的地方就可以了
		client: http.DefaultClient,
	}
}
func (s *service) AuthURL(ctx *gin.Context, state string) (string, error) {
	const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
	return fmt.Sprintf(authURLPattern, s.appID, redirectURL, state), nil
}
func (s *service) VerifyCode(ctx *gin.Context, code string, state string) (domain.WechatInfo, error) {
	accessTokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`,
		s.appID, s.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	httpResp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	var res Result
	//decoder 只读了一遍body
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		//转json结构体时出错
		return domain.WechatInfo{}, err
	}
	//未设置，默认是0，如果不为0的话，说明出错了
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("调用微信接口失败 errcode %d,errmsg %s", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		UnionId: res.UnionId,
		OpenId:  res.OpenId,
	}, nil
}
