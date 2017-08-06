package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"

	"github.com/fysac/golem/mc"
)

const (
	ListenPort = "25565"
)

var Status mc.ServerStatus

func initConfig() {
	bs, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bs, &Status)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	listen, err := net.Listen("tcp", ":"+ListenPort)
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	initConfig()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleConn(conn)
	}
}
