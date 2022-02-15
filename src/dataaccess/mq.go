package dataaccess

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/yangsun0/apiproto/src/dataaccess/apiproto/pb"
	"google.golang.org/protobuf/proto"
)

type MessageQueue struct {
	client *pubsub.Client
	context context.Context
	topic *pubsub.Topic
	sub *pubsub.Subscription
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


func (mq* MessageQueue) Publish(clickEvent *pb.ClickEvent) error {
	data, _ := proto.Marshal(clickEvent)
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
		var clickEvent pb.ClickEvent = pb.ClickEvent{}
		err := proto.Unmarshal(msg.Data, &clickEvent)
		if err != nil {
			fmt.Printf("receive: %v\n", err)
		}else {
			fmt.Println(clickEvent.String())
		}
	})
	if err != nil {
			return fmt.Errorf("receive: %v", err)
	}
	return nil
}
