package main

import (
	"fmt"
	"github.com/cuongvng/blockchain-bar/database"
	"os"
	"time"
)

func main() {
	state, err := database.GetStateFromDisk()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer state.Close()

	block0 := database.CreateNewBlock(
		database.Hash{},
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("cuong", "cuong", 3, ""),
			database.NewTx("cuong", "cuong", 700, "reward"),
		},
	)

	state.AddBlock(block0)
	block0hash, _ := state.Persist()

	block1 := database.CreateNewBlock(
		block0hash,
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("cuong", "babayaga", 2000, ""),
			database.NewTx("cuong", "cuong", 100, "reward"),
			database.NewTx("babayaga", "cuong", 1, ""),
			database.NewTx("babayaga", "caesar", 1000, ""),
			database.NewTx("babayaga", "cuong", 50, ""),
			database.NewTx("cuong", "cuong", 600, "reward"),
		},
	)

	state.AddBlock(block1)
	state.Persist()
}
