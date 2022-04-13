package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goblockchain/common"
	nodes "goblockchain/common"
	"log"
	"net/http"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 20

	BLOCKCHAIN_PORT_RANGE_START      = 5000
	BLOCKCHAIN_PORT_RANGE_END        = 5003
	NEIGHBOR_IP_RANGE_START          = 0
	NEIGHBOR_IP_RANGE_END            = 1
	BLOCKCHIN_NEIGHBOR_SYNC_TIME_SEC = 30
)

func (bc *Blockchain) NodeSyncNewBlock() {
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
}

func (bc *Blockchain) NodeSyncTransaction(t *common.Transaction) {
	for _, n := range bc.neighbors {
		m, _ := json.Marshal(&t)
		buf := bytes.NewBuffer(m)
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", endpoint, buf)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
}

func (bc *Blockchain) NodeSyncChain(n string) []*Block {
	endpoint := fmt.Sprintf("http://%s/blockchain", n)
	resp, _ := http.Get(endpoint)
	if resp.StatusCode == 200 {
		var bcResp Blockchain
		decoder := json.NewDecoder(resp.Body)
		_ = decoder.Decode(&bcResp)
	
		return bcResp.Chain()
	}
	return nil
} 

func (bc *Blockchain) NodeSyncConsensus(){
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/consensus", n)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
}


func (bc *Blockchain) SetNeighbors() {
	bc.neighbors = nodes.FindNeighbors(
		nodes.GetHost(), bc.port,
		NEIGHBOR_IP_RANGE_START, NEIGHBOR_IP_RANGE_END,
		BLOCKCHAIN_PORT_RANGE_START, BLOCKCHAIN_PORT_RANGE_END)
	log.Printf("%v", bc.neighbors)
}

func (bc *Blockchain) SyncNeighbors() {
	bc.muxNeighbors.Lock()
	defer bc.muxNeighbors.Unlock()
	bc.SetNeighbors()
}

func (bc *Blockchain) StartSyncNeighbors() {
	bc.SyncNeighbors()
	_ = time.AfterFunc(time.Second*BLOCKCHIN_NEIGHBOR_SYNC_TIME_SEC, bc.StartSyncNeighbors)
}

func (bc *Blockchain) Run() {
	bc.StartSyncNeighbors()
	bc.ResolveConflicts()
	bc.StartMining()
}
