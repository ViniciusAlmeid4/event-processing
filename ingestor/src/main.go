package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"uuid"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	r.POST("/reading", func(c *gin.Context) {
		reading := Reading{Schema_version: 1}

		if err := c.ShouldBindJSON(&reading); err != nil {
			fmt.Println("[INTERNAL] Error connecting, message:", err)
			c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		reading.Id = uuid.NewV7()

		partition, err := strconv.Atoi(os.Getenv("PARTITION"))
		if err != nil {
			fmt.Println("[INTERNAL] Error getting env partition, message:", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		conn, err := connect("processor.reading", partition)
		if err != nil {
			fmt.Println("[INTERNAL] Error connecting to partition, message:", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = writeReading(conn, reading)
		if err != nil {
			fmt.Println("[INTERNAL] Error writing to kafka, message:", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(
			http.StatusAccepted,
			reading,
		)
	})

	r.GET("/health", func(c *gin.Context) {
		var http_status int = http.StatusOK
		status := Health{Status: true, Kafka_connection: true}
		partition, err := strconv.Atoi(os.Getenv("PARTITION"))
		if err != nil {
			status.Kafka_connection = false
			status.Status = false
			http_status = http.StatusServiceUnavailable
		}

		_, err = connect("processor.reading", partition)
		if err != nil {
			status.Kafka_connection = false
			status.Status = false
			http_status = http.StatusServiceUnavailable
		}

		c.IndentedJSON(
			http_status,
			status,
		)
	})

	r.Run(":" + os.Getenv("APP_PORT"))
}
