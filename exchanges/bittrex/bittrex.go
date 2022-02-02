package bittrex

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"main/database/mongodb"
	"net/http"
	"strings"
)

func InsertPairData() {
	res, err := http.Get("https://api.bittrex.com/v3/markets")
	if err == nil {
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Fatal(err)
		}

		var response []PairData
		if err := json.Unmarshal(data, &response); err != nil {
			log.Fatal(err)
		}

		client := mongodb.GetClient()
		coll := client.Database("cryptohamster").Collection("pairs")

		for i := 0; i <= len(response)-1; i++ {
			filter := bson.M{"symbol": strings.Replace(response[i].Symbol, "-", "", -1)}
			update := bson.M{"$set": bson.M{"baseAsset": response[i].BaseCurrencySymbol}}
			opts := options.Update().SetUpsert(true)
			coll.UpdateOne(context.TODO(), filter, update, opts)
		}
	}
}

func GetPrices() []PriceData {
	res, err := http.Get("https://api.bittrex.com/v3/markets/tickers")
	if err == nil {
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Fatal(err)
		}

		var response []PriceData
		if err := json.Unmarshal(data, &response); err != nil {
			log.Fatal(err)
		}
		for i := 0; i <= len(response)-1; i++ {
			response[i].Symbol = strings.Replace(response[i].Symbol, "-", "", -1)
		}
		return response
	} else {
		return nil
	}
}

type PairData struct {
	Symbol             string `json:"symbol"`
	BaseCurrencySymbol string `json:"baseCurrencySymbol"`
}

type PriceData struct {
	Symbol string `json:"symbol"`
	Price  string `json:"bidRate"`
}
