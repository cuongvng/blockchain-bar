package database

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Hash [32]byte
type State struct{
	Balances map[Account]uint32
	txMempool []Tx
	dbFile *os.File
	hash Hash
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

	f, err := os.OpenFile(filepath.Join(cwd, "database", "tx.db"), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f, Hash{}}

	for scanner.Scan(){
		if err := scanner.Err(); err != nil{
			return nil, err
		}

		var tx Tx
		err = json.Unmarshal(scanner.Bytes(), &tx)
		if err != nil {
			return nil, err
		}

		err := state.Apply(tx)
		if err != nil{
			return nil, err
		}
	}

	err = state.takeSnapshot()
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *State) Add(tx Tx) error {
	err := s.Apply(tx)
	if err != nil{
		return err
	}

	s.txMempool = append(s.txMempool, tx)
	return nil
}

func (s *State) Apply(tx Tx) error {
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

func (s *State) takeSnapshot() error {
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return err
	}

	txsData, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		return err
	}

	s.hash = sha256.Sum256(txsData)
	return nil
}

func (s *State) GetLastestHash() Hash {
	return s.hash
}

func (s *State) SaveToDisk() (Hash, error) {
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i:=0; i<len(mempool); i++{
		txJson, err := json.Marshal(mempool[i])
		if err != nil{
			return Hash{}, err
		}
		_, err = s.dbFile.Write(append(txJson, '\n'))
		if err != nil{
			return Hash{}, err
		}

		err = s.takeSnapshot()
		if err != nil{
			return Hash{}, err
		}
		fmt.Printf("New DB hash: %x\n", s.hash)

		s.txMempool = s.txMempool[1:]
	}
	return s.hash, nil
}

func (s *State) Close() error{
	return s.dbFile.Close()
}