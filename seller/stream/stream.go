package stream

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	"github.com/Shopify/sarama"

	sellerpb "github.com/ihtkas/farm/seller/v1"
)

// KafkaProducer implements  message broker functionalities for seller service
type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// Init creates a new async producer for given kafka broker addresses.
func (p *KafkaProducer) Init(brokerAddrs []string, topic string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerAddrs, config)
	if err != nil {
		return err
	}
	p.producer = producer
	p.topic = topic
	return nil
}

// PublishNewProduct produces new product in the Kafka steam for matcher service.
func (p *KafkaProducer) PublishNewProduct(product *sellerpb.ProductInfo) error {
	producer := p.producer
	bs, err := proto.Marshal(product)
	msg := &sarama.ProducerMessage{Topic: p.topic, Value: sarama.ByteEncoder(bs)}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		glog.Errorln("FAILED to send message: %s\n", err)
		return err
	}
	glog.Info("Message sent to partition %d at offset %d\n", partition, offset)
	return nil
}
