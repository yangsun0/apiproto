package dataaccess

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
)

type MessageQueue struct {
	client *pubsub.Client
	context context.Context
	topic *pubsub.Topic
	sub *pubsub.Subscription
}

type ClickEvent struct {
	Target string	`json:"target"`
	Timestamp int64	`json:"timestamp"`
}


func (mq *MessageQueue) Connect(ctx context.Context) error {
	var projectId string = os.Getenv("GCP_PROJECT_ID")
	mq.context = ctx
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("create client: %v", err)
	}
	mq.client = client
	var topicId string = os.Getenv("GCP_TOPIC_ID")
	mq.topic = client.Topic(topicId)
	subId := topicId + "-sub"
	mq.sub = client.Subscription(subId)

	fmt.Println("mq connect successfully")
	return nil
}

func (mq *MessageQueue) Close() {
	fmt.Println("mq close")
	mq.client.Close()
}


func (mq* MessageQueue) Publish(clickEvent ClickEvent) error {
	data, _ := json.Marshal(&clickEvent)
	fmt.Println(string(data))
	msg := pubsub.Message{Data: data}
	result := mq.topic.Publish(mq.context, &msg)
	_, err := result.Get(mq.context)
	if err != nil {
		return fmt.Errorf("result.Get: %v", err)
	}

	return nil
}

func (mq* MessageQueue) Pull() error {
	err := mq.sub.Receive(mq.context, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Got message: %v\n", string(msg.Data))
		msg.Ack()
		var clickEvent ClickEvent = ClickEvent{}
		err := json.Unmarshal(msg.Data, &clickEvent)
		if err != nil {
			fmt.Printf("receive: %v\n", err)
		}else {
			fmt.Println(clickEvent)
		}
	})
	if err != nil {
			return fmt.Errorf("receive: %v", err)
	}
	return nil
}
