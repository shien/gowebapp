package main

import (
	"backup"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matryer/filedb"
)

type path struct {
	Path string
	Hash string
}

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			log.Fatalln(fatalErr)
		}
	}()
	var (
		interval = flag.Int("interval", 10, "check interval / seconds")
		archive  = flag.String("archive", "archive", "where does arcive save")
		dbpath   = flag.String("db", "./db", "Path to filedb database")
	)
	flag.Parse()

	m := &backup.Monitor{
		Destination: *archive,
		Archiver:    backup.ZIP,
		Paths:       make(map[string]string),
	}

	db, err := filedb.Dial(*dbpath)
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()
	col, err := db.C("paths")
	if err != nil {
		fatalErr = err
		return
	}
	var path path
	col.ForEach(func(_ int, data []byte) bool {
		if err := json.Unmarshal(data, &path); err != nil {
			fatalErr = err
			return true
		}
		m.Paths[path.Path] = path.Hash
		return false
	})
	if fatalErr != nil {
		return
	}
	if len(m.Paths) < 1 {
		fatalErr = errors.New("Path is not found, please add path using backup tool")
		return
	}
	check(m, col)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
Loop:
	for {
		select {
		case <-time.After(time.Duration(*interval) * time.Second):
			check(m, col)
		case <-signalChan:
			// Finish
			fmt.Println()
			log.Printf("Stopping ...")
			break Loop
		}
	}
}

func check(m *backup.Monitor, col *filedb.C) {
	log.Println("Checking ...")
	counter, err := m.Now()
	if err != nil {
		log.Println("Failed to backup", err)
	}
	if counter > 0 {
		log.Printf("Archive %d directories\n", counter)
		var path path
		col.SelectEach(func(_ int, data []byte) (bool, []byte, bool) {
			if err := json.Unmarshal(data, &path); err != nil {
				log.Println("Failed to read JSON data "+"Go to next data", err)
				return true, data, false
			}
			path.Hash, _ = m.Paths[path.Path]
			newData, err := json.Marshal(&path)
			if err != nil {
				log.Println("Failed to write JSON data "+"Go to next data", err)
				return true, data, false
			}
			return true, newData, false
		})
	} else {
		log.Println("No any updates")
	}
}
