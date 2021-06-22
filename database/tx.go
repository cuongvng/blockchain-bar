package database

type Account string

type Tx struct{
	From Account `json:"from"`
	To Account `json:"to"`
	Value uint32 `json:"value"`
	Data string `json:"data"`
}

func NewTx(from Account, to Account, value uint32, data string) Tx{
	return Tx{from, to , value, data}
}

func (tx Tx) IsReward() bool{
	return tx.Data == "reward"
}