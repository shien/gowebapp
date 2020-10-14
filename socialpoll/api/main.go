package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/stretchr/graceful"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var addr string
	var mongo string

	flag.StringVar(&addr, "addr", ":8080", "EndPoint Address")
	flag.StringVar(&mongo, "mongo", "localhost", "MongoDB Address")

	flag.Parse()

	log.Println("Connect MongoDB", mongo)
	mongoinfo := "mongodb://root:example" + "@" + mongo

	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withVars(withData(mongoinfo, withAPIKey(handlePolls)))))
	log.Println("Start web server: ", addr)
	graceful.Run(addr, 1*time.Second, mux)
	log.Println("Stopping...")
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "Invalid Key")
			return
		}
		fn(w, r)
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}

// データベースのセッション管理
func withData(mongoinfo string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := mongo.NewClient(options.Client().ApplyURI(mongoinfo))
		if err != nil {
			log.Println("Failed to get mongo client.")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = c.Connect(ctx)
		defer c.Disconnect(ctx)
		SetVars(r, "c", c.Database("ballots").Collection("polls"))
		f(w, r)
	}
}

func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}
