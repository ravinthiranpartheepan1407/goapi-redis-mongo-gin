package main

import (
	"context"
	"fmt"
	"goapi-ravi/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client
	redisClient *redis.Client
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load env variables", err)
	}

	ctx = context.TODO()

	mongoconnect := options.Client().ApplyURI(config.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoconnect)

	if err != nil {
		panic(err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Mongodb successfully connected")

	redisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisClient.Set(ctx, "test", "Connected GO App with Redis and MongoDB", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis Server Connected")

	server = gin.Default()

}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could load en config", err)
	}

	defer mongoClient.Disconnect(ctx)

	value, err := redisClient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("Key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	router := server.Group("/api")
	router.GET("/check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	log.Fatal(server.Run(":" + config.Port))
}
