package database

import(
	"io/ioutil"
	"encoding/json"
)

type Genesis struct{
	Balances map[Account]uint32 `json:"balances"`
}

func LoadGenesis(path string) (Genesis, error){
	content, err := ioutil.ReadFile(path)
	if err != nil{
		return Genesis{}, err
	}

	var g Genesis;
	err = json.Unmarshal(content, &g)
	if err != nil{
		return Genesis{}, err
	}
	return g, nil
}