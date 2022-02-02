package binance

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"main/database/mongodb"
	"net/http"
)

func InsertPairIndex() {
	res, err := http.Get("https://api.binance.com/api/v3/exchangeInfo")
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

func GetPrices() []PriceData {
	res, err := http.Get("https://api.binance.com/api/v3/ticker/price")
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
		return response
	} else {
		return nil
	}
}

type PairData struct {
	Symbols []struct {
		Symbol    string `json:"symbol"`
		BaseAsset string `json:"baseAsset"`
	} `json:"symbols"`
}

type PriceData struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type coinIndex struct {
	Symbol    string `bson:"symbol"`
	BaseAsset string `bson:"baseAsset"`
	Status    bool
}
