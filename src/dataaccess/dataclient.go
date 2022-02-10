package dataaccess

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type DataClient struct {
	client *firestore.Client
	context context.Context
}

func NewDataClient() *DataClient {
	return &DataClient{}
}

func (dc *DataClient) Connect(ctx context.Context) error {
	var projectId string = os.Getenv("GCP_PROJECT_ID")

	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
			return err
	}

	dc.client = client
	dc.context = ctx
	log.Println("data client connect successfully")
	return nil
}

func (dc *DataClient) Close() {
	if dc.client != nil {
		dc.client.Close()
		log.Println("close data client")
	}
}

func (dc *DataClient) Add(coll string, data interface{}) {
	dc.client.Collection(coll).Add(dc.context, data)
	log.Println("add data")
}

func (dc *DataClient) Read(coll string) {
	iter := dc.client.Collection(coll).Documents(dc.context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalf("Failed to read collection %v, error: %v", coll, err)
		}
		var user User
		doc.DataTo(&user)
		user.Id = doc.Ref.ID
		fmt.Println(user)
	}
	log.Println("read data")
} 

func (dc *DataClient) ReadAllBookmarks(userId string) {
	iter := dc.client.Collection("bookmarks").Where("userId", "==", userId).Documents(dc.context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalf("Failed to read bookmarks, error: %v", err)
		}
		var bookmark Bookmark
		doc.DataTo(&bookmark)
		bookmark.Id = doc.Ref.ID
		fmt.Println(bookmark)
	}
	log.Println("read bookmarks")
} 

type User struct {
	Id string		`firestore:"-"`
	Name string		`firestore:"name"`
	Email string	`firestore:"email"`
}


type Bookmark struct {
	Id string		`firestore:"-"`
	UserId string	`firestore:"userId"`
	Name string		`firestore:"name"`
	Url string		`firestore:"url"`
}