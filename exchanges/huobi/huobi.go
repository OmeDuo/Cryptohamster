package huobi

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main/database/mongodb"
	"net/http"
)

func InsertPairIndex() {
	res, err := http.Get("https://api.huobi.pro/v1/common/symbols")
	if err == nil {
		defer res.Body.Close()
		response := new(PairData)
		json.NewDecoder(res.Body).Decode(response)

		client := mongodb.GetClient()
		coll := client.Database("cryptohamster").Collection("pairs")

		for i := 0; i <= len(response.Symbols)-1; i++ {
			filter := bson.M{"symbol": response.Symbols[i].Symbol}
			update := bson.M{"$set": bson.M{"baseAsset": response.Symbols[i].BaseAsset}}
			opts := options.Update().SetUpsert(true)
			coll.UpdateOne(context.TODO(), filter, update, opts)
		}
	}
}

func GetPrices() Ticker {
	res, err := http.Get("https://api.huobi.pro/market/tickers")
	if err == nil {
		defer res.Body.Close()
		response := new(PriceData)
		json.NewDecoder(res.Body).Decode(response)
		for i := 0; i <= len(response.Ticker)-1; i++ {
			response.Ticker[i].StringPrice = fmt.Sprintf("%f", response.Ticker[i].Price)
		}
		return response.Ticker
	} else {
		return nil
	}
}

type PairData struct {
	Symbols []struct {
		Symbol    string `json:"symbol"`
		BaseAsset string `json:"base-currency"`
	} `json:"data"`
}

type PriceData struct {
	Ticker Ticker `json:"data"`
}

type Ticker []struct {
	Symbol      string  `json:"symbol"`
	Price       float64 `json:"bid"`
	StringPrice string
}
