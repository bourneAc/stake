package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sero-cash/stake/rpc"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Host   string
	PoolPK string
	Vote   string
	Pool   string
	Fee    int64

	PK       string
	SoloVote string
	Period   string
}

var (
	basePrice = big.NewInt(2000000000000000000) //2 SERO
	addition  = big.NewInt(759357240838722)

	oneSero        = big.NewInt(1000000000000000000)
	towSero        = big.NewInt(2000000000000000000)
	twentyWSero, _ = new(big.Int).SetString("200000000000000000000000", 10)
	gas            = big.NewInt(25000)
	gasPrice       = big.NewInt(1000000000)
)

func encodeBig(value *big.Int) string {
	s := hex.EncodeToString(value.Bytes())
	var index int
	for i := 0; i < len(s); i++ {
		if s[i] != '0' {
			index = i
			break
		}
	}
	return "0x" + s[index:]
}

func sum(basePrice, addition *big.Int, n int64) *big.Int {
	return new(big.Int).Add(
		new(big.Int).Mul(basePrice, big.NewInt(n)),
		new(big.Int).Div(
			new(big.Int).Mul(
				new(big.Int).Mul(big.NewInt(n), big.NewInt(n-1)),
				addition,
			),
			big.NewInt(2),
		),
	)
}

func main() {
	fmt.Println(encodeBig(gas))
	var cfg Config
	readConfig(&cfg)

	//fmt.Println(rpc.CurrentPrice(cfg.Host))
	//fmt.Println(sum(rpc.CurrentPrice(cfg.Host), addition, 1000))

	//rpc.RegistStakePool(cfg.Host, rpc.RegistStakePoolTxArg{From: cfg.PoolPK, Value: encodeBig(twentyWSero), Vote: cfg.Vote, Gas: encodeBig(gas), GasPrice: encodeBig(gasPrice), Fee: encodeBig(big.NewInt(cfg.Fee))})

	period := mustParseDuration(cfg.Period)
	for {

		price := rpc.GasPrice(cfg.Host)
		if price.Sign() == 0 {
			price = gasPrice
		}
		price.Add(price, big.NewInt(2))

		availAmount := rpc.GetMaxAvailable(cfg.Host, cfg.PK, "SERO")
		sum := sum(rpc.CurrentPrice(cfg.Host), addition, 1500)

		var amount *big.Int
		if availAmount.Cmp(sum) > 0 {
			amount = sum
		} else {
			amount = availAmount
		}

		if amount.Cmp(towSero) > 0 {
			amount.Sub(amount, towSero)
			tx, _ := rpc.BuyShare(cfg.Host, rpc.BuyShareTxArg{From: cfg.PK, Value: encodeBig(amount), Vote: cfg.SoloVote, Pool: cfg.Pool, Gas: encodeBig(gas), GasPrice: encodeBig(price)})
			log.Printf("txhash : %s, value : %v", tx, amount.Text(10))
		} else {
			log.Printf("amount : %v", amount.Text(10))
		}

		time.Sleep(period)
	}
}

func readConfig(cfg *Config) {

	configFileName := "config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Loading config: %s", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		log.Fatal("Config error: ", err.Error())
	}

}

func mustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}
