package queue

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	_ "strconv"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type config struct {
	connection configConnection
}

type configConnection struct {
	host string
	port string
	user string
	pass string
}

func GetConnection() (*amqp.Connection, error) {
	config, err := getConfig("configs/queue.env")
	if err != nil {
		return nil, err
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.connection.user, config.connection.pass, config.connection.host, config.connection.port))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CloseConnection(conn *amqp.Connection) error {
	return conn.Close()
}

func GetCannel(connection *amqp.Connection) (*amqp.Channel, error) {
	ch, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func CloseCannel(ch *amqp.Channel) error {
	return ch.Close()
}

func GetQueue(chanel *amqp.Channel, name string) (*amqp.Queue, error) {
	queue, err := chanel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum, // использовать кворумную очередь (рекомендуется)
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &queue, nil
}

func Send(channel *amqp.Channel, queue *amqp.Queue, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := channel.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func getConfig(file string) (*config, error) {
	config := &config{}
	err := config.load(file)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *config) load(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		return errors.New("error loading Queue config")
	}

	c.connection.host = os.Getenv("RABBITMQ_HOST")
	c.connection.port = os.Getenv("RABBITMQ_PORT")
	c.connection.user = os.Getenv("RABBITMQ_USER")
	c.connection.pass = os.Getenv("RABBITMQ_PASS")

	return nil
}
