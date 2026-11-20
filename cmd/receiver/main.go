package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OliverVood/courier/cmd/internal/queue"
	"github.com/go-chi/chi/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	Connection *amqp.Connection
	Chanel     *amqp.Channel
	Queue      *amqp.Queue
}

func main() {
	fmt.Println("Launching the receiver management program...")

	var err error

	h := &Handler{}

	h.Connection, err = queue.GetConnection()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to broker...")
	defer func(connection *amqp.Connection) {
		err := queue.CloseConnection(connection)
		if err != nil {
			log.Fatal(err)
		}
	}(h.Connection)

	h.Chanel, err = queue.GetCannel(h.Connection)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to chanel...")
	defer func(cannel *amqp.Channel) {
		err := queue.CloseCannel(cannel)
		if err != nil {
			log.Fatal(err)
		}
	}(h.Chanel)

	h.Queue, err = queue.GetQueue(h.Chanel, "test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to queue...")

	router := chi.NewRouter()
	regRouts(router, h)
	runServer(router)
}

func regRouts(router chi.Router, h *Handler) {
	fmt.Println("Router registration...")

	router.Post("/message", h.store)
}

func runServer(router chi.Router) {
	fmt.Println("Server startup...")
	fmt.Println("Open site: http://localhost:8181")

	log.Fatal(http.ListenAndServe(":8181", router))
}
