package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		appId:    ekit.ToPtr(appId),
		signName: ekit.ToPtr(signName),
		client:   client,
	}

}
func (s Service) Send(ctx context.Context, tpl string, args []string, number ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId //这里要传入的是指针
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr(tpl) //可能会有问题，出问题再来看这里
	req.PhoneNumberSet = s.toPtrSlice(number)
	req.TemplateParamSet = s.toPtrSlice(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "ok" {
			return fmt.Errorf("发送短信失败%S,%s", *status.Code, *status.Message)
		}
	}
	return nil

}
func (s *Service) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data,
		func(idx int, src string) *string {
			return &src
		})
}
