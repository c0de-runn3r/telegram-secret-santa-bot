package main

import (
	"log"
	"os"

	tgClient "main/clients/telegram"
	event_consumer "main/consumer/event-consumer"
	"main/events/telegram"

	"github.com/joho/godotenv"
)

// to glue everything here

const batchSize = 100

func main() {
	token, host := processENV()

	eventsProcessor := telegram.New(
		tgClient.New(host, token),
	)
	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func processENV() (token string, host string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token = os.Getenv("TOKEN")
	host = os.Getenv("HOST")
	return token, host
}
