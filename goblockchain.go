package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block { // * is a pointer, & is a reference
	b := new(Block) // new() return a pointer
	b.nonce = nonce
	b.previousHash = previousHash
	b.timestamp = time.Now().UnixNano()
	b.transactions = transactions
	return b
	// return &Block {
	// 	timestamp: time.Now().UnixNano(),
	// }
}

func (b *Block) Print() { // create method Print in Block
	fmt.Printf("timestamp           %d\n", b.timestamp)
	fmt.Printf("nonce               %d\n", b.nonce)
	fmt.Printf("previousHash        %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	// fmt.Println(string(m))
	sum := sha256.Sum256([]byte(m))
	return sum
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int 			`json:"nonce"`
		PreviousHash [32]byte 		`json:"previous_hash"`
		Timestamp    int64 			`json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce: b.nonce,
		PreviousHash: b.previousHash,
		Timestamp: b.timestamp,
		Transactions: b.transactions,
	})
}

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 75))
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) LastHash() [32]byte {
	return bc.chain[len(bc.chain)-1].Hash()
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value     		 float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	// t := new(Transaction)
	// t.senderAddress = sender
	// t.recipientAddress = recipient
	// t.value = value
	// return t
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() { 
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_address      %s\n", t.senderAddress)
	fmt.Printf(" recipient_address   %s\n", t.recipientAddress)
	fmt.Printf(" sender_address      %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string    `json:"sender_address"`
		RecipientAddress string    `json:"receiver_address"`
		Value     		 float32   `json:"value"`
	}{
		SenderAddress:     t.senderAddress,
		RecipientAddress:  t.recipientAddress,
		Value:             t.value,
	})
}







func init() {
	log.SetPrefix(("GO: "))
}

func main() {
	bc := NewBlockchain()
	bc.AddTransaction("aaa", "bbb", 10)
	bc.CreateBlock(1, bc.LastHash())
	bc.AddTransaction("ccc", "ddd", 20)
	bc.AddTransaction("ddd", "eee", 30)
	bc.CreateBlock(5, bc.LastHash())
	bc.Print()
}
