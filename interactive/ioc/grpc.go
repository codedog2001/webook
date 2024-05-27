package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpc2 "xiaoweishu/webook/interactive/grpc"
	"xiaoweishu/webook/pkg/grpcx"
)

func NewGrpcxServer(intrSvc *grpc2.InteractiveServiceServer) *grpcx.Server {
	s := grpc.NewServer()
	intrSvc.Register(s)
	add := viper.GetString("grpc.server.addr")
	return &grpcx.Server{
		Server: s,
		Addr:   add,
	}
}
