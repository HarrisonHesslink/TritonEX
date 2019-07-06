package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
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
			iter := client.Collection("trades").Documents(ctx)
			var trades []Trade
			for {
				doc, err := iter.Next()
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
				fmt.Println(p)
				fmt.Println(a)
				newTrade := Trade{trtype, time, a, p}
				trades = append(trades, newTrade)
			}

			c.IndentedJSON(http.StatusOK, trades)

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
