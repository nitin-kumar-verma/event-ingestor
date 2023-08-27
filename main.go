package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Event struct {
	UserID    string `json:"userId"`
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
}

const (
	EVENT_KEY      = "queued_events"
	PROCESSING_KEY = "processing_events"
)

func main() {
	app := fiber.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client := redisClient()
	defer client.Close()
	app.Post("/event", func(c *fiber.Ctx) error {
		var event Event
		if err := c.BodyParser(&event); err != nil {
			return err
		}
		if event.UserID == "" || event.Payload == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid Json")
		}
		event.Timestamp = time.Now().UTC().Unix()
		eventJSON, _ := json.Marshal(event)

		ctx := context.Background()

		//Lpush and rpop to implement queue
		err = client.LPush(ctx, EVENT_KEY, string(eventJSON)).Err()
		if err != nil {
			log.Println(err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Event not accepted, please retry")
		}

		return c.SendStatus(fiber.StatusAccepted)
	})

	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

func redisClient() *redis.Client {

	redisAddr, redisAddrExists := os.LookupEnv("REDIS_ADDR")
	redisPassword, redisPassExists := os.LookupEnv("REDIS_PASSWORD")
	redisDB, redisDBExists := os.LookupEnv("REDIS_DB")
	if !(redisAddrExists || redisPassExists || redisDBExists) {
		panic("redis envs not set properly")
	}

	db, err := strconv.Atoi(redisDB)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       db,
	})
	return client
}
