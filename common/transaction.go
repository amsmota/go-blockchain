package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"strings"
)

type BlockTransaction struct {
	SenderAddress    string
	RecipientAddress string
	Value            float32
}

func (t *BlockTransaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_address      %s\n", t.SenderAddress)
	fmt.Printf(" recipient_address   %s\n", t.RecipientAddress)
	fmt.Printf(" value               %.1f\n", t.Value)
}

type Transaction struct {
	SenderPublicKey *ecdsa.PublicKey
	Signature       *Signature
	Tx              BlockTransaction
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderPublicKey  *ecdsa.PublicKey `json:"sender_public_key"`
		Signature        *Signature       `json:"signature"`
		SenderAddress    string           `json:"sender_address"`
		RecipientAddress string           `json:"recipient_address"`
		Value            float32          `json:"value"`
	}{
		SenderPublicKey:  t.SenderPublicKey,
		Signature:        t.Signature,
		SenderAddress:    t.Tx.SenderAddress,
		RecipientAddress: t.Tx.RecipientAddress,
		Value:            t.Tx.Value,
	})
}

func (t *Transaction) UnmarshalJSON(mt []byte) error {
	type ttt struct {
		SenderPublicKey  json.RawMessage `json:"sender_public_key"`
		Signature        *Signature      `json:"signature"`
		SenderAddress    string          `json:"sender_address"`
		RecipientAddress string          `json:"recipient_address"`
		Value            float32         `json:"value"`
	}
	tt := new(ttt)
	if err := json.Unmarshal(mt, &tt); err != nil {
		return err
	}

	var spk *ecdsa.PublicKey
	json.Unmarshal(tt.SenderPublicKey, &spk)
	spk.Curve = elliptic.P256()

	t.SenderPublicKey = spk
	t.Signature = tt.Signature
	t.Tx.SenderAddress = tt.SenderAddress
	t.Tx.RecipientAddress = tt.RecipientAddress
	t.Tx.Value = tt.Value

	return nil
}

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	Value                      *string `json:"value"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPrivateKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil {
		return false
	}
	return true
}
