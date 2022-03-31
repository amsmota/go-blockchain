package main

import (
	"fmt"
	"goblockchain/blockchain"
	"goblockchain/wallet"
	"log"
)

func init() {
	log.SetPrefix(("GO BLOCKCHAIN: "))
}

func main() {
	walletM := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	// Wallet
	t := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(),
		walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)
	// Blockchain
	bc := blockchain.NewBlockchain(walletM.BlockchainAddress())
	added := bc.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0,
		walletA.PublicKey(), t.GenerateSignature())
	fmt.Println("Added: ", added)

	bc.Mining()
	bc.Print()

	fmt.Printf("A %.1f\n", bc.CalculateTotalAmount(walletA.BlockchainAddress()))
	fmt.Printf("B %.1f\n", bc.CalculateTotalAmount(walletB.BlockchainAddress()))
	fmt.Printf("M %.1f\n", bc.CalculateTotalAmount(walletM.BlockchainAddress()))

}