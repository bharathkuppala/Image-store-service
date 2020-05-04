package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

// NewConsumer ,,,
type NewConsumer struct{}

// SetupCustomer ...
func (nc NewConsumer) SetupCustomer() (sarama.Consumer, error) {
	fmt.Println("AM HERE INSIDE SetUp")
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	cust, err := sarama.NewConsumer(KafkaBrokers, config)
	if err != nil {
		log.Println("error in creating new customer", err)
		return nil, err
	}
	return cust, nil
}

// ConsumeClaim ...
func (nc NewConsumer) ConsumeClaim(consumer sarama.Consumer) {
	pCust, err := consumer.ConsumePartition(KafkaTopic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Println("error in creating consumer partition", err)
		return
	}
	// cSignals := make(chan os.Signal, 0)
	// signal.Notify(cSignals, os.Interrupt)
	// signal.Notify(cSignals, os.Kill)

	defer func() {
		custErr := pCust.Close()
		if custErr != nil {
			log.Println("error in closing customer")
			return
		}
	}()

	finishChannel := make(chan bool, 0)

	go func() {
		for {
			select {
			case msg := <-pCust.Messages():
				messagesProduced++
				fmt.Println("Received messages:", string(msg.Key), string(msg.Value))
			case <-time.After(3 * time.Second):
				fmt.Printf("number of messages recieved from broker are %d\n", messagesProduced)
				finishChannel <- true
			}
		}
	}()
	<-finishChannel
}
