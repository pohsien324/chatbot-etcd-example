# Line Bot echo-reply example with etcd
The simply echo-reply example for Line bot. In this example, the webhook server will connect to etcd and get the matched responses for the incoming messages. You can easily create events by inserting the keyword/value records into etcd.

## Preparation

Make sure you have imported the following packages:
1. [line-bot-sdk-go](https://github.com/line/line-bot-sdk-go)
2. [etcd-io/etcd/client](https://github.com/etcd-io/etcd/tree/master/client)

## How to execute?
```{bash}
$ go get github.com/pohsienshih/chatbot-etcd-example/line
```
```{bash}
$ export CHANNEL_SECRET=<yoursecret>
$ export CHANNEL_TOKEN=<yourtoken>
$ export PORT=<the port you want to listen on>

$ cd $GOPATH/src/pohsienshih/chatbot-etcd-example/line
$ go build -o webhook .
$ ./webhook <etcd server ip> <etcd port>
```
> Make sure you already have etcd service.

## Notice
TLS connection for this example is not yet supported. You can expose your service by using [ngrok](https://ngrok.com/).
```{bash}
$ ngrok http <port>
```


