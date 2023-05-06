package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
)

var (
	genesis = new(core.Genesis)
	lc      = &logger.Config{
		EnableMemory:     true,
		DisableStack:     false,
		DisableStorage:   false,
		EnableReturnData: true,
		Debug:            true,
	}
	memDB = state.NewDatabaseWithConfig(rawdb.NewMemoryDatabase(), &trie.Config{Preimages: false})
	s, _  = state.New(common.Hash{}, memDB, nil)
	rc    = runtime.Config{
		State:       s,
		Difficulty:  genesis.Difficulty,
		Time:        genesis.Timestamp,
		Coinbase:    genesis.Coinbase,
		BlockNumber: new(big.Int).SetUint64(genesis.Number),
		EVMConfig: vm.Config{
			Tracer: logger.NewJSONLogger(lc, os.Stdout),
		},
		ChainConfig: params.AllEthashProtocolChanges,
	}
)

type CreateAccountMsg struct {
	Balance int64 `json:"balance"`
}

func HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["status"] = "success"
	response["message"] = "account created"

	var msg CreateAccountMsg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(fmt.Sprintf("%+v", msg))

	addr := common.BigToAddress(big.NewInt(1000))
	s.CreateAccount(addr)
	s.AddBalance(addr, big.NewInt(msg.Balance))

	response["address"] = fmt.Sprintf("0x%x", addr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

type CreateContractMsg struct {
	From  string `json:"from"`
	Value int64  `json:"value"`
	Input string `json:"input"`
	Code  string `json:"code"`
}

func HandleCreateContract(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["status"] = "success"
	response["message"] = "contract created"

	var msg CreateContractMsg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(fmt.Sprintf("%+v", msg))

	rc.GasPrice = big.NewInt(1)
	rc.GasLimit = 100000000
	rc.Origin = common.HexToAddress(msg.From)
	rc.Value = big.NewInt(msg.Value)

	input := common.FromHex(msg.Input)
	code := common.FromHex(msg.Code)

	input = append(code, input...)
	output, addr, gasLeft, err := runtime.Create(input, &rc)

	response["output"] = fmt.Sprintf("%x", output)
	response["address"] = fmt.Sprintf("0x%v", addr)
	response["gasLeft"] = fmt.Sprintf("%v", gasLeft)
	if err != nil {
		response["status"] = "error"
		response["message"] = fmt.Sprintf("contract creation failed, err: %s", err.Error())
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

type CallContractMsg struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int64  `json:"value"`
	Input string `json:"input"`
}

func HandleCallContract(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["status"] = "success"
	response["message"] = "contract called"

	var msg CallContractMsg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(fmt.Sprintf("%+v", msg))

	rc.GasPrice = big.NewInt(1)
	rc.GasLimit = 100000000
	rc.Origin = common.HexToAddress(msg.From)
	rc.Value = big.NewInt(msg.Value)

	input := common.FromHex(msg.Input)
	receiver := common.HexToAddress(msg.To)

	output, gasLeft, err := runtime.Call(receiver, input, &rc)

	response["output"] = fmt.Sprintf("%x", output)
	response["gasLeft"] = fmt.Sprintf("%v", gasLeft)
	if err != nil {
		response["status"] = "error"
		response["message"] = fmt.Sprintf("contract call failed, err: %s", err.Error())
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
