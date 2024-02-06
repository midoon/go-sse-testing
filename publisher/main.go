package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	stockChannel := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go streamUpdatePrice(ctx, stockChannel)

	rmqConn, err := rmqProducerInit()
	if err != nil {
		panic(err)
	}

	defer rmqConn.Close()

	rmqChannelConnection, err := rmqConn.Channel()
	if err != nil {
		fmt.Println("error nih")
		panic(err)
	}

	for stock := range stockChannel {
		message := amqp091.Publishing{
			Headers: amqp091.Table{
				"ContentType": "text/plain",
			},
			Body: []byte(stock),
		}
		err := rmqChannelConnection.PublishWithContext(ctx, "sse", "stock", false, false, message)
		if err != nil {
			panic(err)
		}
	}
}

func rmqProducerInit() (*amqp091.Connection, error) {
	rmqConnection, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	return rmqConnection, err

}

func streamUpdatePrice(ctx context.Context, channel chan<- string) {
	ticker := time.NewTicker(time.Second)
	stocks := []string{"BBRI", "BBCA", "UNTR", "HEXA", "ICBP", "GOTO"}

labelBreak:
	for {
		select {
		case <-ctx.Done():
			break labelBreak
		case <-ticker.C:
			price := rand.Intn(10000)
			stockIndex := rand.Intn(len(stocks))

			channel <- fmt.Sprintf("%s-%d", stocks[stockIndex], price)
		}
	}

	ticker.Stop()
	close(channel)
}
