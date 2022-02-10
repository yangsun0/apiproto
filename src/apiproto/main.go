package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/yangsun0/apiproto/src/dataaccess"
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
	defer cacheClient.Close()
	cacheClient.Connect(ctx)
	cacheClient.Set("hello", "world")
	value := cacheClient.Get("hello")
	fmt.Printf("key: %v, value: %v\n", "hello", value)
}

// func addUser(dataClient *dataaccess.DataClient) {
// 	user := dataaccess.User{
// 		Name: "user3", 
// 		Email: "user3@gmail.com"}
// 	dataClient.Add("users", user)
// }

// func addBookmark(dataClient *dataaccess.DataClient) {
// 	userId :=  "YCEKK4FI1n3JfmEEP81u"
// 	bookmark := dataaccess.Bookmark{
// 		UserId: userId,
// 		Name: "google", 
// 		Url: "https://www.google.com"}
// 	dataClient.Add("bookmarks", bookmark)
// }

// func getBookmarks(dataClient *dataaccess.DataClient) {

// 	userId :=  "YCEKK4FI1n3JfmEEP81u"
// 	dataClient.ReadAllBookmarks(userId)
// }
