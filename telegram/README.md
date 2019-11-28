

# Telegram Bot echo-reply example with etcd
The simply echo-reply example for Telegram bot. In this instance, the webhook server will connect to etcd and get the matched responses for the incoming messages. You can easily create events by inserting the keyword/response records into etcd.

## Preparation

Make sure you have imported the following packages:
1. [go-telegram-bot-api/telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
2. [etcd-io/etcd/client](https://github.com/etcd-io/etcd/tree/master/client)

## How to execute?
```{bash}
$ go get github.com/pohsienshih/chatbot-etcd-example/telegram
```
```{bash}
$ export TELEGRAM_TOKEN=<your access token>

$ cd $GOPATH/src/pohsienshih/chatbot-etcd-example/telegram
$ go build -o webhook .
$ ./webhook <etcd ip> <etcd port>
```
> Make sure you already have etcd service.

## Notice
TLS connection for this example is not yet supported. You can expose your service by using [ngrok](https://ngrok.com/).
```{bash}
$ ngrok http <port>
```


