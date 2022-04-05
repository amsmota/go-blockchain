package main

import (
	"bytes"
	"encoding/json"
	"goblockchain/common"
	jsonUtils "goblockchain/common"
	"goblockchain/wallet"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
)

const tempDir = "wallet_server/templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) PortStr() string {
	return strconv.FormatUint(uint64(ws.port), 10)
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
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

var TheWallet *wallet.Wallet

func (ws *WalletServer) Wallet(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		res.Header().Add("Content-Type", "application/json")
		TheWallet = wallet.NewWallet()
		m, _ := TheWallet.MarshalJSON()
		io.WriteString(res, string(m[:]))
	default:
		res.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Transaction(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		// res.Header().Add("Content-Type", "application/json")
		decoder := json.NewDecoder(req.Body)
		var t common.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(res, string(jsonUtils.JsonStatus("fail")))
			return
		}
		// if !t.Validate() {
		// 	log.Println("ERROR: missing field(s)")
		// 	io.WriteString(res, string(jsonUtils.JsonStatus("fail")))
		// 	return
		// }

		// publicKey := ecdsa.PublicKeyFromString(*t.SenderPublicKey)
		// privateKey := ecdsa.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		// value, err := strconv.ParseFloat(*t.Value, 32)
		// if err != nil {
		// 	log.Println("ERROR: parse error")
		// 	io.WriteString(res, string(jsonUtils.JsonStatus("fail")))
		// 	return
		// }
		// value32 := float32(value)

		res.Header().Add("Content-Type", "application/json")

		value, _ := strconv.ParseFloat(*t.Value, 32)
		transaction := TheWallet.CreateTransaction(*t.RecipientBlockchainAddress, float32(value))
		TheWallet.SignTransaction(transaction)

		m, _ := json.Marshal(transaction)
		buf := bytes.NewBuffer(m)

		response, err := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		if err != nil {
			log.Printf("ERROR: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			io.WriteString(res, string(jsonUtils.JsonStatus("fail")))
			return
		}
		if response.StatusCode >= 200 && response.StatusCode < 300{
			res.Header().Add("status", string(response.StatusCode))
			io.WriteString(res, string(jsonUtils.JsonStatus("success")))
			return
		}
		res.Header().Add("status", string(response.StatusCode))
		io.WriteString(res, string(jsonUtils.JsonStatus("fail")))
		
	default:
		res.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.Transaction)

	log.Println("WalletServer listening on 0.0.0.0:" + ws.PortStr())
	log.Fatal(http.ListenAndServe("0.0.0.0:"+ws.PortStr(), nil))
}
