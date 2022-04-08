package main

import (
	"encoding/json"
	. "goblockchain/blockchain"
	"goblockchain/common"
	"goblockchain/wallet"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*Blockchain = make(map[string]*Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) PortStr() string {
	return strconv.FormatUint(uint64(bcs.port), 10)
}

func (bcs *BlockchainServer) GetBlockchain() *Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = NewBlockchain(minersWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc
		log.Printf("PLEASE REMOVE THOSE LINES BELOW")
		log.Printf("private_key %v", minersWallet.PrivateKeyString())
		log.Printf("public_key %v", minersWallet.PublicKeyString())
		log.Printf("blockchcain_address %v", minersWallet.BlockchainAddress())
		log.Printf("PLEASE REMOVE THOSE LINES ABOVE")
	}
	return bc
}

func (bcs *BlockchainServer) GetChain(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		res.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(res, string(m[:]))
	}
}

func (bcs *BlockchainServer) TransactionPool(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		res.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		tp := bc.CopyTransactionPool()
		m, _ := json.Marshal(tp)
		io.WriteString(res, string(m[:]))
	}
}

func (bcs *BlockchainServer) Transactions(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t common.Transaction
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(res, string(common.JsonStatus("fail")))
			return
		}

		nt := NewTransaction(t.Tx.SenderAddress, t.Tx.RecipientAddress, t.Tx.Value)
		if !VerifyTransaction(t.SenderPublicKey, t.Signature, nt) {
			log.Println("ERROR: Verifiy Transaction")
			io.WriteString(res, string(common.JsonStatus("fail")))
			return
		}
		
		bc := bcs.GetBlockchain()
		bc.AddTransaction(nt)
		m, _ := bc.MarshalJSON()
		io.WriteString(res, string(m[:])) // SEND SOMETHING ELSE
	}
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/blockchain", bcs.GetChain)
	http.HandleFunc("/transactionPool", bcs.TransactionPool)
	http.HandleFunc("/transactions", bcs.Transactions)

	log.Println("BlockchainServer listening on localhost:" + bcs.PortStr())
	log.Fatal(http.ListenAndServe("localhost:"+bcs.PortStr(), nil))
}
