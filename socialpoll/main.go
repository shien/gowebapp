package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client
var ctx context

func dialdb() error {
	db, err := NewClient(options.Client().ApplyURI("mongodb://foo:bar@localhost:27017"))
	log.Println("Dial MongoDB...")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = db.Connect(ctx)
	return err
}

func closedb() {
	db.Disconnect(ctx)
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfing())
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote))
		}
		log.Println("Publisher 停止中です")
		pub.Stop()
		log.Println("Publisher 停止中しました")
		stopchan <- struct{}{}
	}()
	return stopchan
}

func main() {}
