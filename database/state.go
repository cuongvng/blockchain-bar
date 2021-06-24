package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Hash [32]byte
type State struct{
	Balances map[Account]uint32
	txMempool []Tx
	dbFile *os.File
	lastestHash Hash
}

func GetStateFromDisk() (*State, error){
	cwd, err := os.Getwd()
	if err != nil{
		return &State{}, err
	}

	genesis, err := LoadGenesis(filepath.Join(cwd, "database", "genesis.json"))
	if err != nil{
		return &State{}, err
	}

	balances := make(map[Account]uint32)
	for account, balance := range genesis.Balances{
		balances[account] = balance
	}

	f, err := os.OpenFile(filepath.Join(cwd, "database", "block.db"), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f, Hash{}}

	for scanner.Scan(){
		if err := scanner.Err(); err != nil{
			return nil, err
		}

		blockFsJson := scanner.Bytes()
		var blockFs BlockFS
		err = json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, err
		}

		err := state.ApplyBlock(blockFs.value)
		if err != nil{
			return nil, err
		}
		state.lastestHash = blockFs.key
	}
	return state, nil
}

func (s *State) AddBlock(block Block) error {
	for _, tx := range block.txs{
		err := s.AddTx(tx)
		if err != nil{
			return err
		}
	}
	return nil
}
func (s *State) AddTx(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)
	return nil
}

func (s *State) ApplyBlock(block Block) error {
	for _, tx := range block.txs{
		err := s.apply(tx)
		if err != nil{
			return err
		}
	}
	return nil
}
func (s *State) apply(tx Tx) error {
	if tx.IsReward(){
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From] < tx.Value{
		return fmt.Errorf("Insufficient balance to send")
	}

	s.Balances[tx.To] += tx.Value
	s.Balances[tx.From] -= tx.Value

	return nil
}

func (s *State) GetLastestHash() Hash {
	return s.lastestHash
}

func (s *State) Persist() (Hash, error) {
	block := CreateNewBlock(s.lastestHash, uint64(time.Now().Unix()), s.txMempool)
	blockHash, err := block.Hash()
	if err != nil{
		return Hash{}, err
	}
	blockFs := BlockFS{blockHash, block}
	blockFsJson, err := json.Marshal(blockFs)
	if err != nil{
		return Hash{}, err
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	_, err = s.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil{
		return Hash{}, err
	}

	s.lastestHash = blockHash
	s.txMempool = []Tx{}

	return s.lastestHash, nil
}

func (s *State) Close() error{
	return s.dbFile.Close()
}