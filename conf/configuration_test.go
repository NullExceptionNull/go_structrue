package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"testing"
)

func Test(t *testing.T) {
	load, err := ini.Load("conf.ini")
	if err != nil {
		fmt.Println("init 文件加载失败")
	}
	var kafkaconf = new(KafkaConf)

	err = load.Section("kafka").MapTo(kafkaconf)

	if err != nil {
		fmt.Println("init 文件解析失败")
	}
	fmt.Println("\v", &kafkaconf.Address)
}
