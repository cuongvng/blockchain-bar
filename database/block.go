package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}
func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type Header struct{
	parent Hash `json:"parent"`
	time uint64	`json:"time"`
}

type Block struct{
	header Header `json:"hash"`
	txs []Tx `json:"payload"`
}

type BlockFS struct{
	key Hash `json:"hash"`
	value Block `json:"block"`
}

func CreateNewBlock(parent Hash, time uint64, txs []Tx) Block {
	return Block{Header{parent: parent, time: time}, txs}
}

func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil{
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}