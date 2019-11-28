package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/paked/messenger"
	"go.etcd.io/etcd/client"
)

var (
	verifyToken = os.Getenv("VERIFY_TOKEN")
	pageToken   = os.Getenv("PAGE_TOKEN")
	port        = os.Getenv("PORT")
)

func main() {
	var (
		etcdServer     = os.Args[1]
		etcdServerPort = os.Args[2]
	)

	if verifyToken == "" || pageToken == "" {
		log.Fatal("missing arguments")
	}

	if port == "" {
		port = "80"
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

	// Setup Messenger Bot Webhook Server
	client := messenger.New(messenger.Options{
		VerifyToken: verifyToken,
		Token:       pageToken,
	})

	client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		//Get the response from Database
		var response string
		value, err := kapi.Get(context.Background(), m.Text, nil)

		if err != nil {
			// No response match the keyword.
			response = "Sorry, I don't know what you say?"
			fmt.Println(err)
		} else {
			response = value.Node.Value
		}

		r.Text(response, messenger.ResponseType)
	})

	if err := http.ListenAndServe(":"+port, client.Handler()); err != nil {
		log.Fatal(err)
	}
}
