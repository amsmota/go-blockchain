package common

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	SenderPrivateKey *ecdsa.PrivateKey
	SenderPublicKey  *ecdsa.PublicKey
	SenderAddress    string
	RecipientAddress string
	Value            float32
	Signature        *Signature
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_address      %s\n", t.SenderAddress)
	fmt.Printf(" recipient_address   %s\n", t.RecipientAddress)
	fmt.Printf(" value               %.1f\n", t.Value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderPublicKey  *ecdsa.PublicKey `json:"sender_public_key"`
		SenderAddress    string           `json:"sender_address"`
		RecipientAddress string           `json:"recipient_address"`
		Value            float32          `json:"value"`
		Signature        *Signature           `json:"signature"`
	}{
		SenderPublicKey:  t.SenderPublicKey,
		SenderAddress:    t.SenderAddress,
		RecipientAddress: t.RecipientAddress,
		Value:            t.Value,
		Signature:        t.Signature,
	})
}

// func (t *Transaction) UnmarshalJSON(mt []byte) error {
// 	type T2 struct {
// 		SenderAddress    string  `json:"sender_address"`
// 		RecipientAddress string  `json:"recipient_address"`
// 		Value            float32 `json:"value"`
// 	}
// 	var tt T2
// 	if err := json.Unmarshal(mt, &tt); err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}
// 	t.RecipientAddress = tt.RecipientAddress
// 	t.SenderAddress = tt.SenderAddress
// 	t.Value = tt.Value

// 	return nil
// }

type TransactionRequest struct {
	SenderPublicKey  *string `json:"sender_public_key"`
	SenderAddress    *string `json:"sender_address"`
	RecipientAddress *string `json:"recipient_address"`
	Value            *float32 `json:"value"`
	Signature        *string
}

// func (t *TransactionRequest) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		SenderAddress    *string `json:"sender_address"`
// 		RecipientAddress *string `json:"recipient_address"`
// 		Value            *string `json:"value"`
// 	}{
// 		SenderAddress:    t.SenderBlockchainAddress,
// 		RecipientAddress: t.RecipientBlockchainAddress,
// 		Value:            t.Value,
// 	})
// }

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderAddress == nil ||
		tr.RecipientAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil {
		return false
	}
	return true
}
