package kafka

import (
	"testing"
	"time"
)

func TestServer_sendLog(t *testing.T) {

	//kafka.init(*kafkaconf)

	SendToProducerChan("zjhzjhzjh", "test_topic")

	time.Sleep(5 * time.Second)
}
