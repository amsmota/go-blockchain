package transaction

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
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
		SenderAddress    string  `json:"sender_address"`
		RecipientAddress string  `json:"recipient_address"`
		Value            float32 `json:"value"`
	}{
		SenderAddress:    t.SenderAddress,
		RecipientAddress: t.RecipientAddress,
		Value:            t.Value,
	})
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
