package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherWord = "*"

//var transforms = []string{
//	otherWord + "web",
//	otherWord + "srv",
//	otherWord + "book",
//	otherWord,
//	otherWord + "app",
//	otherWord + "site",
//	otherWord + "time",
//	"get" + otherWord,
//	"go" + otherWord,
//	"lets" + otherWord,
//}

func readFile(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var l []string

	s := bufio.NewScanner(f)
	for s.Scan() {
		l = append(l, s.Text())
	}
	return l
}

func main() {
	transforms := readFile("abc")
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
	}
}
