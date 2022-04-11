package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	. "goblockchain/common"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 20
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

type Blockchain struct {
	transactionPool   []*BlockTransaction
	chain             []*Block
	blockchainAddress string
	port              uint16
	mux               sync.Mutex
}

func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.port = port
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*BlockTransaction{}
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 75))
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"blochchain"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) LastHash() [32]byte {
	return bc.chain[len(bc.chain)-1].Hash()
}

func (bc *Blockchain) AddTransaction(t *BlockTransaction) bool {
	sender := t.SenderAddress
	// senderPublicKey := t.SenderPublicKey

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	// VERIFY HERE OR ON THE BC_SERVER?
	if true { //bc.VerifyTransaction(senderPublicKey, t.Signature, t) {
		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("ERROR: Not Enough Gas")
		// 	return false
		// }
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}
	log.Println("ERROR: Verifiy Transaction")
	return false
}

func VerifyTransaction(senderPublicKey *ecdsa.PublicKey, sig *Signature, t *BlockTransaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], sig.R, sig.S)
}

func (bc *Blockchain) CopyTransactionPool() []*BlockTransaction {
	transactions := make([]*BlockTransaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.SenderAddress, t.RecipientAddress, t.Value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*BlockTransaction, dificulty int) bool {
	zeros := strings.Repeat("0", dificulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHash := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHash[:dificulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastHash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	if len(bc.transactionPool) == 0{
		return false
	}

	t := NewTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	bc.AddTransaction(t)
	nonce := bc.ProofOfWork()
	bc.CreateBlock(nonce, bc.LastHash())
	log.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) StartMining() {
	mined := bc.Mining()
	log.Printf("Mining: %t", mined)
	_ = time.AfterFunc(time.Second * MINING_TIMER_SEC, bc.StartMining)
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.Value
			if blockchainAddress == t.RecipientAddress {
				totalAmount += value
			}
			if blockchainAddress == t.SenderAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func NewTransaction(sender string, recipient string, value float32) *BlockTransaction {
	t := new(BlockTransaction)
	t.SenderAddress = sender
	t.RecipientAddress = recipient
	t.Value = value
	return t
}
