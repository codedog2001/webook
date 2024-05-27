package main

import (
	"xiaoweishu/webook/internal/events"
	"xiaoweishu/webook/pkg/grpcx"
)

type App struct {
	consumers []events.Consumer
	server    *grpcx.Server
}
