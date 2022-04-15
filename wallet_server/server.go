package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goblockchain/common"
	jsonUtils "goblockchain/common"
	"goblockchain/wallet"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
)

const tempDir = "wallet_server/templates"

type WalletServer struct {
	port    uint16
	gateway uint16
	wallet  wallet.Wallet
}

func NewWalletServer(port uint16, gateway uint16) *WalletServer {
	wallet := wallet.NewWallet()
	ws := &WalletServer{port, gateway, *wallet}

	//give them some money...
	t := wallet.CreateTransaction(wallet.BlockchainAddress(), 100)
	t.Tx.SenderAddress = "THE BLOCKCHAIN"
	m, _ := json.Marshal(t)
	buf := bytes.NewBuffer(m)
	http.Post(ws.Gateway()+"/transactions", "application/json", buf)

	return ws
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) PortStr() string {
	return strconv.FormatUint(uint64(ws.port), 10)
}

func (ws *WalletServer) Gateway() string {
	return fmt.Sprintf("http://localhost:%d", ws.gateway)
}

func (ws *WalletServer) Index(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(tempDir, "index.html"))
		if err != nil {
			panic(err)
		} else {
			t.Execute(res, "")
		}
	default:
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Wallet(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		res.Header().Add("Content-Type", "application/json")
		m, _ := ws.wallet.MarshalJSON()
		io.WriteString(res, string(m[:]))
	default:
		res.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Transaction(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t common.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(res, string(jsonUtils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(res, string(jsonUtils.JsonStatus("fail")))
			return
		}

		res.Header().Add("Content-Type", "application/json")
		value, _ := strconv.ParseFloat(*t.Value, 32)
		transaction := ws.wallet.CreateTransaction(*t.RecipientBlockchainAddress, float32(value))
		ws.wallet.SignTransaction(transaction)

		m, _ := json.Marshal(transaction)
		buf := bytes.NewBuffer(m)

		response, err := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		if err != nil {
			log.Printf("ERROR: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			io.WriteString(res, string(jsonUtils.JsonStatus("fail")))
			return
		}
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			res.WriteHeader(response.StatusCode)
			io.WriteString(res, string(jsonUtils.JsonStatus("success")))
			return
		}
		res.WriteHeader(response.StatusCode)
		io.WriteString(res, string(jsonUtils.JsonStatus("fail")))

	default:
		res.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Amount(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		response, err := http.Get(ws.Gateway() + "/amounts?address=" + ws.wallet.BlockchainAddress())
		amount, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("ERROR: No Response from Gateway")
			res.Header().Add("Content-Type", "application/json")
			res.WriteHeader(http.StatusInternalServerError)
			io.WriteString(res, string(jsonUtils.JsonStatus("fail")))
			return
		}
		res.WriteHeader(response.StatusCode)
		io.WriteString(res, string(amount))

	default:
		res.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)                  // GET
	http.HandleFunc("/wallet", ws.Wallet)           // POST
	http.HandleFunc("/transaction", ws.Transaction) // POST
	http.HandleFunc("/amount", ws.Amount)           // GET

	log.Println("WalletServer listening on localhost:" + ws.PortStr())
	log.Fatal(http.ListenAndServe("localhost:"+ws.PortStr(), nil))
}
