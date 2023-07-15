package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
  webPort = "80"
  rpcPort = "5001"
  mongoURL = "mongodb://mongo:27017"
  gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
  Models data.Models
}

func main() {
  mongoClient, err := connectToMongo()
  if err != nil {
    log.Panic(err)
  }

  client = mongoClient

  // context required to disconnect
  ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
  defer cancel()

  // close connection
  defer func() {
    if err = client.Disconnect(ctx); err != nil {
      panic(err)
    }
  } ()

  app := Config{
    Models: data.New(mongoClient),
  }

  srv := &http.Server{
    Addr: fmt.Sprintf(":%s", webPort),
    Handler: app.routes(),
  }

  err = srv.ListenAndServe()
  if err != nil {
    log.Panic("Cannot connect")
  }
  // start server
  // app.serve()
}

// func (app *Config) serve() { 
//   srv := &http.Server{
//     Addr: fmt.Sprintf(":%s", webPort),
//     Handler: app.routes(),
//   }

//   err := srv.ListenAndServe()
//   if err != nil {
//     log.Panic()
//   }
// }

func connectToMongo() (*mongo.Client, error) {
  clientOptions := options.Client().ApplyURI(mongoURL)
  
  // TODO: replace with env variable
  clientOptions.SetAuth(options.Credential{
    Username: "admin",
    Password: "password", 
  })


  c, err := mongo.Connect(context.TODO(), clientOptions)
  if err != nil {
    log.Println("Error connecting:", err)
    return nil, err
  } 

  log.Println("Connected!")
  
  return c, nil
}