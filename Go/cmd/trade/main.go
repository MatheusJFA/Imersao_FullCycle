package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/MatheusJFA/Imersao_FullCycle/Go/internal/infrastructure/kafka"
	"github.com/MatheusJFA/Imersao_FullCycle/Go/internal/market/dto"
	"github.com/MatheusJFA/Imersao_FullCycle/Go/internal/market/entity"
	"github.com/MatheusJFA/Imersao_FullCycle/Go/internal/market/transformer"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	orderIn := make(chan *entity.Order)
	orderOut := make(chan *entity.Order)

	waitGroup := &sync.WaitGroup{}

	defer waitGroup.Wait()

	kafkaMessageChannel := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "trader",
		"auto.offset.reset": "latest",
	}

	producer := kafka.NewProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"order"})

	go kafka.Consume(kafkaMessageChannel)

	book := entity.NewBook(orderIn, orderOut, waitGroup)

	go book.Trade()

	go func() {
		for message := range kafkaMessageChannel {
			waitGroup.Add(1)
			fmt.Println(string(message.Value))

			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(message.Value, &tradeInput)

			if err != nil {
				panic(err)
			}

			order := transformer.TransformInput(tradeInput)
			orderIn <- order
		}
	}()

	for result := range orderOut {
		output := transformer.TransformOutput(result)
		outputJSON, err := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(outputJSON))

		if err != nil {
			panic(err)
		}

		producer.Publish(outputJSON, []byte("orders"), "output")
	}

}
