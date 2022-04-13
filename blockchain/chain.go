package blockchain

import (
	"encoding/json"
	"fmt"
	. "goblockchain/common"
	"log"
	"strings"
	"sync"
	"time"
)

type Blockchain struct {
	transactionPool   []*BlockTransaction
	chain             []*Block
	blockchainAddress string
	port              uint16
	muxMining         sync.Mutex

	neighbors    []string
	muxNeighbors sync.Mutex
}

func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.port = port
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) Chain() []*Block {
	return bc.chain
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*BlockTransaction{}

	bc.NodeSyncNewBlock()

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
		Blocks []*Block `json:"blockchain"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) UnmarshalJSON(data []byte) error {
	v := &struct {
		Blocks *[]*Block `json:"blockchain"`
	}{
		Blocks: &bc.chain,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) LastHash() [32]byte {
	return bc.chain[len(bc.chain)-1].Hash()
}

func (bc *Blockchain) CreateTransaction(t *Transaction) bool {
	isTransacted := bc.AddTransaction(t)
	if isTransacted {
		bc.NodeSyncTransaction(t)
	}
	return isTransacted
}

func (bc *Blockchain) AddTransaction(t *Transaction) bool {

	if t.Tx.SenderAddress == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, &t.Tx)
		return true
	}

	if !VerifyTransaction(t.SenderPublicKey, t.Signature, &t.Tx) {
		log.Println("ERROR: Verifiy Transaction")
		return false
	}
	// if bc.CalculateTotalAmount(t.Tx.SenderAddress) < t.Tx.Value {
	// 	log.Println("ERROR: Not Enough Gas")
	// 	return false
	// }

	bc.transactionPool = append(bc.transactionPool, &t.Tx)
	return true
}

func (bc *Blockchain) CopyTransactionPool() []*BlockTransaction {
	transactions := make([]*BlockTransaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.SenderAddress, t.RecipientAddress, t.Value))
	}
	return transactions
}

func (bc *Blockchain) TransactionPool() []*BlockTransaction {
	return bc.transactionPool
}

func (bc *Blockchain) ClearTransactionPool() {
	bc.transactionPool = bc.transactionPool[:0]
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
	bc.muxMining.Lock()
	defer bc.muxMining.Unlock()

	if len(bc.transactionPool) == 0 {
		return false
	}

	t := NewTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	tt := Transaction{Tx: *t}
	bc.AddTransaction(&tt)
	nonce := bc.ProofOfWork()
	bc.CreateBlock(nonce, bc.LastHash())
	log.Println("action=mining, status=success")

	bc.NodeSyncConsensus()

	return true
}

func (bc *Blockchain) StartMining() {
	mined := bc.Mining()
	log.Printf("Mining: %t", mined)
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, bc.StartMining)
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

func (bc *Blockchain) ValidChain(chain []*Block) bool {
	preBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		b := chain[currentIndex]
		if b.previousHash != preBlock.Hash() {
			return false
		}

		if !bc.ValidProof(b.Nonce(), b.PreviousHash(), b.Transactions(), MINING_DIFFICULTY) {
			return false
		}

		preBlock = b
		currentIndex += 1
	}
	return true
}

func (bc *Blockchain) ResolveConflicts() bool {
	var longestChain []*Block = nil
	maxLength := len(bc.chain)

	for _, n := range bc.neighbors {
		chain := bc.NodeSyncChain(n)
		if len(chain) > maxLength && bc.ValidChain(chain) {
			maxLength = len(chain)
			longestChain = chain
		}
	}

	if longestChain != nil {
		bc.chain = longestChain
		log.Printf("Resolve confilicts replaced")
		return true
	}
	log.Printf("Resolve conflicts not replaced")
	return false
}

