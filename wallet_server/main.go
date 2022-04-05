package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix(("GO-WALLET: "))
}

func main() {
	port := flag.Uint("port", 8080, "TCP port for Wallet")
	flag.Parse()
	gateway := flag.String("gateway", "http://127.0.0.1:5000", "Gateway URL")

	app := NewWalletServer(uint16(*port), *gateway)
	app.Run()
}
