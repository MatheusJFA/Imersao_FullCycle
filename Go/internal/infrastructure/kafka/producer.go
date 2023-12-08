package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

type Producer struct {
	configMap *ckafka.ConfigMap
}

func NewProducer(configMap *ckafka.ConfigMap) *Producer {
	return &Producer{
		configMap: configMap,
	}
}

func (p *Producer) Publish(messageValue any, key []byte, topic string) error {
	producer, err := ckafka.NewProducer(p.configMap)

	if err != nil {
		return err
	}

	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          messageValue,
		Key:            key,
	}

	err = producer.Produce(message, nil)

	if err != nil {
		return err
	}

	return nil
}
