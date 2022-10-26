package main

import (
	"context"
	"log"
	"time"

	"os"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-contrib/cors"
)

var MCLI *mongo.Client
var MCTX context.Context
var MCAN context.CancelFunc

const (
	connectTimeout           = 20
	connectionStringTemplate = "mongodb://%s:%s@%s:%s/%s"
)

func main() {

	value, ok := os.LookupEnv("MONGO_SERVER")

	if !ok {
		log.Fatal("No Environment Variable")
		os.Exit(1)
	} else {
		log.Printf("MONGO_SERVER: %s", value)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/health", watcherHealthCheck)
	router.GET("/profiles", watcherReadProfile)
	router.GET("/db/populate", watcherPopulateDB)
	router.GET("/close", watcherCloseDB)
	router.GET("/clearDB", watcherClearDB)
	router.POST("/addProfiles", watcherAddProfiles)

	router.Run(":5000")
}
