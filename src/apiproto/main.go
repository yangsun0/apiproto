package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yangsun0/apiproto/src/dataaccess"
	"github.com/yangsun0/apiproto/src/dataaccess/apiproto/pb"
)

func main() {
	var app_key string = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	fmt.Println("Hello world")
	fmt.Println("api key: ", app_key)
	ctx := context.Background()
	dataClient := dataaccess.NewDataClient()
	err := dataClient.Connect(ctx)
	defer dataClient.Close()
	if err != nil {
		log.Fatal("connect to data failed")
	} 

	cacheClient := dataaccess.NewCacheClient()
	cacheClient.Connect(ctx)
	defer cacheClient.Close()

	var mq dataaccess.MessageQueue = dataaccess.MessageQueue{}
	defer mq.Close()
	err = mq.Connect(ctx)
	if err != nil {
		log.Fatalf("connect to messsage queue failed, %v", err)
	}


	mode := "pull"
	switch mode {
	case "addUser":
		addUser(dataClient)
	case "addBookmark":
		addBookmark(dataClient)
	case "getBookmarks":
		getBookmarks(dataClient)
	case "setCache":
		setCache(cacheClient)
	case "publish":
		publish(&mq)
	case "pull":
		pull(&mq)
	}
}


func addUser(dataClient *dataaccess.DataClient) {
	user := dataaccess.User{
		Name: "user3", 
		Email: "user3@gmail.com"}
	dataClient.Add("users", user)
}

func addBookmark(dataClient *dataaccess.DataClient) {
	userId :=  "YCEKK4FI1n3JfmEEP81u"
	bookmark := dataaccess.Bookmark{
		UserId: userId,
		Name: "google", 
		Url: "https://www.google.com"}
	dataClient.Add("bookmarks", bookmark)
}

func getBookmarks(dataClient *dataaccess.DataClient) {
	userId :=  "YCEKK4FI1n3JfmEEP81u"
	dataClient.ReadAllBookmarks(userId)
}

func setCache(cacheClient *dataaccess.CacheClient) {
	cacheClient.Set("hello", "world")
	value := cacheClient.Get("hello")
	fmt.Printf("key: %v, value: %v\n", "hello", value)
}

func publish(mq *dataaccess.MessageQueue) {
	timestamp := time.Now().Unix()
	clickEvent := pb.ClickEvent{Target: "btn1", Timestamp: timestamp}
	err := mq.Publish(&clickEvent)
	if err != nil {
		fmt.Printf("publish error%v\n", err)
	}
}

func pull(mq *dataaccess.MessageQueue) {
	err := mq.Pull()
	if err != nil {
		fmt.Printf("publish error%v\n", err)
	}
}