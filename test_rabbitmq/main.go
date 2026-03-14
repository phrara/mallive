package main

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	host := "rabbitmq.default.svc.cluster.local"
	port := "5672"
	username := "guest"
	password := "guest"

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)

	fmt.Printf("Testing RabbitMQ connection to %s:%s...\n", host, port)

	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Printf("RabbitMQ connection failed: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("RabbitMQ connection OK!")

	// 测试 channel
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Open channel failed: %v\n", err)
		return
	}
	defer ch.Close()

	// 测试声明队列
	_, err = ch.QueueDeclare(
		"test_queue",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		fmt.Printf("Declare queue failed: %v\n", err)
		return
	}

	fmt.Println("Queue declare OK!")

	// 设置 QoS
	err = ch.Qos(10, 0, false)
	if err != nil {
		fmt.Printf("Set QoS failed: %v\n", err)
		return
	}

	fmt.Println("All tests passed!")

	// 保持连接一段时间以便观察
	time.Sleep(2 * time.Second)
}
