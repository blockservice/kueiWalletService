package ews

import (
	"encoding/json"
	"github.com/ChungkueiBlock/kueiWalletService/internal/log"
	"math"
	"net/http"
	"time"
)

var (
	reties      = 0
	predictJson []predictTableItem
	ethgasJson  *ethgas
)

type GasEstimateResult struct {
	TxMeanSecs float64 `json:"txMeanSecs"`
	Gasprice   float64 `json:"gasprice"`
}

func (api *WalletAPI) EstimateGasprice(speed string) GasEstimateResult {
	if ethgasJson == nil {
		return GasEstimateResult{-1, -1}
	}

	switch speed {
	case "fast":
		return GasEstimateResult{ethgasJson.FastWait * 60, ethgasJson.Fast}
	case "middle":
		return GasEstimateResult{ethgasJson.AvgWait * 60, ethgasJson.Avg}
	case "slow":
		return GasEstimateResult{ethgasJson.SlowWait * 60, ethgasJson.Slow}
	}
	return GasEstimateResult{ethgasJson.AvgWait * 60, ethgasJson.Avg}
}

type TxEstimateResult struct {
	TxMeanSecs float64 `json:"txMeanSecs"`
	MinedProb  string  `json:"minedProb"`
}

func (api *WalletAPI) EstimateTxTime(gasprice, gasoffered float64) TxEstimateResult {
	if predictJson == nil || len(predictJson) == 0 {
		return TxEstimateResult{-1, "no results"}
	}

	if gasprice < predictJson[1].Gasprice {
		gasprice = predictJson[1].Gasprice
	} else if gasprice > predictJson[len(predictJson)-1].Gasprice {
		gasprice = predictJson[len(predictJson)-1].Gasprice
	}

	i := 0
	for ; i < len(predictJson); i++ {
		if predictJson[i].Gasprice == gasprice {
			break
		}
	}

	intercept := 4.2794
	hpa := 0.0329
	hgo := -3.2836
	//wb  := -0.0048
	tx := -0.0004

	intercept2 := 7.5375
	hpa_coef := -0.0801
	txatabove_coef := 0.0003
	high_gas_coef := 0.3532

	sum1 := 0.0
	sum2 := 0.0
	if gasoffered > 1000000 {
		sum1 = intercept + (predictJson[i].HashpowerAccepting * hpa) + (predictJson[i].TxAtabove * tx) + hgo
		sum2 = intercept2 + (predictJson[i].HashpowerAccepting * hpa_coef) + (predictJson[i].TxAtabove * txatabove_coef) + high_gas_coef
	} else {
		sum1 = intercept + (predictJson[i].HashpowerAccepting * hpa) + (predictJson[i].TxAtabove * tx)
		sum2 = intercept2 + (predictJson[i].TxAtabove * txatabove_coef) + (predictJson[i].HashpowerAccepting * hpa_coef)
	}

	minedProb := ""
	factor := math.Exp(-1 * sum1)
	prob := 1 / (1 + factor)
	if prob > 0.95 {
		minedProb = "Very High"
	} else if prob > 0.9 && prob <= 0.95 {
		minedProb = "Medium"
	} else {
		minedProb = "Low"
	}

	expectedWait := math.Exp(sum2)
	if expectedWait < 2 {
		expectedWait = 2
	}

	if gasoffered > 2000000 {
		expectedWait += 100
	}

	// Mean Time to Confirm (Blocks)
	blocksWait := expectedWait
	//% of last 200 blocks accepting this gas price
	//hashpower := predictJson[i].HashpowerAccepting
	//Transactions At or Above in Current Txpool
	//txatabove := predictJson[i].TxAtabove
	//minedprob
	blockInterval := 15.586734693878 //@TODO block平均挖矿时间如何获得?
	txMeanSecs := blocksWait * blockInterval

	return TxEstimateResult{txMeanSecs, minedProb}
}

// Crawling Thread

type predictTableItem struct {
	Gasprice           float64 `json:"gasprice"`
	HashpowerAccepting float64 `json:"hashpower_accepting"`
	TxAtabove          float64 `json:"tx_atabove"`
}
type ethgas struct {
	Fast     float64 `json:"fast"`
	FastWait float64 `json:"fast"`
	Avg      float64 `json:"average"`
	AvgWait  float64 `json:"avgWait"`
	Slow     float64 `json:"safelow"`
	SlowWait float64 `json:"safeLowWait"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func updatePredictTableJson() {
	isloaded := false
	for {
		results := make([]predictTableItem, 0)
		err := fetchPredictJson(&results)
		if err != nil {
			log.Error("fetch predictTableJson fail", "err", err.Error())
		} else {
			log.Info("fetch predictTableJson successfully!")
			isloaded = true
			predictJson = results
		}

		results2 := &ethgas{}
		err = fetchEthgasJson(results2)
		if err != nil {
			isloaded = false
			log.Error("fetch ethgasAPI fail", "err", err.Error())
		} else {
			log.Info("fetch ethgasAPI successfully!")
			ethgasJson = results2
		}

		if isloaded {
			time.Sleep(time.Minute * 30)
		} else {
			reties++
			log.Info("fetch json failed", "retry", reties)
			time.Sleep(time.Second * 30)
		}
	}
}

func fetchPredictJson(results *[]predictTableItem) error {
	r, err := myClient.Get("https://ethgasstation.info/json/predictTable.json")
	if err != nil {
		return err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&results)
	return nil
}

func fetchEthgasJson(results *ethgas) error {
	r, err := myClient.Get("https://ethgasstation.info/json/ethgasAPI.json")
	if err != nil {
		return err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&results)
	return nil
}
