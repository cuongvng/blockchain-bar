package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/scanner"
)

type State struct{
	Balances map[Account]uint32
	txMempool []Tx
	dbFile *os.File
}

func GetStateFromDisk() (*State, error){
	cwd, err := os.Getwd()
	if err != nil{
		return *State{}, err
	}

	genesis, err = LoadGenesis(filepath.Join(cwd, "database", "genesis.json"))
	if err != nil{
		return *State{}, err
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
	state := &State{balances, make([]Tx, 0), f}

	for scanner.Scan(){
		if err := scanner.Err(); err != nil{
			return nil, err
		}

		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)

		err := state.Apply(tx)
		if err != nil{
			return nil, err
		}
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

func (s *State) SaveToDisk() error {
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i:=0; i<len(mempool); i++{
		txJson, err := json.Marshal(mempool[i])
		if err != nil{
			return err
		}
		_, err = s.dbFile.Write(append(txJson, '\n'))
		if err != nil{
			return err
		}
		s.txMempool = s.txMempool[1:]
	}
	return nil
}

func (s *State) Close() {
	s.dbFile.Close()
}