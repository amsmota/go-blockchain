package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	. "goblockchain/common"
)

func NewTransaction(sender string, recipient string, value float32) *BlockTransaction {
	t := new(BlockTransaction)
	t.SenderAddress = sender
	t.RecipientAddress = recipient
	t.Value = value
	return t
}

func VerifyTransaction(senderPublicKey *ecdsa.PublicKey, sig *Signature, t *BlockTransaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], sig.R, sig.S)
}




