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

type LogData struct {
	topic string
	data  string
}

var producer sarama.SyncProducer

var ProductChan chan *LogData

func initProps() *conf.KafkaConf {
	load, err := ini.Load("../conf/conf.ini")
	if err != nil {
		fmt.Println("init 文件加载失败")
	}
	var kafkaConf = new(conf.KafkaConf)

	err = load.Section("kafka").MapTo(kafkaConf)

	if err != nil {
		fmt.Println("init 文件解析失败")
	}
	return kafkaConf
}

func init() {
	props := initProps()
	collector := newDataCollector(props.Address)
	producer = collector
	ProductChan = make(chan *LogData, props.ChanSize)
	fmt.Println("===init kafka producer success===")
	SendMsg()
}

func SendMsg() {
	message := &sarama.ProducerMessage{}
	for data := range ProductChan {
		message.Value = sarama.StringEncoder(data.data)
		message.Topic = data.topic
		message.Timestamp = time.Now()
		sendMessage, offset, err := producer.SendMessage(message)
		fmt.Println(sendMessage, offset)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendToProducerChan(msg string, topic string) {
	logData := new(LogData)
	logData.topic = topic
	logData.data = msg
	ProductChan <- logData
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
