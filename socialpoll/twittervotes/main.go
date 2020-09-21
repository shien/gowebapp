package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	nsq "github.com/bitly/go-nsq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client
var ctx context.Context

func dialdb() error {
	var err error
	db, err = mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@172.20.0.12:27017"))
	log.Println("Dial MongoDB...")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = db.Connect(ctx)
	return err
}

func closedb() error {
	return db.Disconnect(ctx)
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	coll := db.Database("ballots").Collection("polls")
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		log.Println("collection の取得に失敗しました", err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var opts bson.M
		if err = cur.Decode(&opts); err != nil {
			log.Fatal(err)
		}
		if opt, ok := opts["options"].(primitive.A); ok {
			s := []interface{}(opt)
			for _, w := range s {
				if w, ok := w.(string); ok {
					options = append(options, string(w))
				}
			}
		}
	}

	return options, err
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("172.20.0.11:4150", nsq.NewConfig())
	go func() {
		for vote := range votes {
			log.Println("Publish します", []byte(vote))
			pub.Publish("votes", []byte(vote))
		}
		log.Println("Publisher 停止中です")
		pub.Stop()
		log.Println("Publisher 停止しました")
		stopchan <- struct{}{}
	}()
	return stopchan
}

func main() {
	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)

	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("停止します")
		stopChan <- struct{}{}
		closeConn()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := dialdb(); err != nil {
		log.Fatalln("MongoDB への接続に失敗しました", err)
	}
	defer closedb()

	votes := make(chan string)
	publisherStoppedChan := publishVotes(votes)
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				break
			}
			stoplock.Unlock()
		}
	}()
	<-twitterStoppedChan
	close(votes)
	<-publisherStoppedChan
}
