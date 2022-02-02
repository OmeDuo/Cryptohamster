package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"os"
)

var client mongo.Client

func init() {
	for {
		if _, err := os.Stat("./connection_config.json"); err != nil {
			if os.IsNotExist(err) {
				createConfig()
			} else {
				os.Exit(1)
			}
		} else {
			// login
			var input string
			fmt.Print("Use previous config? (y/n default=yes): ")
			fmt.Scan(&input)

			if input == "n" {
				createConfig()
			}
		}

		if SetAndCheckConnection() == true {
			break
		}
	}
	fmt.Println("Connection established")
}

func SetAndCheckConnection() (status bool) {
	data := readConfig()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+data.Credentials.Username+":"+data.Credentials.Password+"@"+data.IPAddress+":"+data.Port+"/?authSource=admin"))
	if err != nil {
		fmt.Println("IP settings or login credentials incorrect!")
		client = nil
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println("IP settings or login credentials incorrect!")
		client = nil
		return false
	} else {
		return true
	}
}

func GetClient() (client *mongo.Client) {
	data := readConfig()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+data.Credentials.Username+":"+data.Credentials.Password+"@"+data.IPAddress+":"+data.Port+"/?authSource=admin"))
	if err != nil {
		fmt.Println("IP settings or login credentials incorrect!")
		client = nil
	}
	return client
}

func readConfig() (data Mongodb) {
	file, _ := ioutil.ReadFile("connection_config.json")
	data = Mongodb{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func createConfig() {
	var ip string
	var port string
	var username string
	var password string

	fmt.Print("Enter MongoDB ip: ")
	fmt.Scan(&ip)
	fmt.Print("Enter MongoDB port: ")
	fmt.Scan(&port)
	fmt.Print("Enter MongoDB user: ")
	fmt.Scan(&username)
	fmt.Print("Enter MongoDB password: ")
	fmt.Scan(&password)

	db := Mongodb{
		IPAddress: ip,
		Port:      port,
		Credentials: struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}(struct {
			Username string
			Password string
		}{Username: username, Password: password}),
	}

	Path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	file, _ := json.MarshalIndent(db, "", " ")
	_ = ioutil.WriteFile(Path+"\\"+"connection_config"+".json", file, 0644)
}

type Mongodb struct {
	IPAddress   string `json:"ipAddress"`
	Port        string `json:"port"`
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"credentials"`
}
