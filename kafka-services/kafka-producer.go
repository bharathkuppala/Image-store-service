package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	utility "github.com/image-store/webservice/utilities"
)

// NewProducer ...
type NewProducer struct{}

// SetupProducer ...
func (np *NewProducer) SetupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.ChannelBufferSize = 4
	config.Producer.Return.Successes = true
	config.Version = sarama.V0_10_0_0

	utility.ValidateConfig(config)

	producer, err := sarama.NewSyncProducer(KafkaBrokers, config)
	if err != nil {
		log.Fatalln("error in creating new sync producer", err)
		return nil, err
	}

	return producer, nil
}

// ProducerMessage ...
func (np *NewProducer) ProducerMessage(producer sarama.SyncProducer, notification string) {
	message := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Value: sarama.StringEncoder(notification),
	}
	partion, offSet, err := producer.SendMessage(message)
	if err != nil {
		log.Println("error in producing the notification:", err)
		return
	}
	log.Println("Produced notification message: ", notification)
	log.Printf("message produced for topic %s is stored in %d partion, %d offset\n", KafkaTopic, partion, offSet)

	defer func() {
		err := producer.Close()
		if err != nil {
			log.Println("error in closing producer")
			return
		}
	}()
}
