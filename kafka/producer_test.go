package kafka

import (
	"fmt"
	"github.com/NullExceptionNull/go_structrue/conf"
	"gopkg.in/ini.v1"
	"sync"
	"testing"
)

func TestServer_sendLog(t *testing.T) {
	load, err := ini.Load("../conf/conf.ini")
	if err != nil {
		fmt.Println("init 文件加载失败")
	}
	var kafkaconf = new(conf.KafkaConf)

	err = load.Section("kafka").MapTo(kafkaconf)

	if err != nil {
		fmt.Println("init 文件解析失败")
	}
	fmt.Println("\v", &kafkaconf.Address)

	kafka := Kafka{Once: sync.Once{}}

	//kafka.init(*kafkaconf)

	kafka.SendLog("golang", "test_topic")
}
