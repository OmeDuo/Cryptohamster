package kucoin

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main/database/mongodb"
	"net/http"
	"strings"
)

func InsertPairIndex() {
	res, err := http.Get("https://api.kucoin.com/api/v1/symbols")
	if err == nil {

		defer res.Body.Close()
		response := new(PairData)
		json.NewDecoder(res.Body).Decode(response)

		client := mongodb.GetClient()
		coll := client.Database("cryptohamster").Collection("pairs")

		for i := 0; i <= len(response.Data)-1; i++ {
			filter := bson.M{"symbol": strings.Replace(response.Data[i].Symbol, "-", "", -1)}
			update := bson.M{"$set": bson.M{"baseAsset": response.Data[i].BaseCurrency}}
			opts := options.Update().SetUpsert(true)
			coll.UpdateOne(context.TODO(), filter, update, opts)
		}

	}
}

func GetPrices() Ticker {
	res, err := http.Get("https://api.kucoin.com/api/v1/market/allTickers")
	if err == nil {

		defer res.Body.Close()
		response := new(PriceData)
		json.NewDecoder(res.Body).Decode(response)

		for i := 0; i <= len(response.Data.Ticker)-1; i++ {
			response.Data.Ticker[i].Symbol = strings.Replace(response.Data.Ticker[i].Symbol, "-", "", -1)
		}
		return response.Data.Ticker
	} else {
		return nil
	}
}

type PairData struct {
	Data []struct {
		Symbol       string `json:"symbol"`
		BaseCurrency string `json:"baseCurrency"`
	} `json:"data"`
}

type PriceData struct {
	Data struct {
		Ticker Ticker `json:"ticker"`
	} `json:"data"`
}

type Ticker []struct {
	Symbol string `json:"symbol"`
	Price  string `json:"buy"`
}
