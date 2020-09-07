package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/NullExceptionNull/go_structrue/conf"
	"github.com/Shopify/sarama"
	"gopkg.in/ini.v1"
	"log"
	"sync"
	"time"
)

type Kafka struct {
	sync.Once
}

var producer sarama.SyncProducer

func initProps() *conf.KafkaConf {
	load, err := ini.Load("../conf/conf.ini")
	if err != nil {
		fmt.Println("init 文件加载失败")
	}
	var kafkaconf = new(conf.KafkaConf)

	err = load.Section("kafka").MapTo(kafkaconf)

	if err != nil {
		fmt.Println("init 文件解析失败")
	}
	return kafkaconf
}

func init() {
	props := initProps()
	collector := newDataCollector(props.Address)
	producer = collector
}

func (kafka *Kafka) SendLog(msg string, topic string) {
	message := sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(msg), Timestamp: time.Now()}
	_, _, err := producer.SendMessage(&message)
	if err != nil {
		fmt.Printf("msg %s ,topic %s 发送失败\n %v", msg, topic, err)
	}
}

type accessLogEntry struct {
	Method       string  `json:"method"`
	Host         string  `json:"host"`
	Path         string  `json:"path"`
	IP           string  `json:"ip"`
	ResponseTime float64 `json:"response_time"`

	encoded []byte
	err     error
}

func (ale *accessLogEntry) ensureEncoded() {
	if ale.encoded == nil && ale.err == nil {
		ale.encoded, ale.err = json.Marshal(ale)
	}
}

func (ale *accessLogEntry) Length() int {
	ale.ensureEncoded()
	return len(ale.encoded)
}

func (ale *accessLogEntry) Encode() ([]byte, error) {
	ale.ensureEncoded()
	return ale.encoded, ale.err
}

func newDataCollector(brokerList []string) sarama.SyncProducer {

	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	return producer
}

//func newAccessLogProducer(brokerList []string) sarama.AsyncProducer {
//
//	// For the access log, we are looking for AP semantics, with high throughput.
//	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
//	config := sarama.NewConfig()
//
//	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
//	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
//	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
//
//	producer, err := sarama.NewAsyncProducer(brokerList, config)
//	if err != nil {
//		log.Fatalln("Failed to start Sarama producer:", err)
//	}
//	// We will just log to STDOUT if we're not able to produce messages.
//	// Note: messages will only be returned here after all retry attempts are exhausted.
//	go func() {
//		for err := range producer.Errors() {
//			log.Println("Failed to write access log entry:", err)
//		}
//	}()
//
//	return producer
//}
