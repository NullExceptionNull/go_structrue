package log_agent

import (
	"fmt"
	"github.com/hpcloud/tail"
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

func Read() {
	tails, err := tail.TailFile(FILE_NAME, *config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	for {
		for line := range tails.Lines {
			fmt.Println(line.Text)
		}
	}

}
