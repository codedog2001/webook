package grpcx

import (
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	*grpc.Server
	Addr string
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	//包装了grpc自带是serve，把监听字段也加入到这里来
	//编程风格更清晰
	return s.Server.Serve(l)
}
