package rpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
)

func GetMaxAvailable(host string, addr string, currency string) (amount *big.Int) {
	resp, e := doPost(host, "exchange_getMaxAvailable", []interface{}{addr, currency})
	if e != nil {
		return
	}
	if resp.Result == nil {
		return
	}
	message := *resp.Result
	amount, _ = new(big.Int).SetString(string(message), 10)
	return
}

type BuyShareTxArg struct {
	From     string `json:"from"`
	Vote     string `json:"vote"`
	Pool     string `json:"pool"`
	Value    string `json:"value"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
}

type RegistStakePoolTxArg struct {
	From     string `json:"from"`
	Vote     string `json:"vote"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
	Fee      string `json:"fee"`
}

func GasPrice(host string) (gasPrice *big.Int) {
	resp, e := doPost(host, "sero_gasPrice", []interface{}{})
	if e != nil {
		fmt.Println(e)
		return big.NewInt(25000000000000)
	}

	var s string
	err := json.Unmarshal(*resp.Result, &s)
	if err != nil {
		fmt.Println(err)
		return
	}
	datas, e := hex.DecodeString(s[2:])
	gasPrice = new(big.Int).SetBytes(datas)
	return
}

func RegistStakePool(host string, arg RegistStakePoolTxArg) {
	resp, e := doPost(host, "stake_registStakePool", []interface{}{arg})
	if e != nil {
		fmt.Println(e)
		return
	}
	message := *resp.Result
	fmt.Println(string(message))
	return
}

func BuyShare(host string, arg BuyShareTxArg) (txHash string, err error) {
	resp, err := doPost(host, "stake_buyShare", []interface{}{arg})
	if err != nil {
		log.Printf("err : %v", err)
		return
	}
	message := *resp.Result
	txHash = string(message)
	return
}

func CurrentPrice(host string) (sharePrice *big.Int) {
	resp, e := doPost(host, "stake_sharePrice", []interface{}{})
	if e != nil {
		fmt.Println(e)
		return big.NewInt(2000000000000000000)
	}

	var s string
	err := json.Unmarshal(*resp.Result, &s)
	if err != nil {
		fmt.Println(err)
		return
	}
	datas, e := hex.DecodeString(s[2:])
	sharePrice = new(big.Int).SetBytes(datas)
	return
}

type JSONRpcResp struct {
	Id     *json.RawMessage       `json:"id"`
	Result *json.RawMessage       `json:"result"`
	Error  map[string]interface{} `json:"error"`
}

func doPost(url string, method string, params interface{}) (*JSONRpcResp, error) {
	client := &http.Client{
		Timeout: MustParseDuration("900s"),
	}
	jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	data, err := json.Marshal(jsonReq)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp *JSONRpcResp
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {

		return nil, err
	}
	if rpcResp.Error != nil {
		return nil, errors.New(rpcResp.Error["message"].(string))

	}
	data, _ = json.Marshal(rpcResp)
	return rpcResp, err
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}
