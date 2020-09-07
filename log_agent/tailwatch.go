package log_agent

import (
	"fmt"
	"github.com/NullExceptionNull/go_structrue/kafka"
	"github.com/hpcloud/tail"
	"sync"
)

const FILE_NAME string = "./log"

var config *tail.Config

func init() {
	config = &tail.Config{
		Location:    &tail.SeekInfo{Offset: 0, Whence: 2}, //0 代表从文件开头开始算起，1 代表从当前位置开始算起，2 代表从文件末尾算起。
		ReOpen:      true,
		MustExist:   false,
		Poll:        true,
		Pipe:        false,
		RateLimiter: nil,
		Follow:      true,
		MaxLineSize: 0,
		Logger:      nil,
	}
}

func run() {
	lines, err := Read()

	if err != nil {
		fmt.Println("tail error")
	}

	for line := range lines.Lines {
		//这里发到kafka
		fmt.Println(line.Text)
		kafka := kafka.Kafka{Once: sync.Once{}}
		kafka.SendLog(line.Text, "test_topic")
	}
}

func Read() (*tail.Tail, error) {
	tails, err := tail.TailFile(FILE_NAME, *config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return nil, err
	}

	return tails, nil
}
