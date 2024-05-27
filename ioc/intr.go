package ioc

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "xiaoweishu/webook/api/proto/gen/intr/v1"
	"xiaoweishu/webook/interactive/service"
	"xiaoweishu/webook/internal/client"
)

// 初始化grpc客户端
func InitIntrClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	type config struct {
		Addr      string `yaml:"addr"`
		Secure    bool
		Threshold int32
	}
	var cfg config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(cc)       //初始化远程客户端
	local := client.NewLocalInteractiveServiceAdapter(svc) //初始化本地客户端
	res := client.NewInteractiveClient(remote, local)
	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg = config{}
		err := viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err != nil {
			panic(err)
		}
		res.UpdateThreshold(cfg.Threshold)
	})
	return res
}
