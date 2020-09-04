package main

import (
	"fmt"
	"net"
)

var conn net.Conn

func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}

	netc, err := net.DialTimeout(netw, addr, 5+time.Second)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

func main() {
	fmt.Println("vim-go")
}
