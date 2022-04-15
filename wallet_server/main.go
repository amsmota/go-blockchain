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
	gateway := flag.Uint("gateway", 5000, "Gateway Port")
	flag.Parse()
	app := NewWalletServer(uint16(*port), uint16(*gateway))
	app.Run()
}
