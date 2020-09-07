package log_agent

import (
	"context"
	"fmt"
	"github.com/NullExceptionNull/go_structrue/kafka"
	"github.com/hpcloud/tail"
)

type LogTask struct {
	Path       string
	Topic      string
	Instance   *tail.Tail
	Ctx        context.Context //退出的context
	CancelFunc context.CancelFunc
}

func (t *LogTask) InitTail(path string) {
	tails, err := tail.TailFile(path, *config)
	t.Instance = tails
	if err != nil {
		fmt.Println("tail file err: ", err)
		return
	}
	fmt.Println("===Init tail success===")
	go t.Run()
}

func (t *LogTask) Run() {

	for {
		select {
		case <-t.Ctx.Done():
			fmt.Println("配置------------->" + t.Path + "-----" + t.Topic + " 已退出")
			return
		case line := <-t.Instance.Lines:
			kafka.SendToProducerChan(line.Text, t.Topic)
		}
	}
}

func NewLogTask(path, topic string) *LogTask {
	logTask := new(LogTask)
	ctx, cancelFunc := context.WithCancel(context.Background())
	logTask.Topic = topic
	logTask.Path = path
	logTask.Ctx = ctx
	logTask.CancelFunc = cancelFunc
	logTask.InitTail(path)
	return logTask
}
