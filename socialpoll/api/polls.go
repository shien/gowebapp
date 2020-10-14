package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

type poll struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title   string             `bson:"title" json:"title"`
	Options []string           `bson:"options" json:"options"`
	Results map[string]int     `bson:"results,omitempty" json:"results,omitempty"`
}

func handlePolls(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlePollsGet(w, r)
		return
	case "POST":
		handlePollsPost(w, r)
		return
	case "DELETE":
		handlePollsDelete(w, r)
		return
	case "OPTIONS":
		w.Header().Add("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
	respondHTTPErr(w, r, http.StatusNotFound)
}

func handlePollsGet(w http.ResponseWriter, r *http.Request) {
	coll := GetVars(r, "c").(*mongo.Collection)

	p := NewPath(r.URL.Path)

	var err error
	var result []poll
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if p.HasID() {
		resultOne, err := getIDResult(ctx, coll, p.ID)
		if err != nil {
			respondErr(w, r, http.StatusInternalServerError, err)
			return
		}
		result = append(result, resultOne)
	} else {
		err = getAllResult(ctx, coll, &result)
		if err != nil {
			respondErr(w, r, http.StatusInternalServerError, err)
			return
		}
	}

	respond(w, r, http.StatusOK, &result)
}

func getAllResult(ctx context.Context, coll *mongo.Collection, resultAll *[]poll) error {
	var cur *mongo.Cursor

	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Println("cur all err")
		return err

	}

	defer cur.Close(ctx)

	if err := cur.All(ctx, resultAll); err != nil {
		log.Println("cur all err")
		return err
	}

	return nil
}

func getIDResult(ctx context.Context, coll *mongo.Collection, id string) (poll, error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed: Object ID From Hex", err)
		return poll{}, err
	}

	var result poll

	err = coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&result)

	if err != nil {
		log.Println("Failed: Could not find ID", err)
		return poll{}, err
	}

	return result, nil
}

func handlePollsPost(w http.ResponseWriter, r *http.Request) {
	coll := GetVars(r, "c").(*mongo.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var p = poll{
		ID: primitive.NewObjectID(),
	}

	if err := decodeBody(r, &p); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "リストから調査項目を読み込めません", err)
		return
	}

	res, err := coll.InsertOne(ctx, p)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "DB への書き込み失敗", err)
		return
	}
	log.Println(res)

	w.Header().Set("Location", "polls/"+p.ID.Hex())
	respond(w, r, http.StatusCreated, nil)
}

func handlePollsDelete(w http.ResponseWriter, r *http.Request) {
	coll := GetVars(r, "c").(*mongo.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	p := NewPath(r.URL.Path)

	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "すべての調査項目を削除することはできません")
		return
	}

	oid, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		log.Println("Failed: Object ID From Hex", err)
		return
	}

	res, err := coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "削除失敗", err)
	}

	respond(w, r, http.StatusOK, res)
}
