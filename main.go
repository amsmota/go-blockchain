package main

import (
	"fmt"
	"goblockchain/wallet"
	"log"
)

func init() {
	log.SetPrefix(("GO BLOCKCHAIN: "))
}

func main() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyString())
	fmt.Println(w.PublicKeyString())
	fmt.Println(w.BlockchainAddress())
}