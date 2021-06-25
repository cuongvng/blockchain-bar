package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}
func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type Header struct{
	Parent Hash `json:"parent"`
	Time uint64	`json:"time"`
}

type Block struct{
	Header Header `json:"header"`
	TXs []Tx `json:"payload"`
}

type BlockFS struct{
	Key Hash `json:"hash"`
	Value Block `json:"block"`
}

func CreateNewBlock(parent Hash, time uint64, txs []Tx) Block {
	return Block{Header{parent, time}, txs}
}

func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil{
		return Hash{}, err
	}
	fmt.Printf("%x", sha256.Sum256(blockJson))
	return sha256.Sum256(blockJson), nil
}