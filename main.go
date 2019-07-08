package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

//Trade
type Trade struct {
	TradeType string `json:"TradeType"`
	TimeStamp string `json:"TimeStamp"`
	Amount    string `json:"Amount"`
	Price     string `json:"Price"`
}

//Order
type Order struct {
	TradeType string `json: "TradeType"`
	TimeStamp string `json: "TimeStamp"`
	Amount    string `json: "Amount"`
	Price     string `json: "Price"`
	Filled    string `json: Filled`
}

//Message
type rawMessageData1 struct {
	message string
}

//Message
type rawMessageData struct {
	message string
	data    json.RawMessage
}

// Json of new Order for buy or sell
type orderJson struct {
	Tradetype string `json: "trade_type"`
	Timestamp string `json: "timestamp"`
	Amount    string `json: "amount"`
	Price     string `json: "price"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("recv: %s", msg)
		conn.WriteMessage(t, msg)
		raw := rawMessageData1{}
		errJSON := conn.ReadJSON(&raw)
		if errJSON != nil {
			fmt.Println("Error reading json.", err)
		}
		if raw.message != "" {
			fmt.Printf("Got message: %#v\n", raw)
			var messages []rawMessageData
			r := json.Unmarshal([]byte(raw.message), &messages)
			if r != nil {
				log.Fatalln("error:", err)
			}
			fmt.Println()
		}

	}
}

func main() {

	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("credentials.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	router := gin.Default()
	router.Use(static.Serve("/assets", static.LocalFile("./assets", true)))
	router.LoadHTMLGlob("views/*")

	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		api.POST("/buy", func(c *gin.Context) {
			newTade := Trade{c.PostForm("trade_type"), c.PostForm("timestamp"), c.PostForm("amount"), c.PostForm("price")}
			_, _, err = client.Collection("trades").Add(ctx, map[string]interface{}{
				"trade_type": newTade.TradeType,
				"timestamp":  newTade.TimeStamp,
				"amount":     newTade.Amount,
				"price":      newTade.Price,
			})

			if err != nil {
				log.Fatalf("Failed adding aturing: %v", err)
			}
		})
		api.POST("/sell", func(c *gin.Context) {
			newTade := Trade{c.PostForm("trade_type"), c.PostForm("timestamp"), c.PostForm("amount"), c.PostForm("price")}
			_, _, err = client.Collection("trades").Add(ctx, map[string]interface{}{
				"trade_type": newTade.TradeType,
				"timestamp":  newTade.TimeStamp,
				"amount":     newTade.Amount,
				"price":      newTade.Price,
			})

			if err != nil {
				log.Fatalf("Failed adding aturing: %v", err)
			}
		})
		api.GET("/get_trades", func(c *gin.Context) {
			iter := client.Collection("trades").OrderBy("timestamp", firestore.Desc)
			iter.Limit(100)
			docum := iter.Documents(ctx)
			var trades []Trade
			for {
				doc, err := docum.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Fatalf("Failed to iterate: %v", err)
				}
				docsnap := doc.Data()
				trtype := docsnap["trade_type"].(string)
				time := docsnap["timestamp"].(string)
				p := docsnap["price"].(string)
				a := docsnap["amount"].(string)
				newTrade := Trade{trtype, time, a, p}
				trades = append(trades, newTrade)
			}

			c.IndentedJSON(http.StatusOK, trades)

		})
		api.GET("/get_orders", func(c *gin.Context) {
			iter := client.Collection("orders").OrderBy("timestamp", firestore.Desc)
			iter.Limit(100)
			docum := iter.Documents(ctx)
			var orders []Order
			for {
				doc, err := docum.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Fatalf("Failed to iterate: %v", err)
				}
				docsnap := doc.Data()
				trtype := docsnap["trade_type"].(string)
				time := docsnap["timestamp"].(string)
				p := docsnap["price"].(string)
				a := docsnap["amount"].(string)
				f := docsnap["filled"].(string)
				newOrder := Order{trtype, time, a, p, f}
				orders = append(orders, newOrder)
			}

			c.IndentedJSON(http.StatusOK, orders)

		})
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "TritonEX | Coming Soon",
			},
		)

	})

	router.Run(":80")

}
