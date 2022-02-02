package pricecontroller

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"main/database/mongodb"
	"main/exchanges/binance"
	"main/exchanges/bittrex"
	"main/exchanges/huobi"
	"main/exchanges/kucoin"
	"time"
)

var coinsIndex []CoinIndex

func Start() {
	
	go binance.InsertPairIndex()
	go kucoin.InsertPairIndex()
	go bittrex.InsertPairData()
	go huobi.InsertPairIndex()

	time.Sleep(5000 * time.Millisecond)
	fmt.Println("PairIndex updated!")

	GetIndex()

	InsertBinancePriceData()
	InsertKucoinPriceData()
	InsertHuobiPriceData()
	InsertBittrexPriceData()

	RemoveDeadCoin()

	time.Sleep(5000 * time.Millisecond)
	fmt.Println("Prices updated!")
}

func GetIndex() {
	client := mongodb.GetClient()
	coll := client.Database("cryptohamster").Collection("pairs")

	cur, err := coll.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.Background()) {
		var coinIndex CoinIndex
		err := cur.Decode(&coinIndex)
		if err != nil {
			log.Fatal(err)
		}
		coinsIndex = append(coinsIndex, coinIndex)
	}
	defer cur.Close(context.Background())
}

func InsertBinancePriceData() {
	pricedata := binance.GetPrices()

	for i := 0; i <= len(pricedata)-1; i++ {
		for j := 0; j <= len(coinsIndex)-1; j++ {

			if pricedata[i].Symbol == coinsIndex[j].Symbol {
				coinsIndex[j].Status = true
			}
		}
	}
}

func InsertKucoinPriceData() {
	pricedata := kucoin.GetPrices()

	for i := 0; i <= len(pricedata)-1; i++ {
		for j := 0; j <= len(coinsIndex)-1; j++ {

			if pricedata[i].Symbol == coinsIndex[j].Symbol {
				coinsIndex[j].Status = true
			}
		}
	}
}

func InsertHuobiPriceData() {
	pricedata := huobi.GetPrices()

	for i := 0; i <= len(pricedata)-1; i++ {
		for j := 0; j <= len(coinsIndex)-1; j++ {

			if pricedata[i].Symbol == coinsIndex[j].Symbol {
				coinsIndex[j].Status = true
			}
		}
	}
}

func RemoveDeadCoin() {
	client := mongodb.GetClient()
	coll := client.Database("cryptohamster").Collection("pairs")
	for i := 0; i <= len(coinsIndex)-1; i++ {
		if coinsIndex[i].Status == false {
			coll.DeleteOne(context.TODO(), bson.M{"_id": coinsIndex[i].Id})
		}
	}
}

func InsertBittrexPriceData() {
	pricedata := bittrex.GetPrices()

	for i := 0; i <= len(pricedata)-1; i++ {
		for j := 0; j <= len(coinsIndex)-1; j++ {

			if pricedata[i].Symbol == coinsIndex[j].Symbol {
				coinsIndex[j].Status = true
			}
		}
	}
}

type CoinIndex struct {
	Id        primitive.ObjectID `bson:"_id"`
	Symbol    string             `bson:"symbol"`
	BaseAsset string             `bson:"baseAsset"`
	Status    bool
}
