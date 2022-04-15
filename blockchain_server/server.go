package main

import (
	"encoding/json"
	"fmt"
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
		tp := bc.TransactionPool()
		m, _ := json.Marshal(tp)
		io.WriteString(res, string(m[:]))
	}
}

func (bcs *BlockchainServer) Transactions(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		res.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		tp := bc.TransactionPool()
		m, _ := json.Marshal(tp)
		io.WriteString(res, string(m[:]))

	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t common.Transaction
		decoder.Decode(&t)
		bc := bcs.GetBlockchain()
		isCreated := bc.CreateTransaction(&t)
		res.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isCreated {
			res.WriteHeader(http.StatusBadRequest)
			m = common.JsonStatus("fail")
		} else {
			res.WriteHeader(http.StatusCreated)
			m = common.JsonStatus("success")
		}
		io.WriteString(res, string(m))

	case http.MethodPut:
		decoder := json.NewDecoder(req.Body)
		var t common.Transaction
		decoder.Decode(&t)
		bc := bcs.GetBlockchain()
		isUpdated := bc.AddTransaction(&t)
		res.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isUpdated {
			res.WriteHeader(http.StatusBadRequest)
			m = common.JsonStatus("fail")
		} else {
			m = common.JsonStatus("success")
		}
		io.WriteString(res, string(m))

	case http.MethodDelete:
		bc := bcs.GetBlockchain()
		bc.ClearTransactionPool()
		io.WriteString(res, string(common.JsonStatus("success")))
	}

}

func (bcs *BlockchainServer) Amounts(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")
		if address == "" {
			// say bye bye
		}

		res.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		amount := bc.CalculateTotalAmount(address)

		io.WriteString(res, fmt.Sprintf("%f", amount))

	default:
		res.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")

	}

}

func (bcs *BlockchainServer) Consensus(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		log.Printf("IGNORING ResolveConflicts")
		// bc := bcs.GetBlockchain()
		// replaced := bc.ResolveConflicts()

		// res.Header().Add("Content-Type", "application/json")
		// if replaced {
		// 	io.WriteString(res, string(common.JsonStatus("success")))
		// } else {
		// 	io.WriteString(res, string(common.JsonStatus("fail")))
		// }
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		res.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/blockchain", bcs.GetChain)       // GET
	http.HandleFunc("/transactions", bcs.Transactions) // GET POST PUT DELETE
	http.HandleFunc("/amounts", bcs.Amounts)           // GET
	http.HandleFunc("/consensus", bcs.Consensus)       // PUT

	log.Println("BlockchainServer listening on localhost:" + bcs.PortStr())
	bcs.GetBlockchain().Run()
	log.Fatal(http.ListenAndServe("localhost:"+bcs.PortStr(), nil))
}
