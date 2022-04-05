package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix(("GO-BLOCKCHAIN: "))
}

func main() {
	port := flag.Uint("port", 5000, "TCP port for BlockchainServer")
	flag.Parse()
	app := NewBlockchainServer(uint16(*port))
	app.Run()
}
