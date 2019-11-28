package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.etcd.io/etcd/client"
)

func main() {

	var (
		etcdServer     = os.Args[1]
		etcdServerPort = os.Args[2]
	)
	token := os.Getenv("TELEGRAM_TOKEN")

	if token == "" {
		log.Fatal("Missing tolen")
	}

	if etcdServer == "" || etcdServerPort == "" {
		log.Fatal("missing database arguments")
	}

	//  Connect and initial the etcd server
	cfg := client.Config{
		Endpoints: []string{"http://" + etcdServer + ":" + etcdServerPort},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)

	fmt.Println("Initial etcd successfully.")

	// Setup telegram bot
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := client.GetUpdatesChan(u)

	// Get the message from client
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//Get the response from Database
		var response string
		value, err := kapi.Get(context.Background(), update.Message.Text, nil)

		if err != nil {
			// No response match the keyword.
			response = "Sorry, I don't know what you say?"
			fmt.Println(err)
		} else {
			response = value.Node.Value
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		client.Send(msg)
	}

}
