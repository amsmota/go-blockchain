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

	// WalletA to WalletB transaction
	t := walletA.CreateTransaction(walletB.BlockchainAddress(), 1.0)
	walletA.SignTransaction(t)

	// Blockchain
	bc := blockchain.NewBlockchain(walletM.BlockchainAddress())

	// Bad guy does this, gives ERROR: Verifiy Transaction
	// t.Value = 100000
	// walletB.SignTransaction(t)

	added := bc.AddTransaction(t)
	fmt.Println("Added: ", added)

	bc.Mining()
	bc.Print()

	fmt.Printf("A %.1f\n", bc.CalculateTotalAmount(walletA.BlockchainAddress()))
	fmt.Printf("B %.1f\n", bc.CalculateTotalAmount(walletB.BlockchainAddress()))
	fmt.Printf("M %.1f\n", bc.CalculateTotalAmount(walletM.BlockchainAddress()))

}
