package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"context"

	"go.etcd.io/etcd/client"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")

	var (
		etcdServer     = os.Args[1]
		etcdServerPort = os.Args[2]
	)

	if port == "" {
		port = "80"
	}

	if etcdServer == "" || etcdServerPort == "" {
		log.Fatal("missing database arguments")
	}

	//  Connect and initial the etcd server
	cfg := client.Config{
		Endpoints:               []string{"http://"+etcdServer+":"+etcdServerPort},
		Transport:               client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		 log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)

	fmt.Println("Initial etcd successfully.")

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:

					//Get the response from Database
					var response string
					value, err := kapi.Get(context.Background(),message.Text , nil)
					
					if err != nil {
						// No response match the keyword.
						response = "Sorry, I don't know what you say?"
						fmt.Println(err)
					} else {
						response = value.Node.Value
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}