package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id   int64  //任务的唯一标识符，通常是一个自增的整数，用于在系统内部唯一地识别每一个任务
	Name string //任务名称
	// Cron 表达式
	Expression string
	Executor   string //定义执行器，本地执行器或者远程执行器
	Cfg        string //配置
	CancelFunc func() //回调函数，成功执行后，释放相关的资源
}

func (j Job) NextTime() time.Time {
	c := cron.NewParser(cron.Second | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := c.Parse(j.Expression)
	//传入当前的时间根据表达式计算出下次执行的时间，比如每一个小时执行一次，
	//那么下次调度的时间就是在现在的时间是上加一个小时
	return s.Next(time.Now())
}
