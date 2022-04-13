package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	. "goblockchain/common"
	"time"
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*BlockTransaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*BlockTransaction) *Block { // * is a pointer, & is a reference
	b := new(Block) // new() return a pointer
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func (b *Block) PreviousHash() [32]byte {
	return b.previousHash
}

func (b *Block) Nonce() int {
	return b.nonce
}

func (b *Block) Transactions() []*BlockTransaction {
	return b.transactions
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
	sum := sha256.Sum256([]byte(m))
	return sum
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64               `json:"timestamp"`
		Nonce        int                 `json:"nonce"`
		PreviousHash string              `json:"previous_hash"`
		Transactions []*BlockTransaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

func (b *Block) UnmarshalJSON(data []byte) error {
	var previousHash string
	v := &struct {
		Timestamp        *int64               `json:"timestamp"`
		Nonce            *int                 `json:"nonce"`
		PreviousHash     *string              `json:"previous_hash"`
		BlockTransaction *[]*BlockTransaction `json:"transactions"`
	}{
		Timestamp:        &b.timestamp,
		Nonce:            &b.nonce,
		PreviousHash:     &previousHash,
		BlockTransaction: &b.transactions,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	ph, _ := hex.DecodeString(*v.PreviousHash)
	copy(b.previousHash[:], ph[:32])
	return nil
}
